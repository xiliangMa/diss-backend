package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ContainerInfo struct {
	Id            string `orm:"pk;description(id)"`
	Name          string `orm:"description(名称)"`
	NameSpaceName string `orm:"description(命名空间)"`
	PodId         string `orm:"description(Pod Id)"`
	PodName       string `orm:"description(Pod 名称)"`
	ImageId       string `orm:"description(imageId)"`
	ImageName     string `orm:"description(image名称)"`
	HostId        string `orm:"description(主机id)"`
	HostName      string `orm:"description(主机名)"`
	Command       string `orm:"default(null);description(命令)"`
	StartedAt      string `orm:"default(null);description(启动时间)"`
	CreatedAt     string `orm:"default(null);description(创建时间)"`
	Status        string `orm:"default(null);description(状态)"`
	Ports         string `orm:"default(null);description(端口)"`
	Ip            string `orm:"default(null);description(ip)"`
	Labels        string `orm:"default(null);description(标签)"`
	Volumes       string `orm:"default(null);description(Volumes)"`
	Mounts        string `orm:"default(null);description(Mounts)"`
}

func init() {
	orm.RegisterModel(new(ContainerInfo))
}

type ContainerInfoInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ContainerInfo) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	var err error
	var ContainerInfogList []*ContainerInfo

	_, err = o.QueryTable(utils.ContainerInfo).Filter("id", this.Id).All(&ContainerInfogList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerInfoErr
		logs.Error("Get ContainerInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(ContainerInfogList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		return this.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddContainerInfoErr
			logs.Error("Add ContainerInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerInfo) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ContainerList []*ContainerInfo = nil
	var ResultData Result
	var total = 0
	var err error

	cond := orm.NewCondition()
	cond = cond.And("id__icontains", this.Id)
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	} else if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	} else if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	} else if this.PodId != "" {
		cond = cond.And("pod_id", this.PodId)
	}
	_, err = o.QueryTable(utils.ContainerInfo).SetCond(cond).All(&ContainerList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerInfoErr
		logs.Error("Get ContainerInfo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if ContainerList != nil {
		total = len(ContainerList)
	}
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

func (this *ContainerInfo) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditContainerInfoErr
		logs.Error("Update ContainerInfo: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
