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

func init() {
	orm.RegisterModel(new(DcokerIds))
}

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
