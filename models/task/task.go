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
	Id              string    `orm:"pk;descripion(任务id)"`
	Name            string    `orm:"descripion(名称)"`
	Spec            string    `orm:"descripion(定时器)"`
	BmtName         string    `orm:"descripion(入侵检测模版)"`
	BmtComman       string    `orm:"descripion(入侵检测命令)"`
	Type            string    `orm:"description(类型 重复执行 单词执行 )"`
	SecurityGroupId string    `orm:"description(安全策略组)"`
	Status          string    `orm:"null;description(状态  未开始、 执行中、完成、 暂停)"`
	CreateTime      time.Time `orm:"description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime      time.Time `orm:"null;description(更新时间);auto_now;type(datetime)"`
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
