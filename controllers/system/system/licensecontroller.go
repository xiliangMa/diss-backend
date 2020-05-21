package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

// @Title License File Import
// @Description 导入授权文件配置
// @Param token header string true "authToken"
// @Param body body models.LicenseConfig false "授权文件配置信息"
// @Success 200 {object} models.Result
// @router /system/licensefileimport [post]
func (this *IntegrationController) AddLicenseFile() {
	licFile := new(models.LicenseConfig)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &licFile)
	if err!=nil{
		logs.Info("parse json err:", err)
	}

	logs.Info("licFile: %#v", licFile.LicenseModule)

	this.Data["json"] = licFile.Add()
	this.ServeJSON(false)
}


// @Title Get License Config
// @Description 获取指定授权id的详情
// @Param token header string true "authToken"
// @Param LicenseUuid query string true "授权文件uuid"
// @Success 200 {object} models.Result
// @router /system/licenseconfig [get]
func (this *IntegrationController) GetLicenseData() {
	licConfig := new(models.LicenseConfig)
	configId := this.GetString(":LicenseUuid")
	licConfig.LicenseUuid = configId
	json.Unmarshal(this.Ctx.Input.RequestBody, &licConfig)

	var ResultData models.Result
	ResultData.Code = http.StatusOK
	ResultData.Data = licConfig.Get()

	this.Data["json"] = ResultData
	this.ServeJSON(false)
}
