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

type PodService struct {
	PodInterface corev1.PodInterface
	Cluster      *models.Cluster
}

func (this *PodService) List() (*v1.PodList, error) {
	return this.PodInterface.List(metav1.ListOptions{})
}

func (this *PodService) Wtach() {
	podWatch, err := this.PodInterface.Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	// 开启 watch 事件
Done:
	for {
		select {
		case event, ok := <-podWatch.ResultChan():
			if event.Object != nil || ok {
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
				switch event.Type {
				case watch.Added:
					pod.Add()
				case watch.Modified:
					pod.Add()
					//result := pod.List(0, 0)
					//if result.Code == http.StatusOK {
					//	data := result.Data.(map[string]interface{})
					//	if data != nil {
					//		items := data["items"].([]*models.Pod)
					//		if items != nil && len(items) == 1 {
					//			dbPod := items[0]
					//			dbPod.Name = name
					//			dbPod.HostName = nodeName
					//			dbPod.ClusterName = clusterName
					//			dbPod.KMetaData = string(KMetaData)
					//			dbPod.KSpec = string(KSpec)
					//			dbPod.KStatus = string(KStatus)
					//			dbPod.Update()
					//		}
					//	}
					//}
				case watch.Deleted:
					pod.Delete()
				case watch.Bookmark:
					//todo
				case watch.Error:
					//todo
				}
			} else {
				// 如果 watch 异常退回重新 watch
				logs.Error("podWatch chan has been close!!!!")
				break Done
			}
		}
	}
}
