package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Groups struct {
	Id          string    `orm:"pk;" description:"(id)"`
	Name        string    `orm:"" description:"(镜像id)"`
	Type        int       `orm:"default(0)" description:"(分组类型 0 主机 1 容器)"`
	AccountName string    `orm:"" description:"(租户)"`
	CreateTime  time.Time `orm:"null;auto_now;type(datetime)" description:"(创建时间)"`
}

type GroupInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *Groups) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddGroupErr
		logs.Error("Add Group failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Groups) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var GroupList []*Groups
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	cond = cond.And("type", this.Type)
	_, err = o.QueryTable(utils.Groups).SetCond(cond).Limit(limit, from).All(&GroupList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetGroupErr
		logs.Error("Get Group List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Groups).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = GroupList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Groups) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditGroupErr
		logs.Error("Update Group: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Groups) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.Groups).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteGroupErr
		logs.Error("Delete Group failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
