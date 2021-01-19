package models

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Cluster struct {
	Id          string    `orm:"pk" description:"(集群id)"`
	Name        string    `orm:"unique;size(32)" description:"(集群名)"`
	FileName    string    `orm:"size(255)" description:"(KubeConfig 文件)"`
	AuthType    string    `orm:"size(32);default(BearerToken)" description:"(认证类型 KubeConfig BearerToken)"`
	BearerToken string    `orm:"" description:"(Token)"`
	MasterUrls  string    `orm:"size(255)" description:"(ApiServer 访问地址)"`
	Status      string    `orm:"size(32);default(Active)" description:"(集群状态 Active Unavailable)"`
	Type        string    `orm:"size(32);default(Kubernetes)" description:"(类型 Kubernetes Openshift Rancher)"`
	IsSync      bool      `orm:"default(false)" description:"(是否同步)"`
	Label       string    `orm:"size(64);default(null)" description:"(标签)"`
	ScopeUrl    string    `orm:"size(512);default()" description:"(scope访问地址)"`
	SocpeStatus string    `orm:"size(64);default()" description:"(scope 操作状态)"`
	AccountName string    `orm:"-" description:"(租户)"`
	SyncStatus  string    `orm:"default(NotSynced)" description:"(同步状态 NotSynced 未同步 Synced 成功 InProcess 同步中 Fail 失败 Clearing 清理中)"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime  time.Time `orm:"auto_now;type(datetime)" description:"(更新时间)"`
}

type ClusterCheck struct {
	Id           string `description:"(集群id)"`
	Name         string `description:"(集群名)"`
	DockerCIS    bool   `description:"(Docker CIS)"`
	KubenetesCIS bool   `description:"(kubernetes CIS)"`
	VirusScan    bool   `description:"(病毒检查)"`
	LeakScan     bool   `description:"(漏洞检查)"`
	Batch        int64  `description:"(批次)"`
}

type ClusterInterface interface {
	Get() *Cluster
	Add() Result
	Update() Result
	List(from, limit int) Result
	ListByAccount(from, limit int) Result
	GetRequiredSyncList() Result
}

func (this *Cluster) Get() *Cluster {
	o := orm.NewOrm()
	C := new(Cluster)

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	err := o.QueryTable(utils.Cluster).SetCond(cond).One(C)
	if err != nil {
		logs.Error("Get cluster failed, code: %d, err: %s", utils.GetClusterErr, err.Error())
		return nil
	}
	return C
}

func (this *Cluster) Add(isForce bool) Result {
	o := orm.NewOrm()
	ResultData := Result{Code: http.StatusOK}
	cond := orm.NewCondition()
	cluster := new(Cluster)

	if this.Name != "" {
		cond = cond.And("name", this.Name)
	}

	if this.MasterUrls != "" {
		cond = cond.Or("master_urls__contains", this.MasterUrls)
	}

	if this.FileName != "" {
		cond = cond.Or("file_name", this.FileName)
	}

	o.QueryTable(utils.Cluster).SetCond(cond).One(cluster)
	if isForce && cluster.Id != "" { // 更新
		cluster.Name = this.Name
		cluster.FileName = this.FileName
		cluster.BearerToken = this.BearerToken
		cluster.MasterUrls = this.MasterUrls
		cluster.Status = Cluster_Status_Active
		cluster.IsSync = Cluster_IsSync
		ResultData = cluster.Update()
		if ResultData.Code != http.StatusOK {
			return ResultData
		}
		this = cluster
	} else { // 添加
		// 根据 master_urls 或者 集群名的唯一性 判断是否重复
		count, err := o.QueryTable(utils.Cluster).SetCond(cond).Count()
		if count != 0 {
			ResultData.Message = "ClusterIsExistErr"
			ResultData.Code = utils.ClusterIsExistErr
			logs.Error("Add Cluster failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
		_, err = o.Insert(this)

		if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddClusterErr
			logs.Error("Add Cluster failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}
	ResultData.Data = this
	return ResultData
}

func (this *Cluster) GetRequiredSyncList() Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	var ClusterList []*Cluster
	var ResultData Result
	var err error
	var total int64
	cond := orm.NewCondition()
	cond = cond.Or("sync_status", Cluster_Sync_Status_Fail).Or("sync_status", Cluster_Sync_Status_NotSynced)
	_, err = o.QueryTable(utils.Cluster).SetCond(cond).All(&ClusterList)
	total, _ = o.QueryTable(utils.Cluster).SetCond(cond).Count()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get Cluster List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ClusterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Cluster) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	var ClusterList []*Cluster
	var ResultData Result
	var err error
	var total int64
	cond := orm.NewCondition()

	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.Status != "" && this.Status != All {
		cond = cond.And("status", this.Status)
	}
	if this.Type != "" && this.Type != All {
		cond = cond.And("type", this.Type)
	}
	if this.Label != "" {
		cond = cond.And("label__contains", this.Label)
	}
	_, err = o.QueryTable(utils.Cluster).SetCond(cond).Limit(limit, from).All(&ClusterList)
	total, _ = o.QueryTable(utils.Cluster).SetCond(cond).Count()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get Cluster List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ClusterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Cluster) Update() Result {
	o := orm.NewOrm()
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditClusterErr
		logs.Error("Update cluster: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Cluster) ListByAccount(from, limit int) Result {
	o := orm.NewOrm()
	var ClusterList []*Cluster
	var cIds []string
	ns := new(NameSpace)
	var ResultData Result
	var err error
	var total int64
	cond := orm.NewCondition()

	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	if this.AccountName != "" && this.AccountName != Account_Admin {
		//根据命名空间查询绑定关系
		ns.AccountName = this.AccountName
		_, cIds = ns.ListByAccountGroupByClusterId()
		if cIds != nil {
			cond = cond.And("id__in", cIds)
			total, err = o.QueryTable(utils.Cluster).SetCond(cond).Limit(limit, from).All(&ClusterList)
		}
	} else {
		_, err = o.QueryTable(utils.Cluster).SetCond(cond).Limit(limit, from).All(&ClusterList)
		total, _ = o.QueryTable(utils.Cluster).SetCond(cond).Count()
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get Cluster List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ClusterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Cluster) Delete() Result {
	o := orm.NewOrm()
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err := o.QueryTable(utils.Cluster).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteClusterErr
		logs.Error("Delete Cluster failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
