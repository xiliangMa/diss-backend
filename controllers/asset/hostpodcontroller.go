package asset

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// HostInfo object api list
type HostPodController struct {
	beego.Controller
}

// @Title GetHostPod
// @Description Get HostPod
// @Param token header string true "auth token"
// @Param name query string "" false "pod name"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /info [post]
func (this *HostPodController) GetHostPodList() {
	name := this.GetString("name")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = models.GetHostPodList(name, from, limit)
	this.ServeJSON(false)

}
