package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ImagePackageVulnerabilities struct {
	PkgUserId                  string    `description:"(用户Id)"`
	PkgImageId                 string    `description:"(镜像id)"`
	PkgName                    string    `description:"(名称)"`
	PkgVersion                 string    `description:"(版本)"`
	PkgType                    string    `description:"(类型)"`
	PkgArch                    string    `description:"(架构)"`
	VulnerabilityId            string    `description:"(漏铜Id)"`
	VulnerabilityNamespaceName string    `description:"(操作系统)"`
	CreatedAt                  time.Time `description:"(创建时间)"`
	Severity                   string    `description:"(安全等级)"`
	Link                       string    `description:"(漏洞详情)"`
}

type FeedDataVulnerabilities struct {
	Id            string    `description:"(Id)"`
	NamespaceName string    `description:"(命名空间名称)"`
	Severity      string    `description:"(安全等级)"`
	Description   string    `description:"(描述)"`
	Link          string    `description:"(漏洞详情)"`
	MetadataJson  string    `description:"(元数据)"`
	Cvss2Vectors  string    `description:"(向量)"`
	Cvss2Score    string    `description:"(分数)"`
	CreatedAt     time.Time `description:"(创建时间)"`
	UpdatedAt     time.Time `description:"(更新时间)"`
}

type ImagePackageVulnerabilitiesInterface interface {
	List(from, limit int) Result
}

func (this *ImagePackageVulnerabilities) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var ImagePackageVulnerabilitiesList []*ImagePackageVulnerabilities = nil
	var ResultData Result
	var err error
	var total int64 = 0
	filter := ""

	countSql := "select " + `"count"(a.pkg_image_id)` + " from " + utils.ImagePackageVulnerabilities + ` as a join ` + utils.FeedDataVulnerabilities +
		` as b on a.vulnerability_id = b."id" and a.vulnerability_namespace_name = b."namespace_name" `
	sql := "select * from " + utils.ImagePackageVulnerabilities + ` as a join ` + utils.FeedDataVulnerabilities +
		` as b on a.vulnerability_id = b."id" and a.vulnerability_namespace_name = b."namespace_name" `

	if this.Severity != "" {
		sql = sql + `and b."severity" = '` + this.Severity + `'`
		countSql = countSql + `and b."severity" = '` + this.Severity + `'`
	}

	if this.PkgUserId != "" {
		filter = filter + `pkg_user_id = '` + this.PkgUserId + "' and "
	}
	if this.PkgImageId != "" {
		filter = filter + `pkg_image_id = '` + this.PkgImageId + "' and "
	}
	if this.VulnerabilityId != "" {
		filter = filter + `vulnerability_id = '` + this.VulnerabilityId + "' and "
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}

	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql).QueryRows(&ImagePackageVulnerabilitiesList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImagePackageVulnerabilitiesListErr
		logs.Error("Get ImagePackageVulnerabilities List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	o.Raw(countSql).QueryRow(&total)

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ImagePackageVulnerabilitiesList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
