package service

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
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
	case models.Tag_ContainerConfig:
		containerConfigList := []models.ContainerConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &containerConfigList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		for _, containerConfig := range containerConfigList {
			if result := containerConfig.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}

	case models.Tag_ImageConfig:
		imageConfigList := []models.ImageConfig{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &imageConfigList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		for _, imageConfig := range imageConfigList {
			if result := imageConfig.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	case models.Tag_HostPs:
		hostPsList := []models.HostPs{}
		s, _ := json.Marshal(ms.Metric)
		if err := json.Unmarshal(s, &hostPsList); err != nil {
			logs.Error("Paraces %s error %s", ms.ResTag, err)
			return err
		}
		if len(hostPsList) != 0 {
			data := hostPsList[0].ListById().Data.(map[string]interface{})
			if data["total"] != 0 {
				for _, v := range data["items"].([]*models.HostPs) {
					v.Delete()
				}
			}
		}
		for _, hostPs := range hostPsList {
			if result := hostPs.Add(); result.Code != http.StatusOK {
				return errors.New(result.Message)
			}
		}
	}

	return nil
}
