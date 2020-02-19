package models

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

func Internal_HostMetricInfo_M(hostname string) Result {
	var ResultData Result
	esclient, err := utils.GetESClient()

	if err != nil {
		logs.Error("esclient connect error:", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	esqueryStr := strings.Replace(ESString("msearch_host_metric"), "!Param@hostname!", hostname, 4)
	mres, err := esclient.API.Msearch(strings.NewReader(esqueryStr), esclient.Msearch.WithIndex("metric*"))

	if err != nil {
		logs.Error("host msearch error: ", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}
	defer mres.Body.Close()

	logs.Info("host metric info is ok , %s", mres.Status())

	var hostInfo map[string]interface{}
	json.NewDecoder(mres.Body).Decode(&hostInfo)

	var hostInfoPure []interface{}
	for _, x := range hostInfo["responses"].([]interface{}) {
		hostInfoRefine1 := x.(map[string]interface{})["hits"]
		hostInfoRefine2 := hostInfoRefine1.(map[string]interface{})["hits"]
		hostInfoRefine2a := hostInfoRefine2.([]interface{})
		if len(hostInfoRefine2a) > 0 {
			hostInfoRefine3 := hostInfoRefine2a[0].(map[string]interface{})["_source"]
			hostInfoPure = append(hostInfoPure, hostInfoRefine3)
		}
	}

	ResultData.Data = hostInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}

// -- 因为结构复杂不使用的方法 ，通过 search - aggregation 方式获取es数据
func GetHostMetricInfo(hostname string) interface{} {
	var ResultData Result

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}
	res, _ := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(ESString("host_metric"))),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	defer res.Body.Close()

	var hostInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&hostInfo)

	hostInfoPure := hostInfo["aggregations"]

	return hostInfoPure
}

func Internal_ContainerListMetricInfo(hostname string) Result {
	var ResultData Result

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	lteTime, gteTime := utils.LteandGteTime()
	esqueryStr := strings.Replace(ESString("container_metric"), "!Param@gteTime!", gteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@lteTime!", lteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@hostname!", hostname, 1)

	res, err := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	if err != nil {
		logs.Error("host msearch error: ", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}
	defer res.Body.Close()

	var containerInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&containerInfo)

	var hostInfoPure []interface{}
	containerInfoRefine1 := containerInfo["aggregations"].(map[string]interface{})
	containerInfoRefine2 := containerInfoRefine1["container_list"].(map[string]interface{})
	for _, x := range containerInfoRefine2["buckets"].([]interface{}) {
		hostInfoPure = append(hostInfoPure, x)
	}

	ResultData.Data = hostInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}

func Internal_ContainerSummaryInfo(hostname string) Result {
	var ResultData Result

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	lteTime, gteTime := utils.LteandGteTime()
	esqueryStr := strings.Replace(ESString("container_summary"), "!Param@gteTime!", gteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@lteTime!", lteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@hostname!", hostname, 1)

	res, err := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())

	if err != nil {
		logs.Error("host msearch error: ", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}
	defer res.Body.Close()

	var containerSummaryInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&containerSummaryInfo)

	containerInfoPure := containerSummaryInfo["aggregations"].(map[string]interface{})

	ResultData.Data = containerInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}

func Internal_ImageListMetricInfo(hostname string) Result {
	var ResultData Result

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	lteTime, gteTime := utils.LteandGteTime()
	esqueryStr := strings.Replace(ESString("dockerimage_metric"), "!Param@gteTime!", gteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@lteTime!", lteTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@hostname!", hostname, 1)

	res, err := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	if err != nil {
		logs.Error("host msearch error: ", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}
	defer res.Body.Close()

	var imageInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&imageInfo)

	var hostInfoPure []interface{}
	containerInfoRefine1 := imageInfo["aggregations"].(map[string]interface{})
	containerInfoRefine2 := containerInfoRefine1["image_list"].(map[string]interface{})
	for _, x := range containerInfoRefine2["buckets"].([]interface{}) {
		hostInfoPure = append(hostInfoPure, x)
	}

	ResultData.Data = hostInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}


func Internal_IntrudeDetectMetricInfo(hostId , fromTime , toTime string) Result {
	var ResultData Result

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	esqueryStr := strings.Replace(ESString("intrude_detect"), "!Param@gteTime!", fromTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@lteTime!", toTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@hostname!", hostId, 1)

	res, err := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex(hostId),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	if err != nil {
		logs.Error("host msearch error: ", err.Error())
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}
	defer res.Body.Close()

	var indDetInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&indDetInfo)

	var intDetPure []interface{}
	intDetPureRefine1 := indDetInfo["hits"].(map[string]interface{})
	intDetPure = intDetPureRefine1["hits"].([]interface{})

	ResultData.Data = intDetPure
	ResultData.Code = http.StatusOK
	return ResultData
}