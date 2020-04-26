package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	msl "github.com/xiliangMa/diss-backend/models/securitylog"
)

// Vulnerabilities Log api list
type VulnerabilitiesLogController struct {
	beego.Controller
}

// @Title GetImageVulnerabilitiesLog
// @Description Get ImageVulnerabilities Log List (暂时不支持风险等级查询)
// @Param token header string true "authToken"
// @Param body body securitylog.ImagePackageVulnerabilities false "镜像漏洞信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilities/image [post]
func (this *VulnerabilitiesLogController) GetImageVulnerabilitiesLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	imagePackageVulnerabilities := new(msl.ImagePackageVulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imagePackageVulnerabilities)
	this.Data["json"] = imagePackageVulnerabilities.List(from, limit)
	this.ServeJSON(false)
}
