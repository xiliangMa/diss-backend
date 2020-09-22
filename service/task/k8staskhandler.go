package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"net/http"
	"time"
)

type K8sTaskHandler struct {
	K8sSyncService *k8s.K8sSyncService
}

/***
 ** 已经废弃
 */
func (this *K8sTaskHandler) SyncAll() {
	var cluster models.Cluster
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*models.Cluster)
		// 清除集群数据
		k8sClearService := k8s.K8sClearService{ClusterList: clusterList, DropCluster: false}
		k8sClearService.ClearAll()

		// 同步集群数据
		syncCheckPoint := time.Now().Unix()
		for _, c := range clusterList {
			k8sSyncService := k8s.NewK8sSyncService(syncCheckPoint, c)
			go k8sSyncService.Sync()
		}
	}
}

func (this *K8sTaskHandler) SyncRequiredCluster() {
	var cluster models.Cluster
	result := cluster.GetRequiredSyncList()

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*models.Cluster)
		// 同步集群数据
		syncCheckPoint := time.Now().Unix()
		for _, c := range clusterList {
			k8sSyncService := k8s.NewK8sSyncService(syncCheckPoint, c)
			go k8sSyncService.Sync()
		}
	} else {
		logs.Info("No cluster require sync.")
	}
}
