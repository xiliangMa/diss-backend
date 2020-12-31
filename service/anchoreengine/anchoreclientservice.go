package anchoreengine

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type AnchoreService struct {
	ImageParams *models.ImageSearchParams
}

func GetAnchoreManager() *models.AnchoreEngineManager {
	if models.AEM == nil {
		models.AEM = models.NewAnchoreEngineManager()
	}
	return models.AEM
}

func (this *AnchoreService) GetImageContent() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	client := GetAnchoreManager().AnchoreClient
	if this.ImageParams == nil || this.ImageParams.ImageId == "" {
		ResultData.Code = utils.NotEnoughParametersErr
		ResultData.Message = "NotEnoughParametersErr"
		return ResultData
	}
	if this.ImageParams.Ctype == "" {
		this.ImageParams.Ctype = models.Image_Content_Type_OS
	}
	respData, _, err := client.ImagesApi.GetImageContentByTypeImageId(nil, this.ImageParams.ImageId, this.ImageParams.Ctype, nil)
	if err != nil {
		ResultData.Code = utils.GetImageContentErr
		ResultData.Message = err.Error()
		logs.Error("Get ImageContent failed from anchore engine, err: %s.", err.Error())
		return ResultData
	}
	data := make(map[string]interface{})
	data["items"] = respData
	data["total"] = len(respData.Content)
	ResultData.Data = data
	return ResultData
}

func (this *AnchoreService) GetImageVuln() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	client := GetAnchoreManager().AnchoreClient
	if this.ImageParams == nil || this.ImageParams.ImageId == "" {
		ResultData.Code = utils.NotEnoughParametersErr
		ResultData.Message = "NotEnoughParametersErr"
		return ResultData
	}
	if this.ImageParams.Vtype == "" {
		this.ImageParams.Vtype = models.Image_Vuln_Type_All
	}
	respData, _, err := client.ImagesApi.GetImageVulnerabilitiesByTypeImageId(nil, this.ImageParams.ImageId, this.ImageParams.Vtype, nil)
	if err != nil {
		ResultData.Code = utils.GetImageVulnErr
		ResultData.Message = err.Error()
		logs.Error("Get ImageVulnerabilities failed from anchore engine, err: %s.", err.Error())
		return ResultData
	}
	data := make(map[string]interface{})
	data["items"] = respData
	data["total"] = len(respData.Vulnerabilities)
	ResultData.Data = data
	return ResultData
}

func (this *AnchoreService) GetImageVulnStatistics() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	client := GetAnchoreManager().AnchoreClient
	if this.ImageParams == nil || this.ImageParams.ImageId == "" {
		ResultData.Code = utils.NotEnoughParametersErr
		ResultData.Message = "NotEnoughParametersErr"
		return ResultData
	}
	if this.ImageParams.Vtype == "" {
		this.ImageParams.Vtype = models.Image_Vuln_Type_All
	}
	respData, _, err := client.ImagesApi.GetImageVulnerabilitiesByTypeImageId(nil, this.ImageParams.ImageId, this.ImageParams.Vtype, nil)
	if err != nil {
		ResultData.Code = utils.GetImageVulnErr
		ResultData.Message = err.Error()
		logs.Error("Get ImageVulnStatistics failed from anchore engine, err: %s.", err.Error())
		return ResultData
	}
	vulnStatistics := make(map[string]interface{})
	vulnLowCount := 0
	vulnMediumCount := 0
	if &respData != nil && len(respData.Vulnerabilities) > 0 {
		for _, vuln := range respData.Vulnerabilities {
			switch vuln.Severity {
			case models.Image_vuln_Severity_Low:
				vulnLowCount = vulnLowCount + 1
			case models.Image_vuln_Severity_Medium:
				vulnMediumCount = vulnMediumCount + 1
			}
		}
	}
	vulnStatistics[models.Image_vuln_Severity_Low] = vulnLowCount
	vulnStatistics[models.Image_vuln_Severity_Medium] = vulnMediumCount
	data := make(map[string]interface{})
	data["items"] = vulnStatistics
	data["total"] = 0
	ResultData.Data = data
	return ResultData
}

func (this *AnchoreService) GetImageMetadata() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	client := GetAnchoreManager().AnchoreClient
	if this.ImageParams == nil || this.ImageParams.ImageDigest == "" {
		ResultData.Code = utils.NotEnoughParametersErr
		ResultData.Message = "NotEnoughParametersErr"
		return ResultData
	}
	if this.ImageParams.Mtype == "" {
		this.ImageParams.Mtype = models.Image_Metadata_Type_Manifest
	}
	respData, _, err := client.ImagesApi.GetImageMetadataByType(nil, this.ImageParams.ImageDigest, this.ImageParams.Mtype, nil)
	if err != nil {
		ResultData.Code = utils.GetImageMetadataErr
		ResultData.Message = err.Error()
		logs.Error("Get ImageMetadata failed from anchore engine, err: %s.", err.Error())
		return ResultData
	}
	data := make(map[string]interface{})
	data["items"] = respData
	data["total"] = 1
	ResultData.Data = data
	return ResultData
}

func (this *AnchoreService) GetImageSensitiveInfo() models.Result {
	ResultData := models.Result{Code: http.StatusOK}
	this.ImageParams.Ctype = ""
	client := GetAnchoreManager().AnchoreClient
	if this.ImageParams == nil || this.ImageParams.ImageId == "" {
		ResultData.Code = utils.NotEnoughParametersErr
		ResultData.Message = "NotEnoughParametersErr"
		return ResultData
	}
	respData, _, err := client.ImagesApi.GetImageContentByTypeImageIdFiles(nil, this.ImageParams.ImageId, nil)
	if err != nil {
		ResultData.Code = utils.GetImageContentErr
		ResultData.Message = err.Error()
		logs.Error("Get ImageSensitiveInfo failed from anchore engine, err: %s.", err.Error())
		return ResultData
	}
	if &respData != nil && len(respData.Content) > 0 {
		for _, file := range respData.Content {
			// 匹配敏感信息
			// todo
			logs.Info(file.Filename)
		}
	}
	data := make(map[string]interface{})
	data["items"] = respData
	data["total"] = 1
	ResultData.Data = data
	return ResultData
}
