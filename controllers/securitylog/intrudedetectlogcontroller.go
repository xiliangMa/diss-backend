package securitypolicy

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// Intrude Detect Log api list
type IntrudeDetectLogController struct {
	web.Controller
}

// @Title GetIntrudeDetectLogInfo
// @Description Get IntrudeDetectLogInfo (查询主机/容器的入侵日志， 主机：targeType = host 容器：targeType = container)
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.IntrudeDetectLog false "入侵检测日志信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /intrudedetect/:hostId [post]
func (this *IntrudeDetectLogController) GetIntrudeDetectLogInfo() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	hostId := this.GetString(":hostId")
	intrudeDetectLog := new(models.IntrudeDetectLog)
	intrudeDetectLog.HostId = hostId
	json.Unmarshal(this.Ctx.Input.RequestBody, &intrudeDetectLog)
	//var securityLogService = ssl.SecurityLogService{nil, intrudeDetectLog}
	//
	//this.Data["json"] = securityLogService.GetIntrudeDetectLogInfo()
	this.Data["json"] = intrudeDetectLog.List(from, limit)
	this.ServeJSON(false)
}

// @Title GetIntrudeDetectLog
// @Description Get IntrudeDetectLog List (1. 根据 TargeType = host 查询主机基线日志 2. 根据 TargeType = container 如果快速查询所有容器日志可以设置 ContianerId =All)
// @Param token header string true "authToken"
// @Param body body models.IntrudeDetectLog false "入侵检测日志信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /idls [post]
func (this *IntrudeDetectLogController) GetIntrudeDetectLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	intrudeDetectLog := new(models.IntrudeDetectLog)
	json.Unmarshal(this.Ctx.Input.RequestBody, &intrudeDetectLog)
	this.Data["json"] = intrudeDetectLog.List1(from, limit)
	this.ServeJSON(false)
}
