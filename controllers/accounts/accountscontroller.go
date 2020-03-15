package accounts

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// Account（租户） object api list
type AccountsController struct {
	beego.Controller
}

// @Title GetAccoounts
// @Description Get Accoounts List
// @Param token header string true "authToken"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *AccountsController) GetHostConfigList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	accounts := new(models.Accounts)
	//json.Unmarshal(this.Ctx.Input.RequestBody, &accounts)
	this.Data["json"] = accounts.List(from, limit)
	this.ServeJSON(false)

}
