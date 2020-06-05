package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 容器命令历史接口列表
type ContainerCmdHistoryController struct {
	beego.Controller
}

// @Title GetContainerCmdHistorys
// @Description Get Container CmdHistory  List
// @Param token header string true "authToken"
// @Param containerId path string "" true "containerId"
// @Param body body models.CmdHistory false "容器命令历史"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ContainerCmdHistoryController) GetContainerCmdHistorys() {
	containerId := this.GetString(":containerId")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	cmdHistory := new(models.CmdHistory)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cmdHistory)
	cmdHistory.ContainerId = containerId
	cmdHistory.Type = models.Cmd_History_Type_Container
	this.Data["json"] = cmdHistory.List(from, limit)
	this.ServeJSON(false)

}
