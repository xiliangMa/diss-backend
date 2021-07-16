package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/utils"
)

type DockerHubService struct {
	ImageConfig *models.ImageConfig
}

type dhNamespaceList struct {
	Namespace []string `json:"namespaces"`
}

type dhRepos struct {
	Next  string `json:"next"`
	Repos []Repo `json:"results"`
}

type Repo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type dhTags struct {
	Next string `json:"next"`
	Tags []struct {
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

func (this *DockerHubService) Imports() (err error) {
	token, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return err
	}
	go func() {
		var repos []Repo
		page := 1
		pageSize := 100
		n := 0
		for {
			repo, reposErr := this.getRepos(token, page, pageSize)
			if reposErr != nil {
				reposErr.Error()
			}
			repos = append(repos, repo.Repos...)
			n += len(repo.Repos)
			if len(repo.Next) == 0 {
				break
			}
			page++
		}

		for _, repo := range repos {
			page = 1
			pageSize = 100
			var tags []string
			for {
				tag, _ := this.getTags(token, repo.Name, page, pageSize)
				for _, t := range tag.Tags {
					tags = append(tags, t.Name)
				}
				if len(tag.Next) == 0 {
					break
				}
				page++
			}
			if len(tags) > 0 {
				for _, t := range tags {
					this.ImageConfig.Name = repo.Namespace + "/" + repo.Name + ":" + t
					cs := CommonService{ImageConfig: this.ImageConfig}
					cs.AddDetail()
				}
			}
		}
	}()
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

func (this *DockerHubService) getRepos(token string, page, pageSize int) (dh *dhRepos, err error) {
	url := fmt.Sprintf("%s/v2/repositories/%s/?page=%d&page_size=%d", this.ImageConfig.Registry.Url, this.ImageConfig.Namespaces, page, pageSize)
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

func (this *DockerHubService) getTags(token string, repo string, page, pageSize int) (dh *dhTags, err error) {
	urls := fmt.Sprintf("%s/v2/repositories/%s/%s/tags/?page=%d&page_size=%d", this.ImageConfig.Registry.Url, this.ImageConfig.Namespaces, repo, page, pageSize)
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
