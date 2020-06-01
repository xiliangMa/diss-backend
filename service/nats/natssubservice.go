package nats

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type NatsSubService struct {
	Message       []byte
	Conn          stan.Conn
	ClientSubject string
	DelTask       *models.Task
}

func (this *NatsSubService) Save() error {
	ms := models.WsData{}
	if err := json.Unmarshal(this.Message, &ms); err != nil {
		logs.Error("Paraces WsData error %s", err)
		return err
	}
	if ms.Data == nil {
		return nil
	}
	switch ms.Type {
	case models.Type_Metric: // 数据上报
		switch ms.Tag {
		case models.Resource_HeartBeat:
			// to do 可以暂时采用点对点的 websocket 上报
			//heartBeat := models.HeartBeat{}
			//s, _ := json.Marshal(ms.Data)
			//if err := json.Unmarshal(s, &heartBeat); err != nil {
			//	logs.Error("Paraces %s error %s", ms.Tag, err)
			//	return err
			//}
			//ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
			//client := &models.Client{Hub: models.WSHub, Conn: this.Conn, Send: make(chan []byte, 256), ClientIp: ip[0], SystemId: heartBeat.SystemId}
			//client.Hub.Register <- client
			//logs.Info("Nats ############################ Agent Heater Beat data, >>> HostId: %s, Type: %s <<<", heartBeat.SystemId, models.Resource_HeartBeat)
			//metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			//this.ReceiveData(metricsResult)
		case models.Resource_HostInfoDynamic:
			// k8s 主机更新主机的外网ip
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostInfo.Id, models.Resource_HostInfoDynamic)
			if result := hostInfo.UpdateDynamic(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostConfigDynamic:
			// k8s 主机更新主机的外网ip
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Resource_HostConfigDynamic)
			if result := hostConfig.UpdateDynamic(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostConfig:
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Resource_HostConfig)
			if err := hostConfig.Inner_AddHostConfig(); err != nil {
				return err
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostInfo:
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, models.Resource_HostInfo)
			if err := hostInfo.Inner_AddHostInfo(); err != nil {
				return err
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_ContainerConfig:
			containerConfigList := []models.ContainerConfig{}
			CheckObject := new(models.ContainerConfig)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerConfigList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, models.Resource_ContainerConfig, size)
			}
			for _, containerConfig := range containerConfigList {
				//if result := containerConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerConfig.Add()
			}
			// 清除脏数据
			if size != 0 {
				CheckObject.HostName = containerConfigList[0].HostName
				CheckObject.SyncCheckPoint = containerConfigList[0].SyncCheckPoint
				agentCheckHandler := synccheck.AgentCheckHadler{CheckObject, nil}
				agentCheckHandler.Check(models.Resource_ContainerConfig)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerInfo:
			containerInfoList := []models.ContainerInfo{}
			CheckObject := new(models.ContainerInfo)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerInfoList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, models.Resource_ContainerInfo, size)
			}
			for _, containerInfo := range containerInfoList {
				//if result := containerInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerInfo.Add()
			}

			// 清除脏数据
			if size != 0 {
				CheckObject.HostName = containerInfoList[0].HostName
				CheckObject.SyncCheckPoint = containerInfoList[0].SyncCheckPoint
				agentCheckHandler := synccheck.AgentCheckHadler{nil, CheckObject}
				agentCheckHandler.Check(models.Resource_ContainerInfo)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ImageConfig:
			imageConfigList := []models.ImageConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageConfigList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, models.Resource_ImageConfig, len(imageConfigList))
				//删除主机下的所有imageconfig
				imageConfigList[0].Delete()
			}

			for _, imageConfig := range imageConfigList {
				//if result := imageConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageConfig.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ImageInfo:
			imageInfoList := []models.ImageInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageInfoList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, models.Resource_ImageInfo, len(imageInfoList))
				//删除主机下的所有imageginfo
				imageInfoList[0].Delete()
			}
			for _, imageInfo := range imageInfoList {
				//if result := imageInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageInfo.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostPs:
			hostPsList := []models.HostPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(hostPsList)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, models.Resource_HostPs, size)
				// 删除该主机下所有的进程
				hostPsList[0].Delete()

			}
			for _, hostPs := range hostPsList {
				//if result := hostPs.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				hostPs.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerPs:
			containerPsList := []models.ContainerPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerPsList)

			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, models.Resource_ContainerPs, len(containerPsList))
				containerPsList[0].Delete()
			}

			for _, containerTop := range containerPsList {
				//if result := containerTop.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerTop.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_DockerBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := models.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_DockerBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			// 上报es
			//esClient, err := utils.GetESClient()
			//if err != nil {
			//	logs.Error("Get ESClient err: %s", err)
			//}
			//
			//respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
			//if err != nil {
			//	logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
			//} else {
			//	logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
			//}
			//defer respones.Body.Close()
		case models.Resource_KubernetesBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := models.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("Nats ############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_KubernetesBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			// 上报es
			//esClient, err := utils.GetESClient()
			//if err != nil {
			//	logs.Error("Get ESClient err: %s", err)
			//}
			//respones, err := esClient.Create(index, benchMarkLog.Id, bytes.NewReader(s))
			//if err != nil {
			//	logs.Error("Add security_log to es fail, benchMarkLog.Id: %s", benchMarkLog.Id)
			//} else {
			//	logs.Info("Add security_log to es success, benchMarkLog.Id: %s", benchMarkLog.Id)
			//}
			//defer respones.Body.Close()
		case models.Resource_HostCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_HostCmdHistory, size)
				// 删除该主机下所有的记录 type = 0
				cmdHistoryList.List[0].Delete()
			}

			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case models.Resource_ContainerCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("Nats ############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_ContainerCmdHistory, size)
				// 删除该主机下所有的记录 type = 1
				cmdHistoryList.List[0].Delete()
			}
			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		}
	case models.Type_Control:
		switch ms.Tag {
		case models.Resource_Task:
			// 获取任务列表接口
			switch ms.RCType {
			case models.Resource_Control_Type_Get:
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Get}
				task := models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &task); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				if result := task.GetUnFinishedTaskList(); result.Code != http.StatusOK {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					logs.Error("Nats ############################ Get un finished task list  fail, >>> HostId: %s, error: <<<", task.Host.Id, result.Message)
					return errors.New(result.Message)
				} else {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					metricsResult.Data = result.Data
					if result.Data != nil {
						data := result.Data.(map[string]interface{})
						total := data["total"]
						if total != 0 {
							logs.Info("Nats ############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, total)
							this.ReceiveData(metricsResult)
						} else {
							logs.Info("Nats ############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, 0)
						}
					}
				}
			case models.Resource_Control_Type_Put:
				//更新任务状态
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Put}
				taskList := []models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &taskList); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				for _, task := range taskList {
					metricsResult.Data = task
					if result := task.Update(); result.Code != http.StatusOK {
						metricsResult.Code = result.Code
						metricsResult.Msg = result.Message
						msg := ""
						if task.Host == nil {
							msg = fmt.Sprintf("Nats ############################ Update task Status: %s, fail, >>> HostName: %s, task id: %s, error: %s <<<", task.Status, task.Container.HostName, task.Id, result.Message)
						} else {
							msg = fmt.Sprintf("Nats ############################ Update task Status: %s, fail, >>> HostId: %s, task id: %s, error: %s <<<", task.Status, task.Host.Id, task.Id, result.Message)
						}
						logs.Error(msg)
						taskLog := models.TaskLog{RawLog: msg, Task: &task, Account: task.Account, Level: models.Log_level_Error}
						taskLog.Add()
						return errors.New(result.Message)
					} else {
						msg := ""
						if task.Host == nil {
							msg = fmt.Sprintf("Nats ############################ Update task Success, Status: %s >>> HostName: %s, Type: %s, task id: %s <<<", task.Status, task.Container.HostName, models.Resource_Task, task.Id)
						} else {
							msg = fmt.Sprintf("Nats ############################ Update task Success, Status: %s >>> HostId: %s, Type: %s, task id: %s <<<", task.Status, task.Host.Id, models.Resource_Task, task.Id)
						}
						logs.Info(msg)
						taskLog := models.TaskLog{RawLog: msg, Task: &task, Account: task.Account, Level: models.Log_level_Info}
						taskLog.Add()
					}
					this.ReceiveData(metricsResult)
				}
			case models.Resource_Control_Type_Delete:
				metricsResult := models.WsData{Code: http.StatusOK, Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Get}
				task := models.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &task); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				if result := task.Delete(); result.Code != http.StatusOK {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					msg := ""
					if task.Host == nil {
						msg = fmt.Sprintf("Nats ############################ Delete task fail, >>> HostName: %s, , task id: %s, error: %s <<<", task.Container.HostName, task.Id, result.Message)
					} else {
						msg = fmt.Sprintf("Nats ############################ Delete task fail, >>> HostId: %s, , task id: %s, error: %s <<<", task.Host.Id, task.Id, result.Message)
					}
					logs.Error(msg)
					taskLog := models.TaskLog{RawLog: msg, Task: &task, Account: task.Account, Level: models.Log_level_Error}
					taskLog.Add()
					return errors.New(result.Message)
				} else {
					msg := ""
					if task.Host == nil {
						msg = fmt.Sprintf("Nats ############################  Delete task success, >>> HostName: %s, Type: %s, task id: %s<<<", task.Container.HostName, models.Resource_Task, task.Id)
					} else {
						msg = fmt.Sprintf("Nats ############################  Delete task success, >>> HostId: %s, Type: %s, task id: %s<<<", task.Host.Id, models.Resource_Task, task.Id)
					}
					logs.Info(msg)
					taskLog := models.TaskLog{RawLog: msg, Task: &task, Account: task.Account, Level: models.Log_level_Info}
					taskLog.Add()
				}
				this.ReceiveData(metricsResult)
			}
		}
	}

	return nil
}

