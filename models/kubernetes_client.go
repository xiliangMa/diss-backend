package models

import (
	"flag"
	"github.com/astaxie/beego/logs"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
	"net/http"
)

type KubernetesClientManager struct {
	ClientHub map[string]ClientGo
}

func NewKubernetesClientManager() *KubernetesClientManager {
	return &KubernetesClientManager{
		ClientHub: make(map[string]ClientGo),
	}
}

type ApiParams struct {
	KubeConfigPath string
	BearerToken    string
	MasterUrl      string
	AuthType       string
}

type ClientGo struct {
	ClientSet  *kubernetes.Clientset
	ErrMessage string
}

func CreateK8sClient(params *ApiParams) ClientGo {
	var config *rest.Config
	var err error
	clientgo := ClientGo{nil, ""}
	if params.AuthType == "BearerToken" {
		// 这里是使用用户名和密码调用APIserver，所以kubeconfig为空
		kubeconfig := flag.String(string(rand.Intn(1000)), "", "BearerToken")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags(params.MasterUrl, *kubeconfig)
		config.Insecure = true
		config.BearerToken = params.BearerToken
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", params.KubeConfigPath)
	}
	if err != nil {
		clientgo.ErrMessage = "BuildConfigFromFlags"
		return clientgo
	}
	// 根据指定的 config 创建 clientset
	clientSet, err := kubernetes.NewForConfig(config)
	clientgo.ClientSet = clientSet
	if err != nil {
		clientgo.ErrMessage = err.Error()
	}

	// 检测是否可用
	if _, err := clientgo.GetNodes(); err != nil {
		clientgo.ErrMessage = err.Error()
	}
	return clientgo
}

func (clientgo *ClientGo) GetNodes() (*v1.NodeList, error) {
	return clientgo.ClientSet.CoreV1().Nodes().List(v12.ListOptions{})
}

func InitClientHub() {
	// 获取所有集群
	cluster := new(Cluster)
	result := cluster.List(0, 0)

	if result.Code == http.StatusOK && result.Data != nil {
		data := result.Data.(map[string]interface{})
		clusterList := data["items"].([]*Cluster)
		for _, c := range clusterList {
			// 构建k8s客户端
			clientGo := CreateK8sClient(BuildApiParams(c))
			if clientGo.ErrMessage == "" {
				// 加入缓存
				KCHub.ClientHub[c.Id] = clientGo
				logs.Info("Init ClientHub add cluster： %s success.", c.Name)
			}
		}
	}
}

func BuildApiParams(c *Cluster) *ApiParams {
	params := new(ApiParams)
	params.AuthType = c.AuthType
	params.BearerToken = c.BearerToken
	params.MasterUrl = c.MasterUrls
	params.KubeConfigPath = c.FileName
	return params
}
