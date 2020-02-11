package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type K8STaskHandler struct {
	Clientgo utils.ClientGo
}

func NewK8STaskHandler(path string) *K8STaskHandler {
	return &K8STaskHandler{
		Clientgo: utils.CreateK8sClient(path),
	}
}

func (this *K8STaskHandler) SyncCluster() {
}

func (this *K8STaskHandler) SyncHostConfig() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
		logs.Error("Sync node err: %s", err.Error())
	} else {
		for _, ns := range nodes.Items {
			uid, _ := uuid.NewV4()
			ob := new(models.HostConfig)
			ob.HostName = ns.ObjectMeta.Name
			ob.Id = uid.String()
			ob.Inner_AddHostConfig()
		}
	}
}

func (this *K8STaskHandler) SyncNameSpace() {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
		logs.Error("Sync namspace err: %s", err.Error())
	} else {
		for _, ns := range nameSpaces.Items {
			ob := new(k8s.NameSpace)

			NID, _ := uuid.NewV4()
			ob.Id = NID.String()
			ob.Name = ns.ObjectMeta.Name
			ob.Add()
			this.SyncPod(ns.ObjectMeta.Name)
		}
	}

}

func (this *K8STaskHandler) SyncPod(namespace string) {
	pods, err := this.Clientgo.GetPodsByNameSpace(namespace)
	if err != nil {
		logs.Error("Sync namespace: %s pods err: %s", namespace, err.Error())
	} else {
		for _, pod := range pods.Items {
			ob := new(k8s.Pod)

			NID, _ := uuid.NewV4()
			ob.Id = NID.String()
			ob.Name = pod.ObjectMeta.Name

			ob.Add()
		}
	}

}

type Data struct {
	items []*k8s.Cluster
	total int
}

func SyncAll() {
	// cluster
	var cluster k8s.Cluster
	result := cluster.List()

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		for _, c := range data["items"].([]*k8s.Cluster) {
			if c.IsSync == k8s.Cluster_IsSync && c.Synced == k8s.Cluster_NoSynced {
				logs.Info("Sync start...., cluster name: %s ", c.Name)
				// 创建k8s客户端
				this := NewK8STaskHandler(c.FileName)

				// 同步主机
				this.SyncHostConfig()

				// 同步 ns &&  ns 下的 pod
				this.SyncNameSpace()
				logs.Info("Sync end...., cluster name: %s ", c.Name)
				// 更新初始化状态
				c.Synced = k8s.Cluster_Synced
				c.Update()
				logs.Info("Update cluster: %s  Synced status: true ", c.Name)
			}
		}

	}
}
