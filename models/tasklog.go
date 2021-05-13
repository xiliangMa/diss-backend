package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TaskLog struct {
	Id         string `orm:"pk;size(64)" description:"(任务id)"`
	Account    string `orm:"default(admin);size(32)" description:"(租户)"`
	Task       string `orm:"" description:"(任务详情json)"`
	RawLog     string `orm:"" description:"(日志)"`
	Level      string `orm:"default(Info);size(32)" description:"(日志级别)"`
	CreateTime int64  `orm:"default(0)" description:"(创建时间)"`
	StartTime  int64  `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime    int64  `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
}

type TaskLog1Interface interface {
	Add() Result
	List(from, limit int) Result
}

func (this *TaskLog) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var taskLogList []*TaskLog = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from ` + utils.TaskLog + ` `
	countSql := `select "count"(id) from ` + utils.TaskLog + ` `
	filter := ""
	fields := []string{}
	if this.Id != "" {
		filter = filter + `id = ? and `
		fields = append(fields, this.Id)
	}
	if this.Account != "" {
		filter = filter + `account = ? and `
		fields = append(fields, this.Account)
	}
	if this.Task != "" {
		filter = filter + `task like ? and `
		fields = append(fields, "%"+this.Task+"%")
	}
	if this.RawLog != "" {
		filter = filter + `raw_log like ? and `
		fields = append(fields, "%"+this.RawLog+"%")
	}
	if this.Level != "" {
		filter = filter + `level = ? and `
		fields = append(fields, this.Level)
	}
	if this.StartTime != 0 && this.EndTime != 0 {
		filter = filter + `create_time BETWEEN ` + fmt.Sprintf("%v", this.StartTime) + ` and ` + fmt.Sprintf("%v", this.EndTime) + ` and `
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	orderBySql := ` ORDER BY create_time desc`
	if from >= 0 && limit > 0 {
		limitSql := orderBySql + ` limit ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql, fields).QueryRows(&taskLogList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetTaskLogErr
		logs.Error("Get TaskLog List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = taskLogList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *TaskLog) Add() Result {
	insetSql := `INSERT INTO ` + utils.TaskLog + ` VALUES(?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	uid, _ := uuid.NewV4()
	_, err := o.Raw(insetSql,
		uid.String(),
		this.Account,
		this.Task,
		this.RawLog,
		this.Level,
		time.Now().UnixNano()).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddTaskLogErr
		logs.Error("Add TaskLog failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
