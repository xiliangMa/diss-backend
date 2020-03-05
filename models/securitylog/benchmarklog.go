package securitylog

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type BenchMarkLog struct {
	Id            string `orm:"pk;description(基线id)"`
	BenchMarkName string `orm:"description(基线模版名称)"`
	Level         string `orm:"description(级别)"`
	ProjectName   string `orm:"description(测试项目)"`
	HostName      string `orm:"description(主机名称)"`
	HostId        string `orm:"description(主机Id)"`
	InternalAddr  string `orm:"description(主机ip 内)"`
	PublicAddr    string `orm:"description(主机ip 外)"`
	OS            string `orm:"description(系统)"`
	UpdateTime    string `orm:"description(更新时间)"`
	FailCount     string `orm:"description(检查失败个数)"`
	WarnCount     string `orm:"description(检查警告个数)"`
	PassCount     string `orm:"description(检查通过个数)"`
	InfoCount     string `orm:"description(检查提示个数)"`
	RawLog        string `orm:"description(结果原始内容)"`
	Type          string `orm:"description(分类)"`
	Result        string `orm:"description(测试结果)"`
}

func init() {
	orm.RegisterModel(new(BenchMarkLog))
}

type BenchMarkLogInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *BenchMarkLog) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result
	var err error

	this.RawLog = ""
	_, err = o.Insert(this)
	if err != nil {
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
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var BenchMarkLogList []*BenchMarkLog = nil
	var ResultData models.Result
	var err error
	cond := orm.NewCondition()

	cond = cond.And("host_id", this.HostId)
	if this.BenchMarkName != "" {
		cond = cond.And("bench_mark_name", this.BenchMarkName)
	}
	_, err = o.QueryTable(utils.BenchMarkLog).SetCond(cond).Limit(limit, from).All(&BenchMarkLogList)

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
