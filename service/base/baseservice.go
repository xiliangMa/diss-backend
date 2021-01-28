package base

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type BaseService struct {
	HostIds      string
	ImageIds     string
	ContainerIds string
	JobId        string
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
