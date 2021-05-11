package kubevuln

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
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
	msg := ""

	// 检测集群客户端
	dymaicClient := models.GetDymaicClient(this.Cluster)
	client := models.GetClient(this.Cluster)
	if dymaicClient.ErrMessage != "" {
		result.Code = utils.CreateK8sClientErr
		if this.IsActive {
			msg = fmt.Sprint("Active kube scan job failed, ClusterName: %s, Err: %s", clusterName, result.Message)
		} else {
			msg = fmt.Sprint("Disable kube scan job failed, ClusterName: %s, Err: %s", clusterName, result.Message)
		}
		logs.Error(msg)
		result.Message = msg
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
	msg := ""
	result, kubeScan := this.CheckKubeScanLog()
	if result.Code != http.StatusOK {
		return result
	}
	taskService := base.TaskService{Task: this.Task}
	// 入库
	result = kubeScan.Add()
	if result.Code != http.StatusOK {
		taskService.Result = &result
		taskService.Task.Status = models.Task_Status_Failed
		taskService.UpdateTaskStatus()
		return result
	}

	// 删除job
	this.IsActive = false
	result = this.ActiveOrDisableKubeScan()
	if result.Code != http.StatusOK {
		msg = fmt.Sprint("Remove kube scan job failed, ClusterId: %s", this.Cluster.Id)
		logs.Error(msg)
		result.Message = msg
		taskService.Result = &result
		taskService.Task.Status = models.Task_Status_Failed
		taskService.UpdateTaskStatus()
		return result
	}

	msg = fmt.Sprintf("Update task success, Status: Finished, ClusterId: %s, Type: %s, TaskId: %s. <<<", this.Cluster.Id, this.Task.Type, this.Task.Id)
	result.Message = msg
	taskService.Result = &result
	taskService.Task.Status = models.Task_Status_Finished
	taskService.UpdateTaskStatus()
	logs.Info("Kube scan success,ClusterId: %s.", this.Cluster.Id)
	return result
}

func (this *KubeVlunService) CheckKubeScanLog() (models.Result, *models.KubeScan) {
	result := models.Result{Code: http.StatusOK}
	msg := ""
	kubeScan := models.KubeScan{}
	// 解析
	err := json.Unmarshal(this.RawLog, &kubeScan)
	if err != nil {
		msg = fmt.Sprintf("Receive kube scan result failed, Err: %s.", err)
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.ParseKubeVulnRawLogErr
		return result, nil
	}
	// 检查task
	if kubeScan.TaskId == "" {
		msg = "Receive kube scan result failed, Task is null."
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.GetTaskErr
		return result, nil
	}
	task := models.Task{}
	task.Id = kubeScan.TaskId
	dbTask := task.Get()
	if dbTask == nil {
		msg = "Receive kube scan result failed, Task is null."
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.GetTaskErr
		return result, nil
	}
	this.Task = dbTask

	// 检查集群
	if kubeScan.ClusterId == "" {
		msg = "Receive kube scan result failed, Cluster is null."
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.GetClusterErr
		return result, nil
	}
	cluster := models.Cluster{}
	cluster.Id = kubeScan.ClusterId
	dbCluster := cluster.Get()
	if dbCluster == nil {
		msg = "Receive kube scan result failed, Cluster is null."
		logs.Error(msg)
		result.Message = msg
		result.Code = utils.GetClusterErr
		return result, nil
	}
	this.Cluster = dbCluster
	return result, &kubeScan
}
