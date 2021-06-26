package base

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/service/registry"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type RegistryService struct {
	Registry *models.Registry
}

func (this *RegistryService) Ping() models.Result {
	var ResultData models.Result
	var err error
	if this.Registry.Type == models.Registry_Type_Harbor || this.Registry.Type == models.Registry_Type_DockerRegistry || this.Registry.Type == models.Registry_Type_JFrogArtifactory {
		err = ping(this.Registry)
	} else if this.Registry.Type == models.Registry_Type_DockerHub {
		dh := registry.DockerHubService{}
		_, err = dh.Auth(this.Registry.Url, this.Registry.User, this.Registry.Pwd)
	} else if this.Registry.Type == models.Registry_Type_AlibabaACR {
		acr := registry.AlibabaACRService{}
		err = acr.NewAuth(this.Registry)
	} else if this.Registry.Type == models.Registry_Type_HuaweiSWR {
		hw := registry.HuaweiSWRService{}
		_, err = hw.Auth(this.Registry.Url, this.Registry.User, this.Registry.Pwd)
	}
	if err != nil {
		ResultData.Message = "Bad credentials"
		ResultData.Code = utils.TestLinkRegistryErr
		logs.Error("Test link failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this.Registry
	return ResultData
}

func ping(registry *models.Registry) (error error) {

	proxy := proxy.ProxyServer{TargetUrl: strings.TrimRight(registry.Url, "/") + "/v2/"}
	resp, err := proxy.Request(registry.User, registry.Pwd)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return
}

func (this *RegistryService) TypeInfos() models.Result {
	var ResultData models.Result
	var infoMap = map[string]*utils.EndpointPattern{}

	infoMap[models.Registry_Type_AlibabaACR] = utils.ACRRegion()

	infoMap[models.Registry_Type_HuaweiSWR] = utils.SWRRegion()

	infoMap[models.Registry_Type_AwsECR] = utils.AWSRegion()

	infoMap[models.Registry_Type_DockerHub] = utils.DockerHubRegion()

	infoMap[models.Registry_Type_GoogleGCR] = utils.GoogleRegion()

	ResultData.Data = infoMap

	ResultData.Code = http.StatusOK
	return ResultData

}
