package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type SensitiveInfo struct {
	Id         string `orm:"pk;" description:"(id)"`
	ImageId    string `description:"(镜像ID)"`
	HostId     string `description:"(主机ID)"`
	HostName   string `description:"(主机名称)"`
	Type       string `description:"(资源类型：host image)"`
	FileName   string `description:"(文件名称)"`
	MD5        string `description:"(文件MD5码)"`
	Permission uint32 `description:"(文件权限)"`
	FileType   string `description:"(文件类型)"`
	Size       int64  `description:"(文件大小)"`
	CreateTime int64  `description:"(添加时间)"`
	Files      []FileInfo
}

//虚拟表，用于接收agent的文件列表
type FileInfo struct {
	Name       string // 文件名称
	MD5        string // 文件MD5码
	Permission uint32 // 文件权限
	Type       string // 文件类型
	Size       int64  // 文件大小
}

type SensitiveInfoInterface interface {
	Add() Result
	List(from, limit int) Result
}

func (this *SensitiveInfo) Add() Result {
	insertSql := `INSERT INTO ` + utils.SensitiveInfo + `(file_name, host_id, host_name, image_id, file_type, m_d5, permission, size, create_time) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insertSql,
		this.FileName,
		this.HostId,
		this.HostName,
		this.ImageId,
		this.FileType,
		this.MD5,
		this.Permission,
		this.Size,
		this.CreateTime).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddSensitiveInfoErr
		logs.Error("Add SensitiveInfo failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *SensitiveInfo) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var SensitiveInfoList []*SensitiveInfo = nil
	var ResultData Result
	var err error
	var total int64 = 0

	resType := "host_id"
	if this.Type == Sc_Type_Image {
		resType = "image_id"
	}

	sql := `select * from ` + utils.SensitiveInfo + ` join
(select distinct on (res_time_filter.` + resType + `) * from
(select ` + resType + `, create_time from ` + utils.SensitiveInfo + ` group by create_time,` + resType + ` order by sensitive_info.create_time desc) as res_time_filter) as dist_res on sensitive_info.` + resType + ` = dist_res.` + resType + ` and sensitive_info.create_time = dist_res.create_time `
	countSql := `select count(id) from ` + utils.SensitiveInfo + ` join
(select distinct on (res_time_filter.` + resType + `) * from
(select ` + resType + `, create_time from ` + utils.SensitiveInfo + ` group by create_time,` + resType + ` order by create_time desc) as res_time_filter) as dist_res on sensitive_info.` + resType + ` = dist_res.` + resType + ` and sensitive_info.create_time = dist_res.create_time `
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

	if this.FileName != "" {
		filter = filter + `name like ? and `
		fields = append(fields, "%"+this.FileName+"%")
	}

	if this.FileType != "" {
		filter = filter + `type = ? and `
		fields = append(fields, this.FileType)
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " order by sensitive_info.create_time desc limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql, fields).QueryRows(&SensitiveInfoList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetWarningInfoListErr
		logs.Error("Get WarningInfo list failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = SensitiveInfoList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
