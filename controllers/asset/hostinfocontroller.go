package asset

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// HostConfig object api list
type HostInfoController struct {
	beego.Controller
}

// @Title GetHostConfigInfo
// @Description Get HostConfig
// @Param token header string true "auth token"
// @Param id query string "" false "host id"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /info [post]
func (this *HostInfoController) GetHostInfoList() {
	id := this.GetString("id")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = models.GetHostInfoList(id, from, limit)
	this.ServeJSON(false)

}


