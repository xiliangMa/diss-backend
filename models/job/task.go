package job

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	msecuritypolicy "github.com/xiliangMa/diss-backend/models/securitypolicy"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Task struct {
	Id             string                          `orm:"pk;" description:"(任务id)"`
	Name           string                          `orm:"" description:"(名称)"`
	Description    string                          `orm:"" description:"(描述)"`
	Spec           string                          `orm:"" description:"(定时器)"`
	Type           string                          `orm:"" description:"(类型 重复执行 单词执行 )"`
	Status         string                          `orm:"null;" description:"(状态: 未开始、执行中、完成、暂停)"`
	SystemTemplate *msecuritypolicy.SystemTemplate `orm:"rel(fk);null;;" description:"(模板类型)"`
	CreateTime     time.Time                       `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime     time.Time                       `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
}

type TaskInterface interface {
	Add() models.Result
	List() models.Result
	Delete() models.Result
	Update() models.Result
}

func (this *Task) Add() models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData models.Result

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

func (this *Task) List(from, limit int) models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var TaskList []*Task
	var ResultData models.Result
	var err error
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	_, err = o.QueryTable(utils.Task).SetCond(cond).RelatedSel().Limit(limit, from).All(&TaskList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetTaskErr
		logs.Error("Get Task List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.SYSTemplate).Count()
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

func (this *Task) Delete() models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData models.Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.Task).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteTaskErr
		logs.Error("Delete Task failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *Task) Update() models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData models.Result

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
