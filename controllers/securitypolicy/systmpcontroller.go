package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/securitypolicy"
)

// System Template object api list
type SystemTemplateController struct {
	beego.Controller
}

// @Title GetSystemTemplateList
// @Description Get System Template List
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.SystemTemplate false "安全策略"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SystemTemplateController) GetSystemTemplateLIst() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	systemTemplate := new(models.SystemTemplate)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplate)
	this.Data["json"] = systemTemplate.List(from, limit)
	this.ServeJSON(false)
}

// @Title DeleteSystemTemplate
// @Description Delete SystemTemplate
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *SystemTemplateController) DeleteSystemTemplate() {
	id := this.GetString(":id")
	systemTemplate := new(models.SystemTemplate)
	systemTemplate.Id = id
	result := systemTemplate.Delete()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddSystemTemplate
// @Description Add SystemTemplate
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.SystemTemplate false "SystemTemplate"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *SystemTemplateController) AddSystemTemplate() {
	systemTemplate := new(models.SystemTemplate)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplate)
	sysTemplateService := securitypolicy.SystemTemplateService{}
	sysTemplateService.SystemTemplate = systemTemplate
	result := sysTemplateService.AddSystemTemplate()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title UpdateSystemTemplate
// @Description Update SystemTemplate
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param id path string "" true "id"
// @Param body body models.SystemTemplate false "SystemTemplate"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *SystemTemplateController) UpdateSystemTemplate() {
	id := this.GetString(":id")
	systemTemplate := new(models.SystemTemplate)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplate)
	systemTemplate.Id = id
	result := systemTemplate.Update()
	this.Data["json"] = result
	this.ServeJSON(false)
}
