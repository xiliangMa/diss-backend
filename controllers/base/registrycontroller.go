package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
)

// 镜像仓库
type RegistryController struct {
	beego.Controller
}

// @Title GetRegistryList
// @Description Get Registry List
// @Param token header string true "authToken"
// @Param body body models.Registry false "镜像仓库信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *RegistryController) GetRegistryList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	registry := new(models.Registry)
	json.Unmarshal(this.Ctx.Input.RequestBody, &registry)

	this.Data["json"] = registry.List(from, limit)
	this.ServeJSON(false)
}

// @Title Add
// @Description Add Registry
// @Param token header string true "authToken"
// @Param body body models.Registry true "Registry"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *RegistryController) Add() {

	registry := new(models.Registry)
	json.Unmarshal(this.Ctx.Input.RequestBody, &registry)

	rs := base.RegistryService{Registry: registry}
	if result := rs.Ping(); result.Data == nil {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	this.Data["json"] = registry.Add()
	this.ServeJSON(false)
}

// @Title GetRegistry
// @Description Registry
// @Param token header string true "authToken"
// @Param id path int 0 true "id"
// @Success 200 {object} models.Result
// @router /:id [post]
func (this *RegistryController) GetRegistry() {
	id, _ := this.GetInt64(":id")

	registry := new(models.Registry)
	registry.Id = id

	this.Data["json"] = registry.Get()
	this.ServeJSON(false)
}

// @Title DelRegistry
// @Description Del Registry
// @Param token header string true "authToken"
// @Param id path int 0 true "id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *RegistryController) DelRegistry() {
	id, _ := this.GetInt64(":id")

	registry := new(models.Registry)
	registry.Id = id

	this.Data["json"] = registry.Delete()
	this.ServeJSON(false)
}

// @Title ping
// @Description Test ping
// @Param token header string true "authToken"
// @Param body body models.Registry true "Registry"
// @Success 200 {object} models.Result
// @router /ping [post]
func (this *RegistryController) Ping() {

	registry := new(models.Registry)
	json.Unmarshal(this.Ctx.Input.RequestBody, &registry)

	rs := base.RegistryService{Registry: registry}

	this.Data["json"] = rs.Ping()

	this.ServeJSON(false)
}

// @Title TypeInfos
// @Description type Infos
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /typeInfos [post]
func (this *RegistryController) TypeInfos() {

	rs := base.RegistryService{}

	this.Data["json"] = rs.TypeInfos()

	this.ServeJSON(false)
}
