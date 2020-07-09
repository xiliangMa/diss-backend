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

type WarningInfo struct {
	Id         string `orm:"pk;size(128)" description:"(Id)"`
	Name       string `orm:"size(256)" description:"(告警名称)"`
	HostName   string `orm:"size(256)" description:"(主机Name agent采集数据)"`
	Type       string `orm:"size(32)" description:"(类型)"`
	Level      string `orm:"size(32)" description:"(告警级别)"`
	Status     string `orm:"size(32)" description:"(状态)"`
	CreateTime string `orm:"size(128)" description:"(发生时间)"`
	UpdateTime string `orm:"size(128)" description:"(更新时间)"`
	StartTime  string `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime    string `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
}

type WarningInfoInterface interface {
	List(from, limit int) Result
	Add() Result
	Update() Result
}

func (this *WarningInfo) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var WarningInfoList []*WarningInfo = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from ` + utils.WarningInfo + ` `
	countSql := `select "count"(id) from ` + utils.WarningInfo + ` `
	filter := ""
	if this.Id != "" {
		filter = filter + `id = '` + this.Id + `' and `
	}
	if this.HostName != "" {
		filter = filter + `docker_event."host_name" like '%` + this.HostName + `%' and `
	}
	if this.Status != "" {
		filter = filter + `status = '` + this.Status + `' and `
	}
	if this.Name != "" {
		filter = filter + `name like '%` + this.Name + `%' and `
	}
	if this.Type != "" {
		filter = filter + `type = '` + this.Type + `' and `
	}
	if this.Level != "" {
		filter = filter + `level = '` + this.Level + `' and `
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
	_, err = o.Raw(resultSql).QueryRows(&WarningInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetWarningInfoListErr
		logs.Error("Get WarningInfo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = WarningInfoList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *WarningInfo) Add() Result {
	insetSql := `INSERT INTO ` + utils.WarningInfo + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insetSql,
		this.Id,
		this.Name,
		this.HostName,
		this.Type,
		this.Level,
		this.Status,
		this.CreateTime,
		this.UpdateTime).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddWarningInfoListErr
		logs.Error("Add WarningInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
