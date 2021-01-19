package k8s

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/scope"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PodService struct {
	ClientGo models.ClientGo
	Cluster  *models.Cluster
	Close    chan bool
}

func (this *PodService) List() (*v1.PodList, error) {
	ns := ""
	if this.Cluster != nil {
		ns = this.Cluster.Name
	}
	return this.ClientGo.ClientSet.CoreV1().Pods(ns).List(nil, metav1.ListOptions{})
}

func (this *PodService) Wtach() {
	podWatch, err := this.ClientGo.ClientSet.CoreV1().Pods("").Watch(nil, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close podWatch, cluster: %s", this.Cluster.Name)
			return
		case event, ok := <-podWatch.ResultChan():
			if event.Object != nil || ok {
				isDelete := false
				object := event.Object.(*v1.Pod)
				id := string(object.UID)
				name := object.Name
				nodeName := object.Spec.NodeName
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)
				KStatus, _ := json.Marshal(object.Status)

				pod := models.Pod{}
				pod.Id = id
				pod.Name = name
				pod.HostName = nodeName
				pod.ClusterName = clusterName
				pod.KMetaData = string(KMetaData)
				pod.KSpec = string(KSpec)
				pod.KStatus = string(KStatus)

				logs.Info("Watch >>> Pod: %s <<<, >>> Cluster: %s <<<, >>> EventType: %s <<<", id, clusterName, event.Type)

				// 根据 scope app 状态动态组件 集群 scope url
				scopeService := scope.ScopeService{Pod: object, Cluster: this.Cluster}

				switch event.Type {
				case watch.Added:
					pod.Add()
				case watch.Modified:
					pod.Add()
				case watch.Deleted:
					isDelete = true
					pod.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
				// 动态初始化 socpe url
				scopeService.UpdatetClusterScopeUrlAndStatus(isDelete)
				// 同步pod下container到数据库
				containerService := ContainerService{Cluster: this.Cluster, Pod: object}
				containerService.InitContainer(event.Type)
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("PodWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.Pod
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove PodWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除 pod 下 container 数据
				k8sClearService := K8sClearService{CurrentCluster: this.Cluster}
				k8sClearService.ClearContainer()

				// 清除数据库数据
				pod := models.Pod{}
				pod.ClusterName = this.Cluster.Name
				pod.Delete()

				// 重启 watch 携程
				podService := PodService{Cluster: this.Cluster, ClientGo: this.ClientGo, Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = podService

				logs.Info("Retry PodWatch, cluster: %s", this.Cluster.Name)
				podService.Wtach()
				break Retry
			}
		}
	}
}
