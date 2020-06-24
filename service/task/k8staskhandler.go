package task

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"net/http"
	"time"
)

type K8STaskHandler struct {
	K8sSyncService *k8s.K8sSyncService
}

func (this *K8STaskHandler) SyncAll() {
	// cluster
	var cluster models.Cluster
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		syncCheckPoint := time.Now().Unix()
		for _, c := range data["items"].([]*models.Cluster) {
			k8sSyncService := k8s.NewK8sSyncService(syncCheckPoint, c)
			go k8sSyncService.Sync()
		}
	}
}
