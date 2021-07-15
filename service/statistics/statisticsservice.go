package statistics

import (
	"net/http"

	"github.com/xiliangMa/diss-backend/models"
)

type StatisticsService struct {
	*models.HostConfig
	*models.ContainerConfig
	*models.DcokerIds
}

func (this *StatisticsService) GetAssetStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})

	//主机数
	hostConfig := new(models.HostConfig)
	hc := make(map[string]int64)
	hc["HostCount"] = hostConfig.Count()

	hostConfig.Status = models.Host_Status_Normal
	hc[models.Host_Status_Normal] = hostConfig.Count()

	hostConfig.Status = models.Host_Status_Abnormal
	hc[models.Host_Status_Abnormal] = hostConfig.Count()
	data["HostConfig"] = hc

	//容器数
	containerConfig := new(models.ContainerConfig)
	cc := make(map[string]int64)
	cc["ContainerCount"] = containerConfig.Count()

	containerConfig.Status = models.Container_Status_Created
	cc[models.Container_Status_Created] = containerConfig.Count()

	containerConfig.Status = models.Container_Status_Running
	cc[models.Container_Status_Running] = containerConfig.Count()

	containerConfig.Status = models.Container_Status_Terminated
	cc[models.Container_Status_Terminated] = containerConfig.Count()

	containerConfig.Status = models.Container_Status_Exited
	cc[models.Container_Status_Exited] = containerConfig.Count()
	data["ContainerConfig"] = cc

	//镜像仓库
	registry := new(models.Registry)
	r := make(map[string]int64)
	r["RegistryCount"] = registry.Count()
	data["Registry"] = r

	//镜像
	imageConfig := new(models.ImageConfig)
	ic := make(map[string]int64)

	imageConfig.Type = models.All
	ic["ImageCount"] = imageConfig.Count()

	imageConfig.Type = "host"
	ic["HostImage"] = imageConfig.Count()

	imageConfig.Type = ""
	ic["RegistryImage"] = imageConfig.Count()
	data["ImageConfig"] = ic

	//集群数
	cluster := new(models.Cluster)
	cmap := make(map[string]int64)
	cmap["ClusterCount"] = cluster.Count()

	cluster.Status = models.Cluster_Status_Active
	cmap[models.Cluster_Status_Active] = cluster.Count()

	cluster.Status = models.Cluster_Status_Unavailable
	cmap[models.Cluster_Status_Unavailable] = cluster.Count()
	data["Cluster"] = cmap

	// Pod
	pod := new(models.Pod)
	pmap := make(map[string]int64)
	pmap["PodCount"] = pod.Count()
	pod.Status = models.Container_Status_Pending
	pmap[models.Container_Status_Pending] = pod.Count()

	pod.Status = models.Container_Status_Running
	pmap[models.Container_Status_Running] = pod.Count()

	pod.Status = models.Container_Status_Succeeded
	pmap[models.Container_Status_Succeeded] = pod.Count()

	pod.Status = models.Container_Status_Failed
	pmap[models.Container_Status_Failed] = pod.Count()

	pod.Status = models.Container_Status_Unknown
	pmap[models.Container_Status_Unknown] = pod.Count()
	data["Pod"] = pmap

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetBenchMarkProportionStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})
	//docker 基线 / k8s 基线
	hostConfig := new(models.HostConfig)
	data["DockerBenchmarkCount"], data["K8sBenchmarkCount"] = hostConfig.GetBenchMarkProportion()
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetBnechMarkSummaryStatistics() models.Result {
	bml := new(models.BenchMarkLog)
	return bml.GetMarkSummary()
}

func (this *StatisticsService) GetIntrudeDetectLogStatistics(timeCycle int) models.Result {
	dcokerIds := new(models.DcokerIds)
	return dcokerIds.GetIntrudeDetectLogStatistics(timeCycle)
}

