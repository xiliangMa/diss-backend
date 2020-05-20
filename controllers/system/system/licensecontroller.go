package system

import (
	"encoding/json"
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
	json.Unmarshal(this.Ctx.Input.RequestBody, &licFile)

	this.Data["json"] = licFile.Add()
	this.ServeJSON(false)
}
