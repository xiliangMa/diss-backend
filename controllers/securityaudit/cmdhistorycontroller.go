package securityaudit

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 命令历史接口列表
type CmdHistoryController struct {
	beego.Controller
}

// @Title GetCmdHistorys
// @Description Get CmdHistory List(主机：Type = Host 容器 Type = Container )
// @Param token header string true "authToken"
// @Param module header string true "moduleCode"
// @Param body body models.CmdHistory false "命令历史"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /cmdhistorys [post]
func (this *CmdHistoryController) GetCmdHistorys() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cmdHistory)
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}
