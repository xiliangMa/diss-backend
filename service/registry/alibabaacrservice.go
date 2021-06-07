package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"regexp"
	"strings"
)

type timeUnix int64

var regRegion = regexp.MustCompile("https://(registry|cr)\\.([\\w\\-]+)\\.aliyuncs\\.com")

const (
	endpointTpl = "cr.%s.aliyuncs.com"
)

type AlibabaACRService struct {
	ImageConfig *models.ImageConfig
}

type aliACRNamespaceResp struct {
	Data struct {
		Namespaces []aliACRNamespace `json:"namespaces"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

type authorizationToken struct {
	Data struct {
		ExpireDate         timeUnix `json:"expireDate"`
		AuthorizationToken string   `json:"authorizationToken"`
		TempUserName       string   `json:"tempUserName"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

type aliACRNamespace struct {
	Namespace       string `json:"namespace"`
	AuthorizeType   string `json:"authorizeType"`
	NamespaceStatus string `json:"namespaceStatus"`
}

type aliReposResp struct {
	Data struct {
		Page     int       `json:"page"`
		Total    int       `json:"total"`
		PageSize int       `json:"pageSize"`
		Repos    []aliRepo `json:"repos"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

type aliRepo struct {
	Summary        string `json:"summary"`
	RegionID       string `json:"regionId"`
	RepoName       string `json:"repoName"`
	RepoNamespace  string `json:"repoNamespace"`
	RepoStatus     string `json:"repoStatus"`
	RepoID         int    `json:"repoId"`
	RepoType       string `json:"repoType"`
	RepoBuildType  string `json:"repoBuildType"`
	GmtCreate      int64  `json:"gmtCreate"`
	RepoOriginType string `json:"repoOriginType"`
	GmtModified    int64  `json:"gmtModified"`
	RepoDomainList struct {
		Internal string `json:"internal"`
		Public   string `json:"public"`
		Vpc      string `json:"vpc"`
	} `json:"repoDomainList"`
	Downloads         int    `json:"downloads"`
	RepoAuthorizeType string `json:"repoAuthorizeType"`
	Logo              string `json:"logo"`
	Stars             int    `json:"stars"`
}

type aliTagResp struct {
	Data struct {
		Total    int         `json:"total"`
		PageSize int         `json:"pageSize"`
		Page     int         `json:"page"`
		Tags     []tagDetail `json:"tags"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

type tagDetail struct {
	ImageUpdate int64  `json:"imageUpdate"`
	ImageID     string `json:"imageId"`
	Digest      string `json:"digest"`
	ImageSize   int    `json:"imageSize"`
	Tag         string `json:"tag"`
	ImageCreate int64  `json:"imageCreate"`
	Status      string `json:"status"`
}

type aliLayers struct {
	Data struct {
		Image struct {
			Layers []struct {
				LayerCMD string `json:"layerCMD"`
			} `json:"layers"`
		} `json:"image"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

type aliManifest struct {
	Data struct {
		Manifest struct {
			FsLayers []struct {
				BlobSum string `json:"blobSum"`
			} `json:"fsLayers"`
		} `json:"manifest"`
	} `json:"data"`
	RequestID string `json:"requestId"`
}

func getRegion(url string) (region string, err error) {
	if url == "" {
		return "", errors.New("empty url")
	}
	rs := regRegion.FindStringSubmatch(url)
	if rs == nil {
		return "", errors.New("invalid registry service url")
	}
	return rs[2], nil
}

func (this *AlibabaACRService) NewAuth(registry *models.Registry) (err error) {
	region, err := getRegion(registry.Url)
	if err != nil {
		return err
	}
	var client *cr.Client
	client, err = cr.NewClientWithAccessKey(region, registry.User, registry.Pwd)
	if err != nil {
		return errors.New("incorrect authentication credentials")
	}
	var tokenRequest = cr.CreateGetAuthorizationTokenRequest()
	domain := fmt.Sprintf(endpointTpl, region)
	tokenRequest.SetDomain(domain)
	tokenResponse, err := client.GetAuthorizationToken(tokenRequest)
	if err != nil {
		return errors.New("incorrect authentication credentials")
	}
	var v authorizationToken
	json.Unmarshal(tokenResponse.GetHttpContentBytes(), &v)
	return nil
}

func (this *AlibabaACRService) getClient(url string, user string, pwd string) (domain string, c *cr.Client, err error) {
	region, _ := getRegion(url)
	var client *cr.Client
	client, err = cr.NewClientWithAccessKey(region, user, pwd)
	d := fmt.Sprintf(endpointTpl, region)
	return d, client, err

}

func (this *AlibabaACRService) Imports() (err error) {
	domain, client, err := this.getClient(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return
	}
	var repositories []aliRepo
	if this.ImageConfig.Namespaces != "" {
		repos, e := this.listReposByNamespace(domain, this.ImageConfig.Namespaces, client)
		if e != nil {
			return
		}

		for _, repo := range repos {
			repositories = append(repositories, repo)
		}
	} else {
		namespaces, e := this.listNamespaces(domain, client)
		if e != nil {
			return
		}
		logs.Info("got namespaces: %v", namespaces)

		for _, ns := range namespaces {
			var repos []aliRepo
			repos, err = this.listReposByNamespace(domain, ns, client)
			if err != nil {
				return
			}

			for _, repo := range repos {
				repositories = append(repositories, repo)
			}
		}
	}

	for _, r := range repositories {
		repo := r

		tags, e := this.getTags(domain, repo, client)
		if e != nil {
			return fmt.Errorf("List tags for repo '%s' error: %v", repo.RepoName, err)
		}

		for _, tag := range tags {
			public := repo.RepoDomainList.Public
			this.ImageConfig.ImageId = "sha256:" + tag.ImageID
			this.ImageConfig.Name = public + "/" + repo.RepoNamespace + "/" + repo.RepoName + ":" + tag.Tag
			if ic := this.ImageConfig.Get(); ic == nil {
				this.ImageConfig.Id = ""
				this.ImageConfig.Size = utils.FormatFileSize(int64(tag.ImageSize))
				this.ImageConfig.CreateTime = tag.ImageCreate * 1e6

				this.ImageConfig.Add()

				imageDetail := models.ImageDetail{}

				imageDetail.ImageId = this.ImageConfig.ImageId
				imageDetail.Name = this.ImageConfig.Name

				if imd := imageDetail.Get(); imd == nil {

					layer, l := this.getImageLayer(domain, repo, tag.Tag, client)
					if l != nil {
						return fmt.Errorf("getImageLayer error: %v", err)
					}

					manifest, m := this.getImageManifest(domain, repo, tag.Tag, client)
					if m != nil {
						return fmt.Errorf("getImageManifest error: %v", err)
					}

					array := make([]string, len(layer.Data.Image.Layers))
					for i, v := range layer.Data.Image.Layers {
						if v.LayerCMD != "" {
							array[i] = v.LayerCMD
						}
					}

					lenx := len(layer.Data.Image.Layers)
					var buffer bytes.Buffer
					re := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
					for i := 0; i < lenx-1; i++ {
						j := lenx - (i + 1)
						if array[j] != "" {
							str := re.ReplaceAllString(array[j], "")
							buffer.WriteString(str + "\n")
						}
					}
					imageDetail.Layers = len(manifest.Data.Manifest.FsLayers)
					imageDetail.Dockerfile = strings.TrimSpace(buffer.String())
					imageDetail.RepoDigests = tag.Digest
					imageDetail.CreateTime = tag.ImageCreate * 1e6
					imageDetail.Size = this.ImageConfig.Size
					if result := imageDetail.Add(); result.Code != http.StatusOK {
						logs.Error("ImageDetail err: %s", errors.New(result.Message))
						return errors.New(result.Message)
					}
				}
			}
		}
	}

	return nil
}

func (this *AlibabaACRService) listNamespaces(domain string, c *cr.Client) (namespaces []string, err error) {

	var nsReq = cr.CreateGetNamespaceListRequest()
	var nsResp = cr.CreateGetNamespaceListResponse()
	nsReq.SetDomain(domain)
	nsResp, err = c.GetNamespaceList(nsReq)
	if err != nil {
		return
	}
	var resp = &aliACRNamespaceResp{}
	err = json.Unmarshal(nsResp.GetHttpContentBytes(), resp)
	if err != nil {
		return
	}
	for _, ns := range resp.Data.Namespaces {
		namespaces = append(namespaces, ns.Namespace)
	}
	return
}

func (this *AlibabaACRService) listReposByNamespace(domain string, namespace string, c *cr.Client) (repos []aliRepo, err error) {
	var reposReq = cr.CreateGetRepoListByNamespaceRequest()
	var reposResp = cr.CreateGetRepoListByNamespaceResponse()
	reposReq.SetDomain(domain)
	reposReq.RepoNamespace = namespace
	var page = 1
	for {
		reposReq.Page = requests.NewInteger(page)
		reposResp, err = c.GetRepoListByNamespace(reposReq)
		if err != nil {
			return
		}
		var resp = &aliReposResp{}
		err = json.Unmarshal(reposResp.GetHttpContentBytes(), resp)
		if err != nil {
			return
		}
		repos = append(repos, resp.Data.Repos...)

		if resp.Data.Total-(resp.Data.Page*resp.Data.PageSize) <= 0 {
			break
		}
		page++
	}
	return
}

func (this *AlibabaACRService) getTags(domain string, repo aliRepo, c *cr.Client) (tags []tagDetail, err error) {
	var tagsReq = cr.CreateGetRepoTagsRequest()
	var tagsResp = cr.CreateGetRepoTagsResponse()
	tagsReq.SetDomain(domain)
	tagsReq.RepoNamespace = repo.RepoNamespace
	tagsReq.RepoName = repo.RepoName
	var page = 1
	for {
		tagsReq.Page = requests.NewInteger(page)
		tagsResp, err = c.GetRepoTags(tagsReq)
		if err != nil {
			return
		}

		var resp = &aliTagResp{}
		json.Unmarshal(tagsResp.GetHttpContentBytes(), resp)
		tags = resp.Data.Tags

		if resp.Data.Total-(resp.Data.Page*resp.Data.PageSize) <= 0 {
			break
		}
		page++
	}

	return
}

func (this *AlibabaACRService) GetNamespaces() models.Result {
	var ResultData models.Result
	domain, client, err := this.getClient(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
	}
	ns, err := this.listNamespaces(domain, client)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = ns
	return ResultData

}

func (this *AlibabaACRService) getImageLayer(domain string, repo aliRepo, tag string, c *cr.Client) (layers *aliLayers, err error) {

	var layerRequest = cr.CreateGetImageLayerRequest()
	var layerResponse = cr.CreateGetImageLayerResponse()
	layerRequest.SetDomain(domain)
	layerRequest.RepoNamespace = repo.RepoNamespace
	layerRequest.RepoName = repo.RepoName
	layerRequest.Tag = tag
	layerResponse, err = c.GetImageLayer(layerRequest)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(layerResponse.GetHttpContentBytes(), &layers)
	return
}

func (this *AlibabaACRService) getImageManifest(domain string, repo aliRepo, tag string, c *cr.Client) (manifest *aliManifest, err error) {

	var manifestRequest = cr.CreateGetImageManifestRequest()
	var manifestResponse = cr.CreateGetImageManifestResponse()
	manifestRequest.SetDomain(domain)
	manifestRequest.RepoNamespace = repo.RepoNamespace
	manifestRequest.RepoName = repo.RepoName
	manifestRequest.Tag = tag
	manifestResponse, err = c.GetImageManifest(manifestRequest)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(manifestResponse.GetHttpContentBytes(), &manifest)
	return
}
