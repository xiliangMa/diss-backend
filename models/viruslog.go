package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

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
	List(from, limit int) Result
}

type DockerVirusInterface interface {
	List(from, limit int) Result
}

/**
 	select * from (
	select a.*, concat_ws(':', concat_ws('/', b."registry", b."repo"), b."tag") as image_name from image_virus as a
	JOIN catalog_image_docker as b ON  a."imageId" = b."imageId" )
	as c  where c."userId" = 'admin' and c."image_name" like '%docker.io/openstack001/av-sample:2%' limit 20 OFFSET 0
*/
func (this *ImageVirus) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var imageVirusList []*ImageVirus = nil
	var ResultData Result
	var err error
	var total int64 = 0

	sql := ` select * from (select a.*, concat_ws(':', concat_ws('/', b."registry", b."repo"), b."tag") as image_name from image_virus as a
	JOIN catalog_image_docker as b ON  a."imageId" = b."imageId" ) as c `
	countSql := `select "count"(c."imageId") from (select a.*, concat_ws(':', concat_ws('/', b."registry", b."repo"), b."tag") as image_name from image_virus as a
	JOIN catalog_image_docker as b ON  a."imageId" = b."imageId" ) as c `
	filter := ""
	if this.ImageDigest != "" {
		filter = filter + `c."imageId" = '` + this.ImageId + "' and "
	}
	if this.UserId != "" {
		filter = filter + `c."userId" = '` + this.UserId + "' and "
	}
	if this.Virus != "" {
		filter = filter + `c."virus" like '%` + this.Virus + "%' and "
	}

	if this.CreatedAt != 0 {
		filter = filter + `c."createdAt" > ` + fmt.Sprintf("%v", this.CreatedAt) + " and "
	}
	if this.ImageName != "" {
		filter = filter + `c."image_name" like '%` + this.ImageName + `%'`
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
	_, err = o.Raw(resultSql).QueryRows(&imageVirusList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetImageVirusErr
		logs.Error("Get ImageVirusErr List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
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

func (this *DockerVirus) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var dockerVirusList []*DockerVirus = nil
	var ResultData Result
	var err error
	var total int64 = 0

	filterSql := ""
	countSql := `select "count"(host_id) from ` + utils.DockerVirus
	sql := "select * from " + utils.DockerVirus

	// 根据 TargeType = host 和 HostId = All 判断是否是查询所有主机日志 如果不是则匹配其它所传入的条件
	// 根据 TargeType = container 和 ContainerId = All 判断是否是查询所有容器日志 如果不是则匹配其它所传入的条件

	if this.TargeType == IDLT_Host && this.HostId == All {
		filterSql = filterSql + "container_id = '" + IDLT_Host + "' and "
	}
	if this.TargeType == IDLT_Docker && this.ContainerId == All {
		filterSql = filterSql + "container_id != '" + IDLT_Host + "' and "
	}

	if (this.ContainerId != "" && this.ContainerId != All) || (this.HostId != "" && this.HostId != All) {
		if this.ContainerId != IDLT_Host && this.TargeType == IDLT_Docker {
			containerId := this.ContainerId
			containerId = string([]byte(this.ContainerId)[:12])
			filterSql = filterSql + "container_id ilike '" + containerId + "' and "
		}

		if this.TargeType == IDLT_Host {
			filterSql = filterSql + "host_id ilike '%" + this.HostId + "%' and "
		}
		if this.Virus != "" {
			filterSql = filterSql + utils.DockerVirus + `."virus" like '%` + this.Virus + "%' and "
		}
	}

	if this.CreatedAt != 0 {
		filterSql = filterSql + utils.ImageVirus + `."createdAt" > ` + fmt.Sprintf("%s", this.CreatedAt) + " and "
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
	_, err = o.Raw(resultSql).QueryRows(&dockerVirusList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetDockerVirusErr
		logs.Error("Get DockerVirus List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dockerVirusList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
