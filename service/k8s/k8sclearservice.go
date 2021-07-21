package k8s

import (
	"fmt"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	sssystem "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
)

type K8sClearService struct {
	ClusterList    []*models.Cluster
	DropCluster    bool
	CurrentCluster *models.Cluster
}

// todo 删除任务
func (this *K8sClearService) ClearAll() {
	for _, param := range this.ClusterList {
		msg := ""
		// 集群检查
		this.CurrentCluster = this.Check(param)

		// 检测成功
		if this.CurrentCluster == nil {
			msg = fmt.Sprintf("Clear ClusterOBJ: %s fail, Error: Not found cluster", param.Id)
			logs.Error(msg)
		} else {
			msg = fmt.Sprintf("Clear ClusterOBJ: %s start......", this.CurrentCluster.Name)
			logs.Info(msg)
			// 更新同步状态未Clearing、并设置为禁止同步
			this.CurrentCluster.IsSync = false
			this.CurrentCluster.SyncStatus = models.Cluster_Sync_Status_Clearing
			if !this.DropCluster {
				this.CurrentCluster.SyncStatus = models.Cluster_Sync_Status_NotSynced
			}
			this.CurrentCluster.Update()
			// 清除Container
			this.ClearContainer()

			// 清除pod
			this.ClearPod()

			// 清除service
			this.ClearService()

			// 清除deployment
			this.ClearDeployment()

			// 清除ns
			this.ClearNs()

			// 清除集群
			if this.DropCluster {
				this.ClearCluster()
			}

			// 清除node // 清除networkpolicy
			if this.DropCluster {
				this.ClearNode()
				this.ClearNetworkPolicy()
			}

			msg = fmt.Sprintf("Clear ClusterOBJ: %s sucess......", this.CurrentCluster.Name)
			logs.Info(msg)
		}
	}
}

func (this *K8sClearService) Check(cluster *models.Cluster) *models.Cluster {
	//检测集群是否存在,并返回数据库中的集群对象
	result := cluster.List(0, 0)
	if result.Data == nil {
		return nil
	}
	clusterList := result.Data.(map[string]interface{})["items"].([]*models.Cluster)
	if len(clusterList) > 0 {
		return clusterList[0]
	}
	return nil
}

func (this *K8sClearService) ClearCluster() {
	msg := fmt.Sprintf("Clear ClusterOBJ, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)
	if this.CurrentCluster.AuthType == models.Api_Auth_Type_KubeConfig {
		if beego.AppConfig.String("RunMode") == "prod" {
			uploadPath := beego.AppConfig.String("system::UploadPath")
			file := fmt.Sprintf("%+v%+v", uploadPath, this.CurrentCluster.FileName)
			err := os.Remove(file)
			if err != nil {
				logs.Error("Remove kubeconfig fail, file: %s, err：%s", file, err.Error())
			}
		} else {
			os.Remove(this.CurrentCluster.FileName)
		}

	}
	this.CurrentCluster.Delete()
}

func (this *K8sClearService) ClearNs() {
	watchType := this.CurrentCluster.Id + `_` + utils.NameSpace
	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(NameSpaceService).Close <- true
	}
	msg := fmt.Sprintf("Clear NameSpace, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	ns := models.NameSpace{}
	ns.ClusterId = this.CurrentCluster.Id
	ns.Delete()
}

func (this *K8sClearService) ClearNetworkPolicy() {
	watchType := this.CurrentCluster.Id + `_` + utils.NetworkPolicy

	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(NetworkPolicyService).Close <- true
	}
	msg := fmt.Sprintf("Clear NetworkPolicy, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	netpol := models.NetworkPolicy{}
	netpol.ClusterName = this.CurrentCluster.Name
	netpol.Delete()
}

func (this *K8sClearService) ClearDeployment() {
	watchType := this.CurrentCluster.Id + `_` + utils.Deployment

	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(DeploymentService).Close <- true
	}
	msg := fmt.Sprintf("Clear Deployment, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	deploy := models.Deployment{}
	deploy.ClusterName = this.CurrentCluster.Name
	deploy.Delete()
}

func (this *K8sClearService) ClearService() {
	watchType := this.CurrentCluster.Id + `_` + utils.Service

	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(SVCService).Close <- true
	}
	msg := fmt.Sprintf("Clear Service, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	svc := models.Service{}
	svc.ClusterId = this.CurrentCluster.Id
	svc.Delete()
}

func (this *K8sClearService) ClearPod() {
	watchType := this.CurrentCluster.Id + `_` + utils.Pod

	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(PodService).Close <- true
	}

	msg := fmt.Sprintf("Clear Pod, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	pod := models.Pod{}
	pod.ClusterName = this.CurrentCluster.Name
	pod.Delete()
}

func (this *K8sClearService) ClearContainer() {
	msg := fmt.Sprintf("Clear Container, ClusterOBJ: %s ", this.CurrentCluster.Name)
	logs.Info(msg)
	cc := models.ContainerConfig{}
	cc.ClusterName = this.CurrentCluster.Name
	cc.Delete()

	logs.Info(msg)
	ci := models.ContainerInfo{}
	ci.ClusterName = this.CurrentCluster.Name
	ci.Delete()
}

func (this *K8sClearService) ClearNode() {
	watchType := this.CurrentCluster.Id + `_` + utils.Host

	if models.GRM != nil && models.GRM.GoRoutineMap != nil && models.GRM.GoRoutineMap[watchType] != nil {
		models.GRM.GoRoutineMap[watchType].(NodeService).Close <- true
	}
	msg := fmt.Sprintf("Clear node, ClusterOBJ %s ", this.CurrentCluster.Name)
	logs.Info(msg)

	hc := models.HostConfig{}
	hc.ClusterId = this.CurrentCluster.Id
	result := hc.Delete()

	hi := models.HostInfo{}
	hi.ClusterId = this.CurrentCluster.Id
	hi.Delete()

	if result.Code == http.StatusOK {
		// 更新主机（基线）授权 删除集群并清除主机时，授权恢复
		licenseService := sssystem.LicenseService{}
		licenseService.GetLicensedHostCount()
	}

}

// todo clear container and host task
func (this *K8sClearService) ClearTask(cluster *models.Cluster) {
}
