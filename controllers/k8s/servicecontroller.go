package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Service 接口
type ServiceController struct {
	beego.Controller
}

// @Title GetService
// @Description Get Service List
// @Param token header string true "authToken"
// @Param body body models.Service false "Service"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ServiceController) GetServicesList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	Service := new(models.Service)
	json.Unmarshal(this.Ctx.Input.RequestBody, &Service)
	this.Data["json"] = Service.List(from, limit)
	this.ServeJSON(false)
}
