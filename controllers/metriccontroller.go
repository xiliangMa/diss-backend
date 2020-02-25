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
// @Param token header string true "authToken"
// @Param hostName query string false "hostName"
// @Success 200 {object} models.Result
// @router /hostinfo [post]
func (this *MetricController) HostMetricInfo() {
	hostName := this.GetString("hostName")
	hostData := models.Internal_HostMetricInfo_M(hostName)

	this.Data["json"] = hostData
	this.ServeJSON(false)
}

// @Title GetContainerList
// @Description Get ContainerList
// @Param token header string true "authToken"
// @Param hostName query string false "hostName"
// @Success 200 {object} models.Result
// @router /containerlist [post]
func (this *MetricController) ContainerList() {
	hostName := this.GetString("hostName")
	containerList := models.Internal_ContainerListMetricInfo(hostName)

	this.Data["json"] = containerList
	this.ServeJSON(false)
}

// @Title ContainerSummary
// @Description Get Container Summary counts
// @Param token header string true "authToken"
// @Param hostName query string false "hostName"
// @Success 200 {object} models.Result
// @router /containersummary [post]
func (this *MetricController) ContainerSummary() {
	hostName := this.GetString("hostName")
	containerSummary := models.Internal_ContainerSummaryInfo(hostName)

	this.Data["json"] = containerSummary
	this.ServeJSON(false)
}

// @Title GetImageList
// @Description Get ImageList
// @Param token header string true "authToken"
// @Param hostName query string false "hostName"
// @Success 200 {object} models.Result
// @router /imagelist [post]
func (this *MetricController) ImageList() {
	hostName := this.GetString("hostName")
	containerList := models.Internal_ImageListMetricInfo(hostName)

	this.Data["json"] = containerList
	this.ServeJSON(false)
}
