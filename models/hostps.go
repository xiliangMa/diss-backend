package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"github.com/xiliangMa/restapi/models"
	"net/http"
	"time"
)

type HostPs struct {
	Id      string `orm:"pk;description(id)"`
	HostId  string `orm:"pk;description(主机id)"`
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
}

func (this *HostPs) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var HostPsList []*HostPs = nil
	var total = 0
	var ResultData Result

	_, err := o.QueryTable(utils.HostPs).Limit(limit, from).All(&HostPsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPsErr
		logs.Error("GetHostPs List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if HostPsList != nil {
		total = len(HostPsList)
	}
	data := make(map[string]interface{})
	data["items"] = HostPsList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostPs) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

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
