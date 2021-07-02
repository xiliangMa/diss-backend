package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"regexp"
	"time"
)

type ImageVulnerabilities struct {
	Id              int                `orm:"pk;auto" description:"(Id)" description:"(id)"`
	RegistryId      int                `orm:"size(32)" description:"(仓库Id)"`
	HostId          string             `orm:"size(128)" description:"(主机Id)"`
	ImageId         string             `orm:"size(512)" description:"(镜像Id)"`
	TaskId          string             `orm:"size(128)" description:"(任务Id)"`
	Target          string             `orm:"size(128)" description:"(扫描镜像名)"`
	Type            string             `orm:"size(32)" description:"(镜像系统类型)"`
	Vulnerabilities []*Vulnerabilities `orm:"reverse(many);" description:"(漏洞列表)"`
	CreateTime      int64              `description:"(感染文件创建时间)"`
}

type Vulnerabilities struct {
	Id                   int                   `orm:"pk;auto" description:"(Id)" description:"(id)"`
	VulnerabilityID      string                `orm:"column(vulnerability_id)" description:"(漏洞Id)"`
	PkgName              string                `orm:"size(128)" description:"(包名)"`
	InstalledVersion     string                `orm:"size(64)" description:"(安装版本)"`
	FixedVersion         string                `orm:"size(32)" description:"(已解决版本)"`
	SeveritySource       string                `orm:"size(32)" description:"(来源)"`
	PrimaryURL           string                `orm:"size(512)" description:"(漏洞地址)"`
	Title                string                `description:"(漏洞标题)"`
	Description          string                `description:"(漏洞描述)"`
	Severity             string                `orm:"size(32)" description:"(等级)"`
	PublishedDate        int64                 `description:"(发布日期)"`
	LastModifiedDate     int64                 `description:"(最后修改日期)"`
	ImageVulnerabilities *ImageVulnerabilities `orm:"rel(fk);null" description:"(镜像漏洞)"`
}

type ImageVulnerabilitiesInterface interface {
	Add() Result
	List(from, limit int) Result
	Delete() Result
}

type VulnerabilitiesInterface interface {
	Add() Result
}

func (this *ImageVulnerabilities) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	o.Begin()
	this.CreateTime = time.Now().UnixNano()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddImageVulnerabilitiesErr
		logs.Error("Add ImageVulnerabilities failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		o.Rollback()
		return ResultData
	}

	var severityStatus = ""
	re := regexp.MustCompile(`(clam|trivy|koala|syft)`)
	for _, vuln := range this.Vulnerabilities {
		if re.FindString(vuln.PkgName) != "" {
			continue
		}
		if severityStatus == "" && vuln.Severity == SEVERITY_High || vuln.Severity == SEVERITY_Critical {
			severityStatus = "Trustee"
		}
		vuln.ImageVulnerabilities = this
		vuln.Add()
	}
	o.Commit()
	task := Task{}
	task.Id = this.TaskId
	t := task.Get()
	if t != nil {
		if severityStatus != "" {
			t.SecurityStatus = "Trustee"
		} else {
			t.SecurityStatus = "NotTrustee"
		}
		t.Update()
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ImageVulnerabilities) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	imageVulnerabilities := new(ImageVulnerabilities)
	var vulnerabilities []*Vulnerabilities
	var ResultData Result
	cond := orm.NewCondition()
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}
	err := o.QueryTable(utils.ImageVulnerabilities).SetCond(cond).One(imageVulnerabilities)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageVulnerabilitiesErr
		logs.Error("Get ImageVulnerabilitiesErr List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	vuln := orm.NewCondition()
	vuln = vuln.And("image_vulnerabilities_id", imageVulnerabilities.Id)
	_, err = o.QueryTable(utils.Vulnerabilities).SetCond(vuln).Limit(limit, from).All(&vulnerabilities)
	total, _ := o.QueryTable(utils.Vulnerabilities).SetCond(vuln).Count()
	imageVulnerabilities.Vulnerabilities = vulnerabilities

	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = imageVulnerabilities

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ImageVulnerabilities) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.ImageId != "" {
		cond = cond.And("image_id", this.ImageId)
	}
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	_, err := o.QueryTable(utils.ImageVulnerabilities).SetCond(cond).Delete()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteImageVulnerabilitiesErr
		logs.Error("Delete Vulnerabilities failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}

func (this *Vulnerabilities) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddVulnerabilitiesErr
		logs.Error("Add Vulnerabilities failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Vulnerabilities) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var vulnerabilities []*Vulnerabilities
	var ResultData Result
	var err error
	cond := orm.NewCondition()

	if this.VulnerabilityID != "" {
		cond = cond.And("vulnerability_id", this.VulnerabilityID)
	}
	if this.PkgName != "" {
		cond = cond.And("pkg_name", this.PkgName)
	}
	if this.Severity != "" {
		cond = cond.And("severity", this.Severity)
	}
	if this.VulnerabilityID != "" {
		cond = cond.And("vulnerability_id", this.VulnerabilityID)
	}

	_, err = o.QueryTable(utils.Vulnerabilities).SetCond(cond).Limit(limit, from).All(&vulnerabilities)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageVulnerabilitiesErr
		logs.Error("Get Vulnerabilities List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	for _, v := range vulnerabilities {
		o.LoadRelated(v, "ImageVulnerabilities", 1, 1, 0)
	}

	total, _ := o.QueryTable(utils.Vulnerabilities).SetCond(cond).Count()
	data := make(map[string]interface{})
	data[Result_Total] = total
	data[Result_Items] = vulnerabilities

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *Vulnerabilities) Count() int64 {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	cond := orm.NewCondition()
	if this.Severity != "" {
		cond = cond.And("severity", this.Severity)
	}
	count, _ := o.QueryTable(utils.Vulnerabilities).SetCond(cond).Count()
	return count
}
