package statistics

import (
	"github.com/astaxie/beego"
	ss "github.com/xiliangMa/diss-backend/service/statistics"
)

// Statistics api list
type StatisticsController struct {
	beego.Controller
}

// @Title GetAssetStatistics
// @Description Get Asset Statistics (资产概况：主机数、容器数、镜像仓库、镜像、集群数、Pod)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /asset [get]
func (this *StatisticsController) GetAssetStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetAssetStatistics()
	this.ServeJSON(false)
}

// @Title GetMirrorRiskStatistics
// @Description Get Mirror Risk Statistics
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /mirrorrisk [get]
func (this *StatisticsController) GetMirrorRiskStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetMirrorRiskStatistics()
	this.ServeJSON(false)
}

// @Title GetBenchMarkProportionStatistics
// @Description Get BnechMark Proportion Statistics (安全基线占比：docker基线、kubernetes基线)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /bmp [get]
func (this *StatisticsController) GetBenchMarkProportionStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetBenchMarkProportionStatistics()
	this.ServeJSON(false)
}

// @Title GetBnechMarkSummaryStatistics
// @Description Get BnechMark Summary Statistics (安全基线摘要统计)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /bms [get]
func (this *StatisticsController) GetBnechMarkSummaryStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetBnechMarkSummaryStatistics()
	this.ServeJSON(false)
}

// @Title GetIntrudeDetectLogStatistics
// @Description Get IntrudeDetect Log Statistics (入侵基线告警)
// @Param timeCycle query int 24 false "timecycle 时间周期"
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /idl [get]
func (this *StatisticsController) GetIntrudeDetectLogStatistics() {
	timeCycle, _ := this.GetInt("timeCycle")
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetIntrudeDetectLogStatistics(timeCycle)
	this.ServeJSON(false)
}

// @Title GetHostBnechMarkSummaryStatistics
// @Description Get Host BnechMark Summary Statistics (主机安全基线摘要统计)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param hostId path string "" true "hostId"
// @Success 200 {object} models.Result
// @router /:hostId/bms/host [get]
func (this *StatisticsController) GetHostBnechMarkSummaryStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetHostBnechMarkSummaryStatistics(this.GetString(":hostId"))
	this.ServeJSON(false)
}

// @Title GetOnlineProportionStatistics
// @Description Get BnechMark Proportion Statistics (主机在线占比：Online / Offline)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /onlineproportion [get]
func (this *StatisticsController) GetOnlineProportionStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetOnlineProportionStatistics()
	this.ServeJSON(false)
}

// @Title GetDissProportionStatistics
// @Description Get Diss Proportion Statistics (安全容器占比：Safe / Unsafe)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /dissproportion [get]
func (this *StatisticsController) GetDissProportionStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetDissProportionStatistics()
	this.ServeJSON(false)
}

// @Title GetWarningStatistics
// @Description Get Diss Proportion Statistics (安全容器占比：Safe / Unsafe)
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Success 200 {object} models.Result
// @router /warning [get]
func (this *StatisticsController) GetWarningStatistics() {
	statisticsService := ss.StatisticsService{nil, nil, nil}
	this.Data["json"] = statisticsService.GetWarningStatistics()
	this.ServeJSON(false)
}
