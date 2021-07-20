package models

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type HostConfig struct {
	Id                  string  `orm:"pk;size(128)" description:"(主机id)"`
	HostName            string  `orm:"size(64)" description:"(主机名)"`
	OS                  string  `orm:"size(32)" description:"(系统)"`
	PG                  string  `orm:"size(32);default(sys-default)" description:"(安全策略组)"`
	InternalAddr        string  `orm:"size(32);default(null);" description:"(主机ip 内)"`
	PublicAddr          string  `orm:"size(32);default(null);" description:"(主机ip 外)"`
	Status              string  `orm:"size(32);default(Abnormal)" description:"(主机状态 正常 Normal 异常 Abnormal)"`
	Diss                string  `orm:"size(32);default(NotInstalled)" description:"(安全容器 Installed NotInstalled)"`
	DissStatus          string  `orm:"size(32);default(Unsafe)" description:"(安全状态 Safe Unsafe)"`
	AccountName         string  `orm:"size(32);default(admin)" description:"(租户)"`
	GroupId             string  `orm:"-" description:"(查询参数：分组Id， 仅仅是查询使用, 返回数据看 Group)"`
	Group               *Groups `orm:"rel(fk);null;on_delete(set_null)" description:"(分组)"`
	Type                string  `orm:"size(32);default(Server);" description:"(类型 服务器: Server 虚拟机: Vm)"`
	IsInK8s             bool    `orm:"default(false);" description:"(是否在k8s集群)"`
	ClusterId           string  `orm:"size(128);default(null);" description:"(集群id)"`
	ClusterName         string  `orm:"size(128);default(null);" description:"(集群名)"`
	Label               string  `orm:"size(32);default(null);" description:"(标签)"`
	TaskList            []*Task `orm:"reverse(many);null" description:"(任务列表)"`
	Job                 []*Job  `orm:"rel(m2m);null;" description:"(job)"`
	IsEnableHeartBeat   bool    `orm:"default(false);" description:"(是否开启心跳上报)"`
	IsEnableDockerEvent bool    `orm:"default(false);" description:"(是否开启容器审计上报)"`
	HeartBeat           int64   `orm:"null;type(0)" description:"(心跳)"`
	KMetaData           string  `orm:"" description:"(源数据)"`
	KSpec               string  `orm:"" description:"(Spec数据)"`
	KStatus             string  `orm:"" description:"(状态数据)"`
	KubernetesVer       string  `orm:"size(64)" description:"(kubernetes 版本)"`
	NodeRole            string  `orm:"size(64)" description:"(集群主机角色)"`
	DockerCISCount      string  `orm:"null;" description:"(docker基线结果个数)"`
	DockerCISUpdateTime int64   `orm:"default(0)" description:"(docker基线更新时间)"`
	KubeCISCount        string  `orm:"null;" description:"(k8s基线结果个数)"`
	KubeCISUpdateTime   int64   `orm:"default(0)" description:"(k8s基线更新时间)"`
	IsLicensed          bool    `orm:"default(false);" description:"(是否已经授权)"`
	LicCount            bool    `orm:"-" description:"(是否获取授权个数操作)"`
	CreateTime          int64   `orm:"default(0)" description:"(上线时间)"`
	OfflineTime         int64   `orm:"default(0)" description:"(离线时间)"`
	WithK8sBench        bool    `orm:"-" description:"(是否获取k8s基线统计)"`
}

type HostConfigInterface interface {
	Add() error
	List(from, limit int) Result

	BaseList(from, limit int) (error, int64, []*HostConfig)
	Update() Result
	Delete() Result
	UpdateDynamic() Result
	Count() int64
	GetBenchMarkProportion() (int64, int64)
	Get() *HostConfig
}

