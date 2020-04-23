package ws

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/bean"
	"github.com/xiliangMa/diss-backend/models/job"
	"github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/models/ws"
	"github.com/xiliangMa/diss-backend/service/synccheck"
	"net/http"
	"strings"
)

type WSMetricsService struct {
	Message []byte
	Conn    *websocket.Conn
}

func (this *WSMetricsService) Save() error {
	ms := ws.WsData{}
	if err := json.Unmarshal(this.Message, &ms); err != nil {
		logs.Error("Paraces WsData error %s", err)
		return err
	}
	switch ms.Type {
	case ws.Type_Metric: // 数据上报
		switch ms.Tag {
		case ws.Resource_HeartBeat:
			heartBeat := bean.HeartBeat{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &heartBeat); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			ip := strings.Split(this.Conn.RemoteAddr().String(), ":")
			client := &Client{Hub: WSHub, Conn: this.Conn, Send: make(chan []byte, 256), ClientIp: ip[0], SystemId: heartBeat.SystemId}
			client.Hub.Register <- client
			logs.Info("############################ Agent Heater Beat data, >>> HostId: %s, Type: %s <<<", heartBeat.SystemId, ws.Resource_HeartBeat)
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Resource_HostConfig:
			hostConfig := models.HostConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostConfig); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s <<<", hostConfig.Id, ws.Resource_HostConfig)
			if err := hostConfig.Inner_AddHostConfig(); err != nil {
				return err
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Resource_HostInfo:
			hostInfo := models.HostInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostInfo); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", hostInfo.Id, ws.Resource_HostInfo)
			if err := hostInfo.Inner_AddHostInfo(); err != nil {
				return err
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Resource_ContainerConfig:
			containerConfigList := []models.ContainerConfig{}
			CheckObject := new(models.ContainerConfig)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerConfigList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostName: %s, Type: %s, Size: %d <<<", containerConfigList[0].HostName, ws.Resource_ContainerConfig, size)
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
				agentCheckHandler.Check(ws.Resource_ContainerConfig)
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_ContainerInfo:
			containerInfoList := []models.ContainerInfo{}
			CheckObject := new(models.ContainerInfo)
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerInfoList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerInfoList[0].HostId, ws.Resource_ContainerInfo, size)
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
				agentCheckHandler.Check(ws.Resource_ContainerInfo)
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_ImageConfig:
			imageConfigList := []models.ImageConfig{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageConfigList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageConfigList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageConfigList[0].HostId, ws.Resource_ImageConfig, len(imageConfigList))
				//删除主机下的所有imageconfig
				imageConfigList[0].Delete()
			}

			for _, imageConfig := range imageConfigList {
				//if result := imageConfig.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageConfig.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_ImageInfo:
			imageInfoList := []models.ImageInfo{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &imageInfoList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(imageInfoList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", imageInfoList[0].HostId, ws.Resource_ImageInfo, len(imageInfoList))
				//删除主机下的所有imageginfo
				imageInfoList[0].Delete()
			}
			for _, imageInfo := range imageInfoList {
				//if result := imageInfo.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				imageInfo.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
		case ws.Resource_HostPs:
			hostPsList := []models.HostPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &hostPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(hostPsList)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", hostPsList[0].HostId, ws.Resource_HostPs, size)
				// 删除该主机下所有的进程
				hostPsList[0].Delete()

			}
			for _, hostPs := range hostPsList {
				//if result := hostPs.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				hostPs.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_ContainerPs:
			containerPsList := []models.ContainerPs{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &containerPsList); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(containerPsList)

			if size != 0 {
				logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s, Size: %d <<<", containerPsList[0].HostId, ws.Resource_ContainerPs, len(containerPsList))
				containerPsList[0].Delete()
			}

			for _, containerTop := range containerPsList {
				//if result := containerTop.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				containerTop.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_DockerBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := securitylog.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, ws.Resource_DockerBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
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
		case ws.Resource_KubernetesBenchMark:
			//index := beego.AppConfig.String("security_log::BenchMarkIndex")
			benchMarkLog := securitylog.BenchMarkLog{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &benchMarkLog); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			logs.Info("############################ Sync agent data, >>>  HostId: %s, Type: %s <<<", benchMarkLog.HostId, ws.Resource_KubernetesBenchMark)
			if result := benchMarkLog.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
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
		case ws.Resource_HostCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>> HostId: %s, ype: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, ws.Resource_HostCmdHistory, size)
				// 删除该主机下所有的记录 type = 0
				cmdHistoryList.List[0].Delete()
			}

			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		case ws.Resource_ContainerCmdHistory:
			cmdHistoryList := models.CmdHistoryList{}
			s, _ := json.Marshal(ms.Data)
			if err := json.Unmarshal(s, &cmdHistoryList.List); err != nil {
				logs.Error("Paraces %s error %s", ms.Tag, err)
				return err
			}
			size := len(cmdHistoryList.List)
			if size != 0 {
				logs.Info("############################ Sync agent data, >>> HostId: %s, Type: %s, Size: %d <<<", cmdHistoryList.List[0].HostId, ws.Resource_ContainerCmdHistory, size)
				// 删除该主机下所有的记录 type = 1
				cmdHistoryList.List[0].Delete()
			}
			for _, cmdHistory := range cmdHistoryList.List {
				//if result := cmdHistory.Add(); result.Code != http.StatusOK {
				//	return errors.New(result.Message)
				//}
				cmdHistory.Add()
			}
			metricsResult := ws.WsData{Type: ws.Type_ReceiveState, Tag: ws.Resource_Received, Data: nil, Config: ""}
			this.ReceiveData(metricsResult)
			return nil
		}
	case ws.Type_Control:
		switch ms.Tag {
		case ws.Resource_Task:
			// 获取任务列表接口
			switch ms.RCType {
			case ws.Resource_Control_Type_Get:
				metricsResult := ws.WsData{Code: http.StatusOK, Type: ws.Type_RequestState, Tag: ws.Resource_Task, RCType: ws.Resource_Control_Type_Get}
				task := job.Task{}
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
						logs.Info("############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, ws.Resource_Task, total)
					} else {
						logs.Info("############################  Get un finished task list, >>> HostId: %s, Type: %s, task size:  %v <<<", task.Host.Id, ws.Resource_Task, 0)
					}
				}
				this.ReceiveData(metricsResult)
			case ws.Resource_Control_Type_Put:
				//更新任务状态
				metricsResult := ws.WsData{Code: http.StatusOK, Type: ws.Type_RequestState, Tag: ws.Resource_Task, RCType: ws.Resource_Control_Type_Put}
				taskList := []job.Task{}
				s, _ := json.Marshal(ms.Data)
				if err := json.Unmarshal(s, &taskList); err != nil {
					logs.Error("Paraces: %s type: %s error: %s  ", ms.Tag, ms.RCType, err)
					return err
				}
				for _, task := range taskList {
					if result := task.Update(); result.Code != http.StatusOK {
						metricsResult.Code = result.Code
						metricsResult.Msg = result.Message
						logs.Error("############################ Update task status  fail, >>> HostId: %s, error: <<<", task.Host.Id, result.Message)
						return errors.New(result.Message)
					} else {
						this.ReceiveData(metricsResult)
						logs.Info("############################ Update task status, >>> HostId: %s, Type: %s, task id:  %v <<<", task.Host.Id, ws.Resource_Task, task.Id)
					}
				}
			}

		}
	}

	return nil
}

func (this *WSMetricsService) ReceiveData(result ws.WsData) {
	data, _ := json.Marshal(result)
	err := this.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logs.Info("############################ Received data from agent fail ############################", err)
	}
}
