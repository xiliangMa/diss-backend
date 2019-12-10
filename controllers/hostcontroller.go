package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Hosts object api list
type HostController struct {
	beego.Controller
}

// @Title GetHost
// @Description Get Hosts
// @Param token header string true "Auth token"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *HostController) HostList() {
	name := this.GetString("name")
	ip := this.GetString("ip")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = models.GetHostList(name, ip, from, limit)
	this.ServeJSON(false)

}

// @Title AddHost
// @Description Add Host
// @Param token header string true "Auth token"
// @Param Host body models.Host true "host object , remove CreateTime and UpdateTime before POST"
// @Success 200 {object} models.Result
// @router /addhost [post]
func (this *HostController) AddHost() {
	var h models.Host
	json.Unmarshal(this.Ctx.Input.RequestBody, &h)

	this.Data["json"] = models.AddHost_Processing(h, 0)
	this.ServeJSON(false)
}

// @Title EditHost
// @Description Edit Host
// @Param token header string true "Auth token"
// @Param Host body models.Host true "host object , remove CreateTime and UpdateTime before POST"
// @Success 200 {object} models.Result
// @router /edithost [post]
func (this *HostController) EditHost() {
	var h models.Host
	json.Unmarshal(this.Ctx.Input.RequestBody, &h)

	this.Data["json"] = models.AddHost_Processing(h, 1)
	this.ServeJSON(false)
}

// @Title GetHostWithContainer
// @Description Get one Host and its containers
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /gethost_container [post]
func (this *HostController) GetHostContainers() {
	hostname := this.GetString("hostname")

	this.Data["json"] = models.GetHostWithContainer_Processing(hostname)
	this.ServeJSON(false)
}

func (this *HostController) GetHost() {
	hostname := this.GetString("hostname")

	this.Data["json"] = models.GetHostWithContainer_Processing(hostname)
	this.ServeJSON(false)
}

// @Title GetHostWithContainer
// @Description Get one Host and its images
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /gethost_images [post]
func (this *HostController) GetHostImages() {
	hostname := this.GetString("hostname")

	this.Data["json"] = models.GetHostWithImage_Processing(hostname)
	this.ServeJSON(false)
}


// @Title DelHost
// @Description Delete Host
// @Param token header string true "Auth token"
// @Param id path int true "host id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *HostController) DeleteHost() {
	id, _ := this.GetInt(":id")

	this.Data["json"] = models.DeleteHost(id)
	this.ServeJSON(false)
}
