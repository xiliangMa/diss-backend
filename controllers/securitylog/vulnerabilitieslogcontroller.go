package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Vulnerabilities Log api list
type VulnerabilitiesLogController struct {
	beego.Controller
}

// @Title GetImageVulnerabilitiesLog
// @Description Get ImageVulnerabilities Log List (暂时不支持风险等级查询)
// @Param token header string true "authToken"
// @Param body body models.ImagePackageVulnerabilities false "镜像漏洞信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilities/image [post]
func (this *VulnerabilitiesLogController) GetImageVulnerabilitiesLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	imagePackageVulnerabilities := new(models.ImagePackageVulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imagePackageVulnerabilities)
	this.Data["json"] = imagePackageVulnerabilities.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetImageVulnerabilityInfo
// @Description 获取镜像漏洞详情 (Cnnvd)
// @Param token header string true "authToken"
// @Param body body models.VulnerabilityInfo false "镜像漏洞详情"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilities/image/info [post]
func (this *VulnerabilitiesLogController) GetImageVulnerabilityInfo() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	imagePackageVulnerabilityInfo := new(models.VulnerabilityInfo)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imagePackageVulnerabilityInfo)
	this.Data["json"] = imagePackageVulnerabilityInfo.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetVulnerabilities
// @Description Get Vulnerabilities Log
// @Param token header string true "authToken"
// @Param body body models.ImageVulnerabilities false "漏洞"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilities [post]
func (this *VulnerabilitiesLogController) GetVulnerabilities() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	imageVulnerabilities := new(models.ImageVulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageVulnerabilities)
	this.Data["json"] = imageVulnerabilities.List(from, limit)
	this.ServeJSON(false)
}
