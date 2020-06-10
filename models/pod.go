package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Pod struct {
	Id             string `orm:"pk;" description:"(pod id)"`
	Name           string `orm:"unique;" description:"(集群名)"`
	AccountName    string `orm:"" description:"(租户)"`
	PodIp          string `orm:"default(null);" description:"(pod ip)"`
	Status         string `orm:"" description:"(pod状态)"`
	GroupId        string `orm:"default(null);" description:"(租户id)"`
	GroupName      string `orm:"default(null);" description:"(租户名)"`
	SyncCheckPoint int64  `orm:"default(0);" description:"(同步检查点)"`
	HostIp         string `orm:"default(null);" description:"(主机ip， 默认内网ip)"`
	HostName       string `orm:"default(null);" description:"(host name)"`
	NameSpaceName  string `orm:"default(null);" description:"(命名空间)"`
	ClusterName    string `orm:"" description:"(集群名)"`
}

type PodInterface interface {
	Add() Result
	Delete() Result
	Update() Result
	Get() Result
	List(from, limit int) Result
	EmptyDirtyData() error
}

func (this *Pod) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var podList []*Pod
	var err error
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)
	if this.Name != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.Pod).SetCond(cond).All(&podList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetPodErr
		logs.Error("Get Pod failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(podList) != 0 {
		updatePod := podList[0]
		// agent 或者 k8s 数据更新 （因为有diss-backend的关系数据，防止覆盖diss-backend的数据，需要替换更新）
		updatePod.HostName = this.HostName
		updatePod.Name = this.Name
		updatePod.PodIp = this.PodIp
		updatePod.Status = this.Status
		updatePod.HostIp = this.HostIp
		updatePod.HostName = this.HostName
		updatePod.NameSpaceName = this.NameSpaceName
		return updatePod.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddPodErr
			logs.Error("Add Pod failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Pod) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var PodList []*Pod = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Name != "" {
		cond = cond.And("name__icontains", this.Name)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	if this.NameSpaceName != "" {
		cond = cond.And("name_space_name", this.NameSpaceName)
	}
	_, err = o.QueryTable(utils.Pod).SetCond(cond).Limit(limit, from).All(&PodList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetPodErr
		logs.Error("Get Pod List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.Pod).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = PodList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Pod) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditPodErr
		logs.Error("Update Pod: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Pod) EmptyDirtyData() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	_, err := o.Raw("delete from "+utils.Pod+" where cluster_name = ? and sync_check_point != ? ", this.ClusterName, this.SyncCheckPoint).Exec()
	if err != nil {
		logs.Error("Empty Dirty Data failed,  model: %s, code: %d, err: %s", utils.Pod, utils.EmptyDirtyDataPodErr, err.Error())
	}
	return err
}

func (this *Pod) Delete() Result {
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
	_, err := o.QueryTable(utils.Pod).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteNameSpaceErr
		logs.Error("Delete NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
