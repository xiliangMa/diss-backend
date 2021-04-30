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
}

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
