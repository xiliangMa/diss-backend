package securitylog

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	msl "github.com/xiliangMa/diss-backend/models/securitylog"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SecurityLogService struct {
	*msl.BenchMarkLog
}

type Match struct {
	key   string
	value string
}

func (this *SecurityLogService) GetSecurityLogInfo() models.Result {
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
	oId["_id"] = this.Id

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