func (this *NatsSubService) DeleteTask() error {
	logs.Info("################ Delete Task <<<start>>> ################")
	task := this.DelTask
	result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Delete, Data: task}
	data, err := json.Marshal(result)

	// 下发删除任务
	subject := ""
	// nats
	if task.Host != nil {
		subject = task.Host.Id
	}
	if task.Container != nil {
		subject = task.Container.HostName
	}
	if subject != "" {
		err = models.Nats.Conn.Publish(subject, data)
		if err == nil {
			logs.Info("Deliver Task to Nats Success, Subject: %s Id: %s, RCType: %s, data: %v", subject, task.Id, models.Resource_Control_Type_Delete, result)
		}
	} else {
		return errors.New(string(utils.ResourceNotFoundErr))
	}
	if err != nil {
		msg := fmt.Sprintf("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		logs.Error(msg)
		taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account, Level: models.Log_level_Error}
		taskLog.Add()
		return err
	}
	logs.Info("################ Delete Task <<<end>>> ################")
	return nil
}

func (this *NatsSubService) ReceiveData(result models.WsData) {
	data, _ := json.Marshal(result)
	err := this.Conn.Publish(this.ClientSubject, data)
	if err != nil {
		logs.Error("Nats ############################ Received data from agent fail ############################", err)
	}
}

func RunClientSub(clientSubject string) {
	serverSubject := models.Subject_Common + `-` + clientSubject
	natsManager := models.Nats
	if natsManager != nil && natsManager.Conn != nil {
		natsManager.Conn.Subscribe(serverSubject, func(m *stan.Msg) {
			natsSubService := NatsSubService{Conn: natsManager.Conn, Message: m.Data, ClientSubject: clientSubject}
			natsSubService.Save()
		}, stan.DurableName(serverSubject))
	}
}
