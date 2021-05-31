package base

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/service/registry"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type RegistryService struct {
	Registry *models.Registry
}

func (this *RegistryService) Ping() models.Result {
	var ResultData models.Result

	path := ""
	if this.Registry.Type == models.Registry_Type_Harbor || this.Registry.Type == models.Registry_Type_DockerRegistry {
		path = "/v2/_catalog"
		return ping(this.Registry, path)
	} else if this.Registry.Type == models.Registry_Type_AlibabaACR {
		acr := registry.AlibabaACRService{}
		err := acr.NewAuth(this.Registry)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.TestLinkRegistryErr
			logs.Error("Test link failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	} else if this.Registry.Type == models.Registry_Type_HuaweiSWR {
		hw := registry.HuaweiSWRService{}
		_, err := hw.Auth(this.Registry.Url, this.Registry.User, this.Registry.Pwd)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.TestLinkRegistryErr
			logs.Error("Test link failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}
	return ResultData
}

func ping(registry *models.Registry, path string) models.Result {
	var ResultData models.Result
	proxy := proxy.ProxyServer{TargetUrl: registry.Url + path}
	resp, err := proxy.Request(registry.User, registry.Pwd)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.TestLinkRegistryErr
		logs.Error("Test link failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if resp.StatusCode == 401 {
		ResultData.Message = resp.Status
		ResultData.Code = utils.TestLinkRegistryErr
		logs.Error("Test link failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = registry
	return ResultData
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
