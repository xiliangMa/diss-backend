package securitypolicy

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

// Virus Log api list
type VirusLogController struct {
	beego.Controller
}

// @Title GetVirusLog
// @Description Get Virus Log List
// @Param token header string true "authToken"
// @Param body body models.VirusScan false "病毒信息列表"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /viruslog [post]
func (this *VirusLogController) GetVirusLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	virusLog := new(models.VirusScan)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &virusLog)
	logs.Warn("err: ", err)
	this.Data["json"] = virusLog.List(from, limit)
	this.ServeJSON(false)
}
