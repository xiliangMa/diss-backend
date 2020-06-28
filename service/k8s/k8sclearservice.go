package k8s

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"os"
)

type K8sClearService struct {
	ClusterList []*models.Cluster
	DropCluster bool
}

// todo 删除任务
func (this *K8sClearService) ClearAll() {
	for _, param := range this.ClusterList {
		msg := ""
		// 集群检查
		cluster := this.Check(param)

		// 检测成功
		if cluster == nil {
			msg = fmt.Sprintf("Clear Cluster: %s fail, Error: Not found cluster", param.Id)
			logs.Error(msg)
		} else {
			msg = fmt.Sprintf("Clear Cluster: %s start......", cluster.Name)
			logs.Info(msg)
			// 更新同步状态未Clearing、并设置为禁止同步
			cluster.IsSync = false
			cluster.SyncStatus = models.Cluster_Sync_Status_Clearing
			if !this.DropCluster {
				cluster.SyncStatus = models.Cluster_Sync_Status_NotSynced
			}
			cluster.Update()

			// 清除Container
			this.ClearContainer(cluster)

			// 清除pod
			this.ClearPod(cluster)

			// 清除ns
			this.ClearNs(cluster)

			// 清除node
			this.ClearNode(cluster)

			// 清除集群
			if this.DropCluster {
				this.ClearCluster(cluster)
			}

			msg = fmt.Sprintf("Clear Cluster: %s sucess......", cluster.Name)
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
	data := result.Data.(map[string]interface{})
	clusetrList := data["items"].([]*models.Cluster)
	total := data["total"]
	if total != 0 {
		return clusetrList[0]
	}
	return nil
}

func (this *K8sClearService) ClearCluster(cluster *models.Cluster) {
	msg := fmt.Sprintf("Clear Cluster, Cluster: %s ", cluster.Name)
	logs.Info(msg)
	if cluster.AuthType == models.Api_Auth_Type_KubeConfig {
		if beego.AppConfig.String("RunMode") == "prod" {
			uploadPath := beego.AppConfig.String("system::UploadPath")
			file := fmt.Sprintf("%+v%+v", uploadPath, cluster.FileName)
			err := os.Remove(file)
			if err != nil {
				logs.Error("Remove kubeconfig fail, file: %s, err：%s", file, err.Error())
			}
		} else {
			os.Remove(cluster.FileName)
		}

	}
	cluster.Delete()
}

func (this *K8sClearService) ClearNs(cluster *models.Cluster) {
	msg := fmt.Sprintf("Clear NameSpace, Cluster: %s ", cluster.Name)
	logs.Info(msg)

	ns := models.NameSpace{}
	ns.ClusterId = cluster.Id
	ns.Delete()
}

func (this *K8sClearService) ClearPod(cluster *models.Cluster) {
	msg := fmt.Sprintf("Clear Pod, Cluster: %s ", cluster.Name)
	logs.Info(msg)

	pod := models.Pod{}
	pod.ClusterName = cluster.Name
	pod.Delete()
}

func (this *K8sClearService) ClearContainer(cluster *models.Cluster) {
	msg := fmt.Sprintf("Clear Container, Cluster: %s ", cluster.Name)
	logs.Info(msg)
	cc := models.ContainerConfig{}
	cc.ClusterName = cluster.Name
	cc.Delete()

	logs.Info(msg)
	ci := models.ContainerInfo{}
	ci.ClusterName = cluster.Name
	ci.Delete()
}

func (this *K8sClearService) ClearNode(cluster *models.Cluster) {
	msg := fmt.Sprintf("Clear node, Cluster %s ", cluster.Name)
	logs.Info(msg)

	hc := models.HostConfig{}
	hc.ClusterId = cluster.Id
	hc.Delete()

	hi := models.HostInfo{}
	hi.ClusterId = cluster.Id
	hi.Delete()
}

// todo clear container and host task
func (this *K8sClearService) ClearTask(cluster *models.Cluster) {
}
