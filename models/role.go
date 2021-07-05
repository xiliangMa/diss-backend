package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Role struct {
	Id          int       `orm:"auto;pk" description:"角色ID"`
	Code        string    `orm:"unique"  description:"角色代码"`
	DispName    string    `orm:"size(64)" description:"显示名"`
	Description string    `orm:"size(64)" description:"描述"`
	Users       []*User   `orm:"reverse(many)" description:"用户列表"` // 以casbin中的关联优先
	Modules     []*Module `orm:"reverse(many)" description:"模块列表"` // 使用casbin中的关联，ORM中不做关联操作
	CreateTime  int64     `orm:"default(0);" description:"(创建时间)"`
	UpdateTime  int64     `orm:"default(0)" description:"(更新时间)"`
}

type RoleInterface interface {
	Add() Result
	List(from, limit int) Result
	RoleList() ([]*Role, int64, error)
	Update() Result
	Delete() Result
	AddUsers() Result
	RemoveUsers() Result
	PolicyList() Result
	AddPolicy() Result
	RemovePolicy() Result
	UpdatePolicy() Result
}

func (this *Role) Add() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	// 检查重名
	roleQuery := Role{DispName: this.DispName}
	roleObj, _ := roleQuery.Get()
	if roleObj != nil && roleObj.Id != 0 {
		ResultData.Message = "Role Name is Exist"
		ResultData.Code = utils.RoleExistErr
		return ResultData
	}

	this.CreateTime = time.Now().UnixNano()
	this.UpdateTime = time.Now().UnixNano()
	if this.Code == "" {
		this.Code = utils.GetRoleString(utils.GenRandomString(8))
	}
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddRoleErr
		logs.Error("Add Role failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) List(from, limit int) Result {
	var ResultData Result

	roleList, total, err := this.RoleList(from, limit)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRoleErr
		logs.Error("Get Role failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = roleList
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *Role) RoleList(from, limit int) (roleLists []*Role, count int64, err error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var roleList []*Role = nil
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.Code != "" {
		cond = cond.And("Code", this.Code)
	}

	_, err = o.QueryTable(utils.Role).SetCond(cond).RelatedSel().Limit(limit, from).OrderBy("-create_time").All(&roleList)

	total, _ := o.QueryTable(utils.Role).SetCond(cond).Count()
	return roleList, total, err
}

func (this *Role) Get() (*Role, error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var role *Role
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.Code != "" {
		cond = cond.And("Code", this.Code)
	}
	if this.DispName != "" {
		cond = cond.And("DispName", this.DispName)
	}

	err := o.QueryTable(utils.Role).SetCond(cond).RelatedSel().One(&role)

	return role, err
}

func (this *Role) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	roleObj := Role{}
	roleObj.Id = this.Id

	roleList, total, _ := roleObj.RoleList(0, 0)
	if total > 0 {
		roleData := roleList[0]
		this.CreateTime = roleData.CreateTime
		this.UpdateTime = time.Now().UnixNano()
		_, err := o.Update(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.EditRoleErr
			logs.Error("Edit Role %s failed, code: %d, err: %s", this.Code, ResultData.Code, ResultData.Message)
			return ResultData
		}
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	} else {
		ResultData.Message = "No RoleList Id."
		ResultData.Code = utils.DeleteRoleErr
		return ResultData
	}

	_, err := o.QueryTable(utils.Role).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteRoleErr
		logs.Error("Delete Role failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *Role) AddUsers() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	if len(this.Users) < 1 {
		ResultData.Code = utils.AddUserToRoleErr
		msg := fmt.Sprintf("Add User To Role failed, No User Info , code: %d", ResultData.Code)
		ResultData.Message = msg
	}

	for _, user := range this.Users {
		_, err := GlobalCasbin.Enforcer.AddRoleForUser(user.Name, utils.GetRoleString(this.Code))
		if err != nil {
			logs.Warn("Add User To Role failed, code: %d, error : %s", ResultData.Code, err)
		}
	}
	GlobalCasbin.Enforcer.LoadPolicy()

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) RemoveUsers() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	if len(this.Users) < 1 {
		ResultData.Code = utils.RemoveUserFromRoleErr
		msg := fmt.Sprintf("Remove User From Role failed, No User Info , code: %d", ResultData.Code)
		ResultData.Message = msg
	}

	for _, user := range this.Users {
		_, err := GlobalCasbin.Enforcer.DeleteRoleForUser(user.Name, utils.GetRoleString(this.Code))
		if err != nil {
			logs.Warn("Remove User From Role failed, No User Info , code: %d, error : %s", ResultData.Code, err)
		}
	}
	GlobalCasbin.Enforcer.LoadPolicy()

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) PolicyList() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	perm := GlobalCasbin.Enforcer.GetPermissionsForUser(utils.GetRoleString(this.Code))
	modulelist := []string{}
	for _, po := range perm {
		modulelist = append(modulelist, po[1])
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = modulelist
	return ResultData
}

