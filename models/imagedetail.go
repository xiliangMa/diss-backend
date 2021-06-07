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
	Id           string        `orm:"pk;" description:"(数据库主键)"`
	ImageId      string        `description:"(镜像id)"`
	Name         string        `description:"(镜像名)"`
	HostId       string        `description:"(主机id)"`
	HostName     string        `description:"(主机名称)"`
	RepoTags     string        `description:"(RepoTags)"`
	RepoDigests  string        `description:"(RepoDigests)"`
	Os           string        `description:"(系统)"`
	Size         string        `description:"(大小)"`
	Layers       int           `description:"(Layers)"`
	Dockerfile   string        `description:"(Dockerfile内容)"`
	CreateTime   int64         `description:"(创建时间)"`
	ModifyTime   int64         `description:"(更新时间)"`
	PackagesJson string        `description:"(软件包列表)"`
	Packages     []PackageInfo `orm:"-" description:"(镜像中软件包列表)"`
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
	insetSql := `INSERT INTO ` + utils.ImageDetail + ` VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
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
		this.PackagesJson).Exec()

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
	fields := []string{}
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

	o.Raw(countSql, fields).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageDetailList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ImageDetail) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	removeSQL := "Delete from " + utils.ImageDetail + " "
	cond := ""
	pre := ""

	if this.HostId != "" {
		cond = cond + "host_id = '" + this.HostId + "' "
	}
	if this.ImageId != "" {
		if cond != "" {
			pre = " And "
		}
		cond = cond + pre + " image_id = '" + this.ImageId + "' "
	}

	if cond != "" {
		cond = " Where " + cond
	}
	removeSQL = removeSQL + cond

	_, err := o.Raw(removeSQL).Exec()

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
	sql := "select * from " + utils.ImageDetail + " where "
	cond := ""
	if this.ImageId != "" {
		cond = cond + " image_id = '" + this.ImageId + "' "
	}

	if cond != "" {
		sql += " and " + cond
	}

	if this.Name != "" {
		cond = cond + " name = '" + this.Name + "' "
	}

	sql += cond

	err = o.Raw(sql).QueryRow(&object)
	if err != nil {
		return nil
	}
	return object
}
