package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
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
	cluster.Status = models.Cluster_Status_Active
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
	resultData.Code = http.StatusOK

	clusterName := this.Cluster.Name
	clusterId := this.Cluster.Id
	clientGo := models.ClientGo{}
	if _, ok := models.KCM.ClientHub[clusterId]; !ok {
		models.KCM.ClientHub[clusterId] = models.CreateK8sClient(models.BuildApiParams(this.Cluster))
	}
	clientGo = models.KCM.ClientHub[clusterId]

	if clientGo.ErrMessage == "" {
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_InProcess
		this.Cluster.Update()
		logs.Info("ClusterOBJ:  %s, Watch start.", clusterName)

		// 异常捕获更新状态
		defer func() {
			if err := recover(); err != nil {
				// 更新集群的同步状态
				this.Cluster.SyncStatus = models.Cluster_Watch_Status_Fail
				this.Cluster.Update()
				logs.Error("ClusterOBJ: %s, id: %s , Watch fail. err: %s", clusterName, clusterId, err)
			}
		}()
		watchType := ""
		// Namespace
		nameSpaceService := NameSpaceService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.NameSpace
		models.GRM.GoRoutineMap[watchType] = nameSpaceService
		go nameSpaceService.Wtach()

		// Pod
		podService := PodService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.Pod
		models.GRM.GoRoutineMap[watchType] = podService
		go podService.Wtach()

		// Deployengt
		deployService := DeploymentService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.Deployment
		models.GRM.GoRoutineMap[watchType] = deployService
		go deployService.Wtach()

		// Node
		nodeService := NodeService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.Host
		models.GRM.GoRoutineMap[watchType] = nodeService
		go nodeService.Wtach()

		// Service
		svcService := SVCService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.Service
		models.GRM.GoRoutineMap[watchType] = svcService
		go svcService.Wtach()

		// NetworkPolicy
		netpolService := NetworkPolicyService{Cluster: this.Cluster, ClientGo: clientGo, Close: make(chan bool)}
		watchType = this.Cluster.Id + `_` + utils.NetworkPolicy
		models.GRM.GoRoutineMap[watchType] = netpolService
		go netpolService.Wtach()

	} else {
		resultData.Code = utils.CreateK8sClientErr
		resultData.Message = clientGo.ErrMessage
	}

	if resultData.Code == http.StatusOK {
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_Synced
	} else {
		this.Cluster.SyncStatus = models.Cluster_Sync_Status_Fail
	}
	this.Cluster.Update()

	return resultData
}

func ClearCluster(clusterList []*models.Cluster) {
	// 清除集群数据
	k8sClearService := K8sClearService{ClusterList: clusterList, DropCluster: false}
	k8sClearService.ClearAll()
}
