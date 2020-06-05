package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 主机命令历史接口列表
type HostCmdHistoryController struct {
	beego.Controller
}

// @Title HostCmdHistory
// @Description Get HostCmdHistory List
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param body body models.CmdHistory false "主机命令历史信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *HostCmdHistoryController) GetHostCmdHistoryList() {
	hostId := this.GetString(":hostId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cmdHistory)
	cmdHistory.HostId = hostId
	cmdHistory.Type = models.Cmd_History_Type_Host
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}
