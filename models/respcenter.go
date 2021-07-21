package models

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
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
	Status        string `orm:"size(256)" description:"(状态)"`
	Analysis      string `orm:"size(256)" description:"(建议)"`
	ProcessNote   string `orm:"size(256)" description:"(处理说明)"`
	Msg           string `orm:"size(1024)" description:"(消息)"`
	WarningInfoId string `orm:"size(128)" description:"(外键id)" `
	Mode          string `orm:"size(128)" description:"(方式)" `
	Code          int64  `orm:"" description:"(错误码)"`
	Action        string `orm:"-" description:"(处理方式：resume、start)"`
}

type RespCenterInterface interface {
	List(from, limit int) Result
	Add()
	Update() Result
	GetRespCenter() Result
	Get() *RespCenter
	Delete() Result
}

func (this *RespCenter) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var respCenterList []*RespCenter
	var ResultData Result
	var total int64

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
	_, _ = o.Raw(resultSql).QueryRows(&respCenterList)

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = respCenterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *RespCenter) Add() Result {
	insertSql := `INSERT INTO ` + utils.RespCenter + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)`
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
		this.Analysis,
		this.ProcessNote,
		this.Msg,
		this.WarningInfoId, this.Mode, this.Code).Exec()

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
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	this.UpdateTime = time.Now().UnixNano()
	_, err := o.Update(this)
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

func (this *RespCenter) GetRespCenter() Result {
	var respCenterList []*RespCenter
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	_, err := o.Raw("select * from "+utils.RespCenter+" where warning_info_id = ? order by create_time", this.Id).QueryRows(&respCenterList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRespCenterErr
		logs.Error("Get RespCenter failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	data := make(map[string]interface{})
	data[Result_Items] = respCenterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data

	return ResultData
}

func (this *RespCenter) Get() *RespCenter {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	rc := new(RespCenter)
	var err error
	cond := orm.NewCondition()
	if this.WarningInfoId != "" {
		cond = cond.And("warning_info_id", this.WarningInfoId)
	}

	err = o.QueryTable(utils.RespCenter).SetCond(cond).RelatedSel().One(rc)
	if err != nil {
		return nil
	}
	return rc
}

func (this *RespCenter) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id__in", strings.Split(this.Id, ","))
	} else {
		ResultData.Message = "No RespCenter Id"
		ResultData.Code = utils.DeleteRespCenterErr
		return ResultData
	}

	_, err := o.QueryTable(utils.RespCenter).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteRespCenterErr
		logs.Error("Delete RespCenter failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
