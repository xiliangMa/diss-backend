package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RespCenter struct {
	Id            string `orm:"pk;size(128)" description:"(Id)"`
	HostId        string `orm:"size(256)" description:"(主机Id agent采集数据)"`
	HostName      string `orm:"size(256)" description:"(主机Name agent采集数据)"`
	Account       string `orm:"size(256)" description:"(租户)"`
	Cluster       string `orm:"size(256)" description:"(集群名)"`
	Namespace     string `orm:"size(256)" description:"(命名空间)"`
	ContainerId   string `orm:"size(256)" description:"(容器id)"`
	ContainerName string `orm:"size(256)" description:"(容器名称)"`
	PodName       string `orm:"size(256)" description:"(pod名称)"`
	ImageName     string `orm:"size(256)" description:"(镜像名称)"`
	CreateTime    int64  `orm:"" description:"(发生时间)" json:"ct"`
	UpdateTime    int64  `orm:"" description:"(更新时间)" json:"ut"`
	Status        string `orm:"size(128)" description:"(状态)"`
	Proposal      string `orm:"size(256)" description:"(隔离原因)"`
	ProcessNote   string `orm:"size(256)" description:"(处理说明)"`
	WarningInfoId string `orm:"size(128)" description:"(外键id)" `
	Action        string `orm:"-" description:"(处理方式：resume、other)"`
}

type RespCenterInterface interface {
	List(from, limit int) Result
	Add()
	Update() Result
}

func (this *RespCenter) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResponseCenterList []*RespCenter = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from ` + utils.RespCenter + ` `
	countSql := `select "count"(id) from ` + utils.RespCenter + ` `
	filter := ""
	if this.HostName != "" {
		filter = filter + `host_name like '%` + this.HostName + `%' and `
	}
	if this.ContainerName != "" {
		filter = filter + `container_name like '%` + this.ContainerName + `%' and `
	}
	if this.Status != "" {
		filter = filter + `status = '` + this.Status + `' and `
	}
	if this.ImageName != "" {
		filter = filter + `image_name like '%` + this.ImageName + `%' and `
	}
	if this.Account != "" {
		filter = filter + `account like '%` + this.Account + `%' and `
	}
	if this.Cluster != "" {
		filter = filter + `cluster = '` + this.Cluster + `' and `
	}
	if this.Namespace != "" {
		filter = filter + `namespace = '` + this.Namespace + `' and `
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
	_, err = o.Raw(resultSql).QueryRows(&ResponseCenterList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRespCenterErr
		logs.Error("Get RespCenter list failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = ResponseCenterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *RespCenter) Add() Result {
	insertSql := `INSERT INTO ` + utils.RespCenter + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
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
		this.HostId,
		this.HostName,
		this.Account,
		this.Cluster,
		this.Namespace,
		this.ContainerId,
		this.ContainerName,
		this.PodName,
		this.ImageName,
		this.CreateTime,
		this.UpdateTime,
		this.Status,
		this.Proposal,
		this.ProcessNote,
		this.WarningInfoId).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddRespCenterErr
		logs.Error("Add RespCenter failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *RespCenter) Update() Result {
	updateSql := `UPDATE ` + utils.RespCenter + ` SET status=?,proposal=?,process_note=? WHERE id=? or warning_info_id=?`
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	this.UpdateTime = time.Now().UnixNano()
	_, err := o.Raw(updateSql,
		this.Status,
		this.Proposal,
		this.ProcessNote,
		this.Id, this.Id).Exec()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.UpdateRespCenterErr
		logs.Error("Update RespCenter failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
