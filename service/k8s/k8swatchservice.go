package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type K8sWatchService struct {
	Cluster *models.Cluster
}

/**
 * watch 集群下所有的资源
 * namespace、pod、node
 */
func (this *K8sWatchService) WatchAll() {
	var cluster models.Cluster
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*models.Cluster)

		// 清除集群
		ClearCluster(clusterList)

		// watch 集群
		for _, c := range clusterList {
			this.Cluster = c
			this.WatchCluster()
		}
	}
}

func (this *K8sWatchService) WatchCluster() models.Result {
	var resultData models.Result
	clusterName := this.Cluster.Name
	clusterId := this.Cluster.Id
	clientGo := this.CreateK8sClient()

	if clientGo.ErrMessage == "" {
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_InProcess
		this.Cluster.Update()
		logs.Info("Cluster:  %s, Watch start.", clusterName)

		// 异常捕获更新状态
		defer func() {
			if err := recover(); err != nil {
				// 更新集群的同步状态
				this.Cluster.SyncStatus = models.Cluster_Watch_Status_Fail
				this.Cluster.Update()
				logs.Error("Cluster: %s, id: %s , Watch fail. err: %s", clusterName, clusterId, err)
			}
		}()

		// Watch Namespace
		nameSpaceService := NameSpaceService{Cluster: this.Cluster, NameSpaceInterface: clientGo.ClientSet.CoreV1().Namespaces()}
		go nameSpaceService.Wtach()

		// Pod Namespace
		podService := PodService{Cluster: this.Cluster, PodInterface: clientGo.ClientSet.CoreV1().Pods("")}
		go podService.Wtach()

		// Node Namespace
		nodeService := NodeService{Cluster: this.Cluster, NodeInterface: clientGo.ClientSet.CoreV1().Nodes()}
		go nodeService.Wtach()

	}

	return resultData
}

func ClearCluster(clusterList []*models.Cluster) {
	// 清除集群数据
	k8sClearService := K8sClearService{ClusterList: clusterList, DropCluster: false}
	k8sClearService.ClearAll()
}

func (this *K8sWatchService) CreateK8sClient() ClientGo {
	// 构建k8s客户端
	params := new(ApiParams)
	params.AuthType = this.Cluster.AuthType
	params.BearerToken = this.Cluster.BearerToken
	params.MasterUrl = this.Cluster.MasterUrls
	params.KubeConfigPath = this.Cluster.FileName
	return CreateK8sClient(params)
}
