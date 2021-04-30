package base

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type BaseService struct {
	HostIds           string
	ImageIds          string
	ContainerIds      string
	ClusterIds        string
	JobId             string
	HostListInCluster []*models.HostConfig
}

func (this *BaseService) CheckJobIsExist() (models.Result, *models.Job) {
	ResultData := models.Result{Code: http.StatusOK}
	if this.JobId == "" {
		ResultData.Code = utils.GetJobErr
		ResultData.Message = "Job id is null."
		logs.Error("Check failed, Err: %s.", ResultData.Message)
		return ResultData, nil
	}
	object := models.Job{Id: this.JobId}
	data := object.Get()
	if data == nil {
		ResultData.Code = utils.GetJobErr
		ResultData.Message = "JobNotFoundErr"
		logs.Error("Check failed, Err: %s.", ResultData.Message)
		return ResultData, nil
	}
	return ResultData, data
}

func (this *BaseService) CheckHostLicense() (models.Result, []*models.HostConfig) {
	hostList := []*models.HostConfig{}
	ResultData := models.Result{Code: http.StatusOK}

	if len(this.HostListInCluster) < 0 {
		// 标准主机
		if this.HostIds == "" {
			logs.Warn("Check failed, Err: HostIds is null.")
			return ResultData, nil
		}

		for _, hostId := range strings.Split(this.HostIds, ",") {
			object := new(models.HostConfig)
			object.Id = hostId
			data := object.Get()
			if data == nil {
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil
			}
			if !data.IsLicensed {
				ResultData.Code = utils.LicenseHostErr
				ResultData.Message = "LicenseHostErr"
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil
			}
			hostList = append(hostList, data)
		}
	} else {
		// 集群主机
		for _, object := range this.HostListInCluster {
			data := object.Get()
			if data == nil {
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil
			}
			if !data.IsLicensed {
				ResultData.Code = utils.LicenseHostErr
				ResultData.Message = "LicenseHostErr"
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil
			}
			hostList = append(hostList, data)
		}
	}

	return ResultData, hostList
}

func (this *BaseService) CheckImageIsExist() (models.Result, []*models.ImageConfig) {
	list := []*models.ImageConfig{}
	ResultData := models.Result{Code: http.StatusOK}
	if this.ImageIds == "" {
		logs.Warn("Check failed, Err: ImageIds is null.")
		return ResultData, nil
	}

	for _, id := range strings.Split(this.ImageIds, ",") {
		object := new(models.ImageConfig)
		if id != "" {
			object.Id = id
			data := object.Get()
			if data == nil {
				ResultData.Code = utils.GetImageConfigErr
				ResultData.Message = "GetImageConfigErr"
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil
			}
			list = append(list, data)
		}

	}
	return ResultData, list
}

func (this *BaseService) CheckContainerIsExist() (models.Result, []*models.ContainerConfig) {
	list := []*models.ContainerConfig{}
	ResultData := models.Result{Code: http.StatusOK}
	if this.ContainerIds == "" {
		logs.Warn("Check failed, Err: ContainerIds is null.")
		return ResultData, nil
	}

	for _, id := range strings.Split(this.ContainerIds, ",") {
		object := new(models.ContainerConfig)
		object.Id = id
		data := object.Get()
		if data == nil {
			ResultData.Code = utils.GetContainerConfigErr
			ResultData.Message = "GetContainerConfigErr"
			logs.Error("Check failed, Err: %s.", ResultData.Message)
			return ResultData, nil
		}
		list = append(list, data)
	}
	return ResultData, list
}

func (this *BaseService) CheckClusterIsExist() (models.Result, []*models.Cluster, []*models.HostConfig) {
	clusterList := []*models.Cluster{}
	hostList := []*models.HostConfig{}
	ResultData := models.Result{Code: http.StatusOK}
	if this.ClusterIds == "" {
		logs.Warn("Check failed, Err: ClusterIds is null.")
		return ResultData, nil, nil
	}

	for _, id := range strings.Split(this.ClusterIds, ",") {
		object := new(models.Cluster)
		if id != "" {
			object.Id = id
			data := object.Get()
			if data == nil {
				ResultData.Code = utils.GetClusterErr
				ResultData.Message = "GetClusterErr"
				logs.Error("Check failed, Err: %s.", ResultData.Message)
				return ResultData, nil, nil
			}
			clusterList = append(clusterList, data)
			// 获取集群主机
			host := models.HostConfig{ClusterId: data.Id}
			_, _, hostData := host.BaseList(0, 0)
			hostList = append(hostList, hostData...)
		}

	}
	return ResultData, clusterList, hostList
}
