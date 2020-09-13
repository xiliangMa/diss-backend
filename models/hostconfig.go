package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type HostConfig struct {
	Id                string    `orm:"pk;size(128)" description:"(主机id)"`
	HostName          string    `orm:"size(64)" description:"(主机名)"`
	OS                string    `orm:"size(32)" description:"(系统)"`
	PG                string    `orm:"size(32);default(sys-default)" description:"(安全策略组)"`
	InternalAddr      string    `orm:"size(32);default(null);" description:"(主机ip 内)"`
	PublicAddr        string    `orm:"size(32);default(null);" description:"(主机ip 外)"`
	Status            string    `orm:"size(32);default(Normal)" description:"(主机状态 正常 Normal 异常 Abnormal)"`
	Diss              string    `orm:"size(32);default(Installed)" description:"(安全容器 Installed NotInstalled)"`
	DissStatus        string    `orm:"size(32);default(Safe)" description:"(安全状态 Safe Unsafe)"`
	AccountName       string    `orm:"size(32);default(admin)" description:"(租户)"`
	GroupId           string    `orm:"-" description:"(查询参数：分组Id， 仅仅是查询使用, 返回数据看 Group)"`
	Group             *Groups   `orm:"rel(fk);null;on_delete(set_null)" description:"(分组)"`
	Type              string    `orm:"size(32);default(Server);" description:"(类型 服务器: Server 虚拟机: Vm)"`
	IsInK8s           bool      `orm:"default(false);" description:"(是否在k8s集群)"`
	ClusterId         string    `orm:"size(128);default(null);" description:"(集群id)"`
	ClusterName       string    `orm:"size(128);default(null);" description:"(集群名)"`
	Label             string    `orm:"size(32);default(null);" description:"(标签)"`
	Job               []*Job    `orm:"rel(m2m);null;" description:"(job)"`
	IsEnableHeartBeat bool      `orm:"default(false);" description:"(是否开启心跳上报)"`
	HeartBeat         time.Time `orm:"null;type(datetime)" description:"(心跳)"`
	KMetaData         string    `orm:"" description:"(源数据)"`
	KSpec             string    `orm:"" description:"(Spec数据)"`
	KStatus           string    `orm:"" description:"(状态数据)"`
	KubernetesVer     string    `orm:"size(64)" description:"(kubernetes 版本)"`
	NodeRole          string    `orm:"size(64)" description:"(集群主机角色)"`
	DockerCISCount    string    `orm:"null;" description:"(docker基线结果个数)"`
	KubeCISCount      string    `orm:"null;" description:"(k8s基线结果个数)"`
}

type HostConfigInterface interface {
	Add() error
	List(from, limit int) Result
	Update() Result
	Delete() Result
	UpdateDynamic() Result
	Count() int64
	GetBnechMarkProportion() (int64, int64)
	Get() *HostConfig
}

func (this *HostConfig) Add() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
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
		updateHostConfig.ClusterId = this.ClusterId
		updateHostConfig.InternalAddr = this.InternalAddr
		if this.PublicAddr != "" {
			updateHostConfig.PublicAddr = this.PublicAddr
		}
		updateHostConfig.OS = this.OS
		updateHostConfig.AccountName = Account_Admin
		result := updateHostConfig.Update()
		if result.Code != http.StatusOK {
			return errors.New(result.Message)
		}
	} else {
		// 插入数据
		this.AccountName = Account_Admin
		//添加默认数据
		this.PG = "Sys-Default"
		_, err = o.Insert(this)
		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			logs.Error("DB Metrics data --- Add %s failed, err: %s", Resource_HostConfig, err.Error())
			return err
		}
	}

	return nil
}

func (this *HostConfig) Get() *HostConfig {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	HostConfig := new(HostConfig)
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	err = o.QueryTable(utils.HostConfig).SetCond(cond).RelatedSel().One(HostConfig)
	if err != nil {
		logs.Error("GetHostConfig failed, code: %d, err: %s", err.Error(), utils.GetHostConfigErr)
		return nil
	}
	return HostConfig
}

func (this *HostConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var HostConfigList []*HostConfig = nil
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Diss != "" && this.Diss != All {
		cond = cond.And("diss", this.Diss)
	}
	if this.DissStatus != "" && this.DissStatus != All {
		cond = cond.And("diss_status", this.DissStatus)
	}
	if this.Label != "" {
		cond = cond.And("label__contains", this.Label)
	}
	if this.GroupId != "" {
		cond = cond.And("Group", this.GroupId)
	}
	if this.HostName != "" {
		cond = cond.And("host_name__contains", this.HostName)
	}
	if this.AccountName != "" && this.AccountName != Account_Admin {
		cond = cond.And("account_name", this.AccountName)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).Limit(limit, from).OrderBy("-host_name").RelatedSel().All(&HostConfigList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostConfigErr
		logs.Error("GetHostConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.HostConfig).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["items"] = HostConfigList
	data["total"] = total
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *HostConfig) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostConfigErr
		logs.Error("Update HostConfig: %s failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostConfig) UpdateDynamic() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	hostConfig := new(HostConfig)
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if err := o.QueryTable(utils.HostConfig).SetCond(cond).One(hostConfig); err != nil {
		ResultData.Code = utils.HostConfigNotFoundErr
		ResultData.Message = err.Error()
		logs.Warn("Not Get HostConfig: %s, code: %d, message: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	hostConfig.PublicAddr = this.PublicAddr
	_, err := o.Update(hostConfig)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostConfigDynamicErr
		logs.Error("Update HostInfo Dynamic failed, HostName: %s, failed, code: %d, err: %s", this.HostName, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostConfig) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	count, _ := o.QueryTable(utils.HostConfig).SetCond(cond).Count()
	return count
}

// docker基线 / k8s 基线
func (this *HostConfig) GetBnechMarkProportion() (int64, int64) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	dockerBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Count()
	k8sBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Filter("is_in_k8s", false).Count()
	return dockerBenchMarkCount, k8sBenchMarkCount
}

// Online / Offline
func (this *HostConfig) GetOnlineProportion() (int64, int64) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	all, _ := o.QueryTable(utils.HostConfig).Count()
	onlineCount, _ := o.QueryTable(utils.HostConfig).Filter("status", Host_Status_Normal).Count()
	return onlineCount, all - onlineCount
}

func (this *HostConfig) Delete() Result {
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
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err := o.QueryTable(utils.HostConfig).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostConfigErr
		logs.Error("Delete HostConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
