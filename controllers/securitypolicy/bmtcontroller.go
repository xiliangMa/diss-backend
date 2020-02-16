package securitypolicy

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Bench Mark Template object api list
type BMTController struct {
	beego.Controller
}

// @Title GetBMTList
// @Description Get BenchMarkTemplate List
// @Param token header string true "auth token"
// @Param name query string "" false "name"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *BMTController) GetBMTList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	benchMarkTemplate := new(models.BenchMarkTemplate)
	benchMarkTemplate.Name = this.GetString("name")
	this.Data["json"] = benchMarkTemplate.List(from, limit)
	this.ServeJSON(false)

}
