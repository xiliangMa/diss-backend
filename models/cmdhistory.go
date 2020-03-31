package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type CmdHistory struct {
	Id          string `orm:"pk;" description:"(id)"`
	HostId      string `orm:"" description:"(主机id)"`
	ContainerId string `orm:"" description:"(容器id)"`
	User        string `orm:"" description:"(用户)"`
	Command     string `orm:"" description:"(命令)"`
	CreateTime  string `orm:"null;" description:"(更新时间)"`
	Type        int8   `orm:"default(0);" description:"(类型 0 host 1 container)"`
}

type CmdHistoryList struct {
	List []*CmdHistory
}

type CmdHistoryInterface interface {
	Add()
	MultiAdd()
	Delete()
	Edit()
	Get()
	List()
}

func (this *CmdHistory) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error
	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddCmdHistoryErr
		logs.Error("Add CmdHistory failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *CmdHistoryList) MultiAdd() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error
	_, err = o.InsertMulti(len(this.List), this.List)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddCmdHistoryErr
		logs.Error("Add CmdHistory failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *CmdHistory) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var imageList []*CmdHistory
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	cond = cond.And("type", this.Type)
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.ContainerId != "" {
		cond = cond.And("container_id", this.ContainerId)
	}

	if this.Command != "" {
		cond = cond.And("command__contains", this.Command)
	}

	_, err = o.QueryTable(utils.CmdHistory).SetCond(cond).Limit(limit, from).All(&imageList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetCmdHistoryErr
		logs.Error("Get CmdHistory List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.CmdHistory).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *CmdHistory) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	cond = cond.And("type", this.Type)
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}

	_, err := o.QueryTable(utils.CmdHistory).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteCmdHistoryErr
		logs.Error("Delete CmdHistory failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
