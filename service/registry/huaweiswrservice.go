package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
)

type HuaweiSWRService struct {
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

type namespace struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
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

type hWRepo struct {
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

func (this *HuaweiSWRService) ListNamespaces(imageConfig *models.ImageConfig) models.Result {
	var ResultData models.Result

	token, err := this.Auth(imageConfig.Registry.Url, imageConfig.Registry.User, imageConfig.Registry.Pwd)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
		return ResultData
	}

	urls := fmt.Sprintf("swr-api.%s/v2/manage/namespaces", imageConfig.Registry.Url)

	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, _ := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)

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

func (this *HuaweiSWRService) Imports(imageConfig *models.ImageConfig) (error error) {
	token, err := this.Auth(imageConfig.Registry.Url, imageConfig.Registry.User, imageConfig.Registry.Pwd)

	if err != nil {
		return err
	}

	path := "/v2/manage/repos"
	if imageConfig.Namespaces != "" {
		path += "?namespace=" + imageConfig.Namespaces
	}

	urls := fmt.Sprintf("swr-api.%s"+path, imageConfig.Registry.Url)

	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, _ := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var hw []*hWRepo
		repos, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(repos, &hw)
		for _, repo := range hw {
			for _, tag := range repo.Tags {
				imageConfig.Name = repo.Path + ":" + tag
				imageConfig.Size = utils.FormatFileSize(repo.Size)
				if ic := imageConfig.Get(); ic == nil {
					imageConfig.Id = ""
					imageConfig.Add()
				}
			}
		}
	}
	return nil
}
