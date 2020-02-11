package asset

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// HostConfig object api list
type HostConfigController struct {
	beego.Controller
}

// @Title GetHostConfig
// @Description Get HostConfig
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *HostConfigController) GetHostConfigList() {
	name := this.GetString("name")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = models.GetHostConfigList(name, from, limit)
	this.ServeJSON(false)

}
