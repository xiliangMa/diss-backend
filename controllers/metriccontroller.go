package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Metric object api list [inner]
type MetricController struct {
	beego.Controller
}

// @Title GetContainerList
// @Description Get ContainerList
// @Param token header string true "Auth token"
// @Param name query string false "host name"
// @Success 200 {object} models.Result
// @router /containerlist [post]
func (this *MetricController) ContainerList() {
	hostname := this.GetString("name")
	containerList := models.GetContainerListMetricInfo(hostname)

	this.Data["json"] = containerList
	this.ServeJSON(false)
}

// @Title HostInfo
// @Description HostMetricBasicInfo
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /hostinfo [post]
func (this *MetricController) HostInfo() {
	hostname := this.GetString("hostname")
	hostInfo := models.GetHostMetricInfo_M(hostname)

	this.Data["json"] = hostInfo
	this.ServeJSON(false)
}

// @Title GetHost
// @Description Get one Host
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /gethost [post]
func (this *MetricController) GetHost() {
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
func (this *MetricController) DeleteHost() {
	id, _ := this.GetInt(":id")
	this.Data["json"] = models.DeleteHost(id)
	this.ServeJSON(false)

}
