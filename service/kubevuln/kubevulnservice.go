package kubevuln

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"os"
)

type KubeVlunService struct {
	Cluster  *models.Cluster
	Task     *models.Task
	IsActive bool
	RawLog   []byte
}

func (this *KubeVlunService) ActiveOrDisableKubeScan() models.Result {
	result := models.Result{Code: http.StatusOK}
	clusterName := this.Cluster.Name

	// 检测集群客户端
	dymaicClient := models.GetDymaicClient(this.Cluster)
	client := models.GetClient(this.Cluster)
	if dymaicClient.ErrMessage != "" {
		result.Code = utils.CreateK8sClientErr
		result.Message = dymaicClient.ErrMessage
		logs.Error("Kube scan fail, ClusterName: %s, Err: %s", clusterName, result.Message)
		return result
	}

	// 读取yml
	f, err := os.Open(utils.GetKubeScanYml())
	if err != nil {
		result.Code = utils.OpenFileErr
		result.Message = err.Error()
		logs.Error("Open kube scan yml fail, ClusterName: %s, Err: %s", clusterName, result.Message)
		return result
	}

	// 创建 kube job
	kubernetesHandler := models.KubernetesHandler{ClientGo: &client, DymaicClient: &dymaicClient, IsActive: this.IsActive, File: f, Task: this.Task, Cluster: this.Cluster}
	return kubernetesHandler.CreateOrDeleteResourceByYml()
}

func (this *KubeVlunService) ReceiveKubeScanLog() models.Result {
	result := models.Result{Code: http.StatusOK}
	msg := ""
	//task := new(models.Task)
	kubeScan := models.KubeScan{}
	// 1. 解析
	err := json.Unmarshal(this.RawLog, &kubeScan)
	if err != nil {
		msg = fmt.Sprintf("Kube scan: Parse kube vuln raw log failed,  Err: %s.", err)
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.ParseKubeVulnRawLogErr
		return result
	}

	// 2. 获取任务
	//if kubeScan.TaskId != "" {
	//	task.Id = kubeScan.TaskId
	//	task = task.Get()
	//	if task == nil {
	//		result.Code = utils.GetTaskErr
	//		return result
	//	}
	//}

	// 3. 获取集群
	if kubeScan.ClusterId != "" {
		cluster := new(models.Cluster)
		cluster.Id = kubeScan.ClusterId
		this.Cluster = cluster.Get()
		if this.Cluster == nil {
			result.Code = utils.GetClusterErr
			return result
		}
	}

	// 4. 入库
	result = kubeScan.Add()
	if result.Code != http.StatusOK {
		//this.saveTaskLog(task, kubeScan, result)
		return result
	}

	// 5. 删除job
	if this.Cluster != nil {
		this.IsActive = false
		logs.Debug("Remove kube scan job, ClusterId: %s.", this.Cluster.Id)
		go this.ActiveOrDisableKubeScan()
	}
	return result
}

func updateKubeScanYml() {
	// 更新 KUBEHUNTER_HTTP_DISPATCH_URL
	// 更新 TASK_ID
	// 更新 CLUSTER_ID

}

func (this *KubeVlunService) saveTaskLog(task *models.Task, kubeScan models.KubeScan, result models.Result) {
	if task == nil && task.Id != "" {
		return
	}
	msg := fmt.Sprint("Kube Scan success, Clsuter Id: %s Task Id: %s.", kubeScan.ClusterId, task.Id)
	taskLog := models.TaskLog{}
	taskRawInfo, _ := json.Marshal(task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Level = models.Log_level_Info
	if result.Code != http.StatusOK {
		task.Status = models.Task_Status_Failed
		task.Update()
		msg = fmt.Sprint("Kube Scan failed, Clsuter Id: %s Task Id: %s Err: %s", kubeScan.ClusterId, task.Id, result.Message)
		logs.Error(msg)
	}
	taskLog.RawLog = msg
	taskLog.Add()
}
