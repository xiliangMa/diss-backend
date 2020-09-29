package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type K8sTaskHandler struct {
}

/**
 * 安全检查任务
 */
func (this *K8sTaskHandler) CheckClusterStatusTask() {
	// 获取集群列表
	cluster := new(models.Cluster)
	result := cluster.List(0, 0)
	if result.Code != http.StatusOK {
		logs.Error("Cluster status check fail, err: %s", result.Message)
		return
	}
	if result.Data == nil {
		return
	}

	items := result.Data.(map[string]interface{})[models.Result_Items]
	if items == nil {
		return
	}
	clusterList := items.([]*models.Cluster)

	// 检测链接状态
	//如果链接不可用， 设置集群状态为不可用
	for _, cluster := range clusterList {
		// 获取集群连接client
		client := models.CreateK8sClient(models.BuildApiParams(cluster))
		if client.ErrMessage != "" {
			if cluster.Status != models.Cluster_Status_Unavailable {
				cluster.Status = models.Cluster_Status_Unavailable
				cluster.Update()
			}
			logs.Error("Cluster Unavailable, please check your cluster.")
			if _, ok := models.KCM.ClientHub[cluster.Id]; ok {
				delete(models.KCM.ClientHub, cluster.Id)
			}
			return
		}
		_, err := client.GetNodes()
		if err != nil {
			if cluster.Status != models.Cluster_Status_Unavailable {
				cluster.Status = models.Cluster_Status_Unavailable
				cluster.Update()
			}
			logs.Error("Cluster Unavailable, please check your cluster.")
			if _, ok := models.KCM.ClientHub[cluster.Id]; ok {
				delete(models.KCM.ClientHub, cluster.Id)
			}
			return
		}
		models.KCM.ClientHub[cluster.Id] = client
		if cluster.Status != models.Cluster_Status_Active {
			cluster.Status = models.Cluster_Status_Active
			cluster.Update()
		}

	}

}