func (this *StatisticsService) GetHostBnechMarkSummaryStatistics(hostId string) models.Result {
	bml := new(models.BenchMarkLog)
	bml.HostId = hostId
	return bml.GetHostMarkSummary()
}
func (this *StatisticsService) GetDissProportionStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})
	// safe / unsafe
	hostConfig := new(models.HostConfig)
	data["Safe"], data["UnSafe"] = hostConfig.GetDissCountProportion()
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetOnlineProportionStatistics() models.Result {
	var ResultData models.Result
	data := make(map[string]interface{})
	// Online / Offline
	hostConfig := new(models.HostConfig)
	data["OnlineCount"], data["OfflineCount"] = hostConfig.GetOnlineProportion()
	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetMirrorRiskStatistics() models.Result {

	var ResultData models.Result
	data := make(map[string]interface{})

	h := make(map[string]int64)
	vu := new(models.Vulnerabilities)

	vu.Severity = models.SEVERITY_Critical
	h[models.SEVERITY_Critical] = vu.Count()

	vu.Severity = models.SEVERITY_High
	h[models.SEVERITY_High] = vu.Count()

	vu.Severity = models.SEVERITY_Medium
	h[models.SEVERITY_Medium] = vu.Count()

	vu.Severity = models.SEVERITY_Low
	h[models.SEVERITY_Low] = vu.Count()

	vu.Severity = models.SEVERITY_Unknown
	h[models.SEVERITY_Unknown] = vu.Count()

	vr := make(map[string]int64)
	virusRecord := new(models.VirusScanRecord)

	virusRecord.Severity = models.SEVERITY_Critical
	vr[models.SEVERITY_Critical] = virusRecord.Count()

	virusRecord.Severity = models.SEVERITY_High
	vr[models.SEVERITY_High] = virusRecord.Count()

	virusRecord.Severity = models.SEVERITY_Medium
	vr[models.SEVERITY_Medium] = virusRecord.Count()

	virusRecord.Severity = models.SEVERITY_Low
	vr[models.SEVERITY_Low] = virusRecord.Count()

	si := make(map[string]int64)
	sensitiveInfo := new(models.SensitiveInfo)

	sensitiveInfo.Severity = models.SEVERITY_Critical
	si[models.SEVERITY_Critical] = sensitiveInfo.Count()

	sensitiveInfo.Severity = models.SEVERITY_High
	si[models.SEVERITY_High] = sensitiveInfo.Count()

	sensitiveInfo.Severity = models.SEVERITY_Medium
	si[models.SEVERITY_Medium] = sensitiveInfo.Count()

	sensitiveInfo.Severity = models.SEVERITY_Low
	si[models.SEVERITY_Low] = sensitiveInfo.Count()

	// 漏洞库
	data["vuln"] = h
	// 病毒库
	data["virus"] = vr
	//敏感信息
	data["sensitive"] = si

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *StatisticsService) GetWarningStatistics() models.Result {

	var ResultData models.Result
	data := make(map[string]interface{})

	wi := make(map[string]int64)
	warningInfo := new(models.WarningInfo)

	warningInfo.Type = models.ReverseShell
	wi[models.ReverseShell] = warningInfo.Count()

	warningInfo.Type = models.ContainerEscalation
	wi[models.ContainerEscalation] = warningInfo.Count()

	warningInfo.Type = models.Shell
	wi[models.Shell] = warningInfo.Count()

	warningInfo.Type = models.PrivilegeEscalation
	wi[models.PrivilegeEscalation] = warningInfo.Count()

	warningInfo.Type = models.VulnExploit
	wi[models.VulnExploit] = warningInfo.Count()

	warningInfo.Type = models.SuspiciousConnection
	wi[models.SuspiciousConnection] = warningInfo.Count()

	warningInfo.Type = models.SuspiciousProcess
	wi[models.SuspiciousProcess] = warningInfo.Count()

	warningInfo.Type = models.NetAttack
	wi[models.NetAttack] = warningInfo.Count()

	warningInfo.Type = models.K8sAudit
	wi[models.K8sAudit] = warningInfo.Count()

	warningInfo.Type = models.TrojanVirus
	wi[models.TrojanVirus] = warningInfo.Count()

	warningInfo.Type = models.SuspiciousLogin
	wi[models.SuspiciousLogin] = warningInfo.Count()

	warningInfo.Type = models.TypeScanning
	wi[models.TypeScanning] = warningInfo.Count()

	warningInfo.Type = models.BruteForce
	wi[models.BruteForce] = warningInfo.Count()

	warningInfo.Type = models.OTHER
	wi[models.OTHER] = warningInfo.Count()

	data["warningInfo"] = wi

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData

}
