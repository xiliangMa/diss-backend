package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostPs struct {
	Id      string        `orm:"pk;" description:"(id)"`
	HostId  string        `orm:"" description:"(主机id)"`
	PID     string        `orm:"" description:"(PID)"`
	User    string        `orm:"" description:"(用户)"`
	CPU     string        `orm:"" description:"(CPU)"`
	Mem     string        `orm:"" description:"(内存)"`
	Time    string        `orm:"" description:"(时间)"`
	Start   string        `orm:"" description:"(运行时长 非mac)"`
	Started string        `orm:"" description:"(运行时长 mac)"`
	Command orm.TextField `orm:"" description:"(Command)"`
}

func init() {
	orm.RegisterModel(new(HostPs))
}

type HostPsInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *HostPs) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var hostPsList []*HostPs = nil
	var ResultData Result
	cond := orm.NewCondition()
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.Command != "" {
		cond = cond.And("command__contains", this.Command)
	}

	_, err := o.QueryTable(utils.HostPs).SetCond(cond).Limit(limit, from).All(&hostPsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPsErr
		logs.Error("GetHostPs List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.HostPs).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["items"] = hostPsList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostPs) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddHostPsErr
		logs.Error("Add HostPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostPs) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostPsErr
		logs.Error("Update HostPs: %s failed, code: %d, err: %s", this.Id, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostPs) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	// 根据agent同步时 依据 host_id 删除该主机上所有的容器历史记录
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	_, err := o.QueryTable(utils.HostPs).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostPsErr
		logs.Error("Delete HostPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
