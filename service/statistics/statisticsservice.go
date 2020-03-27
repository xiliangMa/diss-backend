package statistics

import (
	"github.com/xiliangMa/diss-backend/models"
	msecuritylog "github.com/xiliangMa/diss-backend/models/securitylog"
	"net/http"
)

type StatisticsService struct {
	*models.HostConfig
	*models.ContainerConfig
	*msecuritylog.DcokerIds
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

func (this *StatisticsService) GetBnechMarkProportionStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})
	//docker 基线 / k8s 基线
	hostConfig := new(models.HostConfig)
	data["DockerBenchmarkCount"], data["K8sBenchmarkCount"] = hostConfig.GetBnechMarkProportion()
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetBnechMarkSummaryStatistics() models.Result {
	bml := new(msecuritylog.BenchMarkLog)
	return bml.GetMarkSummary()
}

func (this *StatisticsService) GetIntrudeDetectLogStatistics(timeCycle int) models.Result {
	dcokerIds := new(msecuritylog.DcokerIds)
	return dcokerIds.GetIntrudeDetectLogStatistics(timeCycle)
}

func (this *StatisticsService) GetHostBnechMarkSummaryStatistics(hostId string) models.Result {
	bml := new(msecuritylog.BenchMarkLog)
	bml.HostId = hostId
	return bml.GetHostMarkSummary()
}
