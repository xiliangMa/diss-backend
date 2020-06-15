package models

import (
	//"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserEvent struct {
	Id          string    `orm:"pk;size(128)" description:"(pod id)"`
	UserName    string    `orm:"size(32)" description:"(用户名)"`
	AccountName string    `orm:"size(32)" description:"(租户)"`
	RawLog      string    `orm:"default(null);" description:"(操作详情)"`
	ModelType   string    `orm:"size(32)" description:"(模块类型)"`
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	StartTime   string    `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime     string    `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
}

type UserEventInterface interface {
	Add() Result
	Delete() Result
	List(from, limit int) Result
}

func (this *UserEvent) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var userEventList []*UserEvent = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from ` + string(utils.UserEvent) + ` `
	countSql := `select "count"(id) from ` + string(utils.UserEvent) + ` `
	filter := ""
	if this.Id != "" {
		filter = filter + `id = '` + this.Id + `' and `
	}
	if this.UserName != "" {
		filter = filter + `user_name = '` + this.UserName + `' and `
	}
	if this.AccountName != "" {
		filter = filter + `account_name = '` + this.AccountName + `' and `
	}
	if this.ModelType != "" {
		filter = filter + `model_type = '` + this.ModelType + `' and `
	}
	if this.RawLog != "" {
		filter = filter + `raw_log like '%` + this.RawLog + `%' and `
	}

	if this.StartTime != "" && this.EndTime != "" {
		filter = filter + `create_time BETWEEN '` + this.StartTime + `' and '` + this.EndTime + `' and `
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
	_, err = o.Raw(resultSql).QueryRows(&userEventList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetUserEventListErr
		logs.Error("Get UserEvent List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = userEventList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *UserEvent) Add() Result {
	insetSql := `INSERT INTO ` + string(utils.UserEvent) + ` VALUES(?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	uid, _ := uuid.NewV4()
	this.Id = uid.String()
	_, err := o.Raw(insetSql,
		this.Id,
		this.UserName,
		this.ModelType,
		this.AccountName,
		this.RawLog,
		this.CreateTime).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddUserEventErr
		logs.Error("Add UserEvent failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
