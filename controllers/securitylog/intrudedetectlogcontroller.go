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
// @Param token header string true "auth token"
// @Param hostid query string "" false "hostid"
// @Param targetype query string "host" false "targetype"
////@Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Param starttime query string "" false "starttime"
// @Param totime query string "" false "totime"
// @Success 200 {object} models.Result
// @router /intrudedetect [post]
func (this *IntrudeDetectLogController) GetIntrudeLogList() {
	starttime := this.GetString("starttime")
	totime := this.GetString("totime")
	hostid := this.GetString("hostid")
	targetype := this.GetString("targetype")
	limit := this.GetString("limit")

	intrudelog := models.Internal_IntrudeDetectMetricInfo(hostid, targetype, starttime, totime, limit)

	this.Data["json"] = intrudelog
	this.ServeJSON(false)

}