func (this *Role) AddPolicy() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	if len(this.Modules) < 1 {
		ResultData.Code = utils.AddPolicyErr
		msg := fmt.Sprintf("Add AddPolicy failed, No Modules Info , code: %d", ResultData.Code)
		ResultData.Message = msg
	}

	for _, module := range this.Modules {
		_, err := GlobalCasbin.Enforcer.AddPolicy(utils.GetRoleString(this.Code), module.Code, "-")

		if err != nil {
			logs.Warn("Add AddPolicy failed, code: %d, error : %s", ResultData.Code, err)
		}
	}
	GlobalCasbin.Enforcer.LoadPolicy()

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) RemovePolicy() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	if len(this.Modules) < 1 {
		ResultData.Code = utils.RemovePolicyErr
		msg := fmt.Sprintf("Remove Policy failed, No Modules Info , code: %d", ResultData.Code)
		ResultData.Message = msg
	}

	for _, module := range this.Modules {
		_, err := GlobalCasbin.Enforcer.RemovePolicy(utils.GetRoleString(this.Code), module.Code, "")
		if err != nil {
			logs.Warn("Add Policy failed, code: %d, error : %s", ResultData.Code, err)
		}
	}
	GlobalCasbin.Enforcer.LoadPolicy()

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Role) UpdatePolicy() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	if len(this.Modules) < 1 {
		ResultData.Code = utils.RemovePolicyErr
		msg := fmt.Sprintf("Update Policy failed, No Modules Info, code: %d", ResultData.Code)
		ResultData.Message = msg
	}

	// 更新修改时间
	roleObj := Role{}
	roleObj.Code = utils.GetRoleString(this.Code)
	roleList, total, _ := roleObj.RoleList(0, 0)
	if total > 0 {
		roleData := roleList[0]
		roleData.UpdateTime = time.Now().UnixNano()
		_, err := o.Update(roleData)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.EditRoleErr
			logs.Error("Edit Role %s failed, code: %d, err: %s", this.Code, ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	// 移除之前的全部权限
	modules := GlobalCasbin.Enforcer.GetPermissionsForUser(utils.GetRoleString(this.Code))
	for _, mo := range modules {
		if this.Code != utils.GetRoleString(System_Role) && mo[1] != Permission_AuthManage {
			_, err := GlobalCasbin.Enforcer.RemovePolicy(utils.GetRoleString(this.Code), mo[1], "")
			if err != nil {
				logs.Warn("Remove Policy failed, code: %d, error : %s", ResultData.Code, err)
			}
		}
	}

	// 添加新指定的权限
	for _, module := range this.Modules {
		_, err := GlobalCasbin.Enforcer.AddPolicy(utils.GetRoleString(this.Code), module.Code, "-")
		GlobalCasbin.Enforcer.LoadPolicy()
		if err != nil {
			logs.Warn("Add AddPolicy failed, code: %d, error : %s", ResultData.Code, err)
		}
	}
	GlobalCasbin.Enforcer.LoadPolicy()

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
