package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Task struct {
	Id             string           `orm:"pk;" description:"(任务id)"`
	Account        string           `orm:"default(admin)" description:"(租户)"`
	Name           string           `orm:"" description:"(名称)"`
	Description    string           `orm:"" description:"(描述)"`
	Spec           string           `orm:"" description:"(定时器)"`
	Type           string           `orm:"" description:"(类型 重复执行 单次执行 )"`
	Status         string           `orm:"null;" description:"(状态: 未开始、执行中、完成、暂停)"`
	Batch          int64            `orm:"default(0);" description:"(任务批次)"`
	SystemTemplate *SystemTemplate  `orm:"rel(fk);null;" description:"(系统模板)"`
	Host           *HostConfig      `orm:"rel(fk);null;on_delete(do_nothing)"description:"(主机)"`
	Container      *ContainerConfig `orm:"rel(fk);null;on_delete(do_nothing)" description:"(容器)"`
	CreateTime     time.Time        `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime     time.Time        `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
	Job            *Job             `orm:"rel(fk);null;" description:"(job)"`
	ClusterId      string           `orm:"size(256)" description:"(集群Id)"`
	IsOne          bool             `orm:"-" description:"(是否取单条记录)"`
	RunCount       int64            `orm:"" description:"(运行次数)"`
	Action         string           `orm:"size(256)" description:"(操作类型标记)"`
}

type TaskLogInterface interface {
	Add() Result
	List(from, limit int) Result
}

type TaskInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
	Update() Result
	GetCurrentBatchTaskList() (error, []*Task)
	GetUnFinishedTaskList() Result
}

func (this *Task) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddTaskErr
		logs.Error("Add Task failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Task) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var TaskList []*Task
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.Batch != 0 {
		cond = cond.And("batch", this.Batch)
	}
	if this.Host != nil && this.Host.Id != "" {
		cond = cond.And("host_id", this.Host.Id)
	}
	if this.Container != nil && this.Container.Id != "" {
		cond = cond.And("container_id", this.Container.Id)
	}
	if this.Account != "" {
		cond = cond.And("account", this.Account)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	if this.Status != "" && this.Status != All {
		cond = cond.And("status", this.Status)
	}
	if this.Job != nil && this.Job.Id != "" {
		cond = cond.And("job__id", this.Job.Id)
		if this.IsOne {
			limit = 1
		}
	}
	_, err = o.QueryTable(utils.Task).SetCond(cond).RelatedSel().Limit(limit, from).OrderBy("-update_time").All(&TaskList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetTaskErr
		logs.Error("Get Task List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Task).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = TaskList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Task) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()
	var err error

	if this.Id != "" {
		cond = cond.And("id", this.Id)
		_, err = o.QueryTable(utils.Task).SetCond(cond).Delete()
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteTaskErr
		logs.Error("Delete Task failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *Task) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditTaskErr
		logs.Error("Update Task: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Task) GetCurrentBatchTaskList() (error, []*Task) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var TaskList []*Task
	var err error
	cond := orm.NewCondition()
	cond = cond.And("batch", this.Batch)
	_, err = o.QueryTable(utils.Task).SetCond(cond).RelatedSel().All(&TaskList)
	if err != nil {
		logs.Error("Get Task List failed, code: %d, err: %s", utils.GetTaskErr, err.Error())
		return err, nil
	}
	return nil, TaskList
}

func (this *Task) GetUnFinishedTaskList() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var TaskList []*Task
	var ResultData Result
	var err error
	var total int64
	cond := orm.NewCondition()

	if this.Host != nil && this.Host.Id != "" {
		cond = cond.And("host_id", this.Host.Id)
	}
	cond = cond.AndCond(cond.And("status", Task_Status_Pending).Or("status", Task_Status_Running).Or("status", Task_Status_Deliver_Failed))

	total, err = o.QueryTable(utils.Task).SetCond(cond).RelatedSel().All(&TaskList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetTaskErr
		logs.Error("Get Task List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = TaskList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
