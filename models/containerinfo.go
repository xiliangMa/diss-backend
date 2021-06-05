package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type ContainerInfo struct {
	Id             string `orm:"pk;" description:"(id)"`
	Name           string `orm:"" description:"(名称)"`
	NameSpaceName  string `orm:"" description:"(命名空间)"`
	PodId          string `orm:"" description:"(Pod Id)"`
	PodName        string `orm:"" description:"(Pod 名称)"`
	ImageId        string `orm:"" description:"(imageId)"`
	ImageName      string `orm:"" description:"(image名称)"`
	HostId         string `orm:"" description:"(主机id)"`
	HostName       string `orm:"" description:"(主机名)"`
	ClusterName    string `orm:"" description:"(集群名)"`
	SyncCheckPoint int64  `orm:"default(0);" description:"(同步检查点)"`
	Command        string `orm:"default(null);" description:"(命令)"`
	StartedAt      int64  `orm:"default(0);" description:"(启动时间)"`
	CreatedAt      int64  `orm:"default(0);" description:"(创建时间)"`
	Status         string `orm:"default(null);" description:"(状态)"`
	Ports          string `orm:"default(null);" description:"(端口)"`
	Ip             string `orm:"default(null);" description:"(ip)"`
	Labels         string `orm:"default(null);" description:"(标签)"`
	Volumes        string `orm:"default(null);" description:"(Volumes)"`
	Mounts         string `orm:"default(null);" description:"(Mounts)"`
}

type ContainerInfoInterface interface {
	Add() Result
	Delete() Result
	Edit() Result
	Get() Result
	List() Result
	EmptyDirtyDataForAgent() error
	EmptyDirtyDataForK8s() error
}

func (this *ContainerInfo) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddContainerInfoErr
		logs.Error("Add ContainerInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerInfo) List() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ContainerList []*ContainerInfo = nil
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	}
	if this.PodId != "" {
		cond = cond.And("pod_id", this.PodId)
	}
	_, err = o.QueryTable(utils.ContainerInfo).SetCond(cond).All(&ContainerList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerInfoErr
		logs.Error("Get ContainerInfo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ContainerInfo).SetCond(cond).Count()
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
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditContainerInfoErr
		logs.Error("Update ContainerInfo: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerInfo) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.PodId != "" {
		cond = cond.And("pod_id", this.PodId)
	}

	_, err := o.QueryTable(utils.ContainerInfo).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteContainerInfoErr
		logs.Error("Delete ContainerInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *ContainerInfo) EmptyDirtyDataForAgent() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.ContainerInfo+" where host_name = ? and sync_check_point != ? and pod_id = '' ", this.HostName, this.SyncCheckPoint).Exec()

	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.ContainerInfo, utils.EmptyDirtyDataContinerConfigErr, err.Error())
	}
	return err
}

func (this *ContainerInfo) EmptyDirtyDataForK8s() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.ContainerInfo+" where cluster_name = ? and sync_check_point != ? and pod_id = '' ", this.ClusterName, this.SyncCheckPoint).Exec()

	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.ContainerInfo, utils.EmptyDirtyDataContinerConfigErr, err.Error())
	}
	return err
}
