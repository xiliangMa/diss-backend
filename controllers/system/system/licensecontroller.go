package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

// @Title License File
// @Description 新建授权文件配置
// @Param token header string true "authToken"
// @Param body body models.LicenseFile false "授权文件配置信息"
// @Success 200 {object} models.Result
// @router /system/licensefile [post]
func (this *IntegrationController) AddLicenseFile() {
	licFile := new(models.LicenseFile)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &licFile)
	if err!=nil{
		logs.Info("parse json err:", err)
	}

	logs.Info("licFile: %#v", licFile.LicenseModule)

	this.Data["json"] = licFile.Add()
	this.ServeJSON(false)
}
