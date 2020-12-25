package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type HostPackage struct {
	Id     string `orm:"pk;size(128)" description:"(id)"`
	Name   string `orm:"" description:"(包名)"`
	Type   string `orm:"size(32)" description:"(java、rpm、dpkg、jar、system)"`
	HostId string `orm:"size(128)" description:"(主机id)"`
}

type HostPackageList struct {
	List []*HostPackage
}

type HostPackageInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
	GetPackageCountByType() Result
}

func (this *HostPackage) Add() Result {
	insetSql := `INSERT INTO host_package VALUES(?, ?, ?, ?)`
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result

	_, err := o.Raw(insetSql,
		this.Id,
		this.Name,
		this.Type,
		this.HostId).Exec()

	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddHostPackageErr
		logs.Error("Add AHostPackage failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *HostPackage) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var ResultData Result
	deleteSql := `DELETE FROM host_package WHERE host_id = ?`

	_, err := o.Raw(deleteSql, this.HostId).Exec()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostPackageErr
		logs.Error("Delete HostPackage failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *HostPackage) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var hostPackageList []*HostPackage
	var ResultData Result
	var err error
	var total = 0

	sql := `SELECT * FROM ` + utils.HostPackage
	countSql := `SELECT "count"(host_id) FROM ` + utils.HostPackage
	filter := ""

	if this.Type != "" {
		filter = filter + ` type = '` + this.Type + `' and `
	}
	if this.HostId != "" {
		filter = filter + `host_id = '` + this.HostId + `' and `
	}
	if this.Name != "" {
		filter = filter + `name like '%` + this.Name + `%' and `
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

	_, err = o.Raw(resultSql).QueryRows(&hostPackageList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPackageErr
		logs.Error("Get HostPackage List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = hostPackageList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *HostPackage) GetPackageCountByType() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var packageCount []orm.Params
	ResultData := Result{Code: http.StatusOK}
	var err error
	var total int64 = 0

	sql := `SELECT count(id), type from ` + utils.HostPackage
	countSql := `SELECT count(id) from ` + utils.HostPackage
	groupBySql := " GROUP BY type"
	filter := " where "

	if this.HostId != "" {
		filter = filter + ` host_id = '` + this.HostId + `' `
		sql = sql + filter
		countSql = countSql + filter
	}
	sql = sql + groupBySql
	_, err = o.Raw(sql).Values(&packageCount)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostPackageErr
		logs.Error("Get HostPackage count group type failed, code: %d, err: %s", utils.GetHostPackageErr, err.Error())
		return ResultData
	}
	o.Raw(countSql).QueryRow(&total)

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = packageCount

	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
