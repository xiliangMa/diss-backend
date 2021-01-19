package securitypolicy

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// System Template Group object api list
type SystemTemplateGroupController struct {
	web.Controller
}

// @Title GetSystemTemplateGroupList
// @Description Get System Template List
// @Param token header string true "authToken"
// @Param body body models.SystemTemplateGroup false "安全策略组"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *SystemTemplateGroupController) GetSystemTemplateGroupLIst() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	SystemTemplateGroup := new(models.SystemTemplateGroup)
	json.Unmarshal(this.Ctx.Input.RequestBody, &SystemTemplateGroup)
	this.Data["json"] = SystemTemplateGroup.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddSystemTemplateGroup
// @Description Add SystemTemplateGroup
// @Param token header string true "authToken"
// @Param body body models.SystemTemplateGroup false "安全策略组"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *SystemTemplateGroupController) AddSystemTemplateGroup() {
	systemTemplateGroup := new(models.SystemTemplateGroup)
	json.Unmarshal(this.Ctx.Input.RequestBody, &systemTemplateGroup)
	this.Data["json"] = systemTemplateGroup.Add()
	this.ServeJSON(false)
}

// @Title DeleteSystemTemplateGroup
// @Description Delete SystemTemplateGroup
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *SystemTemplateGroupController) DeleteSystemTemplateGroup() {
	id := this.GetString(":id")
	SystemTemplateGroup := new(models.SystemTemplateGroup)
	SystemTemplateGroup.Id = id
	result := SystemTemplateGroup.Delete()
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title UpdateSystemTemplateGroup
// @Description Update SystemTemplateGroup
// @Param token header string true "authToken"
// @Param id path string "" true "id"
// @Param body body models.SystemTemplateGroup false "安全策略组"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *SystemTemplateGroupController) UpdateSystemTemplateGroup() {
	id := this.GetString(":id")
	SystemTemplateGroup := new(models.SystemTemplateGroup)
	json.Unmarshal(this.Ctx.Input.RequestBody, &SystemTemplateGroup)
	SystemTemplateGroup.Id = id
	result := SystemTemplateGroup.Update()
	this.Data["json"] = result
	this.ServeJSON(false)
}
