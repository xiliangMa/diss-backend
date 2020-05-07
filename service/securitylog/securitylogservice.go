package securitylog

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type SecurityLogService struct {
	*models.BenchMarkLog
	*models.IntrudeDetectLog
}

func (this *SecurityLogService) GetHostBenchMarkLogInfo() models.Result {
	var ResultData models.Result

	esClient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}
	index := beego.AppConfig.String("security_log::BenchMarkIndex")
	//esqueryStr := `{
	//"size": 1,
	//"query": {
	// "bool": {
	//   "must":[
	//       {"match":{"HostId": "` + this.HostId + `"}},
	//       {"match":{"Id": "` + this.Id + `"}}
	//   ]
	// }
	//}
	//}`
	//oHostId := make(map[string]interface{})
	//oHostId["HostId"] = this.HostId
	oId := make(map[string]interface{})
	oId["_id"] = this.BenchMarkLog.Id

	//mustHostId := make(map[string]interface{})
	//mustHostId["match"] = oHostId

	mustId := make(map[string]interface{})
	mustId["match"] = oId

	var must []map[string]interface{}
	must = append(must, mustId)
	//must = append(must, mustHostId)

	esqueryStr := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"size": 10,
	}
	jsonBody, _ := json.Marshal(esqueryStr)

	res, err := esClient.API.Search(esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(bytes.NewReader(jsonBody)),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty())
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

func (this *SecurityLogService) GetIntrudeDetectLogInfo() models.Result {
	var ResultData models.Result
	targetType := this.IntrudeDetectLog.TargeType
	containerid := this.IntrudeDetectLog.ContainerId
	startTime := this.IntrudeDetectLog.StartTime
	toTime := this.IntrudeDetectLog.ToTime
	hostId := strings.ToLower(this.IntrudeDetectLog.HostId)
	limit := strconv.Itoa(this.IntrudeDetectLog.Limit)

	esclient, err := utils.GetESClient()
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}
	matchMode := "should"
	containerFilterStr := ""
	if targetType == "container" {
		targetType = "host"
		matchMode = "must_not"
		if containerid != "" {
			shortId := string([]byte(containerid)[:12])
			containerFilterStr = strings.Replace(models.ContainerFilterPattern, "!Param@containerId!", shortId, 1)
		}
	}

	esqueryStr := strings.Replace(models.ESString("intrude_detect"), "!Param@gteTime!", startTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@lteTime!", toTime, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@hostname!", hostId, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@targetTypeM!", matchMode, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@targetType!", targetType, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Filter@container!", containerFilterStr, 1)
	esqueryStr = strings.Replace(esqueryStr, "!Param@limit!", limit, 1)

	// fmt.Println("esqueryStr\n",  esqueryStr)
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

	if indDetInfo != nil && indDetInfo["hits"] != nil {
		intDetPureRefine1 := indDetInfo["hits"].(map[string]interface{})
		intDetPure = intDetPureRefine1["hits"].([]interface{})
	}

	ResultData.Data = intDetPure
	ResultData.Code = http.StatusOK
	return ResultData
}
