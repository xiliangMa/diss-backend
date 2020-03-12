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
// @Description Get Asset Statistics (资产概况：主机数、容器数)
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /asset [post]
func (this *StatisticsController) GetAssetStatistics() {
	statisticsService := ss.StatisticsService{nil, nil}
	this.Data["json"] = statisticsService.GetAssetStatistics()
	this.ServeJSON(false)

}

// @Title GetBnechMarkProportionStatistics
// @Description Get BnechMark Proportion Statistics (安全基线占比：docker基线、kubernetes基线)
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /bmp [post]
func (this *StatisticsController) GetBnechMarkProportionStatistics() {
	statisticsService := ss.StatisticsService{nil, nil}
	this.Data["json"] = statisticsService.GetBnechMarkProportionStatistics()
	this.ServeJSON(false)

}

// @Title GetBnechMarkSummaryStatistics
// @Description Get BnechMark Summary Statistics (安全基线摘要统计)
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /bms [post]
func (this *StatisticsController) GetBnechMarkSummaryStatistics() {
	statisticsService := ss.StatisticsService{nil, nil}
	this.Data["json"] = statisticsService.GetBnechMarkSummaryStatistics()
	this.ServeJSON(false)
}
