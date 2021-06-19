package ws

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
	"github.com/xiliangMa/diss-backend/service/kubevuln"
	"net/http"
	"strings"
)

type WSDeliverService struct {
	*models.Hub
	Bath                 int64
	CurrentBatchTaskList []*models.Task
	DelTask              *models.Task
}

func (this *WSDeliverService) DeliverTaskToNats() {
	logs.Info("################ Deliver Task start, ################")
	for _, task := range this.CurrentBatchTaskList {
		msg := ""
		subject := ""
		// 根据是host还是container类型的任务获得topic
		if task.Host != nil {
			hostCount := task.Host.Count()
			if hostCount > 0 {
				subject = task.Host.Id
			}
		} else if task.Container != nil {
			containerCount := task.Container.Count()
			if containerCount > 0 && task.Container.HostName != "" {
				subject = task.Container.HostName
			}
		} else if task.Image != nil {
			subject = task.Image.HostId
		} else if task.ClusterOBJ != nil {
			subject = models.Subject_Cluster
		}

		if subject != "" {
			switch subject {
			case models.Subject_Cluster:
				if task.SystemTemplate.Type == models.TMP_Type_KubernetesVulnScan {
					// 强制删除之前存在的job
					kubeVlunService := kubevuln.KubeVlunService{IsActive: false, Cluster: task.ClusterOBJ, Task: task}
					result := kubeVlunService.ActiveOrDisableKubeScan()
					if result.Code != http.StatusOK && !strings.Contains(result.Message, "not found") {
						return
					}
					kubeVlunService.IsActive = true
					result = kubeVlunService.ActiveOrDisableKubeScan()
					if result.Code == http.StatusOK {
						task.Status = models.Task_Status_Pending
						msg = fmt.Sprintf("Update task success, Status: %s, ClusterId: %s, Type: %s, TaskId: %s. ", task.Status, task.ClusterOBJ.Id, task.Type, task.Id)
						result.Message = msg
						logs.Info("Deliver Task to cluster Success, TaskId: %s, ClusterName: %s.", task.Id, task.ClusterOBJ.Name)
						taskService := base.TaskService{Task: task, Result: &result}
						taskService.UpdateTaskStatus()

						// 设置running状态
						task.Status = models.Task_Status_Running
						msg = fmt.Sprintf("Update task success, Status: %s, ClusterId: %s, Type: %s, TaskId: %s. ", task.Status, task.ClusterOBJ.Id, task.Type, task.Id)
						result.Message = msg
						taskService = base.TaskService{Task: task, Result: &result}
						taskService.UpdateTaskStatus()
					} else {
						task.Status = models.Task_Status_Deliver_Failed
						msg = fmt.Sprintf("Update task failed, Status: %s, ClusterId: %s, Type: %s, TaskId: %s. ", task.Status, task.ClusterOBJ.Id, task.Type, task.Id)
						result.Message = msg
						logs.Error("Deliver Task to cluster failed, TaskId: %s, ClusterName: %s.", task.Id, task.ClusterOBJ.Name)
						taskService := base.TaskService{Task: task, Result: &result}
						taskService.UpdateTaskStatus()
					}

				}
			default:
				result := models.NatsData{Type: models.Type_Control, Tag: models.Resource_Task, Data: task, RCType: models.Resource_Control_Type_Post}
				data, _ := json.MarshalIndent(result, "", "  ")
				err := models.Nats.Conn.Publish(subject, data)
				logs.Debug("Send task data: %s.", string(data))
				if err == nil {
					logs.Info("Deliver Task to Nats Success, Subject: %s Id: %s, data: %v", subject, task.Id, result)
				} else {
					//更新 task 状态
					task.Status = models.Task_Status_Deliver_Failed
					task.Update()
					logs.Error("Deliver Task to Nats Fail,  Subject: %s Id: %s, err: %s", subject, task.Id, err.Error())
				}

				if task.Host != nil {
					hc := new(models.HostConfig)
					hc.Id = task.Host.Id
					hostConfig := hc.Get()
					hostConfig.TaskStatus = task.Status
					hostConfig.Update()
				} else if task.Image != nil {
					ic := new(models.ImageConfig)
					ic.Id = task.Image.Id
					imageConfig := ic.Get()
					imageConfig.TaskStatus = task.Status
					imageConfig.Update()
				} else if task.Container != nil {
					cc := new(models.ContainerConfig)
					cc.Id = task.Container.Id
					containerConfig := cc.Get()
					containerConfig.TaskStatus = task.Status
					containerConfig.Update()
				}
			}

		}
	}
	logs.Info("################ Deliver Task end, ################")
}

func (this *WSDeliverService) DeliverTask() {
	logs.Info("################ Deliver Task start, ################")
	for _, task := range this.CurrentBatchTaskList {
		if _, ok := this.Hub.DissClient[task.Host.Id]; ok {
			client := this.Hub.DissClient[task.Host.Id]
			result := models.NatsData{Type: models.Type_Control, Tag: models.Resource_Task, Data: task, RCType: models.Resource_Control_Type_Post}
			data, _ := json.Marshal(result)
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err == nil {
				msg := fmt.Sprintf("Deliver Task Success, Id: %s, data: %v", task.Id, result)
				logs.Info(msg)
				taskRawInfo, _ := json.Marshal(task)
				taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Info}
				taskLog.Add()
			} else {
				//更新 task 状态
				task.Status = models.Task_Status_Deliver_Failed
				task.Update()
				msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, err.Error())
				logs.Error(msg)
				taskRawInfo, _ := json.Marshal(task)
				taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
				taskLog.Add()
			}
		} else {
			//更新 task 状态
			task.Status = models.Task_Status_Deliver_Failed
			task.Update()
			errMsg := "Agent not connect"
			msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, errMsg)
			logs.Error(msg)
			taskRawInfo, _ := json.Marshal(task)
			taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
			taskLog.Add()
		}
	}
	logs.Info("################ Deliver Task end, ################")
}
