package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 用户接口口列表
type UserController struct {
	beego.Controller
}

// @Title GetUserEvents
// @Description Get User event List
// @Param token header string true "authToken"
// @Param body body models.UserEvent false "用户事件"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /events [post]
func (this *UserController) UserEventList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	userEvent := new(models.UserEvent)
	json.Unmarshal(this.Ctx.Input.RequestBody, &userEvent)
	this.Data["json"] = userEvent.List(from, limit)
	this.ServeJSON(false)
}
