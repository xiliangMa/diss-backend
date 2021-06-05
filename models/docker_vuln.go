package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type DockerVulnerabilities struct {
	Id         int    `orm:"pk;auto" description:"(Id)"`
	HostId     string `orm:"size(64)" description:"(主机ID)"`
	HostName   string `orm:"size(128)" description:"(主机名称)"`
	TaskId     string `orm:"size(64)" description:"(任务Id)"`
	Docker     string `orm:"" description:"(Docker版本信息, Docker库自定义数据结构)"`
	CveIds     string `orm:"" description:"(任务Id)"`
	CreateTime int64  `orm:"" description:"(创建时间)"`
}

type DockerVulnerabilitiesScanInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
}

func (this *DockerVulnerabilities) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	createTime := time.Now().UnixNano()
	this.CreateTime = createTime

	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddDockerVulnErr
		logs.Error("Add DockerScan failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		o.Rollback()
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *DockerVulnerabilities) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var DockerVulnList []*DockerVulnerabilities
	var ResultData Result
	var err error
	cond := orm.NewCondition()
	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.TaskId != "" {
		cond = cond.And("task_id", this.TaskId)
	}
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}

	_, err = o.QueryTable(utils.DockerVulnerabilities).SetCond(cond).Limit(limit, from).OrderBy("-id").All(&DockerVulnList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDockerVulnErr
		logs.Error("Get DockerScan List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.DockerVulnerabilities).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = DockerVulnList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *DockerVulnerabilities) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}
	_, err := o.QueryTable(utils.DockerVulnerabilities).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteDockerVulnErr
		logs.Error("Delete DeleteDockerVulnErr failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
