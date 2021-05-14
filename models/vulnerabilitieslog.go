package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type ImagePackageVulnerabilities struct {
	PkgUserId                  string `description:"(用户Id)"`
	PkgImageId                 string `description:"(镜像id)"`
	PkgImageName               string `description:"(镜像名)"`
	PkgName                    string `description:"(名称)"`
	PkgVersion                 string `description:"(版本)"`
	PkgPath                    string `description:"(Path)"`
	PkgType                    string `description:"(类型)"`
	PkgArch                    string `description:"(架构)"`
	VulnerabilityId            string `description:"(漏铜Id)"`
	VulnerabilityNamespaceName string `description:"(操作系统)"`
	CreatedAt                  int64  `description:"(创建时间)"`
	Severity                   string `description:"(安全等级)"`
	Link                       string `description:"(漏洞详情)"`
}

type FeedDataVulnerabilities struct {
	Id            string `description:"(Id)"`
	NamespaceName string `description:"(命名空间名称)"`
	Severity      string `description:"(安全等级)"`
	Description   string `description:"(描述)"`
	Link          string `description:"(漏洞详情)"`
	MetadataJson  string `description:"(元数据)"`
	Cvss2Vectors  string `description:"(向量)"`
	Cvss2Score    string `description:"(分数)"`
	CreatedAt     int64  `description:"(创建时间)"`
	UpdatedAt     int64  `description:"(更新时间)"`
}

type ImagePackageVulnerabilitiesInterface interface {
	List(from, limit int) Result
}

/**
   	select * from (
	select a.*, b."link", b."severity", concat_ws(':', concat_ws('/', c."registry", c."repo"), c."tag") as pkg_image_name  from image_package_vulnerabilities as a
	join feed_data_vulnerabilities as b on a.vulnerability_id = b."id" and a.vulnerability_namespace_name = b."namespace_name"
	join catalog_image_docker as c on a.pkg_image_id = c."imageId"
	) as d
	where d."severity" = 'Low' and d."pkg_user_id" = 'admin' and d."vulnerability_id" = 'CVE-2016-2781' and d."pkg_image_name" like '%docker.io/mysql:8.0.17%' limit 20 OFFSET 0
*/
func (this *ImagePackageVulnerabilities) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var ImagePackageVulnerabilitiesList []*ImagePackageVulnerabilities = nil
	var ResultData Result
	var err error
	var total int64 = 0
	filter := ""
	fields := []string{}
	countSql := `select "count"(d.pkg_image_id) from (select a.*, b."link", b."severity", concat_ws(':', concat_ws('/', c."registry", c."repo"), c."tag") as pkg_image_name  from ` + utils.ImagePackageVulnerabilities +
		` as a join ` + utils.FeedDataVulnerabilities +
		` as b on a.vulnerability_id = b."id" and a.vulnerability_namespace_name = b."namespace_name" ` +
		` join catalog_image_docker as c on a.pkg_image_id = c."imageId") as d `
	sql := `select * from (select a.*, b."link", b."severity", concat_ws(':', concat_ws('/', c."registry", c."repo"), c."tag") as pkg_image_name  from ` + utils.ImagePackageVulnerabilities +
		` as a join ` + utils.FeedDataVulnerabilities +
		` as b on a.vulnerability_id = b."id" and a.vulnerability_namespace_name = b."namespace_name" ` +
		` join catalog_image_docker as c on a.pkg_image_id = c."imageId") as d `

	if this.Severity != "" {
		filter = filter + `d."severity" = ? and `
		fields = append(fields, this.Severity)
	}

	if this.PkgUserId != "" {
		filter = filter + `d."pkg_user_id" = ? and `
		fields = append(fields, this.PkgUserId)
	}
	if this.PkgImageId != "" {
		filter = filter + `d."pkg_image_id" = ? and `
		fields = append(fields, this.PkgImageId)
	}
	if this.VulnerabilityId != "" {
		filter = filter + `d."vulnerability_id" like ? and `
		fields = append(fields, "%"+this.VulnerabilityId+"%")
	}
	if this.PkgImageName != "" {
		filter = filter + `d."pkg_image_name" like ? `
		fields = append(fields, "%"+this.PkgImageName+"%")
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
	_, err = o.Raw(resultSql, fields).QueryRows(&ImagePackageVulnerabilitiesList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImagePackageVulnerabilitiesListErr
		logs.Error("Get ImagePackageVulnerabilities List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	o.Raw(countSql, fields).QueryRow(&total)

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
