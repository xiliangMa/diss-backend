package task

import (
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/utils"
)

type K8STaskHandler struct {
	Clientgo utils.ClientGo
}

func (this *K8STaskHandler) SyncCluster() {
}

func (this *K8STaskHandler) SyncHostConfig() {
	nodes, err := this.Clientgo.GetNodes()
	if err != nil {
	}
	for _, ns := range nodes.Items {
		ob := new(models.HostConfig)
		ob.HostName = ns.ObjectMeta.Name
		ob.Inner_AddHostConfig()
	}
}

func (this *K8STaskHandler) SyncNameSpace() {
	nameSpaces, err := this.Clientgo.GetNameSpaces()
	if err != nil {
	}
	for _, ns := range nameSpaces.Items {
		ob := new(k8s.NameSpace)

		NID, _ := uuid.NewV4()
		ob.Id = NID.String()
		ob.Name = ns.ObjectMeta.Name

		ob.Add()
	}
}

func (this *K8STaskHandler) SyncPod() {
	pods, err := this.Clientgo.GetPodsByNameSpace("default")
	if err != nil {
	}
	for _, pod := range pods.Items {
		ob := new(k8s.Pod)

		NID, _ := uuid.NewV4()
		ob.Id = NID.String()
		ob.Name = pod.ObjectMeta.Name

		ob.Add()
	}
}
