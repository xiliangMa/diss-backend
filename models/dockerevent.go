package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type DockerEvent struct {
	Id        string `orm:"pk;size(256)" description:"(容器id)"`
	HostId    string `orm:"size(256)" description:"(主机Id agent采集数据)"`
	HostName  string `orm:"size(256)" description:"(主机Name agent采集数据)"`
	From      string `orm:"size(256)" description:"(镜像来源)"`
	Type      string `orm:"size(256)" description:"(类型)"`
	Action    string `orm:"size(256)" description:"(执行操作)"`
	Actor     string `orm:"" description:"(操作明细)"`
	Status    string `orm:"size(256)" description:"(状态)"`
	Scope     string `orm:"size(256)" description:"(范围)"`
	Time      int64  `orm:"" description:"(时间)"`
	TimeNano  int64  `orm:"" description:"(精确时间)"`
	StartTime int64  `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime   int64  `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
}

type DockerEventInterface interface {
	List(from, limit int) Result
	Add() Result
	GetLatestTime() Result
}

func (this *DockerEvent) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var dockerEventList []*DockerEvent = nil
	var ResultData Result
	var err error
	var total int64 = 0
	var orderSql = " order by time_nano desc"

	sql := ` select * from docker_event `
	countSql := `select "count"(id) from docker_event `
	filter := ""
	if this.Id != "" {
		filter = filter + `id = '` + this.Id + `' and `
	}
	if this.HostId != "" {
		filter = filter + `host_id = '` + this.HostId + `' and `
	}
	if this.HostName != "" {
		filter = filter + `docker_event."host_name" like '%` + this.HostName + `%' and `
	}
	if this.From != "" {
		filter = filter + `docker_event."from" like '%` + this.From + `%' and `
	}
	if this.Status != "" {
		filter = filter + `status = '` + this.Status + `' and `
	}
	if this.Actor != "" {
		filter = filter + `actor like '%` + this.Actor + `%' and `
	}
	if this.Action != "" {
		filter = filter + `action = '` + this.Action + `' and `
	}
	if this.Scope != "" {
		filter = filter + `scope = '` + this.Scope + `' and `
	}

	if this.StartTime != 0 && this.EndTime != 0 {
		//startTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.Local)
		//endTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.EndTime, time.Local)
		filter = filter + `time BETWEEN ` + fmt.Sprintf("%v", this.StartTime) + ` and ` + fmt.Sprintf("%v", this.EndTime) + ` and `
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql + orderSql
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

func (this *DockerEvent) Add() Result {
	insetSql := `INSERT INTO docker_event VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insetSql,
		this.Id,
		this.HostId,
		this.HostName,
		this.From,
		this.Type,
		this.Action,
		this.Actor,
		this.Status,
		this.Scope,
		this.Time,
		this.TimeNano).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddDockerEventErr
		logs.Error("Add DockerEvent failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *DockerEvent) GetLatestTime() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	var DockerEvent *DockerEvent

	sql := `SELECT * FROM ` + utils.DockerEvent
	filter := ""

	if this.HostId != "" {
		filter = filter + `host_id = '` + this.HostId + `' and `
	}

	if filter != "" {
		sql = sql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	resultSql := sql

	limitSql := " ORDER BY time DESC limit " + strconv.Itoa(1) + " OFFSET " + strconv.Itoa(0)
	resultSql = resultSql + limitSql

	err := o.Raw(resultSql).QueryRow(&DockerEvent)

	if err != nil {
		ResultData.Code = utils.GetLatestTimeForDockerEventErr
		ResultData.Message = string(utils.GetLatestTimeForDockerEventErr)
		logs.Error("Get LatestTimeForDockerEventErr failed, code: %d, err: %s", utils.GetLatestTimeForDockerEventErr, err.Error())
		return ResultData
	}
	data := make(map[string]interface{})
	data["items"] = DockerEvent

	ResultData.Code = http.StatusOK
	ResultData.Data = data

	return ResultData
}