func (this *HostConfig) Add() error {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var err error

	if hostConfig := this.Get(); hostConfig != nil {
		updateHostConfig := hostConfig
		// agent 或者 k8s 数据更新 （因为有diss-backend的关系数据，防止覆盖diss-backend的数据，需要替换更新）

		if this.HostName != "" {
			updateHostConfig.HostName = this.HostName
		}
		if this.InternalAddr != "" {
			updateHostConfig.InternalAddr = this.InternalAddr
		}
		if this.PublicAddr != "" {
			updateHostConfig.PublicAddr = this.PublicAddr
		}

		if this.OS != "" {
			updateHostConfig.OS = this.OS
		}
		if this.PG != "" {
			updateHostConfig.PG = this.PG
		}

		if this.ClusterId != "" {
			updateHostConfig.ClusterId = this.ClusterId
		}
		if this.ClusterName != "" {
			updateHostConfig.ClusterName = this.ClusterName
		}
		if this.NodeRole != "" {
			updateHostConfig.NodeRole = this.NodeRole
		}
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
	hostConfig := new(HostConfig)
	var err error
	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.InternalAddr != "" {
		cond = cond.And("internal_addr", this.InternalAddr)
	}
	err = o.QueryTable(utils.HostConfig).SetCond(cond).RelatedSel().One(hostConfig)
	if err != nil {
		return nil
	}
	return hostConfig
}

func (this *HostConfig) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	err, total, list := this.BaseList(from, limit)
	if err != nil {
		ResultData.Code = utils.GetHostConfigErr
	}
	data := make(map[string]interface{})
	data[Result_Items] = list
	data[Result_Total] = total
	ResultData.Data = data
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *HostConfig) BaseList(from, limit int) (error, int64, []*HostConfig) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var list []*HostConfig
	var TaskList []*Task
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
	if this.PublicAddr != "" {
		cond = cond.And("public_addr__contains", this.PublicAddr)
	}
	if this.InternalAddr != "" {
		cond = cond.And("internal_addr__contains", this.InternalAddr)
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
	if this.ClusterName != "" {
		cond = cond.And("cluster_name", this.ClusterName)
	}
	if this.Status != "" {
		cond = cond.And("status", this.Status)
	}
	if this.IsLicensed {
		cond = cond.And("is_licensed", this.IsLicensed)
	}
	if this.LicCount {
		cond = cond.And("is_licensed", true)
	}
	_, err = o.QueryTable(utils.HostConfig).SetCond(cond).Limit(limit, from).OrderBy("-host_name").RelatedSel().All(&list)
	total, _ := o.QueryTable(utils.HostConfig).SetCond(cond).Count()
	if err != nil {
		return err, 0, nil
	}

	for _, l := range list {
		cond = orm.NewCondition()
		cond = cond.And("host_id", l.Id)
		if this.ClusterId != "" {
			cond = cond.And("cluster_id", this.ClusterId)
		}
		_, err = o.QueryTable(utils.Task).SetCond(cond).RelatedSel().Limit(1, 0).OrderBy("-update_time").All(&TaskList)
		l.TaskList = TaskList

		if this.WithK8sBench && l.ClusterId != "" {
			bml := BenchMarkLog{}
			bml.HostId = l.Id
			bml.Type = BMLT_K8s
			bmlData, _ := bml.Get()
			if bmlData != nil {
				k8sMarkSummary := MarkSummary{}
				k8sMarkSummary.FailCount = bmlData.FailCount
				k8sMarkSummary.PassCount = bmlData.PassCount
				k8sMarkSummary.InfoCount = bmlData.InfoCount
				k8sMarkSummary.WarnCount = bmlData.WarnCount
				k8sMarkSummaryJson, _ := json.Marshal(k8sMarkSummary)
				l.KubeCISCount = string(k8sMarkSummaryJson)
			}
		}

	}

	return nil, total, list
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
	if this.LicCount {
		cond = cond.And("IsLicensed", true)
	}

	if this.Status != "" {
		cond = cond.And("status", this.Status)
	}
	count, _ := o.QueryTable(utils.HostConfig).SetCond(cond).Count()
	return count
}

// docker基线 / k8s 基线
func (this *HostConfig) GetBenchMarkProportion() (int64, int64) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	dockerBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Count()
	k8sBenchMarkCount, _ := o.QueryTable(utils.HostConfig).Filter("is_in_k8s", true).Count()
	return dockerBenchMarkCount, k8sBenchMarkCount
}

// 已安全容器 / 未安装安全容器
func (this *HostConfig) GetDissCountProportion() (int64, int64) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	safeCount, _ := o.QueryTable(utils.HostConfig).Filter("diss_status", Diss_status_Safe).Count()
	unSafeCount, _ := o.QueryTable(utils.HostConfig).Filter("diss_status", Diss_Status_Unsafe).Count()
	return safeCount, unSafeCount
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

func (this *HostConfig) UpdateHostCISCount(benchMarkSummary MarkSummary) {
	hostConfig := this.Get()
	if hostConfig != nil {
		kubeCISCount := map[string]int{}
		kubeCISCount[CIS_FailCount] = benchMarkSummary.FailCount
		kubeCISCount[CIS_InfoCount] = benchMarkSummary.InfoCount
		kubeCISCount[CIS_PassCount] = benchMarkSummary.WarnCount
		kubeCISCount[CIS_WarnCount] = benchMarkSummary.PassCount
		benchCountJson, _ := json.Marshal(kubeCISCount)
		hostConfig.KubeCISCount = string(benchCountJson)
		hostConfig.Update()
	}
}

// Restore kube result summary for renew cluster node or reopen license host
func (this *HostConfig) RestoreKubeBenchSummary() {
	kubeCountStr := this.KubeCISCount
	benchmarkSummary := MarkSummary{}
	_ = json.Unmarshal([]byte(kubeCountStr), &benchmarkSummary)
	if benchmarkSummary.WarnCount == 0 && benchmarkSummary.InfoCount == 0 && benchmarkSummary.PassCount == 0 && benchmarkSummary.FailCount == 0 {
		benchmarkQuery := BenchMarkLog{}
		benchmarkQuery.HostId = this.Id
		benchmarkQuery.Type = BMLT_K8s
		benchmarkData := benchmarkQuery.List(0, 1)
		if benchmarkData.Data != nil {
			benchmarkMap := benchmarkData.Data.(map[string]interface{})
			if benchmarkMap != nil {
				benchmarkList := benchmarkMap["items"].([]*BenchMarkLog)
				if len(benchmarkList) > 0 {
					benchmarkSummary := MarkSummary{}
					benchmarkSummary.FailCount = benchmarkList[0].FailCount
					benchmarkSummary.InfoCount = benchmarkList[0].InfoCount
					benchmarkSummary.PassCount = benchmarkList[0].PassCount
					benchmarkSummary.WarnCount = benchmarkList[0].WarnCount
					this.KubeCISUpdateTime = benchmarkList[0].UpdateTime
					this.UpdateHostCISCount(benchmarkSummary)
				}
			}
		}
	}
}
