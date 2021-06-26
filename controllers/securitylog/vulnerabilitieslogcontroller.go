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

// @Title GetVulnerabilitiesScan
// @Description Get VulnerabilitiesScan
// @Param token header string true "authToken"
// @Param body body models.ImageVulnerabilities false "漏洞"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilitiesscan [post]
func (this *VulnerabilitiesLogController) GetVulnerabilitiesScan() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	imageVulnerabilities := new(models.ImageVulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &imageVulnerabilities)

	this.Data["json"] = imageVulnerabilities.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetVulnerabilities
// @Description Get Vulnerabilities
// @Param token header string true "authToken"
// @Param body body models.Vulnerabilities false "漏洞"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /vulnerabilities [post]
func (this *VulnerabilitiesLogController) GetVulnerabilities() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	vulnerabilities := new(models.Vulnerabilities)
	json.Unmarshal(this.Ctx.Input.RequestBody, &vulnerabilities)

	this.Data["json"] = vulnerabilities.List(from, limit)
	this.ServeJSON(false)
}
