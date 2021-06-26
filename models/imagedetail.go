package models

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type ImageDetail struct {
	Id            string        `orm:"pk;" description:"(数据库主键)"`
	ImageId       string        `description:"(镜像id)"`
	Name          string        `description:"(镜像名)"`
	HostId        string        `description:"(主机id)"`
	HostName      string        `description:"(主机名称)"`
	RepoTags      string        `description:"(RepoTags)"`
	RepoDigests   string        `description:"(RepoDigests)"`
	Os            string        `description:"(系统)"`
	Size          string        `description:"(大小)"`
	Layers        int           `description:"(Layers)"`
	Dockerfile    string        `description:"(Dockerfile内容)"`
	CreateTime    int64         `description:"(创建时间)"`
	ModifyTime    int64         `description:"(更新时间)"`
	PackagesJson  string        `description:"(软件包列表)"`
	ImageConfigId string        `description:"(镜像Id)"`
	Packages      []PackageInfo `orm:"-" description:"(镜像中软件包列表)"`
}

type PackageInfo struct {
	Name       string   // 软件包名
	Maintainer string   // 软件包发行方
	Licenses   []string // 软件包授权方式
	Type       string   // 软件包类型
	Version    string   // 软件包版本
}

type ImageDetailInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
	Get() *ImageDetail
}

func (this *ImageDetail) Add() Result {
	insetSql := `INSERT INTO ` + utils.ImageDetail + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	packagesJson, _ := json.Marshal(this.Packages)
	if string(packagesJson) != "" {
		this.PackagesJson = string(packagesJson)
	}

	if this.Id == "" {
		uid, _ := uuid.NewV4()
		this.Id = uid.String()
	}
	_, err := o.Raw(insetSql,
		this.Id,
		this.ImageId,
		this.Name,
		this.HostId,
		this.HostName,
		this.RepoTags,
		this.RepoDigests,
		this.Os,
		this.Size,
		this.Layers,
		this.Dockerfile,
		this.CreateTime,
		this.ModifyTime,
		this.PackagesJson, this.ImageConfigId).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageDetailErr
		logs.Error("Add ImageDetail failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageDetail) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var imageDetailList []*ImageDetail = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from ` + utils.ImageDetail + ` `
	countSql := `select count(id) from ` + utils.ImageDetail + ` `
	filter := ""
	var fields []string
	if this.Id != "" {
		filter = filter + `id = ? and `
		fields = append(fields, this.Id)
	}
	if this.Name != "" {
		filter = filter + `name like ? and `
		fields = append(fields, "%"+this.Name+"%")
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
	if this.RepoTags != "" {
		filter = filter + `repo_tags like ? and `
		fields = append(fields, "%"+this.RepoTags+"%")
	}
	if this.Os != "" {
		filter = filter + `contaienr_id = ? and `
		fields = append(fields, this.Os)
	}
	if this.CreateTime != 0 {
		filter = filter + `create_at  > this.CreateTime and `
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
	_, err = o.Raw(resultSql, fields).QueryRows(&imageDetailList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageDetailErr
		logs.Error("Get ImageDetail List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	_ = o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})

	data[Result_Total] = total
	data[Result_Items] = imageDetailList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *ImageDetail) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	removeSQL := "delete from " + utils.ImageDetail + " "
	var fields []string
	filter := ""

	if this.ImageId != "" {
		filter = filter + `image_id = ? and `
		fields = append(fields, this.ImageId)
	}

	if this.ImageConfigId != "" {

		ids := strings.Split(this.ImageConfigId, ",")
		var placeholder []string
		for i := 0; i < len(ids); i++ {
			placeholder = append(placeholder, "?")
		}
		filter = filter + `image_config_id in (` + strings.Join(placeholder, ",") + `) and `
		fields = append(fields, ids...)
	}

	if this.HostId != "" {
		filter = filter + `host_id = ? and `
		fields = append(fields, this.HostId)
	}

	if filter != "" {
		removeSQL = removeSQL + " where " + filter
	}

	removeSQL = strings.TrimSuffix(strings.TrimSpace(removeSQL), "and")
	_, err := o.Raw(removeSQL, fields).Exec()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteImageDetailErr
		logs.Error("Delete ImageDetail failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
func (this *ImageDetail) Get() *ImageDetail {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	object := new(ImageDetail)
	var err error
	var fields []string
	filter := ""
	sql := "select * from " + utils.ImageDetail + " where "

	if this.ImageId != "" {
		filter = filter + `image_id = ? and `
		fields = append(fields, this.ImageId)
	}

	if this.Name != "" {
		filter = filter + `name = ? and `
		fields = append(fields, this.Name)
	}

	if this.ImageConfigId != "" {
		filter = filter + `image_config_id = ? and `
		fields = append(fields, this.ImageConfigId)
	}

	if filter != "" {
		sql = sql + " where " + filter
	}

	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	err = o.Raw(sql, fields).QueryRow(&object)
	if err != nil {
		return nil
	}
	return object
}
