package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"os"
)

// 集群接口
type ClusterController struct {
	beego.Controller
}

// @Title AddCluster
// @Description Add Cluster
// @Param token header string true "authToken"
// @Param k8sFile formData file true "k8sFile"
// @Param clusterName formData string true "clusterName"
// @Param type formData string true "类型 Kubernetes Openshift Rancher"
// @Param isForce formData bool true true "force"
// @Success 200 {object} models.Result
// @router /add [post]
func (this *ClusterController) UploadK8sFile() {
	ClusterType := this.GetString("type")
	key := "k8sFile"
	isForce, _ := this.GetBool("isForce", true)
	clusterName := this.GetString("clusterName")
	f, h, _ := this.GetFile(key)
	result, fpath := css.Check(h)
	defer f.Close()
	if result.Code == http.StatusOK {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.UploadK8sFileErr
			result.Message = "UploadK8sFileErr"
			logs.Error("Upload k8s file  fail, err: %s", err.Error())
		} else {
			//检查连接是否可用
			if code := css.TestK8sFile(fpath); code != http.StatusOK {
				result.Code = code
				result.Message = "CheckK8sFileTestErr"
				os.Remove(fpath)
			} else {
				logs.Info("Upload k8s file success, file name: %s", h.Filename)
				// 添加集群记录
				css.AddCluster(clusterName, fpath, ClusterType)
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
				if code := css.TestK8sFile(fpath); code != http.StatusOK {
					result.Code = code
					result.Message = "CheckK8sFileTestErr"
					os.Remove(fpath)
				} else {
					logs.Info("Force update k8s file success, file name: %s", h.Filename)
					// to do 强制更新后文件名相同、内容不一样
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
