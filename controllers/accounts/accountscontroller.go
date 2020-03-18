package accounts

import (
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
)

// Account（租户） object api list
type AccountsController struct {
	beego.Controller
}

// @Title GetAccoounts
// @Description Get Accoounts List
// @Param token header string true "authToken"
// @Param user header string "admin" true "diss api 系统的登入用户"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *AccountsController) GetHostConfigList() {
	accountName := models.Account_Admin
	user := this.Ctx.Input.Header("user")
	if user != models.Account_Admin && user != "" {
		accountUsers := models.AccountUsers{}
		accountUsers.UserName = user
		err, account := accountUsers.GetAccountByUser()
		accountName = account
		if err != nil {
			this.Data["json"] = models.Result{Code: utils.NoAccountUsersErr, Data: nil, Message: err.Error()}
			this.ServeJSON(false)
		}
	}
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	accounts := new(models.Accounts)
	accounts.Name = accountName
	this.Data["json"] = accounts.List(from, limit)
	this.ServeJSON(false)

}
