package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ContainerConfig struct {
	Id            string `orm:"pk;" description:"(id)"`
	Name          string `orm:"" description:"(容器名)"`
	NameSpaceName string `orm:"" description:"(命名空间)"`
	PodId         string `orm:"" description:"(pod id)"`
	PodName       string `orm:"" description:"(pod 名)"`
	HostName      string `orm:"" description:"(主机名)"`
	Status        string `orm:"" default(null);size(1000);description:"(状态)"`
	Command       string `orm:"" default(null);size(1000);description:"(命令)"`
	ImageName     string `orm:"" default(null);description:"(镜像名)"`
	Age           string `orm:"null;" description:"(运行时长)"`
	CreateTime    string `orm:"null;" description:"(创建时间);"`
	UpdateTime    string `orm:"null;" description:"(更新时间);"`
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
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error
	var containerConfigList []*ContainerConfig

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.ContainerConfig).SetCond(cond).All(&containerConfigList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(containerConfigList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		return this.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddContainerConfigErr
			logs.Error("Add ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using(utils.DS_Default)
	var ContainerList []*ContainerConfig = nil
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	if this.ImageName != "" {
		cond = cond.And("image_name", this.ImageName)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	}
	if this.PodId != "" {
		cond = cond.And("pod_id", this.PodId)
	}
	if this.Status != "" && this.Status != Container_Status_All {
		switch this.Status {
		case Container_Status_Run:
			cond = cond.AndCond(cond.And("status__contains", "Up").Or("status", Pod_Container_Statue_Running).Or("status", Pod_Container_Statue_Terminated))
		case Container_Status_Pause:
			cond = cond.AndNotCond(cond.And("status__contains", "Up").Or("status", Pod_Container_Statue_Running).Or("status", Pod_Container_Statue_Terminated))

		}
	}
	_, err = o.QueryTable(utils.ContainerConfig).SetCond(cond).Limit(limit, from).All(&ContainerList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ContainerConfig).SetCond(cond).Count()
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

func (this *ContainerConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditContainerConfigErr
		logs.Error("Update ContainerConfig: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
