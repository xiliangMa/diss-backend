package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Groups struct {
	Id          string        `orm:"pk;" description:"(id)"`
	FirstLevel  string        `orm:"unique" description:"(一级分组)"`
	SecondLevel string        `orm:"null" description:"(二级分组)"`
	ThirdLevel  string        `orm:"null" description:"(三级分组)"`
	Type        string        `orm:"default(Host)" description:"(All Host Container)"`
	AccountName string        `orm:"default(admin)" description:"(租户 默认 admin)"`
	CreateTime  time.Time     `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime  time.Time     `orm:"auto_now;type(datetime)" description:"(更新时间)"`
	HostConfig  []*HostConfig `orm:"reverse(many);null" description:"(主机列表)"`
}

type GroupInterface interface {
	Add() Result
	Delete() Result
	Update() Result
	Get() Result
	List(from, limit int) Result
}

func (this *Groups) Add() Result {
	o := orm.NewOrm()
	var ResultData Result
	uuid, _ := uuid.NewV4()
	this.Id = uuid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddGroupErr
		logs.Error("Add Group failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return this.List(0, 0)
}

func (this *Groups) List(from, limit int) Result {
	o := orm.NewOrm()
	var groupList []*Groups
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.FirstLevel != "" {
		cond = cond.And("first_level__contains", this.FirstLevel)
	}
	if this.SecondLevel != "" {
		cond = cond.And("second_level__contains", this.SecondLevel)
	}
	if this.ThirdLevel != "" {
		cond = cond.And("third_level__contains", this.ThirdLevel)
	}
	if this.Type != "" && this.Type != All {
		cond = cond.And("type", this.Type)
	}

	_, err = o.QueryTable(utils.Groups).SetCond(cond).Limit(limit, from).All(&groupList)
	for _, group := range groupList {
		o.LoadRelated(group, "HostConfig")
	}
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetGroupErr
		logs.Error("Get Group List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Groups).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = groupList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Groups) Update() Result {
	o := orm.NewOrm()
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditGroupErr
		logs.Error("Update Group: %s failed, code: %d, err: %s", this.FirstLevel, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Groups) Delete() Result {
	o := orm.NewOrm()
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
