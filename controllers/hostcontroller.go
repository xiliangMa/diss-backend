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
// @Param name query string false "host name"
// @Param ip query string false "host ip"
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
// @Param Host body models.Host true "host object"
// @Success 200 {object} models.Result
// @router /addhost [post]
func (this *HostController) AddHost() {
	var h models.Host
	json.Unmarshal(this.Ctx.Input.RequestBody, &h)

	existhost := models.GetHostInternal(h.HostName)
	if existhost == nil {
		this.Data["json"] = models.AddHost(&h)
	} else {
		existhost["Message"] = "Host Exist"
		this.Data["json"] = existhost
	}

	this.ServeJSON(false)
}

// @Title GetHost
// @Description Get one Host
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /gethost [post]
func (this *HostController) GetHost() {
	hostname := this.GetString("hostname")
	this.Data["json"] = models.GetHost(hostname)
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
