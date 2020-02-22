package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostConfig struct {
	Id         string `orm:"pk;description(主机id)"`
	HostName   string `orm:"description(主机名)"`
	OS         string `orm:"description(系统)"`
	PG         string `orm:"description(安全策略组)"`
	Status     int8   `orm:"description(主机状态)"`
	Diss       int8   `orm:"description(安全容器)"`
	DissStatus int8   `orm:"description(安全状态)"`
	TenantId   string `orm:"description(租户id)"`
	Group      string `orm:"description(分组)"`
	Type       int8   `orm:"default(0);description(类型 服务器: 0 虚拟机: 1)"`
	IsInK8s    bool   `orm:"default(false);description(是否在k8s集群)"`
	ClusterId  string `orm:"default(null);description(集群id)"`
}

type HostInfo struct {
	Id            string `orm:"pk;description(主机id)"`
	HostName      string `orm:"description(主机名称)"`
	InternalAddr  string `orm:"default(null);description(主机ip 内)"`
	PublicAddr    string `orm:"default(null);description(主机ip 外)"`
	CpuCore       int64  `orm:"description(cpu)"`
	Mem           int64  `orm:"description(内存)"`
	Disk          string `orm:"description(磁盘)"`
	OS            string `orm:"description(系统)"`
	OSVer         string `orm:"description(系统版本)"`
	Kernel        string `orm:"description(内核)"`
	Architecture  string `orm:"description(架构)"`
	Mac           string `orm:"description(mac)"`
	DockerRuntime string `orm:"description(容器运行时)"`
	KubernetesVer string `orm:"description(kubernetes 版本)"`
	KubeletVer    string `orm:"description(kubelet 版本)"`
	Kubeproxy     string `orm:"description(kubeproxy 版本)"`
	DockerStatus  string `orm:"default(false);description(容器状态)"`
}

func init() {
	orm.RegisterModel(new(HostConfig), new(HostInfo))
}

func (this *HostConfig) Inner_AddHostConfig() error {
	o := orm.NewOrm()
	o.Using("default")
	var err error
	var hostConfigList []*HostConfig
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).All(&hostConfigList)
	if err != nil {
		return err
	}
	if len(hostConfigList) != 0 {
		updateHostConfig := hostConfigList[0]
		// agent 或者 k8s 数据更新 （因为有diss-backend的关系数据，防止覆盖diss-backend的数据，需要替换更新）
		updateHostConfig.HostName = this.HostName
		updateHostConfig.IsInK8s = this.IsInK8s
		updateHostConfig.OS = this.OS
		updateHostConfig.Status = this.Status

		resilt := updateHostConfig.Update()
		if resilt.Code != http.StatusOK {
			return errors.New(resilt.Message)
		}
	} else {
		// 插入数据
		_, err = o.Insert(this)
		if err != nil {
			logs.Error("DB Metrics data --- Add %s failed, err: %s", Tag_HostConfig, err.Error())
			return err
		}
	}

	return nil
}

func (this *HostInfo) Inner_AddHostInfo() error {
	o := orm.NewOrm()
	o.Using("default")
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
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		resilt := this.Update()
		if resilt.Code != http.StatusOK {
			return errors.New(resilt.Message)
		}
	} else {
		// 插入数据
		_, err = o.Insert(this)
		if err != nil {
			logs.Error("DB Metrics data --- Add %s failed, err: %s", Tag_HostInfo, err.Error())
			return err
		}
	}

	return nil
}
