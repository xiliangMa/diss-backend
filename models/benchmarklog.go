package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type BenchMarkLog struct {
	Id            string `orm:"pk;" description:"(基线id)"`
	BenchMarkName string `orm:"" description:"(基线模版名称)"`
	BenchMarkType string `orm:"" description:"(基线模版类型，DockerBenchMark 和 KubernetesBenchMark)"`
	Level         string `orm:"" description:"(级别 info warn  fail pass)"`
	ProjectName   string `orm:"" description:"(测试项目)"`
	HostName      string `orm:"" description:"(主机名称)"`
	HostId        string `orm:"" description:"(主机Id)"`
	InternalAddr  string `orm:"" description:"(主机ip 内)"`
	PublicAddr    string `orm:"" description:"(主机ip 外)"`
	OS            string `orm:"" description:"(系统)"`
	UpdateTime    int64  `orm:"default(0)" description:"(更新时间)"`
	FailCount     int    `orm:"" description:"(检查失败个数, kubeCIS)"`
	WarnCount     int    `orm:"" description:"(检查警告个数, dockerCIS和kubeCIS)"`
	PassCount     int    `orm:"" description:"(检查通过个数, dockerCIS和kubeCIS)"`
	InfoCount     int    `orm:"" description:"(检查提示个数, dockerCIS和kubeCIS)"`
	NoteCount     int    `orm:"" description:"(检查Note个数, dockerCIS)"`
	RawLog        string `orm:"" description:"(结果原始内容)"`
	Type          string `orm:"" description:"(分类)"`
	Result        string `orm:"" description:"(测试结果)"`
	IsInfo        bool   `orm:"-" description:"(是否取日志原始内容)"`
	TaskId        string `orm:"size(64)" description:"(任务ID)"`
}

type BenchMarkLogInterface interface {
	Add() Result
	List(from, limit int) Result
	GetMarkSummary() Result
	GetHostMarkSummary() Result
}

func (this *BenchMarkLog) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	var err error

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

func (this *BenchMarkLog) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var BenchMarkLogList []*BenchMarkLog = nil
	var ResultData Result
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

	if this.Level != "" && this.Level != BML_Level_ALL {
		cond = cond.And("level", this.Level)
	}

	if this.Type != "" {
		cond = cond.And("type", this.Type)
	}

	if this.Result != "" && this.Result != BML_Result_ALL {
		cond = cond.And("result", this.Result)
	}

	if this.BenchMarkName != "" && this.BenchMarkName != BML_Template_ALL {
		cond = cond.And("bench_mark_name", this.BenchMarkName)
	}

	if this.BenchMarkType != "" && this.BenchMarkType != BML_Template_ALL {
		cond = cond.And("bench_mark_type", this.BenchMarkType)
	}

	if this.UpdateTime != 0 {
		cond = cond.And("update_time__gte", this.UpdateTime)
	}

	if this.InternalAddr != "" {
		cond = cond.And("internal_addr__contains", this.InternalAddr)
	}

	if this.TaskId != "" {
		cond = cond.And("task_id", this.TaskId)
	}

	isInfo := this.IsInfo
	if isInfo {
		_, err = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Limit(limit, from).OrderBy("-update_time").All(&BenchMarkLogList)
	} else {
		fields := []string{"id", "bench_mark_name", "level", "project_name", "host_name", "host_id", "internal_addr", "public_addr", "o_s", "update_time", "fail_count", "warn_count", "pass_count", "info_count", "type", "result", "task_id"}
		_, err = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Limit(limit, from).OrderBy("-update_time").
			All(&BenchMarkLogList, fields...)
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetBenchMarkLogErr
		logs.Error("Get BenchMarkLog List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}
	var total int64
	if this.UpdateTime == 0 {
		total, _ = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Count()
	} else {
		total, _ = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Filter("update_time__gte", this.UpdateTime).Count()
	}
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = BenchMarkLogList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

type MarkSummary struct {
	FailCount int
	WarnCount int
	PassCount int
	InfoCount int
}

func (this *BenchMarkLog) GetMarkSummary() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	// docker 基线统计
	dockerMarkSummary := new(MarkSummary)
	o.Raw(utils.GetMarkSummarySql(BMLT_Docker)).QueryRow(&dockerMarkSummary)

	// k8s 基线统计
	k8sMarkSummary := new(MarkSummary)
	o.Raw(utils.GetMarkSummarySql(BMLT_K8s)).QueryRow(&k8sMarkSummary)

	data := make(map[string]interface{})
	data[BMLT_Docker] = dockerMarkSummary
	data[BMLT_K8s] = k8sMarkSummary
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *BenchMarkLog) GetHostMarkSummary() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	// host 基线统计
	hostMarkSummary := new(MarkSummary)
	o.Raw(utils.GetHostMarkSummarySql(this.HostId)).QueryRow(&hostMarkSummary)
	data := make(map[string]interface{})
	data[BMLT_Host_All] = hostMarkSummary
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}
