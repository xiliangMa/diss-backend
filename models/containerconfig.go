package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ContainerConfig struct {
	Id            string    `orm:"pk;description(id)"`
	Name          string    `orm:"unique;description(容器名)"`
	NameSpaceId   string    `orm:"description(命名空间id)"`
	NameSpaceName string    `orm:"description(命名空间)"`
	PodId         string    `orm:"description(pod id)"`
	PodName       string    `orm:"description(pod 名)"`
	Status        string    `orm:"default(null);description(状态)"`
	Command       string    `orm:"default(null);description(命令)"`
	ImageName     string    `orm:"default(null);description(镜像名)"`
	CreateTime    time.Time `orm:"null;description(创建时间);type(datetime)"`
	UpdateTime    time.Time `orm:"null;description(更新时间);type(datetime)"`
}

func init() {
	orm.RegisterModel(new(ContainerConfig))
}

type ContainerConfigInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ContainerConfig) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddContainerConfigErr
		logs.Error("Add ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerConfig) List() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ContainerList []*ContainerConfig
	var ResultData Result

	_, err := o.QueryTable(utils.ContainerConfig).All(&ContainerList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ContainerConfig).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ContainerList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
