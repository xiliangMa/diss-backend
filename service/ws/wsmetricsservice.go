package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"net/http"
	"strings"
	"time"
)

type WSMetricsService struct {
	Message []byte
	Conn    *websocket.Conn
}

func (this *WSMetricsService) Save() error {
	ms := models.WsData{}
	if err := json.Unmarshal(this.Message, &ms); err != nil {
		logs.Error("Paraces WsData error %s", err)
		return err
	}
	switch ms.Type {
	case models.Type_Metric: // 数据上报
		switch ms.Tag {
		case models.Resource_HeartBeat:
			heartBeat := models.HeartBeat{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &heartBeat); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
			client := &models.Client{Hub: models.WSHub, Conn: this.Conn, Send: make(chan []byte, 256), ClientIp: ip[0], SystemId: heartBeat.SystemId}

			// 开启 nats 订阅
			if models.WSHub != nil {
				if _, ok := models.WSHub.DissClient[client.SystemId]; !ok {
					nats.RunClientSub(client.SystemId)
					logs.Info("Run nats client sub success , HostId: %s", client.SystemId)

					hostObj := new(models.HostConfig)
					hostObj.Id = client.SystemId
					host := hostObj.Get()
					if host != nil {
						nats.RunClientSub(host.HostName)
						logs.Info("Run nats client sub success , HostName: %s", host.HostName)
					}
				}
			}
			client.Hub.Register <- client

			//更新主机心跳（更新时间
			host := models.HostConfig{Id: heartBeat.SystemId}
			listData := host.List(0, 0).Data
			if listData == nil {
				return nil
			}
			data := listData.(map[string]interface{})
			items := data["items"].([]*models.HostConfig)
			if len(items) != 0 {
				currentHost := items[0]
				currentHost.HeartBeat = time.Now()
				currentHost.IsEnableHeartBeat = true
				currentHost.Update()
			}
			logs.Info("############################ Agent Heater Beat data, >>> HostId: %s, Type: %s <<<", heartBeat.SystemId, models.Resource_HeartBeat)
			metricsResult := models.WsData{Type: models.Type_ReceiveState, Tag: models.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case models.Resource_HostConfig:
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, models.Resource_HostConfig)
			if err := hostConfig.Add(); err != nil {
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
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, models.Resource_HostInfo)
			if err := hostInfo.Add(); err != nil {
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
				logs.Info("############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, models.Resource_ContainerConfig, size)
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
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, models.Resource_ContainerInfo, size)
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
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, models.Resource_ImageConfig, len(imageConfigList))
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
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, models.Resource_ImageInfo, len(imageInfoList))
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
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, models.Resource_HostPs, size)
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
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, models.Resource_ContainerPs, len(containerPsList))
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
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_DockerBenchMark)
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
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, models.Resource_KubernetesBenchMark)
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
				logs.Info("############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_HostCmdHistory, size)
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
				logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, models.Resource_ContainerCmdHistory, size)
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
					logs.Error("############################ Get un finished task list  fail, >>> HostId: %s, error: <<<", task.Host.Id, result.Message)
					return errors.New(result.Message)
				} else {
					metricsResult.Code = result.Code
					metricsResult.Msg = result.Message
					metricsResult.Data = result.Data
					if result.Data != nil {
						data := result.Data.(map[string]interface{})
						total := data["total"]
						logs.Info("############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, total)
					} else {
						logs.Info("############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, models.Resource_Task, 0)
					}
				}
				this.ReceiveData(metricsResult)
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
						msg := fmt.Sprintf("############################ Update task KStatus: %s, fail, >>> HostId: %s, error: <<<", task.Status, task.Host.Id, result.Message)
						logs.Error(msg)
						taskRawInfo, _ := json.Marshal(task)
						taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Error}
						taskLog.Add()
						return errors.New(result.Message)
					} else {
						msg := fmt.Sprintf("############################ Update task KStatus: %s, >>> HostId: %s, Type: %s, task id:  %v <<<", task.Status, task.Host.Id, models.Resource_Task, task.Id)
						logs.Info(msg)
						taskRawInfo, _ := json.Marshal(task)
						taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Info}
						taskLog.Add()
					}
					this.ReceiveData(metricsResult)
				}
			}

		}
	}

	return nil
}

func (this *WSMetricsService) ReceiveData(result models.WsData) {
	data, _ := json.Marshal(result)
	err := this.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logs.Error("############################ Received data from agent fail ############################", err)
	}
}
