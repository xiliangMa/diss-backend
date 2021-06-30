package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 用户接口列表
type UserController struct {
	beego.Controller
}

// @Title GetUserEvents
// @Description Get User event List
// @Param token header string true "authToken"
// @Param body body models.UserEvent false "用户事件"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /events [post]
func (this *UserController) UserEventList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	userEvent := new(models.UserEvent)
	json.Unmarshal(this.Ctx.Input.RequestBody, &userEvent)
	this.Data["json"] = userEvent.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddUser
// @Description Add User
// @Param token header string true "authToken"
// @Param body body models.User false "用户信息"
// @Success 200 {object} models.Result
// @router / [post]
func (this *UserController) AddUser() {
	user := new(models.User)
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)

	this.Data["json"] = user.Add()
	this.ServeJSON(false)
}

// @Title GetUsers
// @Description Get Users
// @Param token header string true "authToken"
// @Param body body models.User false "用户信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /userlist [post]
func (this *UserController) UserList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	user := new(models.User)
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	this.Data["json"] = user.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateUser
// @Description Update User
// @Param token header string true "authToken"
// @Param userId path string "" true "userId"
// @Param body body models.User true "用户信息"
// @Success 200 {object} models.Result
// @router /:userId [put]
func (this *UserController) UpdateUser() {
	userId, _ := this.GetInt(":userId")
	user := new(models.User)
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	user.Id = userId
	this.Data["json"] = user.Update()
	this.ServeJSON(false)
}

// @Title DeleteUser
// @Description Delete User
// @Param token header string true "authToken"
// @Param userId path string "" true "userId"
// @Success 200 {object} models.Result
// @router /:userId [delete]
func (this *UserController) DeleteUser() {
	userId, _ := this.GetInt(":userId")
	user := new(models.User)
	user.Id = userId

	result := user.Delete()

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title AddUser
// @Description Add User
// @Param token header string true "authToken"
// @Param body body models.User false "用户信息"
// @Success 200 {object} models.Result
// @router /role [post]
func (this *UserController) AddRoleForUser() {
	user := new(models.User)
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)

	this.Data["json"] = user.AddRole()
	this.ServeJSON(false)
}

// @Title UpdateUser
// @Description Update User
// @Param token header string true "authToken"
// @Param userId path string "" true "userId"
// @Param body body models.User true "用户信息"
// @Success 200 {object} models.Result
// @router /role [put]
func (this *UserController) UpdateRoleForUser() {
	userId, _ := this.GetInt(":userId")
	user := new(models.User)
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	user.Id = userId
	this.Data["json"] = user.UpdateRole()
	this.ServeJSON(false)
}

// @Title DeleteUser
// @Description Delete User
// @Param token header string true "authToken"
// @Param userId path string "" true "userId"
// @Success 200 {object} models.Result
// @router /user [delete]
func (this *UserController) RemoveRoleUser() {
	userId, _ := this.GetInt(":userId")
	user := new(models.User)
	user.Id = userId

	result := user.RemoveRole()

	this.Data["json"] = result
	this.ServeJSON(false)
}
