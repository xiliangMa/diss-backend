package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type VirusScan struct {
	Id            int    `description:"(id)"`
	Name          string `description:"(本次扫描任务名)"`
	TaskId        string `description:"(任务ID)"`
	HostId        string `description:"(主机ID)"`
	HostName      string `description:"(主机名称)"`
	ImageId       string `description:"(镜像ID)"`
	ImageName     string `description:"(镜像名称)"`
	ContainerId   string `description:"(容器ID)"`
	ContainerName string `description:"(容器名称)"`
	InternalAddr  string `description:"(内网IP)"`
	PublicAddr    string `description:"(外网IP)"`
	Type          string `description:"(扫描类型，HostVirusScan, ImageVirusScan, ContainerVirusScan)"`
	CreatedAt     int64  `description:"(创建时间)"`
	Records       []VirusRecord
	FileName      string
}

type VirusRecord struct {
	Id          int    `description:"(Id)"`
	VirusScanId int    `description:"(VirusScan Id)"`
	Filename    string `description:"(感染文件名称)"`
	Virus       string `description:"(感染病毒名称)"`
	Database    string `description:"(命中病毒库)"`
	Type        string `description:"(感染文件类型)"`
	Size        int64  `description:"(感染文件大小)"`
	Owner       string `description:"(感染文件所属用户)"`
	Permission  uint32 `description:"(感染文件权限)"`
	ModifyTime  int64  `description:"(感染文件最近修改时间)"`
	CreateTime  int64  `description:"(感染文件创建时间)"`
	Severity    string `orm:"size(32)" description:"(等级)"`
}

type VirusScanRecord struct {
	Id            int    `description:"(id)"`
	Name          string `description:"(本次扫描任务名)"`
	TaskId        string `description:"(任务ID)"`
	HostId        string `description:"(主机ID)"`
	HostName      string `description:"(主机名称)"`
	ImageId       string `description:"(镜像ID)"`
	ImageName     string `description:"(镜像名称)"`
	ContainerId   string `description:"(容器ID)"`
	ContainerName string `description:"(容器名称)"`
	InternalAddr  string `description:"(内网IP)"`
	PublicAddr    string `description:"(外网IP)"`
	Type          string `description:"(扫描类型，HostVirusScan, ImageVirusScan, ContainerVirusScan)"`
	CreatedAt     int64  `description:"(创建时间)"`
	VirusScanId   int    `description:"(VirusScan Id)"`
	Filename      string `description:"(感染文件名称)"`
	Virus         string `description:"(感染病毒名称)"`
	Database      string `description:"(命中病毒库)"`
	FileType      string `description:"(感染文件类型)"`
	Size          int64  `description:"(感染文件大小)"`
	Owner         string `description:"(感染文件所属用户)"`
	Permission    uint32 `description:"(感染文件权限)"`
	ModifyTime    int64  `description:"(感染文件最近修改时间)"`
	CreateTime    int64  `description:"(感染文件创建时间)"`
	Severity      string `description:"(等级)"`
}

type ImageVirus struct {
	ImageId     string `orm:"column(imageId)" description:"(镜像Id)"`
	ImageName   string `orm:"" description:"(镜像名)"`
	ImageDigest string `orm:"column(imageDigest)" description:"(镜像digest)"`
	UserId      string `orm:"column(userId)" description:"(用户id)"`
	FileName    string `description:"(文件名)"`
	Virus       string `description:"(病毒)"`
	FileHash    string `description:"(文件hash)"`
	FileSize    int64  `description:"(文件大小)"`
	CreatedAt   int64  `description:"(创建时间)"`
	LastUpdated int64  `description:"(更新时间)"`
}

type ContainerVirus struct {
	HostId      string `description:"(主机id)"`
	ContainerId string `description:"(容器id)"`
	FileName    string `description:"(文件名)"`
	Virus       string `description:"(病毒)"`
	FileHash    string `description:"(文件hash)"`
	FileSize    int64  `description:"(文件大小)"`
	CreatedAt   int64  `description:"(创建时间)"`
	LastUpdated int64  `description:"(更新时间)"`
	TargeType   string `description:"(类型)"`
}

type VirusScanInterface interface {
	Add() Result
	List(from, limit int) Result
	VirusList(from, limit int) ([]*VirusScanRecord, int64, error)
}

