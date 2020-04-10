package base

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
// @Description Get Groups List（获取租户下的分组 主机/容器分组, 暂不支持主机、容器对象参数传入查询）
// @Param token header string true "authToken"
// @Param accountName query string "admin" true "diss api 系统登入用户的所属租户"
// @Param body body models.Groups false "分组信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *GroupsController) GetGroupsList() {
	accountName := this.GetString("accountName")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	groups := new(models.Groups)
	json.Unmarshal(this.Ctx.Input.RequestBody, &groups)
	groups.AccountName = accountName
	this.Data["json"] = groups.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddGroup
// @Description Add Group
// @Param token header string true "authToken"
// @Param body body models.Groups false "分组信息"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *GroupsController) AddGroup() {
	groups := new(models.Groups)
	json.Unmarshal(this.Ctx.Input.RequestBody, &groups)
	this.Data["json"] = groups.Add()
	this.ServeJSON(false)
}

// @Title DeleteGroup
// @Description Delete Group
// @Param token header string true "authToken"
// @Param grouprId path string "" true "grouprId"
// @Success 200 {object} models.Result
// @router /:grouprId [delete]
func (this *GroupsController) DeleteGroup() {
	grouprId := this.GetString(":grouprId")
	group := new(models.Groups)
	group.Id = grouprId
	this.Data["json"] = group.Delete()
	this.ServeJSON(false)
}

// @Title GetContainers
// @Description Get Groups List（获取分组下的容器）
// @Param token header string true "authToken"
// @Param user query string "admin" true "diss api 系统的登入用户"
// @Param body body models.ContainerConfig false "分组信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /containers [post]
func (this *GroupsController) GetContainersList() {
	accountName := models.Account_Admin
	user := this.GetString("user")
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
	containerConfig := new(models.ContainerConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &containerConfig)
	containerConfig.AccountName = accountName
	this.Data["json"] = containerConfig.List(from, limit, true)
	this.ServeJSON(false)

}
