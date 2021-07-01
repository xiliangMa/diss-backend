package base

import (
	"encoding/json"
	"fmt"
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
	var err error

	if this.ImageConfig.Registry.Type == models.Registry_Type_Harbor || this.ImageConfig.Registry.Type == models.Registry_Type_DockerRegistry {
		err = this.generalType(this.ImageConfig.Registry.Url)
	} else {
		if this.ImageConfig.Namespaces == "" {
			result.Message = "NamespacesNotEmpty"
			result.Code = utils.GetNamespacesErr
			return result
		}
	}
	switch this.ImageConfig.Registry.Type {
	case models.Registry_Type_DockerHub:
		dh := registry.DockerHubService{ImageConfig: this.ImageConfig}
		err = dh.Imports()
	case models.Registry_Type_AlibabaACR:
		acr := registry.AlibabaACRService{ImageConfig: this.ImageConfig}
		err = acr.Imports()
	case models.Registry_Type_HuaweiSWR:
		hw := registry.HuaweiSWRService{ImageConfig: this.ImageConfig}
		err = hw.Imports()
	case models.Registry_Type_JFrogArtifactory:
		url := fmt.Sprintf("%s/artifactory/api/docker/%s", this.ImageConfig.Registry.Url, this.ImageConfig.Namespaces)
		err = this.generalType(url)
	case models.Registry_Type_AwsECR:
		ae := registry.AwsECRService{ImageConfig: this.ImageConfig}
		err = ae.Imports()
	}

	if err != nil {
		result.Message = err.Error()
		result.Code = utils.ImportImageErr
		logs.Error("Import Image failed, code: %d, err: %s", result.Code, result.Message)
		return result
	}

	result.Code = http.StatusOK
	return result
}

func (this *ImageConfigService) GetNamespaces() models.Result {
	result := models.Result{}
	switch this.ImageConfig.Registry.Type {
	case models.Registry_Type_AlibabaACR:
		ics := registry.AlibabaACRService{ImageConfig: this.ImageConfig}
		return ics.ListNamespaces()
	case models.Registry_Type_DockerHub:
		dh := registry.DockerHubService{ImageConfig: this.ImageConfig}
		return dh.ListNamespaces()
	case models.Registry_Type_HuaweiSWR:
		hw := registry.HuaweiSWRService{ImageConfig: this.ImageConfig}
		return hw.ListNamespaces()
	case models.Registry_Type_JFrogArtifactory:
		hw := registry.JFrogArtifactoryService{ImageConfig: this.ImageConfig}
		return hw.ListRepositories()
	case models.Registry_Type_AwsECR:
		ae := registry.AwsECRService{ImageConfig: this.ImageConfig}
		return ae.ListRepositories()
	}

	return result
}

func (this *ImageConfigService) generalType(url string) (error error) {
	imageConfig := this.ImageConfig
	proxy := proxy.ProxyServer{}
	proxy.TargetUrl = fmt.Sprintf("%s/v2/_catalog", url)
	resp, err := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)

	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, respErr := ioutil.ReadAll(resp.Body)
		if respErr != nil {
			return respErr
		}
		var cc map[string]interface{}
		json.Unmarshal(body, &cc)

		if cc["repositories"] != nil {
			for _, imageName := range cc["repositories"].([]interface{}) {

				name := imageName.(string)
				proxy.TargetUrl = fmt.Sprintf("%s/v2/%s/tags/list", url, name)
				tags, _ := proxy.Request(imageConfig.Registry.User, imageConfig.Registry.Pwd)
				if tags.StatusCode == 200 {
					defer tags.Body.Close()
					t, _ := ioutil.ReadAll(tags.Body)
					var tagObj map[string]interface{}
					json.Unmarshal(t, &tagObj)

					if tagObj["tags"] != nil {
						for _, tag := range tagObj["tags"].([]interface{}) {
							imageConfig.Name = name + ":" + tag.(string)
							cs := registry.CommonService{ImageConfig: imageConfig}
							cs.AddDetail()
						}
					}
				}
			}
		}
	}
	return
}
