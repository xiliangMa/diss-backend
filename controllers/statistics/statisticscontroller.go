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
// @Description Get Asset Statistics (主机数、容器数)
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /asset [post]
func (this *StatisticsController) GetAssetStatistics() {
	statisticsService := ss.StatisticsService{nil, nil}
	this.Data["json"] = statisticsService.GetAssetStatistics()
	this.ServeJSON(false)

}

// @Title GetBnechMarkProportionStatistics
// @Description Get BnechMark Proportion Statistics (docker基线、kubernetes基线)
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /bmp [post]
func (this *StatisticsController) GetBnechMarkProportionStatistics() {
	statisticsService := ss.StatisticsService{nil, nil}
	this.Data["json"] = statisticsService.GetBnechMarkProportionStatistics()
	this.ServeJSON(false)

}
