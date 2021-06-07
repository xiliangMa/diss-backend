package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type HuaweiSWRService struct {
	ImageConfig *models.ImageConfig
}

type hwNamespaceList struct {
	Namespace []hwNamespace `json:"namespaces"`
}

type hwNamespace struct {
	ID           int64  `json:"id" orm:"column(id)"`
	Name         string `json:"name"`
	CreatorName  string `json:"creator_name,omitempty"`
	DomainPublic int    `json:"-"`
	Auth         int    `json:"auth"`
	DomainName   string `json:"-"`
	UserCount    int64  `json:"user_count"`
	ImageCount   int64  `json:"image_count"`
}

type Auth struct {
	Identity Identity               `json:"identity"`
	Scope    map[string]interface{} `json:"scope"`
}

type Identity struct {
	Methods  []string `json:"methods"`
	Password password `json:"password"`
}

type password struct {
	User map[string]interface{} `json:"user"`
}

type hwRepo struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Namespace   string   `json:"namespace"`
	NumDownload int      `json:"num_download"`
	NumImages   int      `json:"num_images"`
	Size        int64    `json:"size"`
	Tags        []string `json:"Tags"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type hwTags struct {
	ImageId  string `json:"image_id"`
	Manifest string `json:"manifest"`
	Digest   string `json:"digest"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

type hwManifest struct {
	Layers []struct {
		Digest string `json:"digest"`
	} `json:"layers"`
}

func (ns hwNamespace) metadata() map[string]interface{} {
	var metadata = make(map[string]interface{})
	metadata["id"] = ns.ID
	metadata["creator_name"] = ns.CreatorName
	metadata["domain_public"] = ns.DomainPublic
	metadata["auth"] = ns.Auth
	metadata["domain_name"] = ns.DomainName
	metadata["user_count"] = ns.UserCount
	metadata["image_count"] = ns.ImageCount

	return metadata
}

func (this *HuaweiSWRService) Auth(url string, user string, pwd string) (token string, error error) {

	authUrl := fmt.Sprintf("iam.%s/v3/auth/tokens", url)

	var u = make(map[string]interface{})
	u["name"] = user
	u["password"] = pwd

	var domain = make(map[string]interface{})
	domain["name"] = user
	u["domain"] = domain

	var project = make(map[string]interface{})
	project["name"] = "cn-north-4"

	var scope = make(map[string]interface{})
	scope["project"] = project

	var params = make(map[string]interface{})

	auth := Auth{
		Identity: Identity{
			Methods: []string{"password"},
			Password: password{
				User: u,
			},
		},
		Scope: scope,
	}

	params["auth"] = auth

	proxy := proxy.ProxyServer{TargetUrl: authUrl, Body: params, Method: "POST"}
	resp, err := proxy.Request(user, pwd)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return "", errors.New(resp.Status)
	}

	value := resp.Header.Get("X-Subject-Token")

	return value, nil
}

func (this *HuaweiSWRService) ListNamespaces() models.Result {
	var ResultData models.Result
	token, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
		return ResultData
	}

	urls := fmt.Sprintf("swr-api.%s/v2/manage/namespaces", this.ImageConfig.Registry.Url)
	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, _ := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, e := ioutil.ReadAll(resp.Body)

		var namespaces []string
		var namespacesData hwNamespaceList
		e = json.Unmarshal(body, &namespacesData)
		if e != nil {
			ResultData.Message = e.Error()
			ResultData.Code = utils.GetNamespacesErr
			return ResultData
		}

		for _, namespaceData := range namespacesData.Namespace {
			namespaces = append(namespaces, namespaceData.Name)
		}
		ResultData.Code = http.StatusOK
		ResultData.Data = namespaces
	}
	return ResultData
}

func (this *HuaweiSWRService) Imports() (error error) {
	token, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return err
	}
	repos, reposErr := this.getRepos(token)
	if reposErr != nil {
		return reposErr
	}
	for _, repo := range repos {
		tags, _ := this.getTags(token, repo.Name)
		for _, tag := range tags {
			this.ImageConfig.Name = tag.Path
			this.ImageConfig.ImageId = "sha256:" + tag.ImageId
			var mf = &hwManifest{}
			json.Unmarshal([]byte(tag.Manifest), mf)
			if ic := this.ImageConfig.Get(); ic == nil {
				this.ImageConfig.Id = ""
				timeTemplate1 := "2006-01-02T15:04:05Z"
				stamp, _ := time.ParseInLocation(timeTemplate1, tag.Created, time.Local)
				this.ImageConfig.CreateTime = stamp.UnixNano()
				this.ImageConfig.Size = utils.FormatFileSize(tag.Size)
				this.ImageConfig.Add()

				imageDetail := models.ImageDetail{}

				imageDetail.ImageId = this.ImageConfig.ImageId
				imageDetail.Name = this.ImageConfig.Name

				if imd := imageDetail.Get(); imd == nil {
					imageDetail.Layers = len(mf.Layers)
					imageDetail.CreateTime = this.ImageConfig.CreateTime
					imageDetail.RepoDigests = tag.Digest
					imageDetail.Dockerfile = ""
					imageDetail.Size = this.ImageConfig.Size
					if result := imageDetail.Add(); result.Code != http.StatusOK {
						logs.Error("ImageDetail err: %s", errors.New(result.Message))
						return errors.New(result.Message)
					}
				}
			}
		}
	}
	return
}

func (this *HuaweiSWRService) getRepos(token string) (hw []*hwRepo, err error) {

	path := "/v2/manage/repos"
	if this.ImageConfig.Namespaces != "" {
		path += "?namespace=" + this.ImageConfig.Namespaces
	}

	urls := fmt.Sprintf("swr-api.%s"+path, this.ImageConfig.Registry.Url)

	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, err := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		repos, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(repos, &hw)
	}
	return
}

func (this *HuaweiSWRService) getTags(token string, repo string) (hw []hwTags, err error) {

	path := "/v2/manage/namespaces"
	if this.ImageConfig.Namespaces != "" {
		path = path + "/" + this.ImageConfig.Namespaces + "/repos/" + repo + "/tags"
	}

	urls := fmt.Sprintf("swr-api.%s"+path, this.ImageConfig.Registry.Url)
	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, err := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		tags, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(tags, &hw)
	}
	return
}
