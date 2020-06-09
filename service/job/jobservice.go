package job

import (
	"github.com/xiliangMa/diss-backend/models"
)

type JobService struct {
	JobParm *models.Job
}

func (this *JobService) GetCheckList() []*models.SecurityCheck {
	joblist := this.JobParm.Internal_Get(this.JobParm.Id)
	secCheckList := []*models.SecurityCheck{}
	if len(joblist) > 0 {
		job := joblist[0]

		// 通过job内容 构建checklist
		targetType := job.SystemTemplate.Type //目标类型为检查模板的类型
		resType := models.SC_Type_Host        // 资源类型包括 主机 容器 镜像
		if (targetType == models.TMP_Type_IDS_Docker) || (targetType == models.TMP_Type_DockerVS) {
			resType = models.Sc_Type_Container
		}
		if (resType == models.SC_Type_Host) && (job.HostConfig != nil) {
			for _, checkHost := range job.HostConfig {
				seccheck := new(models.SecurityCheck)
				seccheck.Type = resType
				switch targetType {
				case models.TMP_Type_BM_Docker, models.TMP_Type_BM_K8S:
					seccheck.BenchMarkCheck = true
				case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
					seccheck.VirusScan = true
				case models.TMP_Type_LS:
					seccheck.LeakScan = true
				}
				seccheck.Host = checkHost
				seccheck.CronType = job.Type
				secCheckList = append(secCheckList, seccheck)
			}
		} else if (resType == models.Sc_Type_Container) && (job.ContainerConfig != nil) {
			for _, checkContainer := range job.ContainerConfig {
				seccheck := new(models.SecurityCheck)
				seccheck.Type = resType
				switch targetType {
				case models.TMP_Type_BM_Docker, models.TMP_Type_BM_K8S:
					seccheck.BenchMarkCheck = true
				case models.TMP_Type_DockerVS, models.TMP_Type_HostVS:
					seccheck.VirusScan = true
				case models.TMP_Type_LS:
					seccheck.LeakScan = true
				}
				seccheck.Container = checkContainer
				seccheck.CronType = job.Type
				secCheckList = append(secCheckList, seccheck)
			}
		}
	}

	return secCheckList
}
