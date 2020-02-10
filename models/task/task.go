package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Task struct {
	Id         int       `orm:"auto;descripion(任务id)"`
	Type       string    `orm:"description(类型)"`
	Status     string    `orm:"null;description(状态)"`
	CreateTime time.Time `orm:"description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"null;description(更新时间);auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Task))
}

type TaskInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *Task) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddTaskErr
		logs.Error("Add Task failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
