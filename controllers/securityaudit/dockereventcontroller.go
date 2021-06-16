package securityaudit

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// docker审计 接口列表
type DockerEventController struct {
	beego.Controller
}

// @Title GetDockerEvents
// @Description Get DockerEven List
// @Param token header string true "authToken"
// @Param body body models.DockerEvent false "docker event 审计"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /dockerevents [post]
func (this *DockerEventController) GetDockerEvents() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	dockerEvent := new(models.DockerEvent)
	json.Unmarshal(this.Ctx.Input.RequestBody, &dockerEvent)
	this.Data["json"] = dockerEvent.List(from, limit)
	this.ServeJSON(false)

}
