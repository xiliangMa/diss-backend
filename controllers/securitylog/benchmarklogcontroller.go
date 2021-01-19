package securitypolicy

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// Bench Mark Log api list
type BenchMarkLogController struct {
	web.Controller
}

// @Title GetBenchMarkLog
// @Description Get BenchMarkLog List (暂不支持租户查询)
// @Param token header string true "authToken"
// @Param body body models.BenchMarkLog false "基线日志信息"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /bmls [post]
func (this *BenchMarkLogController) GetBenchMarkLogList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	benchMarkLog := new(models.BenchMarkLog)
	json.Unmarshal(this.Ctx.Input.RequestBody, &benchMarkLog)
	this.Data["json"] = benchMarkLog.List(from, limit)
	this.ServeJSON(false)

}
