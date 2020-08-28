package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

// NetworkPolicy 接口
type NetworkPolicyController struct {
	beego.Controller
}

// @Title GetNetworkPolicy
// @Description Get NetworkPolicy List
// @Param token header string true "authToken"
// @Param body body models.NetworkPolicy false "网络策略"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *NetworkPolicyController) GetNetworkPolicysList() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	this.Data["json"] = NetworkPolicy.List(from, limit)
	this.ServeJSON(false)
}

// @Title AddNetworkPolicy
// @Description Add NetworkPolicy
// @Param token header string true "authToken"
// @Param body body models.NetworkPolicy true "网络策略"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *NetworkPolicyController) AddNetworkPolicy() {
	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	netpolService := k8s.NetworkPolicyService{NetworkPolicy: NetworkPolicy, ClientGo: models.KCM.ClientHub[NetworkPolicy.ClusterId]}
	object, err := netpolService.Create()
	if err != nil {
		logs.Error("Add NetworkPolicy fail, err: %s", err)
		this.Data["json"] = models.Result{Code: utils.AddNetworkPolicyErr, Message: err.Error()}
		this.ServeJSON(false)
	} else {
		NetworkPolicy.Id = string(object.UID)
		this.Data["json"] = NetworkPolicy.Add()
		this.ServeJSON(false)
	}

}

// @Title UpdateNetworkPolicy
// @Description Update NetworkPolicy
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Param body body models.NetworkPolicy true "网络策略"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *NetworkPolicyController) UpdateNetworkPolicy() {
	id := this.GetString(":id")
	NetworkPolicy := new(models.NetworkPolicy)
	json.Unmarshal(this.Ctx.Input.RequestBody, &NetworkPolicy)
	NetworkPolicy.Id = id
	this.Data["json"] = NetworkPolicy.Update()
	this.ServeJSON(false)
}

// @Title DeleteNetworkPolicy
// @Description Delete NetworkPolicy
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *NetworkPolicyController) DeleteNetworkPolicy() {
	id := this.GetString(":id")
	NetworkPolicy := new(models.NetworkPolicy)
	NetworkPolicy.Id = id
	clientGo := models.ClientGo{}
	if data := NetworkPolicy.List(0, 0).Data; data != nil {
		items := data.(map[string]interface{})["items"].([]*models.NetworkPolicy)
		NetworkPolicy = items[0]
		if _, ok := models.KCM.ClientHub[NetworkPolicy.ClusterId]; ok {
			clientGo = models.KCM.ClientHub[NetworkPolicy.ClusterId]
		} else {
			cluster := new(models.Cluster)
			if data := cluster.List(0, 0).Data; data != nil {
				items := data.(map[string]interface{})["items"].([]*models.Cluster)
				cluster = items[0]
				clientGo = models.CreateK8sClient(models.BuildApiParams(cluster))
			}
		}

		if clientGo.ErrMessage != "" {
			logs.Error("Delete network policy fail from kubernetes, err: %s", clientGo.ErrMessage)
			this.Data["json"] = models.Result{Code: utils.DeleteNetworkPolicyErr}
			this.ServeJSON(false)
		}
		netpolService := k8s.NetworkPolicyService{ClientGo: clientGo, NetworkPolicy: NetworkPolicy}
		err := netpolService.Delete()
		if err != nil {
			logs.Error("Delete network policy fail from kubernetes, err: %s", err)
			this.Data["json"] = models.Result{Code: utils.DeleteNetworkPolicyErr}
			this.ServeJSON(false)
		}
	}
	NetworkPolicy.Delete()
	this.Data["json"] = models.Result{Code: http.StatusOK}
	this.ServeJSON(false)
}
