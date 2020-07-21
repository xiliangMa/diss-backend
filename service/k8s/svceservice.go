package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type SVCService struct {
	ServiceInterface corev1.ServiceInterface
	Cluster          *models.Cluster
}

func (this *SVCService) List() (*v1.ServiceList, error) {
	return this.ServiceInterface.List(metav1.ListOptions{})
}

func (this *SVCService) Wtach() {
	serviceWatch, err := this.ServiceInterface.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Retry:
	for {
		select {
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
					//result := service.List(0, 0)
					//if result.Code == http.StatusOK {
					//	data := result.Data.(map[string]interface{})
					//	if data != nil {
					//		items := data["items"].([]*models.Service)
					//		if items != nil && len(items) == 1 {
					//			dbservice := items[0]
					//			dbservice.Name = name
					//			dbservice.ClusterName = clusterName
					//			dbservice.KMetaData = string(KMetaData)
					//			dbservice.KSpec = string(KSpec)
					//			dbservice.KStatus = string(KStatus)
					//			dbservice.Update()
					//		}
					//	}
					//}
				case watch.Deleted:
					service.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Error("serviceWatch chan has been close!!!!")
				break Retry
			}
		}
	}
}
