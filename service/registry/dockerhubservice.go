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

type DockerHubService struct {
	ImageConfig *models.ImageConfig
}

type dhNamespaceList struct {
	Namespace []string `json:"namespaces"`
}

type dhRepos struct {
	Results []struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"results"`
}

type dhTags struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func (this *DockerHubService) Auth(url string, user string, pwd string) (token string, error error) {

	urls := fmt.Sprintf("%s/v2/users/login", url)

	var u = make(map[string]interface{})
	u["username"] = user
	u["password"] = pwd

	proxy := proxy.ProxyServer{TargetUrl: urls, Body: u, Method: "POST"}
	resp, err := proxy.Request(user, pwd)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var cc map[string]interface{}
	json.Unmarshal(body, &cc)

	value := cc["token"].(string)
	return "JWT " + value, nil
}

func (this *DockerHubService) Imports() (error error) {
	token, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return err
	}
	repos, reposErr := this.getRepos(token)
	if reposErr != nil {
		return reposErr
	}
	for _, repo := range repos.Results {
		tags, _ := this.getTags(token, repo.Name)
		for _, t := range tags.Results {
			this.ImageConfig.Name = repo.Namespace + "/" + repo.Name + ":" + t.Name
			cs := CommonService{ImageConfig: this.ImageConfig}
			cs.AddDetail()
		}
	}
	return
}

func (this *DockerHubService) ListNamespaces() models.Result {
	var ResultData models.Result
	token, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
		return ResultData
	}

	url := fmt.Sprintf("%s/v2/repositories/namespaces/", this.ImageConfig.Registry.Url)
	proxy := proxy.ProxyServer{TargetUrl: url, Token: token}
	resp, _ := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, e := ioutil.ReadAll(resp.Body)

		var na dhNamespaceList
		json.Unmarshal(body, &na)

		if e != nil {
			ResultData.Message = e.Error()
			ResultData.Code = utils.GetNamespacesErr
			return ResultData
		}
		ResultData.Data = na.Namespace
		ResultData.Code = http.StatusOK
	}
	return ResultData
}

func (this *DockerHubService) getRepos(token string) (dh *dhRepos, err error) {

	path := "/v2/repositories"
	if this.ImageConfig.Namespaces != "" {
		path = path + "/" + this.ImageConfig.Namespaces + "/?page_size=10000"
	}

	url := fmt.Sprintf("%s"+path, this.ImageConfig.Registry.Url)
	proxy := proxy.ProxyServer{TargetUrl: url, Token: token}
	resp, err := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		repos, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(repos, &dh)
	}
	return
}

func (this *DockerHubService) getTags(token string, repo string) (dh *dhTags, err error) {

	path := "/v2/repositories/"
	if this.ImageConfig.Namespaces != "" {
		path = path + "/" + this.ImageConfig.Namespaces + "/" + repo + "/tags"
	}

	urls := fmt.Sprintf("%s"+path, this.ImageConfig.Registry.Url)
	proxy := proxy.ProxyServer{TargetUrl: urls, Token: token}
	resp, err := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		tags, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(tags, &dh)
	}
	return
}
