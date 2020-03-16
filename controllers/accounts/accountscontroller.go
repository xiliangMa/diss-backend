package accounts

import (
	"encoding/json"
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

// @Title GetGroups
// @Description Get Groups List（获取租户下的分组 主机/容器分组）
// @Param token header string true "authToken"
// @Param body body models.Groups false "分组信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /:accountName/groups [post]
func (this *AccountsController) GetGroupsList() {
	accountName := this.GetString(":accountName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	groups := new(models.Groups)
	json.Unmarshal(this.Ctx.Input.RequestBody, &groups)
	groups.AccountName = accountName
	this.Data["json"] = groups.List(from, limit)
	this.ServeJSON(false)

}
