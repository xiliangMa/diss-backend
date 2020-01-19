package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type HostConfig struct {
	Id         string `orm:"pk;description(主机id)"`
	HostName   string `orm:"description(主机名)"`
	OS         string `orm:"description(系统)"`
	PG         string `orm:"description(安全策略组)"`
	Status     uint8  `orm:"description(主机状态)"`
	Diss       uint8  `orm:"description(安全容器)"`
	DissStatus uint8  `orm:"description(安全状态)"`
	TenantId   string `orm:"description(租户id)"`
	Group      string `orm:"description(分组)"`
	Type       uint8  `orm:"description(系统类型，服务器或者虚拟机)"`
}

type HostInfo struct {
	Id            string `orm:"pk;description(主机id)"`
	HostName      string `orm:"description(主机名称)"`
	InternalAddr  string `orm:"description(主机ip 内)"`
	PublicAddr    string `orm:"description(主机ip 外)"`
	CpuCore       uint8  `orm:"description(cpu)"`
	Mem           uint64 `orm:"description(内存)"`
	Disk          string `orm:"description(磁盘)"`
	OS            string `orm:"description(系统)"`
	OSVer         string `orm:"description(系统版本)"`
	Kernel        string `orm:"description(内核)"`
	Architecture  string `orm:"description(架构)"`
	Mac           string `orm:"description(mac)"`
	DockerRuntime string `orm:"description(容器运行时)"`
	KubernetVer   string `orm:"description(kubernetes 版本)"`
	KubeletVer    string `orm:"description(kubelet 版本)"`
	Kubeproxy     string `orm:"description(kubeproxy 版本)"`
}

func init() {
	orm.RegisterModel(new(HostConfig), new(HostInfo))
}

func (hostConfig *HostConfig) Inner_AddHostConfig() error {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(hostConfig)
	if err != nil {
		logs.Error("DB Metrics data --- Add %s failed, err: %s", Tag_HostConfig, err.Error())
		return err
	}
	return nil
}

func (hostInfo *HostInfo) Inner_AddHostInfo() error {
	o := orm.NewOrm()
	o.Using("default")
	_, err := o.Insert(hostInfo)
	if err != nil {
		logs.Error("DB Metrics data --- Add %s failed, err: %s", Tag_HostInfo, err.Error())
		return err
	}
	return nil
}
