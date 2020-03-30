package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models/k8s"
)

// 命名空间接口
type NSController struct {
	beego.Controller
}

// @Title GetNameSpaces
// @Description Get NameSpace List
// @Param token header string true "authToken"
// @Param body body k8s.NameSpace false "命名空间"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *NSController) GetNameSpaceList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	ns := new(k8s.NameSpace)
	json.Unmarshal(this.Ctx.Input.RequestBody, &ns)
	this.Data["json"] = ns.List(from, limit)
	this.ServeJSON(false)
}