func (this *VirusScan) Add() Result {
	insetSql := `INSERT INTO virus_scan(name, task_id, host_id, host_name, image_id, image_name, container_id, container_name, internal_addr, public_addr, type, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) Returning id`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	var id int
	o.Raw(insetSql,
		this.Name,
		this.TaskId,
		this.HostId,
		this.HostName,
		this.ImageId,
		this.ImageName,
		this.ContainerId,
		this.ContainerName,
		this.InternalAddr,
		this.PublicAddr,
		this.Type,
		this.CreatedAt).QueryRow(&id)

	if id <= 0 {
		ResultData.Code = utils.AddVirusLogErr
		logs.Error("Add VirusScan failed, code: %d ", ResultData.Code)
		return ResultData
	}
	this.Id = id

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *VirusScan) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var virusLogList []*VirusScanRecord = nil
	var ResultData Result

	virusLogList, total, err := this.VirusList(from, limit)
	if err != nil {
		ResultData.Code = utils.GetImageVirusErr
		ResultData.Message = err.Error()
	}

	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = virusLogList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *VirusScan) VirusList(from, limit int) ([]*VirusScanRecord, int64, error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var virusLogList []*VirusScanRecord = nil
	var err error
	var total int64 = 0

	sql := `select *, virus_record.type as file_type from virus_record join 
 (select distinct on(created_at) *  from virus_scan order by created_at desc) as virus_log
    on virus_record.virus_scan_id = virus_log.id  `
	countSql := `select count(virus_record.id) from virus_record join
 (select distinct on(created_at) *  from virus_scan order by created_at desc)  as virus_log
    on virus_record.virus_scan_id = virus_log.id `
	filter := ""
	var fields []string
	if this.Id != 0 {
		filter = filter + `virus_log.id = ? and `
		fields = append(fields, string(rune(this.Id)))
	}
	if this.HostId != "" {
		filter = filter + `host_id = ? and `
		fields = append(fields, this.HostId)
	}
	if this.HostName != "" {
		filter = filter + `host_name like ? and `
		fields = append(fields, "%"+this.HostName+"%")
	}
	if this.ImageId != "" {
		filter = filter + `image_id = ? and `
		fields = append(fields, this.ImageId)
	}
	if this.ImageName != "" {
		filter = filter + `image_name like ? and `
		fields = append(fields, "%"+this.ImageName+"%")
	}
	if this.ContainerId != "" {
		filter = filter + `contaienr_id = ? and `
		fields = append(fields, this.ContainerId)
	}
	if this.ContainerName != "" {
		filter = filter + `container_name like ? and `
		fields = append(fields, "%"+this.ContainerName+"%")
	}
	if this.Type != "" {
		filter = filter + `virus_log.type = ? and `
		fields = append(fields, this.Type)
	}
	if this.CreatedAt != 0 {
		filter = filter + `virus_log.create_at  > this.CreatedAt and `
	}
	if this.FileName != "" {
		filter = filter + `virus_record.filename  like ? and `
		fields = append(fields, "%"+this.FileName+"%")
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
	_, err = o.Raw(resultSql, fields).QueryRows(&virusLogList)
	if err != nil {
		logs.Error("Get VirusLogErr List failed, code: %d, err: %s", utils.GetImageVirusErr, err.Error())
		return nil, 0, err
	}

	o.Raw(countSql, fields).QueryRow(&total)
	return virusLogList, total, nil
}

type VirusRecordInterface interface {
	List(from, limit int)
	Add() Result
	Count() int64
}

func (this *VirusRecord) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var virusLogList []*VirusScanRecord
	var ResultData Result
	var err error
	var total int64 = 0
	var fields []string

	sql := `select vr.*,vs.image_id from virus_scan vs  join virus_record vr on vs.id = vr.virus_scan_id where vs.type = 'ImageVirusScan'  `
	countSql := `select count(vr.id) from virus_scan vs  join virus_record vr on vs.id = vr.virus_scan_id where vs.type = 'ImageVirusScan' `
	filter := ""

	if this.Filename != "" {
		filter = filter + `vr.filename like ? and `
		fields = append(fields, "%"+this.Filename+"%")
	}

	if this.Severity != "" {
		filter = filter + `vr.severity = ? and `
		fields = append(fields, this.Severity)
	}

	if filter != "" {
		sql = sql + "where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql

	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql, fields).QueryRows(&virusLogList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageVirusErr
		logs.Error("Get VirusLogErr List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = virusLogList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *VirusRecord) Add() Result {
	insetSql := `INSERT INTO virus_record(virus_scan_id, filename, virus, database, type, size, owner, permission, modify_time, create_time,severity) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insetSql,
		this.VirusScanId,
		this.Filename,
		this.Virus,
		this.Database,
		this.Type,
		this.Size,
		this.Owner,
		this.Permission,
		this.ModifyTime,
		this.CreateTime, this.Severity).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddVirusLogErr
		logs.Error("Add VirusScan failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *VirusRecord) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var total int64
	var fields []string
	filter := ""

	countSql := `select count(vr.id) from virus_scan vs  join virus_record vr on vs.id = vr.virus_scan_id `

	if this.Type != "" {
		filter = filter + `vs.type = ? and `
		fields = append(fields, this.Type)
	}

	if this.Severity != "" {
		filter = filter + `vr.severity = ? and `
		fields = append(fields, this.Severity)
	}

	if filter != "" {
		countSql = countSql + " where " + filter
	}

	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")

	o.Raw(countSql, fields).QueryRow(&total)
	return total
}
