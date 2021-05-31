package base

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/service/registry"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
)

type ImageConfigService struct {
	ImageConfig *models.ImageConfig
}

func (this *ImageConfigService) BatchImportImage() models.Result {
	result := models.Result{}
	if this.ImageConfig.Registry.Type == models.Registry_Type_Harbor || this.ImageConfig.Registry.Type == models.Registry_Type_DockerRegistry {
		GeneralType(this.ImageConfig)
	} else if this.ImageConfig.Registry.Type == models.Registry_Type_AlibabaACR {
		acr := registry.AlibabaACRService{}
		err := acr.Imports(this.ImageConfig)
		if err != nil {
			result.Message = err.Error()
			result.Code = utils.ImportImageErr
			logs.Error("Import Image failed, code: %d, err: %s", result.Code, result.Message)
			return result
		}
	} else if this.ImageConfig.Registry.Type == models.Registry_Type_HuaweiSWR {
		hw := registry.HuaweiSWRService{}
		err := hw.Imports(this.ImageConfig)
		if err != nil {
			result.Message = err.Error()
			result.Code = utils.ImportImageErr
			logs.Error("Import Image failed, code: %d, err: %s", result.Code, result.Message)
			return result
		}
	}

	result.Code = http.StatusOK
	return result
}

func (this *ImageConfigService) GetNamespaces() models.Result {
	result := models.Result{}
	switch this.ImageConfig.Registry.Type {
	case models.Registry_Type_AlibabaACR:
		ics := registry.AlibabaACRService{}
		return ics.GetNamespaces(this.ImageConfig)
	case models.Registry_Type_HuaweiSWR:
		hw := registry.HuaweiSWRService{}
		return hw.ListNamespaces(this.ImageConfig)
	}

	return result
}

func GeneralType(imageConfig *models.ImageConfig) {
	proxy := proxy.ProxyServer{TargetUrl: imageConfig.Registry.Url + "/v2/_catalog"}
	resp, _ := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var cc map[string]interface{}
		json.Unmarshal(body, &cc)

		if cc["repositories"] != nil {
			for _, imageName := range cc["repositories"].([]interface{}) {

				in := imageName.(string)
				proxy.TargetUrl = imageConfig.Registry.Url + "/v2/" + in + "/tags/list"
				tags, _ := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)

				if tags.StatusCode == 200 {
					defer tags.Body.Close()
					t, _ := ioutil.ReadAll(tags.Body)
					var tagObj map[string]interface{}
					json.Unmarshal(t, &tagObj)

					if tagObj["tags"] != nil {
						for _, tag := range tagObj["tags"].([]interface{}) {
							imageConfig.Name = in + ":" + tag.(string)
							if ic := imageConfig.Get(); ic == nil {
								imageConfig.Id = ""
								imageConfig.Add()
							}
						}
					}
				}
			}
		}
	}
}
