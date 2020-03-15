package securitylog

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"time"
)

// 针对 es 解析使用的结构体
type IntrudeDetectLog struct {
	HostId      string
	TargeType   string
	ContainerId string
	StartTime   string
	ToTime      string
	Limit       int
}

// 入侵检测日志（IntrudeDetectLog） 保存于 timescaledb
type DcokerIds struct {
	HostId       string    `orm:"pk" description:"(主机id)"`
	HostName     string    `description:"(主机名)"`
	MachineId    string    `description:"(Machine_id)"`
	ContainerId  string    `description:"(容器id)"`
	Time         time.Time `description:"(日志生成时间)"`
	Priority     string    `description:"(安全等级)"`
	Rule         string    `description:"(规则)"`
	Output       string    `description:"(事件信息)"`
	OutputFields string    `description:"(Output json)"`
	CreatedAt    int       `description:"(日志保存时间)"`
}

type DcokerIdsInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
	GetIntrudeDetectLogStatistics()
}

//func init() {
//	orm.RegisterModel(new(DcokerIds))
//}

func (this *DcokerIds) GetIntrudeDetectLogStatistics(timeCycle int) models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using(utils.DS_Security_Log)
	dcokerIdsCountList := make(map[int]*int)
	var ResultData models.Result
	var err error
	now := time.Now()
	for i := 0; i < timeCycle; i++ {
		count := 0
		m := i
		c1, _ := time.ParseDuration("-1h")
		c2, _ := time.ParseDuration("-" + strconv.Itoa(m+1) + "h")
		hStart := now.Add(c2).Unix()
		hEnd := now.Unix()
		if i != 0 {
			c1, _ = time.ParseDuration("-" + strconv.Itoa(i) + "h")
			hEnd = now.Add(c1).Unix()
		}
		err = o.Raw("select count(host_id) as count from docker_ids where created_at > ? and created_at < ?", hStart, hEnd).QueryRow(&count)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.GetIntrudeDetectLogErr
			logs.Error("Get IntrudeDetectLog List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
		dcokerIdsCountList[i] = &count
	}
	//QueryRows(&dcokerIdsList)
	//_, err = o.QueryTable(utils.DcokerIds).SetCond(cond).Limit(1, from).OrderBy("-time").All(&dcokerIdsList)

	data := make(map[string]interface{})
	data["items"] = dcokerIdsCountList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *IntrudeDetectLog) List(from, limit int) models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using(utils.DS_Security_Log)
	var dcokerIdsList []*DcokerIds = nil
	var ResultData models.Result
	var err error
	containerId := this.ContainerId
	//cond := orm.NewCondition()
	st, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.UTC)
	tt, _ := time.ParseInLocation("2006-01-02T15:04:05", this.ToTime, time.UTC)

	//cond = cond.And("host_id", this.HostId)
	//if this.HostId != "" {
	//	cond = cond.And("host_id", this.HostId)
	//}
	//if this.ContainerId != "" {
	//	cond = cond.And("container_id", this.ContainerId)
	//}
	//if this.StartTime != "" {
	//	loc, _ := time.LoadLocation("Asia/Shanghai")
	//	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", this.StartTime, loc)
	//	cond = cond.And("created_at__gte", tt.Unix())
	//}
	//if this.ToTime != "" {
	//	loc, _ := time.LoadLocation("Asia/Shanghai")
	//	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", this.ToTime, loc)
	//	cond = cond.And("created_at__lte", tt.Unix())
	//}
	if this.TargeType == models.IDLT_Docker {
		containerId = string([]byte(this.ContainerId)[:12])
	}

	total, err := o.Raw("select * from docker_ids where host_id = ? and container_id = ? and created_at > ? and created_at < ?", this.HostId, containerId, st.Unix(), tt.Unix()).QueryRows(&dcokerIdsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetIntrudeDetectLogErr
		logs.Error("Get IntrudeDetectLo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	//total, _ := o.QueryTable(utils.DcokerIds).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dcokerIdsList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
