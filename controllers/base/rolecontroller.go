package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

// 角色接口列表
type RoleController struct {
	beego.Controller
}

// @Title GetRoles
// @Description Get Roles
// @Param token header string true "authToken"
// @Param body body models.Role false "角色信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /rolelist [post]
func (this *RoleController) RoleList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)
	this.Data["json"] = role.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddRole
// @Description Add Role
// @Param token header string true "authToken"
// @Param body body models.Role false "角色及用户信息"
// @Success 200 {object} models.Result
// @router /role [post]
func (this *RoleController) AddRole() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.Add()
	this.ServeJSON(false)
}

// @Title AddUser
// @Description Add Role
// @Param token header string true "authToken"
// @Param body body models.Role false "角色及用户信息"
// @Success 200 {object} models.Result
// @router /binduser [post]
func (this *RoleController) AddUserToRole() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.AddUser()
	this.ServeJSON(false)
}

// @Title RemoveUser
// @Description Remove Role Of User
// @Param token header string true "authToken"
// @Param body body models.Role false "角色及用户信息"
// @Success 200 {object} models.Result
// @router /user [delete]
func (this *RoleController) RemoveUserOfRole() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.RemoveUser()
	this.ServeJSON(false)
}

// @Title PolicyList
// @Description list Policy for Role
// @Param token header string true "authToken"
// @Param body body models.Role false "Role信息"
// @Success 200 {object} models.Result
// @router /policylist [post]
func (this *RoleController) PolicyList() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.PolicyList()
	this.ServeJSON(false)
}

// @Title AddPolicy
// @Description Add Policy for Role
// @Param token header string true "authToken"
// @Param body body models.Role false "Policy信息"
// @Success 200 {object} models.Result
// @router /policy [post]
func (this *RoleController) AddPolicy() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.AddPolicy()
	this.ServeJSON(false)
}

// @Title RemovePolicy
// @Description Remove Policy for Role
// @Param token header string true "authToken"
// @Param body body models.Role false "Policy信息"
// @Success 200 {object} models.Result
// @router /policy [delete]
func (this *RoleController) RemovePolicy() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.RemovePolicy()
	this.ServeJSON(false)
}

// @Title UpdatePolicy
// @Description Update Policy for Role
// @Param token header string true "authToken"
// @Param body body models.Role false "Policy信息"
// @Success 200 {object} models.Result
// @router /policy [put]
func (this *RoleController) UpdatePolicy() {
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)

	this.Data["json"] = role.UpdatePolicy()
	this.ServeJSON(false)
}

// @Title UpdateRole
// @Description Update Role
// @Param token header string true "authToken"
// @Param roleId path string "" true "roleId"
// @Param body body models.Role true "主机配置信息"
// @Success 200 {object} models.Result
// @router /:roleId [put]
func (this *RoleController) UpdateRole() {
	roleId, _ := this.GetInt(":roleId")
	role := new(models.Role)
	json.Unmarshal(this.Ctx.Input.RequestBody, &role)
	role.Id = roleId
	this.Data["json"] = role.Update()
	this.ServeJSON(false)
}

// @Title DeleteRole
// @Description Delete Role
// @Param token header string true "authToken"
// @Param roleId path string "" true "roleId"
// @Success 200 {object} models.Result
// @router /:roleId [delete]
func (this *RoleController) DeleteRole() {
	roleId, _ := this.GetInt(":roleId")
	role := new(models.Role)
	role.Id = roleId

	result := role.Delete()

	this.Data["json"] = result
	this.ServeJSON(false)
}
