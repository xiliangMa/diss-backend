package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type HostPs struct {
	Id      string `orm:"pk;description(id)"`
	HostId  string `orm:"description(主机id)"`
	PID     string `orm:"description(PID)"`
	User    string `orm:"description(用户)"`
	CPU     string `orm:"description(CPU)"`
	Mem     string `orm:"description(内存)"`
	Time    string `orm:"description(时间)"`
	Start   string `orm:"description(运行时长 非mac)"`
	Started string `orm:"description(运行时长 mac)"`
	Command string `orm:"description(Command)"`
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
	ListById()
}

func (this *HostPs) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var hostPsList []*HostPs = nil
	var total = 0
	var ResultData Result

	_, err := o.QueryTable(utils.HostPs).Limit(limit, from).All(&hostPsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPsErr
		logs.Error("GetHostPs List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if hostPsList != nil {
		total = len(hostPsList)
	}
	data := make(map[string]interface{})
	data["items"] = hostPsList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostPs) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil {
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
	o.Using("default")
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
	o.Using("default")
	var ResultData Result
	_, err := o.Delete(&HostPs{Id: this.Id, HostId: this.HostId})
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostPsErr
		logs.Error("Delete HostPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *HostPs) ListById() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var hostPsList []*HostPs = nil
	var total = 0
	var ResultData Result

	_, err := o.QueryTable(utils.HostPs).Filter("host_id", this.HostId).All(&hostPsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPsErr
		logs.Error("GetHostPs List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if hostPsList != nil {
		total = len(hostPsList)
	}
	data := make(map[string]interface{})
	data["items"] = hostPsList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}
