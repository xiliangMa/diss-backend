package statistics

import (
	"github.com/astaxie/beego"
	ss "github.com/xiliangMa/diss-backend/service/statistics"
)

// Package statistics api list
type PackageStatisticsController struct {
	beego.Controller
}

// @Title GetPackacgeStatistics
// @Description Get Package Statistics
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param hostId query string "" false "hostId"
// @Success 200 {object} models.Result
// @router /hostpackage [get]
func (this *PackageStatisticsController) GetPackageStatistics() {
	hostId := this.GetString("hostId")
	HPStatisticsService := ss.PackageStatisticsService{hostId}
	this.Data["json"] = HPStatisticsService.GetHostPackageStatistics()
	this.ServeJSON(false)
}

// @Title GetDBImageStatistics
// @Description Get DB Image Statistics
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param hostId query string "" false "hostId"
// @Success 200 {object} models.Result
// @router /dbpackage [get]
func (this *PackageStatisticsController) GetDBImageStatistics() {
	hostId := this.GetString("hostId")
	HPStatisticsService := ss.PackageStatisticsService{hostId}
	this.Data["json"] = HPStatisticsService.GetDBImageStatistics()
	this.ServeJSON(false)
}
