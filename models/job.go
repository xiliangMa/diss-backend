package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Job struct {
	Id             string          `orm:"pk;" description:"(job id)"`
	Account        string          `orm:"default(admin)" description:"(租户)"`
	Name           string          `orm:"" description:"(名称)"`
	Description    string          `orm:"" description:"(描述)"`
	Spec           string          `orm:"" description:"(定时器)"`
	Type           string          `orm:"" description:"(类型 重复执行 单次执行 )"`
	Status         string          `orm:"null;" description:"(状态: 执行中、启用、禁用)"`
	SystemTemplate *SystemTemplate `orm:"rel(fk);null;" description:"(系统模板)"`
	Task           []*Task         `orm:"reverse(many);null" description:"(任务列表)"`
	//HostConfig      []*HostConfig      `orm:"reverse(many);null" description:"(主机列表)"`
	//ContainerConfig []*ContainerConfig `orm:"reverse(many);null" description:"(容器列表)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime time.Time `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
}

type JobInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
}

func (this *Job) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var JobList []*Job
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Account != "" {
		cond = cond.And("account", this.Account)
	}
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	_, err = o.QueryTable(utils.Job).SetCond(cond).RelatedSel().Limit(limit, from).OrderBy("-create_time").All(&JobList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetJobErr
		logs.Error("Get Job List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Job).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = JobList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Job) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	uid, _ := uuid.NewV4()
	this.Id = uid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddJobErr
		logs.Error("Add Job failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Job) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.Job).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteImageInfoErr
		logs.Error("Delete Job failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
