package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DockerEvent struct {
	Id        string `orm:"pk;size(128)" description:"(容器id)"`
	HostId    string `orm:"size(256)" description:"(主机Id agent采集数据)"`
	From      string `orm:"size(256)" description:"(镜像来源)"`
	Type      string `orm:"size(32)" description:"(类型)"`
	Action    string `orm:"size(32)" description:"(执行操作)"`
	Actor     string `orm:"" description:"(操作明细)"`
	Status    string `orm:"size(32)" description:"(状态)"`
	Scope     string `orm:"size(64)" description:"(范围)"`
	Time      int64  `orm:"" description:"(时间)"`
	TimeNano  int64  `orm:"" description:"(精确时间)"`
	StartTime string `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime   string `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
}

type DockerEventInterface interface {
	List(from, limit int) Result
}

func (this *DockerEvent) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var dockerEventList []*DockerEvent = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from docker_event `
	countSql := `select "count"(id) from docker_event `
	filter := ""
	if this.Id != "" {
		filter = filter + `id = '` + this.Id + `' and `
	}
	if this.HostId != "" {
		filter = filter + `host_id = '` + this.HostId + `' and `
	}
	if this.From != "" {
		filter = filter + `docker_event."from" like '%` + this.From + `%' and `
	}
	if this.Status != "" {
		filter = filter + `status = '` + this.Status + `' and `
	}
	if this.Actor != "" {
		filter = filter + `actor like '%` + this.Action + `%' and `
	}
	if this.Action != "" {
		filter = filter + `action = '` + this.Action + `' and `
	}
	if this.Scope != "" {
		filter = filter + `scope = '` + this.Scope + `' and `
	}

	if this.StartTime != "" && this.EndTime != "" {
		startTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.Local)
		endTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.EndTime, time.Local)
		filter = filter + `time BETWEEN ` + fmt.Sprintf("%v", startTime.Unix()) + ` and '` + fmt.Sprintf("%v", endTime.Unix()) + `' and `
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql).QueryRows(&dockerEventList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDockerEventListErr
		logs.Error("Get DockerEvent List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dockerEventList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
