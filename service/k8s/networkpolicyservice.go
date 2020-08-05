package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

type NetworkPolicyService struct {
	NetworkPolicyInterface networkingv1.NetworkPolicyInterface
	Cluster                *models.Cluster
	Close                  chan bool
}

func (this *NetworkPolicyService) List() (*v1.NetworkPolicyList, error) {
	return this.NetworkPolicyInterface.List(metav1.ListOptions{})
}

func (this *NetworkPolicyService) Wtach() {
	netpolWatch, err := this.NetworkPolicyInterface.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close NetworkPolicyWatch, cluster: %s", this.Cluster.Name)
			netpolWatch.Stop()
			return
		case event, ok := <-netpolWatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.NetworkPolicy)
				id := string(object.UID)
				name := object.Name
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)

				netpol := models.NetworkPolicy{}
				netpol.Id = id
				netpol.Name = name
				netpol.AccountName = models.Account_Admin
				netpol.ClusterName = clusterName
				netpol.KMetaData = string(KMetaData)
				netpol.KSpec = string(KSpec)

				logs.Info("Watch >>> NetworkPolicy: %s <<<, >>> Cluster: %s <<<,  >>> NameSpace: %s <<<, >>> EventType: %s <<<", id, clusterId, name, event.Type)
				switch event.Type {
				case watch.Added:
					netpol.Add()
				case watch.Modified:
					netpol.Delete()
					netpol.Add()
				case watch.Deleted:
					netpol.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("NetworkPolicyWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.NetworkPolicy
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove NetworkPolicyWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除数据库数据
				netpol := models.NetworkPolicy{}
				netpol.ClusterName = this.Cluster.Name
				netpol.Delete()

				// 重启 watch 携程
				k8sWatchService := K8sWatchService{Cluster: this.Cluster}
				clientGo := k8sWatchService.CreateK8sClient()
				deployService := NetworkPolicyService{Cluster: this.Cluster, NetworkPolicyInterface: clientGo.ClientSet.NetworkingV1().NetworkPolicies(""), Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = deployService

				logs.Info("Retry NetworkPolicyWatch, cluster: %s", this.Cluster.Name)
				go deployService.Wtach()
				break Retry
			}
		}
	}
}
