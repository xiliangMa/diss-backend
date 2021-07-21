package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 模块接口列表
type ModuleController struct {
	beego.Controller
}

// @Title AddModule
// @Description Add Module
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.Module false "模块信息"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ModuleController) AddModule() {
	module := new(models.Module)
	json.Unmarshal(this.Ctx.Input.RequestBody, &module)

	this.Data["json"] = module.Add()
	this.ServeJSON(false)
}

// @Title GetModules
// @Description Get Modules
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.Module false "模块信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /modulelist [post]
func (this *ModuleController) ModuleList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	module := new(models.Module)
	json.Unmarshal(this.Ctx.Input.RequestBody, &module)
	this.Data["json"] = module.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateModule
// @Description Update Module
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param moduleId path string "" true "moduleId"
// @Param body body models.Module true "用户信息"
// @Success 200 {object} models.Result
// @router /:moduleId [put]
func (this *ModuleController) UpdateModule() {
	moduleId, _ := this.GetInt(":moduleId")
	module := new(models.Module)
	json.Unmarshal(this.Ctx.Input.RequestBody, &module)
	module.Id = moduleId
	this.Data["json"] = module.Update()
	this.ServeJSON(false)
}

// @Title DeleteModule
// @Description Delete Module
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param userId path string "" true "moduleId"
// @Success 200 {object} models.Result
// @router /:moduleId [delete]
func (this *ModuleController) DeleteModule() {
	moduleId, _ := this.GetInt(":moduleId")
	module := new(models.Module)
	module.Id = moduleId

	result := module.Delete()

	this.Data["json"] = result
	this.ServeJSON(false)
}
