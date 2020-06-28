package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CmdHistory struct {
	Id            string    `orm:"pk;size(64)" description:"(id)"`
	HostId        string    `orm:"size(64)" description:"(主机id)"`
	HostName      string    `orm:"size(64)" description:"(主机名)"`
	ContainerId   string    `orm:"size(256)" description:"(容器id)"`
	ContainerName string    `orm:"size(256)" description:"(容器名)"`
	User          string    `orm:"size(32)" description:"(用户)"`
	Command       string    `orm:"" description:"(命令)"`
	CreateTime    time.Time `orm:"null;" description:"(更新时间)"`
	Type          string    `orm:"default(Host);size(32)" description:"(类型 Host Container)"`
	StartTime     string    `orm:"-;default(null)" description:"(开始时间)"`
	EndTime       string    `orm:"-;default(null)" description:"(结束时间)"`
}

type CmdHistoryList struct {
	List []*CmdHistory
}

type CmdHistoryInterface interface {
	Add() Result
	List() Result
	Delete() Result
}

func (this *CmdHistory) Add() Result {
	insetSql := `INSERT INTO cmd_history VALUES(?, ?, ?, ? , ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insetSql, this.Id,
		this.HostId,
		this.HostName,
		this.ContainerId,
		this.ContainerName,
		this.User,
		this.Command,
		this.CreateTime,
		this.Type).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddCmdHistoryErr
		logs.Error("Add CmdHistory failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *CmdHistory) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	deleteSql := `DELETE FROM cmd_history WHERE TYPE = ? AND host_id = ?`

	_, err := o.Raw(deleteSql, this.Type, this.HostId).Exec()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteCmdHistoryErr
		logs.Error("Delete CmdHistory failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *CmdHistory) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var cmdHistoryList []*CmdHistory
	var ResultData Result
	var err error
	var total = 0

	sql := `SELECT * FROM ` + utils.CmdHistory
	countSql := `SELECT "count"(id) FROM ` + utils.CmdHistory
	filter := ""

	if this.Type != "" {
		filter = filter + ` type = '` + this.Type + `' and `
	}
	if this.User != "" {
		filter = filter + ` cmd_history."user" = '` + this.User + `' and `
	}
	if this.HostId != "" {
		filter = filter + `host_id = '` + this.HostId + `' and `
	}
	if this.HostName != "" {
		filter = filter + `host_name like '%` + this.HostName + `%' and `
	}
	if this.ContainerId != "" {
		filter = filter + `container_id = '` + this.ContainerId + `' and `
	}
	if this.ContainerId != "" {
		filter = filter + `container_name like '%` + this.ContainerName + `%' and `
	}
	if this.Command != "" {
		filter = filter + `command like '%` + this.Command + `%' and `
	}

	if this.StartTime != "" && this.EndTime != "" {
		filter = filter + `create_time  BETWEEN '` + this.StartTime + `' and '` + this.EndTime + `' and `
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

	_, err = o.Raw(resultSql).QueryRows(&cmdHistoryList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetCmdHistoryErr
		logs.Error("Get CmdHistory List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = cmdHistoryList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
