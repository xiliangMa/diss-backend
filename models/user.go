package models

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type User struct {
	Id         int    `orm:"auto;pk" description:"用户ID"`
	Name       string `orm:"unique" description:"用户名"`
	DispName   string `orm:"size(64)" description:"显示名"`
	Password   string `description:"用户密码"`
	Role       *Role  `orm:"rel(fk);null" description:"角色"`
	CreateTime int64  `orm:"default(0);" description:"(创建时间)"`
	UpdateTime int64  `orm:"default(0)" description:"(更新时间)"`
}

type UserInterface interface {
	List(from, limit int) Result
	Add() Result
	Update() Result
	Delete() Result
	UserList() ([]*User, int64, error)
}

func (this *User) Add() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	userQuery := User{}
	userQuery.Name = this.Name

	_, count, _ := userQuery.UserList(0, 1)
	if count > 0 {
		ResultData.Message = fmt.Sprintf("User already exist, code: %d", utils.UserExistErr)
		ResultData.Code = utils.UserExistErr
		logs.Error(ResultData.Message)
		return ResultData
	}

	this.CreateTime = time.Now().UnixNano()
	password := utils.MD5(this.Password)
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))
	this.Password = passwordBase64
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddUserErr
		logs.Error("Add User failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *User) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	userObj := User{}
	userObj.Id = this.Id

	userList, total, _ := userObj.UserList(0, 1)
	if total > 0 {
		userData := userList[0]
		this.CreateTime = userData.CreateTime
		this.UpdateTime = time.Now().UnixNano()
		_, err := o.Update(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.EditUserErr
			logs.Error("Edit User %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
			return ResultData
		}

	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *User) List(from, limit int) Result {
	var ResultData Result

	userList, total, err := this.UserList(from, limit)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetUserErr
		logs.Error("Get User failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = userList
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *User) UserList(from, limit int) (userLists []*User, count int64, err error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var userList []*User = nil
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}
	if this.Password != "" {
		cond = cond.And("password", this.Password)
	}

	_, err = o.QueryTable(utils.User).SetCond(cond).Limit(limit, from).OrderBy("-create_time").All(&userList)
	for _, user := range userList {
		userRole, _ := Enforcer.GetRolesForUser(user.Name)
		if len(userRole) > 0 {
			roleQuery := Role{Code: utils.GetRoleString(userRole[0])}
			roleObj, count, _ := roleQuery.RoleList(0, 1)
			if count > 0 {
				user.Role = roleObj[0]
			}
		}
	}

	total, _ := o.QueryTable(utils.User).SetCond(cond).Count()
	return userList, total, err
}

func (this *User) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	} else {
		ResultData.Message = "No User Id."
		ResultData.Code = utils.DeleteUserErr
		return ResultData
	}

	_, err := o.QueryTable(utils.User).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteUserErr
		logs.Error("Delete User failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
