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
// @Param hostid path string "" true "hostid"
// @Param targetype query string "host" true "targetype"
// @Param containerid query string "" false "containerid"
////@Param from query int 0 true "from"
// @Param limit query int 20 true "limit"
// @Param starttime query string "" true "starttime"
// @Param totime query string "" true "totime"
// @Success 200 {object} models.Result
// @router /intrudedetect/:hostid [post]
func (this *IntrudeDetectLogController) GetIntrudeLogList() {
	starttime := this.GetString("starttime")
	totime := this.GetString("totime")
	hostid := this.GetString(":hostid")
	targetype := this.GetString("targetype")
	containerid := this.GetString("containerid")
	limit := this.GetString("limit")

	intrudelog := models.Internal_IntrudeDetectMetricInfo(hostid, targetype, containerid, starttime, totime, limit)

	this.Data["json"] = intrudelog
	this.ServeJSON(false)

}
