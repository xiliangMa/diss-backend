package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/scope"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type NameSpaceService struct {
	ClientGo models.ClientGo
	Cluster  *models.Cluster
	Close    chan bool
}

func (this *NameSpaceService) List() (*v1.NamespaceList, error) {
	return this.ClientGo.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
}

func (this *NameSpaceService) Wtach() {
	nswatch, err := this.ClientGo.ClientSet.CoreV1().Namespaces().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close namespaceWatch, cluster: %s", this.Cluster.Name)
			nswatch.Stop()
			return
		case event, ok := <-nswatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.Namespace)
				id := string(object.UID)
				name := object.Name
				clusterId := this.Cluster.Id
				clusterName := this.Cluster.Name

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)
				KStatus, _ := json.Marshal(object.Status)

				ns := models.NameSpace{}
				ns.Id = id
				ns.Name = name
				ns.AccountName = models.Account_Admin
				ns.ClusterId = clusterId
				ns.ClusterName = clusterName
				ns.KMetaData = string(KMetaData)
				ns.KSpec = string(KSpec)
				ns.KStatus = string(KStatus)

				logs.Info("Watch >>> Namespace: %s <<<, >>> Cluster: %s <<<, >>> EventType: %s <<<", id, clusterId, event.Type)
				switch event.Type {
				case watch.Added:
					ns.Add(false)
				case watch.Modified:
					ns.Delete()
					ns.Add(true)
				case watch.Deleted:
					scopeService := scope.ScopeService{Namespace: object, Cluster: this.Cluster}
					scopeService.CheckScopeIsDisable()
					ns.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("NamespaceWatch chan has been close!!!!, cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.NameSpace
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove NamespaceWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除数据库数据
				ns := models.NameSpace{}
				ns.ClusterId = this.Cluster.Id
				ns.Delete()

				// 重启 watch 携程
				nameSpaceService := NameSpaceService{Cluster: this.Cluster, ClientGo: this.ClientGo, Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = nameSpaceService

				logs.Info("Retry NamespaceWatch, cluster: %s", this.Cluster.Name)
				go nameSpaceService.Wtach()
				break Retry
			}
		}
	}
}
