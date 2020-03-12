package statistics

import (
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type StatisticsService struct {
	*models.HostConfig
	*models.ContainerConfig
}

func (this *StatisticsService) GetAssetStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})
	data["ContainerCount"] = 0
	data["HostCount"] = 0

	//主机数
	hostConfig := new(models.HostConfig)
	data["HostCount"] = hostConfig.Count()
	//容器数
	containerConfig := new(models.ContainerConfig)
	data["ContainerCount"] = containerConfig.Count()

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}
