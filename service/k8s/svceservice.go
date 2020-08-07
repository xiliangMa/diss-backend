package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type SVCService struct {
	ClientGo ClientGo
	Cluster  *models.Cluster
	Close    chan bool
}

func (this *SVCService) List() (*v1.ServiceList, error) {
	return this.ClientGo.ClientSet.CoreV1().Services("").List(metav1.ListOptions{})
}

func (this *SVCService) Wtach() {
	serviceWatch, err := this.ClientGo.ClientSet.CoreV1().Services("").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
		case <-this.Close:
			logs.Info("Close serviceWatch, cluster: %s", this.Cluster.Name)
			serviceWatch.Stop()
			return
		case event, ok := <-serviceWatch.ResultChan():
			if event.Object != nil || ok {
				object := event.Object.(*v1.Service)
				id := string(object.UID)
				name := object.Name
				clusterName := this.Cluster.Name
				clusterId := this.Cluster.Id

				KMetaData, _ := json.Marshal(object.ObjectMeta)
				KSpec, _ := json.Marshal(object.Spec)
				KStatus, _ := json.Marshal(object.Status)

				service := models.Service{}
				service.Id = id
				service.Name = name
				service.ClusterName = clusterName
				service.ClusterId = clusterId
				service.KMetaData = string(KMetaData)
				service.KSpec = string(KSpec)
				service.KStatus = string(KStatus)

				logs.Info("Watch >>> Service: %s <<<, >>> Cluster: %s <<<, >>> EventType: %s <<<", id, clusterName, event.Type)
				switch event.Type {
				case watch.Added:
					service.Add()
				case watch.Modified:
					service.Add()
				case watch.Deleted:
					service.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Warn("ServiceWatch chan has been close!!!! cluster: %s", this.Cluster.Name)

				// 清除全局 GRM（携程对象）
				watchType := this.Cluster.Id + `_` + utils.Service
				delete(models.GRM.GoRoutineMap, watchType)
				logs.Info("Remove ServiceWatch from global GRM object, cluster: %s", this.Cluster.Name)

				// 清除数据库数据
				svc := models.Service{}
				svc.ClusterName = this.Cluster.Name
				svc.Delete()

				// 重启 watch 携程
				svcService := SVCService{Cluster: this.Cluster, ClientGo: this.ClientGo, Close: make(chan bool)}
				models.GRM.GoRoutineMap[watchType] = svcService

				logs.Info("Retry ServiceWatch, cluster: %s", this.Cluster.Name)
				go svcService.Wtach()

				break Retry
			}
		}
	}
}
