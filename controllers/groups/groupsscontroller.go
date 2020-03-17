package accounts

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

// Groups（分组） object api list
type GroupsController struct {
	beego.Controller
}

// @Title GetGroups
// @Description Get Groups List（获取租户下的分组 主机/容器分组）
// @Param token header string true "authToken"
// @Param user header string "admin" true "diss api 系统的登入用户"
// @Param body body models.Groups false "分组信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *GroupsController) GetGroupsList() {
	accountName := models.Account_Admin
	if this.Ctx.Input.Header("user") != models.Account_Admin {
		accountUsers := models.AccountUsers{}
		accountUsers.UserName = this.Ctx.Input.Header("user")
		err,  account := accountUsers.GetAccountByUser()
		accountName = account
		if err != nil {
			this.Data["json"] = models.Result{Code: utils.NoAccountUsersErr, Data: nil, Message: err.Error()}
			this.ServeJSON(false)
		}
	}
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	groups := new(models.Groups)
	json.Unmarshal(this.Ctx.Input.RequestBody, &groups)
	groups.AccountName = accountName
	this.Data["json"] = groups.List(from, limit)
	this.ServeJSON(false)

}
