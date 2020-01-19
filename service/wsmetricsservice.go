package service

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
)

type WSMetricsService struct {
	Message []byte
}

func (wsmh *WSMetricsService) Save() error {
	ms := models.MetricsResult{}
	if err := json.Unmarshal(wsmh.Message, &ms); err != nil {
		logs.Error("Paraces MetricsResult error %s", err)
		return err
	}
	switch ms.ResTag {
	case models.Tag_HostConfig:
		hostConfig := models.HostConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostConfig); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		if err := hostConfig.Inner_AddHostConfig(); err != nil {
			return err
		}
	case models.Tag_HostInfo:
		hostInfo := models.HostInfo{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostInfo); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		if err := hostInfo.Inner_AddHostInfo(); err != nil {
			return err
		}
	}

	return nil
}

func (wsmh *WSMetricsService) Update() error {
	return nil
}
