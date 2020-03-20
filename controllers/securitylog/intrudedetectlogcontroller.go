package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	msl "github.com/xiliangMa/diss-backend/models/securitylog"
)

// Intrude Detect Log api list
type IntrudeDetectLogController struct {
	beego.Controller
}

// @Title GetIntrudeDetectLogInfo
// @Description Get IntrudeDetectLogInfo (查询主机/容器的入侵日志， 主机：targeType = host 容器：targeType = container)
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body securitylog.IntrudeDetectLog false "入侵检测日志信息"
// @Success 200 {object} models.Result
// @router /intrudedetect/:hostId [post]
func (this *IntrudeDetectLogController) GetIntrudeDetectLogInfo() {
	hostId := this.GetString(":hostId")
	intrudeDetectLog := new(msl.IntrudeDetectLog)
	intrudeDetectLog.HostId = hostId
	json.Unmarshal(this.Ctx.Input.RequestBody, &intrudeDetectLog)
	//var securityLogService = ssl.SecurityLogService{nil, intrudeDetectLog}
	//
	//this.Data["json"] = securityLogService.GetIntrudeDetectLogInfo()
	this.Data["json"] = intrudeDetectLog.List(0, intrudeDetectLog.Limit)
	this.ServeJSON(false)
}


// @Title GetIntrudeDetectLog
// @Description Get IntrudeDetectLog List
// @Param token header string true "authToken"
// @Param body body securitylog.IntrudeDetectLog false "入侵检测日志信息"
// @Success 200 {object} models.Result
// @router /idls [post]
func (this *IntrudeDetectLogController) GetIntrudeDetectLogList() {
	intrudeDetectLog := new(msl.IntrudeDetectLog)
	json.Unmarshal(this.Ctx.Input.RequestBody, &intrudeDetectLog)
	this.Data["json"] = intrudeDetectLog.List(0, intrudeDetectLog.Limit)
	this.ServeJSON(false)
}