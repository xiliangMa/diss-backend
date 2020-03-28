package securitylog

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type BenchMarkLog struct {
	Id            string `orm:"pk;" description:"(基线id)"`
	BenchMarkName string `orm:"" description:"(基线模版名称)"`
	Level         string `orm:"" description:"(级别 info warn  fail pass)"`
	ProjectName   string `orm:"" description:"(测试项目)"`
	HostName      string `orm:"" description:"(主机名称)"`
	HostId        string `orm:"" description:"(主机Id)"`
	InternalAddr  string `orm:"" description:"(主机ip 内)"`
	PublicAddr    string `orm:"" description:"(主机ip 外)"`
	OS            string `orm:"" description:"(系统)"`
	UpdateTime    string `orm:"" description:"(更新时间)"`
	FailCount     int    `orm:"" description:"(检查失败个数)"`
	WarnCount     int    `orm:"" description:"(检查警告个数)"`
	PassCount     int    `orm:"" description:"(检查通过个数)"`
	InfoCount     int    `orm:"" description:"(检查提示个数)"`
	RawLog        string `orm:"" description:"(结果原始内容)"`
	Type          string `orm:"" description:"(分类)"`
	Result        string `orm:"" description:"(测试结果)"`
}

type BenchMarkLogInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
	GetMarkSummary()
	GetHostMarkSummary()
}

func (this *BenchMarkLog) Add() models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData models.Result
	var err error

	//this.RawLog = ""
	_, err = o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddBenchMarkLogErr
		logs.Error("Add BenchMarkLog failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *BenchMarkLog) List(from, limit int) models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var BenchMarkLogList []*BenchMarkLog = nil
	var ResultData models.Result
	var err error
	cond := orm.NewCondition()

	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}

	if this.HostName != "" {
		cond = cond.And("host_name", this.HostName)
	}

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	if this.Level != "" {
		cond = cond.And("level", this.Level)
	}

	if this.BenchMarkName != "" {
		cond = cond.And("bench_mark_name", this.BenchMarkName)
	}
	_, err = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Limit(limit, from).OrderBy("-update_time").All(&BenchMarkLogList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetBenchMarkLogErr
		logs.Error("Get BenchMarkLog List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.BenchMarkLog).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = BenchMarkLogList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

type MarkSummary struct {
	FailCount int
	WarnCount int
	PassCount int
	InfoCount int
}

func (this *BenchMarkLog) GetMarkSummary() models.Result {
	var ResultData models.Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	// docker 基线统计
	dockerMarkSummary := new(MarkSummary)
	o.Raw(utils.GetMarkSummarySql(models.BMLT_Docker)).QueryRow(&dockerMarkSummary)

	// k8s 基线统计
	k8sMarkSummary := new(MarkSummary)
	o.Raw(utils.GetMarkSummarySql(models.BMLT_K8s)).QueryRow(&k8sMarkSummary)

	data := make(map[string]interface{})
	data[models.BMLT_Docker] = dockerMarkSummary
	data[models.BMLT_K8s] = k8sMarkSummary
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *BenchMarkLog) GetHostMarkSummary() models.Result {
	var ResultData models.Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	// host 基线统计
	hostMarkSummary := new(MarkSummary)
	o.Raw(utils.GetHostMarkSummarySql()).QueryRow(&hostMarkSummary)
	data := make(map[string]interface{})
	data[models.BMLT_Host_All] = hostMarkSummary
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}
