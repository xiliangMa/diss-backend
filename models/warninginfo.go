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

type WarningInfo struct {
	Id            string `orm:"pk;size(128)" description:"(Id)"`
	Name          string `orm:"size(256)" description:"(告警名称)"`
	HostId        string `orm:"size(256)" description:"(主机Id agent采集数据)"`
	HostName      string `orm:"size(256)" description:"(主机Name agent采集数据)"`
	Cluster       string `orm:"size(256)" description:"(集群名)"`
	Account       string `orm:"size(256)" description:"(租户)"`
	Type          string `orm:"size(32)" description:"(类型 如：基线检测、病毒检查、入侵检测、镜像安全等)"`
	Info          string `orm:"size(1024)" description:"(告警详情，json，请自定义内部结构)"`
	Level         string `orm:"size(32)" description:"(告警级别)"`
	Status        string `orm:"size(256)" description:"(状态)"`
	CreateTime    int64  `orm:"" description:"(发生时间)"`
	UpdateTime    int64  `orm:"" description:"(更新时间)"`
	StartTime     int64  `orm:"-" description:"(开始时间, 注意时间格式为 local 时间)"`
	EndTime       int64  `orm:"-" description:"(结束时间, 注意时间格式为 local 时间)"`
	Proposal      string `orm:"size(256)" description:"(建议)"`
	Analysis      string `orm:"size(256)" description:"(关联分析)"`
	Mode          string `orm:"size(128)" description:"(方式)"`
	ContainerId   string `orm:"-"`
	ContainerName string `orm:"-"`
	ImageName     string `orm:"-"`
	Action        string `orm:"-" description:"(处理方式：isolation、pause、stop、kill)"`
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
	fields := []string{}
	if this.Id != "" {
		filter = filter + `id = ? and `
		fields = append(fields, this.Id)
	}
	if this.HostName != "" {
		filter = filter + `host_name like ? and `
		fields = append(fields, "%"+this.HostName+"%")
	}
	if this.HostId != "" {
		filter = filter + `host_id like ? and `
		fields = append(fields, "%"+this.HostId+"%")
	}
	if this.Status != "" {
		filter = filter + `status = ? and `
		fields = append(fields, this.Status)
	}
	if this.Name != "" {
		filter = filter + `name like ? and `
		fields = append(fields, "%"+this.Name+"%")
	}
	if this.Account != "" {
		filter = filter + `account like ? and `
		fields = append(fields, "%"+this.Account+"%")
	}
	if this.Cluster != "" {
		filter = filter + `cluster = ? and `
		fields = append(fields, this.Cluster)
	}
	if this.Type != "" {
		filter = filter + `type = ? and `
		fields = append(fields, this.Type)
	}
	if this.Level != "" {
		filter = filter + `level = ? and `
		fields = append(fields, this.Level)
	}

	if this.StartTime != 0 && this.EndTime != 0 {
		//startTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.Local)
		//endTime, _ := time.ParseInLocation("2006-01-02T15:04:05", this.EndTime, time.Local)
		//filter = filter + `create_time BETWEEN  '` + fmt.Sprintf("%v", startTime.Unix()) + `' and '` + fmt.Sprintf("%v", endTime.Unix()) + `' and `
		startTime := strconv.FormatInt(this.StartTime, 10)
		endTime := strconv.FormatInt(this.EndTime, 10)

		filter = filter + `create_time BETWEEN  ` + fmt.Sprintf("%v", startTime) + ` and ` + fmt.Sprintf("%v", endTime) + ` and `
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " order by create_time desc limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql, fields).QueryRows(&WarningInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetWarningInfoListErr
		logs.Error("Get WarningInfo list failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = WarningInfoList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *WarningInfo) Add() Result {
	insertSql := `INSERT INTO ` + utils.WarningInfo + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	if this.Id == "" {
		uid, _ := uuid.NewV4()
		this.Id = uid.String()
	}

	if this.CreateTime == 0 {
		this.CreateTime = time.Now().UnixNano()
	}
	if this.UpdateTime == 0 {
		this.UpdateTime = time.Now().UnixNano()
	}
	_, err := o.Raw(insertSql,
		this.Id,
		this.Name,
		this.HostId,
		this.HostName,
		this.Cluster,
		this.Account,
		this.Type,
		this.Info,
		this.Level,
		this.Status,
		this.CreateTime,
		this.UpdateTime,
		this.Proposal,
		this.Analysis, this.Mode).Exec()

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

func (this *WarningInfo) Update() Result {
	var err error
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	this.UpdateTime = time.Now().UnixNano()
	if this.Status != "" {
		_, err = o.Raw(`UPDATE `+utils.WarningInfo+` SET status=? WHERE id=?`, this.Status, this.Id).Exec()
	}

	if this.Proposal != "" {
		_, err = o.Raw(`UPDATE `+utils.WarningInfo+` SET proposal=? WHERE id=?`, this.Proposal, this.Id).Exec()
	}

	if this.Proposal != "" && this.Status != "" {
		_, err = o.Raw(`UPDATE `+utils.WarningInfo+` SET proposal=?,status=? WHERE id=?`, this.Proposal, this.Status, this.Id).Exec()
	}

	if this.Mode != "" {
		_, err = o.Raw(`UPDATE `+utils.WarningInfo+` SET mode=? WHERE id=?`, this.Mode, this.Id).Exec()
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.UpdateWarningInfoErr
		logs.Error("Update WarningInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
