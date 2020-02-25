package securitypolicy

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Intrude Detect Log api list
type IntrudeDetectLogController struct {
	beego.Controller
}

// @Title GetIntrudeLogList
// @Description Get IntrudeLog List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param targeType query string "host" true "targeType"
// @Param containerId query string "" false "containerId"
////@Param from query int 0 true "from"
// @Param limit query int 20 true "limit"
// @Param starTime query string "" true "starTime"
// @Param toTime query string "" true "toTime"
// @Success 200 {object} models.Result
// @router /intrudedetect/:hostId [post]
func (this *IntrudeDetectLogController) GetIntrudeLogList() {
	startTime := this.GetString("startTime")
	toTime := this.GetString("toTime")
	hostId := this.GetString(":hostId")
	targeType := this.GetString("targeType")
	containerId := this.GetString("containerId")
	limit := this.GetString("limit")

	intrudelog := models.Internal_IntrudeDetectMetricInfo(hostId, targeType, containerId, startTime, toTime, limit)

	this.Data["json"] = intrudelog
	this.ServeJSON(false)

}
