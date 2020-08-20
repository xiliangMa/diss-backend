package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/k8s"
	"github.com/xiliangMa/diss-backend/service/securitycheck"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"os"
	"time"
)

// 集群接口
type ClusterController struct {
	beego.Controller
}

// @Title AddCluster
// @Description Add Cluster "authType=KubeConfig, 需要 上传KubeConfig文件; authType=BearerToken， 需要设置 masterUrl、bearerToken 参数"
// @Param token header string true "authToken"
// @Param authType formData string true "default: BearerToken、KubeConfig"
// @Param masterUrl formData string false "ApiServer 访问地址"
// @Param bearerToken formData string false "ApiServer 访问token"
// @Param k8sFile formData file false "KubeConfig 文件"
// @Param clusterName formData string true "clusterName"
// @Param type formData string true "类型 Kubernetes Openshift Rancher"
// @Param isForce formData bool false true "force"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *ClusterController) AddCluster() {
	authType := this.GetString("authType")
	masterUrl := this.GetString("masterUrl")
	bearerToken := this.GetString("bearerToken")
	clusterType := this.GetString("type")
	isForce, _ := this.GetBool("isForce", true)
	clusterName := this.GetString("clusterName")

	cluster := models.Cluster{}
	uid, _ := uuid.NewV4()
	cluster.Id = uid.String()
	cluster.Name = clusterName
	cluster.Type = clusterType
	cluster.Status = models.Cluster_Status_Active
	cluster.SyncStatus = models.Cluster_Sync_Status_NotSynced
	cluster.IsSync = models.Cluster_IsSync
	cluster.AuthType = authType
	if authType == models.Api_Auth_Type_BearerToken {
		cluster.MasterUrls = masterUrl
		cluster.BearerToken = bearerToken
	}

	params := models.ApiParams{}
	params.AuthType = authType
	params.BearerToken = bearerToken
	params.MasterUrl = masterUrl

	result := models.Result{Code: http.StatusOK}
	fpath := ""

	switch authType {
	//BearerToken 方式
	case models.Api_Auth_Type_BearerToken:
		if result = css.TestClient(params); result.Code == http.StatusOK {
			if result = cluster.Add(isForce); result.Code == http.StatusOK {
				logs.Info("Add Cluster success, MasterUrl: %s", cluster.MasterUrls)
				if result.Code == http.StatusOK {
					// 启动watch
					k8sWatchService := k8s.K8sWatchService{Cluster: &cluster}
					k8sWatchService.WatchCluster()
				}
			}
		}
	case models.Api_Auth_Type_KubeConfig:
		// KubeConfig 方式
		key := "k8sFile"
		f, h, _ := this.GetFile(key)
		result, fpath = css.Check(h)
		defer f.Close()
		if result.Code == http.StatusOK {
			err := this.SaveToFile(key, fpath)
			if err != nil {
				result.Code = utils.UploadK8sFileErr
				result.Message = "UploadK8sFileErr"
				logs.Error("Upload KubeConfig file  fail, err: %s", err.Error())
			} else {
				//检查连接是否可用
				params.KubeConfigPath = fpath
				if result = css.TestClient(params); result.Code != http.StatusOK {
					os.Remove(fpath)
				} else {
					logs.Info("Upload KubeConfig file success, file name: %s", h.Filename)
					// 添加集群记录
					cluster.FileName = fpath
					result = cluster.Add(isForce)

					if result.Code == http.StatusOK {
						// 启动watch
						k8sWatchService := k8s.K8sWatchService{Cluster: &cluster}
						k8sWatchService.WatchCluster()
					}

				}
			}
		} else {
			// 强制更新
			if result.Code == utils.CheckK8sFileIsExistErr && isForce {
				// 更新返回值
				result.Code = http.StatusOK
				result.Message = ""

				// 更新上传文件path
				fpath = fpath + h.Filename

				// 删除旧文件
				os.Remove(fpath)

				err := this.SaveToFile(key, fpath)
				if err != nil {
					result.Code = utils.UploadK8sFileErr
					result.Message = err.Error()
				} else {
					//检查连接是否可用
					params.KubeConfigPath = fpath
					if result = css.TestClient(params); result.Code != http.StatusOK {
						os.Remove(fpath)
					} else {
						cluster.FileName = fpath
						logs.Info("Force update k8s file success, file name: %s", h.Filename)
						result = cluster.Add(isForce)
						// to do 强制更新后文件名相同、内容不一样
					}
				}
			}
		}
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title GetClusters
// @Description Get Cluster List(不支持租户查询)
// @Param token header string true "authToken"
// @Param body body models.Cluster false "集群"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *ClusterController) GetClusters() {
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	cluster := new(models.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	this.Data["json"] = cluster.List(from, limit)
	this.ServeJSON(false)
}

// @Title UpdateCluster
// @Description Update Cluster
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Param body body models.Cluster true "集群"
// @Success 200 {object} models.Result
// @router /:id [put]
func (this *ClusterController) UpdateCluster() {
	id := this.GetString(":id")
	cluster := new(models.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	cluster.Id = id
	this.Data["json"] = cluster.Update()
	this.ServeJSON(false)
}

// @Title DeleteCluster
// @Description delete Cluster
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *ClusterController) DeleteCluster() {
	id := this.GetString(":id")
	clusterList := []*models.Cluster{}
	cluster := new(models.Cluster)
	cluster.Id = id
	clusterList = append(clusterList, cluster)
	k8sClearService := k8s.K8sClearService{ClusterList: clusterList, DropCluster: true}
	go k8sClearService.ClearAll()
	this.Data["json"] = models.Result{Code: http.StatusOK}
	this.ServeJSON(false)
}

// @Title SyncCluster
// @Description Sync Cluster
// @Param token header string true "authToken"
// @Param id path string "" true "Id"
// @Success 200 {object} models.Result
// @router /:id/sync [post]
func (this *ClusterController) SyncCluster() {
	id := this.GetString(":id")
	cluster := new(models.Cluster)
	json.Unmarshal(this.Ctx.Input.RequestBody, &cluster)
	cluster.Id = id
	syncCheckPoint := time.Now().Unix()
	K8sSyncService := k8s.NewK8sSyncService(syncCheckPoint, cluster)
	this.Data["json"] = K8sSyncService.Sync()
	this.ServeJSON(false)
}


// @Title ClusterSecurityCheck
// @Description ClusterSecurityCheck
// @Param token header string true "authToken"
// @Param docker query bool false "docker cis"
// @Param kubernetes query bool false "kubernetes cis"
// @Param id path string "" true "Id"
// @Success 200 {object} models.Result
// @router /:id/securitycheck [post]
func (this *ClusterController) ClusterSecurityCheck() {
	id := this.GetString(":id")
	docker := this.GetString("docker")
	kubernetes := this.GetString("kubernetes")
	cluster := new(models.Cluster)
	cluster.Id = id
	batch := time.Now().Unix()
	securityCheckService := securitycheck.SecurityCheckService{SecurityCheckList: checkList, Batch: batch}
	this.Data["json"] = securityCheckService.DeliverTask(true)
	this.ServeJSON(false)
	this.Data["json"] = K8sSyncService.Sync()
	this.ServeJSON(false)
}
