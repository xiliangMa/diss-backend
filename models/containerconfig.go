package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type ContainerConfig struct {
	Id             string      `orm:"pk;" description:"(id)"`
	Name           string      `orm:"" description:"(容器名)"`
	NameSpaceName  string      `orm:"" description:"(命名空间)"`
	PodId          string      `orm:"default(null)" description:"(pod id)"`
	PodName        string      `orm:"default(null)" description:"(pod 名)"`
	HostName       string      `orm:"" description:"(主机名)"`
	HostId         string      `orm:"" description:"(主机id)"`
	AccountName    string      `orm:"" description:"(租户)"`
	ClusterName    string      `orm:"" description:"(集群名)"`
	SyncCheckPoint int64       `orm:"default(0);" description:"(同步检查点)"`
	Status         string      `orm:"default(null);" description:"(状态)"`
	Command        string      `orm:"default(null);" description:"(命令)"`
	ImageName      string      `orm:"default(null);" description:"(镜像名)"`
	Age            string      `orm:"null;" description:"(运行时长)"`
	CreateTime     int64       `orm:"default(0);" description:"(创建时间);"`
	UpdateTime     int64       `orm:"default(0);" description:"(更新时间);"`
	TaskList       []*Task     `orm:"reverse(many);null" description:"(任务列表)"`
	Job            []*Job      `orm:"rel(m2m);null;" description:"(job)"`
	HostConfig     *HostConfig `orm:"rel(fk);null;" description:"(主机列表)"`
}

type ContainerConfigInterface interface {
	Add() Result
	Delete() Result
	Edit() Result
	Get() *ContainerConfig
	List(from, limit int, groupSearch bool) Result
	Count() int64
	EmptyDirtyDataForAgent() error
	EmptyDirtyDataForK8s() error
	GetContainerConfigList() []*ContainerConfig
}

func (this *ContainerConfig) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

	if this.Id != "" {
		if this.Get() == nil {
			_, err = o.Insert(this)
			if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
				ResultData.Message = err.Error()
				ResultData.Code = utils.AddContainerConfigErr
				logs.Error("Add ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
				return ResultData
			}
		} else {
			this.Update()
		}
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerConfig) Get() *ContainerConfig {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	containerConfig := new(ContainerConfig)
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	err = o.QueryTable(utils.ContainerConfig).SetCond(cond).RelatedSel().One(containerConfig)
	if err != nil {
		return nil
	}
	return containerConfig
}

func (this *ContainerConfig) List(from, limit int, groupSearch bool) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ContainerList []*ContainerConfig
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}
	if this.ImageName != "" {
		cond = cond.And("image_name__contains", this.ImageName)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name__contains", this.NameSpaceName)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name__contains", this.ClusterName)
	}
	if this.Status != "" && this.Status != Container_Status_All {
		switch this.Status {
		case Container_Status_Run:
			cond = cond.AndCond(cond.And("status__contains", "Up").Or("status", Container_Status_Running).Or("status", Container_Status_Terminated))
		case Container_Status_Pause:
			cond = cond.AndNotCond(cond.And("status__contains", "Up").Or("status", Container_Status_Running).Or("status", Container_Status_Terminated))
		}
	}

	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	// 分组条件 只能查询pod 为空的主机
	if groupSearch == true {
		cond = cond.And("pod_id", "")
	} else {
		if this.PodId != "" {
			cond = cond.And("pod_id", this.PodId)
		}
	}
	_, err = o.QueryTable(utils.ContainerConfig).SetCond(cond).Limit(limit, from).All(&ContainerList)
	total, _ := o.QueryTable(utils.ContainerConfig).SetCond(cond).Count()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerConfigErr
		logs.Error("Get ContainerConfig List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	for _, containerConfig := range ContainerList {
		o.LoadRelated(containerConfig, "TaskList", 1, 1, 0, "-update_time")
		if containerConfig.AccountName == "" {
			containerConfig.AccountName = this.AccountName
		}
	}
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = ContainerList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
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

func (this *ContainerConfig) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	cond := orm.NewCondition()
	if this.Status != "" {
		cond = cond.And("status", strings.ToLower(this.Status))
	}
	count, _ := o.QueryTable(utils.ContainerConfig).SetCond(cond).Count()
	return count
}

func (this *ContainerConfig) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.PodId != "" {
		cond = cond.And("pod_id", this.PodId)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if !o.QueryTable(utils.ContainerConfig).SetCond(cond).Exist() {
		ResultData.Message = "ContainerConfigNotFoundErr"
		ResultData.Code = utils.ContainerConfigNotFoundErr
		logs.Error("Delete ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	_, err := o.QueryTable(utils.ContainerConfig).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteContainerConfigErr
		logs.Error("Delete ContainerConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *ContainerConfig) EmptyDirtyDataForAgent() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.ContainerConfig+" where host_name = ? and sync_check_point != ? and pod_id = '' ", this.HostName, this.SyncCheckPoint).Exec()
	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.ContainerConfig, utils.EmptyDirtyDataContinerConfigErr, err.Error())
	}
	return err
}

func (this *ContainerConfig) EmptyDirtyDataForK8s() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.ContainerConfig+" where cluster_name = ? and sync_check_point != ? and pod_id = '' ", this.ClusterName, this.SyncCheckPoint).Exec()

	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.ContainerConfig, utils.EmptyDirtyDataContinerConfigErr, err.Error())
	}
	return err
}

func (this *ContainerConfig) GetContainerConfigList() []*ContainerConfig {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var err error
	var containerConfigList []*ContainerConfig

	cond := orm.NewCondition()
	if this.ImageName != "" {
		cond = cond.And("image_name", this.ImageName)
	}

	_, err = o.QueryTable(utils.ContainerConfig).SetCond(cond).All(&containerConfigList)
	if err != nil {
		logs.Error("Get ContainerConfig failed, code: %d, err: %s", utils.GetContainerConfigErr, err.Error())
	}
	return containerConfigList
}
