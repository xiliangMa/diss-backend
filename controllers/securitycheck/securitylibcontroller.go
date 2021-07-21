package securitycheck

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
)

type SecurityLibController struct {
	beego.Controller
}

// @Title AddKubeVulnLib
// @Description Add VulnerabilityLib Item
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.VulnerabilityLib false "漏洞库资料"
// @Success 200 {object} models.Result
// @router /vulnlib [post]
func (this *SecurityLibController) AddVulnLib() {
	vulnLib := new(models.VulnerabilityLib)

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &vulnLib)
	logs.Warn(err)
	this.Data["json"] = vulnLib.Add()

	this.ServeJSON(false)
}

// @Title UpdateVulnerabilityLib
// @Description Update VulnerabilityLib
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "Id"
// @Param body body models.VulnerabilityLib false "漏洞库资料"
// @Success 200 {object} models.Result
// @router /vulnlib/:id [put]
func (this *SecurityLibController) UpdateVulnLib() {
	id, _ := this.GetInt64(":id")
	vulnLib := new(models.VulnerabilityLib)
	json.Unmarshal(this.Ctx.Input.RequestBody, &vulnLib)
	vulnLib.Id = id
	result := vulnLib.Update()

	natsManager := models.Nats
	natsPubService := nats.NatsPubService{Conn: natsManager.Conn}
	natsPubService.Type = models.RuleType_DockerVulnerability
	natsPubService.RuleDefinePub()

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title GetVulnerabilityLib
// @Description Get VulnerabilityLib
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.VulnerabilityLib false "漏洞库资料"
// @Success 200 {object} models.Result
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @router /vulnliblist [post]
func (this *SecurityLibController) GetVulnLib() {
	vulenerLib := new(models.VulnerabilityLib)
	json.Unmarshal(this.Ctx.Input.RequestBody, &vulenerLib)

	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	this.Data["json"] = vulenerLib.List(from, limit)
	this.ServeJSON(false)
}

// @Title DeleteVulnerabilityLib
// @Description Delete VulnerabilityLib
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "id"
// @Success 200 {object} models.VulnerabilityLib
// @router /vulnlib/:id [delete]
func (this *SecurityLibController) DeleteVulnerabilityLib() {
	id, _ := this.GetInt64(":id")
	vulnLib := new(models.VulnerabilityLib)
	vulnLib.Id = id

	this.Data["json"] = vulnLib.Delete()
	this.ServeJSON(false)
}
