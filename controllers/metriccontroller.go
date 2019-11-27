package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// [以下为内部接口] Metric object api list
type MetricController struct {
	beego.Controller
}

// @Title HostInfo
// @Description HostMetricBasicInfo
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /hostinfo [post]
func (this *MetricController) ContainerSummary() {
	hostname := this.GetString("hostname")
	hostInfo := models.GetHostMetricInfo_M(hostname)

	this.Data["json"] = hostInfo
	this.ServeJSON(false)
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

// @Title ContainerSummary
// @Description Get Container Summary counts
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /containersummary [post]
func (this *MetricController) GetHost() {
	hostname := this.GetString("hostname")
	containerSummary := models.GetContainerSummaryInfo(hostname)

	this.Data["json"] = containerSummary
	this.ServeJSON(false)
}
