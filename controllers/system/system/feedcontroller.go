package system

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
)

// Feeds（订阅）object api list
type FeedsController struct {
	web.Controller
}

// @Title GetFeeds
// @Description Get Feed List
// @Param token header string true "authToken"
// @Param body body models.Feeds false "订阅"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router /feeds [post]
func (this *FeedsController) GetFeedList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	feeds := new(models.Feeds)
	json.Unmarshal(this.Ctx.Input.RequestBody, &feeds)
	this.Data["json"] = feeds.List(from, limit)
	this.ServeJSON(false)
}
