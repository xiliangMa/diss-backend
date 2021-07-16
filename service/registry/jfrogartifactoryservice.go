package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/utils"
)

type JFrogArtifactoryService struct {
	ImageConfig *models.ImageConfig
}

type repositories struct {
	Repository []struct {
		Key string `json:"key"`
	} `json:"repository"`
}

func (this *JFrogArtifactoryService) ListRepositories() models.Result {
	var ResultData models.Result

	urls := fmt.Sprintf("%s/artifactory/api/repositories?packageType=%s", this.ImageConfig.Registry.Url, "docker")
	proxy := proxy.ProxyServer{TargetUrl: urls}
	resp, _ := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, e := ioutil.ReadAll(resp.Body)
		var repositoriesData []string
		var repo repositories
		json.Unmarshal(body, &repo.Repository)
		if e != nil {
			ResultData.Message = e.Error()
			ResultData.Code = utils.GetNamespacesErr
			return ResultData
		}
		for _, r := range repo.Repository {
			repositoriesData = append(repositoriesData, r.Key)
		}
		ResultData.Data = repositoriesData
		ResultData.Code = http.StatusOK
	}
	return ResultData
}
