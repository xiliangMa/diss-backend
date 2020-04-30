package securitylog

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type ImageVirus struct {
	ImageDigest string `orm:"column(imageDigest)" description:"(镜像digest)"`
	UserId      string `orm:"column(userId)" description:"(用户id)"`
	FileName    string `description:"(文件名)"`
	Virus       string `description:"(病毒)"`
	FileHash    string `description:"(文件hash)"`
	FileSize    int64  `description:"(文件大小)"`
	CreatedAt   int64  `description:"(创建时间)"`
	LastUpdated int64  `description:"(更新时间)"`
}

type DockerVirus struct {
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

type ImageVirusInterface interface {
	List(from, limit int) models.Result
}

type DockerVirusInterface interface {
	List(from, limit int) models.Result
}

func (this *ImageVirus) List(from, limit int) models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var imageVirusList []*ImageVirus = nil
	var tempList []*ImageVirus = nil
	var ResultData models.Result
	var err error
	var total int64 = 0

	sql := "select * from " + utils.ImageVirus
	filter := ""
	if this.ImageDigest != "" {
		filter = filter + utils.ImageVirus + `."imageDigest" = '` + this.ImageDigest + "' and "
	}
	if this.UserId != "" {
		filter = filter + utils.ImageVirus + `."userId" = '` + this.UserId + "' and "
	}
	if this.Virus != "" {
		filter = filter + utils.ImageVirus + `."virus" like '%` + this.Virus + "%' and "
	}
	if filter != "" {
		sql = sql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql).QueryRows(&imageVirusList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageVirusErr
		logs.Error("Get ImageVirusErr List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ = o.Raw(sql).QueryRows(&tempList)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageVirusList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *DockerVirus) List(from, limit int) models.Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var imageVirusList []*ImageVirus = nil
	var ResultData models.Result
	var err error
	var total int64 = 0

	filterSql := ""
	countSql := `select "count"(host_id) from ` + utils.DockerVirus
	sql := "select * from " + utils.DockerVirus

	// 根据 TargeType = host 和 HostId = All 判断是否是查询所有主机日志 如果不是则匹配其它所传入的条件
	// 根据 TargeType = container 和 ContainerId = All 判断是否是查询所有容器日志 如果不是则匹配其它所传入的条件

	if this.TargeType == models.IDLT_Host && this.HostId == models.All {
		filterSql = filterSql + "container_id = '" + models.IDLT_Host + "' and "
	}
	if this.TargeType == models.IDLT_Docker && this.ContainerId == models.All {
		filterSql = filterSql + "container_id != '" + models.IDLT_Host + "' and "
	}

	if (this.ContainerId != "" && this.ContainerId != models.All) || (this.HostId != "" && this.HostId != models.All) {
		if this.ContainerId != models.IDLT_Host && this.TargeType == models.IDLT_Docker {
			containerId := this.ContainerId
			containerId = string([]byte(this.ContainerId)[:12])
			filterSql = filterSql + "container_id = '" + containerId + "' and "
		}

		if this.TargeType == models.IDLT_Host {
			filterSql = filterSql + "host_id = '" + this.HostId + "' and "
		}
		if this.Virus != "" {
			filterSql = filterSql + utils.DockerVirus + `."virus" like '%` + this.Virus + "%' and "
		}
	}
	if filterSql != "" {
		sql = sql + ` where ` + filterSql
		countSql = countSql + ` where ` + filterSql
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}
	_, err = o.Raw(resultSql).QueryRows(&imageVirusList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDockerVirusErr
		logs.Error("Get DockerVirus List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = imageVirusList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
