package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type KubeScan struct {
	Id        int    `orm:"pk;auto" description:"(Id)"`
	ClusterId string `orm:"size(64)" description:"(集群Id)" json:"cluster_id"`
	TaskId    string `orm:"size(64)" description:"(任务Id)" json:"task_id"`
	//Nodes           []Nodes               `orm:"" description:"(节点列表)"`
	//Services        []Services            `orm:"" description:"(服务列表)"`
	Vulnerabilities []*KubeVulnerabilities `orm:"reverse(many);" description:"(漏洞列表)"`
	Severity        string                 `orm:"-" description:"(等级，查询参数)"`
}

// todo add nodes and services info
//type Nodes struct {
//	Type     string
//	Location string
//}
//
//type Services struct {
//	Service  string
//	Location string
//}

type KubeVulnerabilities struct {
	Id            int       `orm:"pk;auto" description:"(Id)"`
	Location      string    `orm:"" description:"(POD、服务、IP)"`
	Vid           string    `orm:"size(32)" description:"(漏洞ID)"`
	Category      string    `orm:"size(128)" description:"(类别)"`
	Severity      string    `orm:"size(32)" description:"(级别)"`
	Vulnerability string    `orm:"size(128)" description:"(漏洞名)"`
	Description   string    `orm:"" description:"(漏洞描述)"`
	Evidence      string    `orm:"" description:"(证据)"`
	AvdReference  string    `orm:"" description:"(参考)" json:"avd_reference"`
	Hunter        string    `orm:"size(64)" description:"(扫描组件)"`
	KubeVuln      *KubeScan `orm:"rel(fk);null" description:"(集群扫描记录)"`
}

type KubeScanInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
}

type KubeVulnerabilitiesInterface interface {
	Add() Result
}

func (this *KubeScan) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	o.Begin()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddKubeScanErr
		logs.Error("Add KubeScan failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		o.Rollback()
		return ResultData
	}
	for _, kubeVuln := range this.Vulnerabilities {
		kubeVuln.KubeVuln = this
		kubeVuln.Add()
	}
	o.Commit()
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *KubeScan) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var KubeScanList []*KubeScan
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.TaskId != "" {
		cond = cond.And("task_id", this.TaskId)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	_, err = o.QueryTable(utils.KubeScan).SetCond(cond).Limit(limit, from).All(&KubeScanList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetKubeScanErr
		logs.Error("Get KubeScan List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if this.Severity != "" {
		for _, kubeScan := range KubeScanList {
			cond2 := orm.NewCondition()
			cond2 = cond2.And("kube_vuln_id", kubeScan.Id)
			cond2 = cond2.And("severity", this.Severity)
			_, err = o.QueryTable(utils.KubeVulnerabilities).SetCond(cond2).Limit(limit, from).All(&kubeScan.Vulnerabilities)
		}
	} else {
		for _, kubeScan := range KubeScanList {
			o.LoadRelated(kubeScan, "Vulnerabilities", true)
		}
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetKubeScanErr
		logs.Error("Get KubeScan List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.KubeScan).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = KubeScanList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *KubeScan) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	if this.TaskId != "" {
		cond = cond.And("task_id", this.TaskId)
	}
	_, err := o.QueryTable(utils.KubeScan).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteKubeScanErr
		logs.Error("Delete KubeScan failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *KubeVulnerabilities) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddKubeVulnerabilitiesErr
		logs.Error("Add KubeVulnerabilities failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
