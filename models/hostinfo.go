package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostInfo struct {
	Id                    string `orm:"pk;size(128)" description:"(主机id)"`
	HostName              string `orm:"size(64)" description:"(主机名称)"`
	InternalAddr          string `orm:"size(32);default(null);" description:"(主机ip 内)"`
	PublicAddr            string `orm:"size(32);default(null);" description:"(主机ip 外)"`
	CpuCore               int64  `orm:"" description:"(cpu)"`
	Mem                   string `orm:"size(32)" description:"(内存)"`
	Disk                  string `orm:"size(32)" description:"(磁盘)"`
	OS                    string `orm:"size(32)" description:"(系统)"`
	OSVer                 string `orm:"size(32)" description:"(系统版本)"`
	Kernel                string `orm:"size(32)" description:"(内核)"`
	Architecture          string `orm:"size(32)" description:"(架构)"`
	Mac                   string `orm:"" description:"(mac)"`
	DockerRuntime         string `orm:"size(128)" description:"(容器运行时)"`
	KubernetesVer         string `orm:"size(64)" description:"(kubernetes 版本)"`
	KubeletVer            string `orm:"size(64)" description:"(kubelet 版本)"`
	Kubeproxy             string `orm:"size(64)" description:"(kubeproxy 版本)"`
	DockerStatus          string `orm:"size(32);default(false);" description:"(容器状态)"`
	ImageCount            int    `orm:"default(0);" description:"(镜像数)"`
	ContainerCount        int    `orm:"default(0);" description:"(容器数)"`
	ContainerRunningCount int    `orm:"default(0);" description:"(容器Running数)"`
	ContainerPausedCount  int    `orm:"default(0);" description:"(容器Paused数)"`
	ContainerStoppedCount int    `orm:"default(0);" description:"(容器Stopped数)"`
	ClusterId             string `orm:"size(128);default(null);" description:"(集群id)"`
	ClusterName           string `orm:"size(128);default(null);" description:"(集群名)"`
	KMetaData             string `orm:"" description:"(源数据)"`
	KSpec                 string `orm:"" description:"(Spec数据)"`
	KStatus               string `orm:"" description:"(状态数据)"`
}

type HostInfoInterface interface {
	Add() error
	Update() Result
	List() Result
	Delete() Result
	UpdateDynamic() Result
}

func (this *HostInfo) Add() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var err error
	var hostInfoList []*HostInfo
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err = o.QueryTable(utils.HostInfo).SetCond(cond).All(&hostInfoList)
	if err != nil {
		return err
	}

	if len(hostInfoList) != 0 {
		// agent 或者 k8s 数据更新 （因为有diss-backend的关系数据，防止覆盖diss-backend的数据，需要替换更新）
		updateHostInfo := hostInfoList[0]
		updateHostInfo.ClusterName = this.ClusterName
		updateHostInfo.HostName = this.HostName
		updateHostInfo.OS = this.OS
		updateHostInfo.ClusterId = this.ClusterId
		updateHostInfo.InternalAddr = this.InternalAddr
		if this.PublicAddr != "" {
			updateHostInfo.PublicAddr = this.PublicAddr
		}
		if this.InternalAddr != "" {
			updateHostInfo.InternalAddr = this.InternalAddr
		}
		if this.DockerRuntime != "" {
			updateHostInfo.DockerRuntime = this.DockerRuntime
		}
		if this.DockerStatus != "" {
			updateHostInfo.DockerStatus = this.DockerStatus
		}
		if this.ImageCount != 0 {
			updateHostInfo.ImageCount = this.ImageCount
		}
		if this.ContainerCount != 0 {
			updateHostInfo.ContainerCount = this.ContainerCount
		}
		if this.ContainerRunningCount != 0 {
			updateHostInfo.ContainerRunningCount = this.ContainerRunningCount
		}
		if this.ContainerPausedCount != 0 {
			updateHostInfo.ContainerPausedCount = this.ContainerPausedCount
		}
		if this.ContainerStoppedCount != 0 {
			updateHostInfo.ContainerStoppedCount = this.ContainerStoppedCount
		}
		result := updateHostInfo.Update()
		if result.Code != http.StatusOK {
			return errors.New(result.Message)
		}
	} else {
		// 插入数据
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			logs.Error("DB Metrics data --- Add %s failed, err: %s", Resource_HostInfo, err.Error())
			return err
		}
	}

	return nil
}

func (this *HostInfo) List() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var HostInfoList []*HostInfo = nil
	var ResultData Result
	cond := orm.NewCondition()
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err := o.QueryTable(utils.HostInfo).SetCond(cond).All(&HostInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostInfoErr
		logs.Error("GetHostInfoList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.HostInfo).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["items"] = HostInfoList
	data["total"] = total

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostInfo) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostInfoErr
		logs.Error("Update HostInfo: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostInfo) UpdateDynamic() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	hostInfo := new(HostInfo)
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if err := o.QueryTable(utils.HostInfo).SetCond(cond).One(hostInfo); err != nil {
		ResultData.Code = utils.HostConfigNotFoundErr
		ResultData.Message = err.Error()
		logs.Warn("Not Get HostInfo: %s, code: %d, message: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	hostInfo.PublicAddr = this.PublicAddr
	hostInfo.ImageCount = this.ImageCount
	hostInfo.ContainerCount = this.ContainerCount
	hostInfo.ContainerRunningCount = this.ContainerRunningCount
	hostInfo.ContainerPausedCount = this.ContainerPausedCount
	hostInfo.ContainerStoppedCount = this.ContainerStoppedCount
	_, err := o.Update(hostInfo)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostInfoDynamicErr
		logs.Error("Update HostInfo Dynamic failed, HostName: %s, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostInfo) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	_, err := o.QueryTable(utils.HostInfo).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostinfoErr
		logs.Error("Delete HostInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
