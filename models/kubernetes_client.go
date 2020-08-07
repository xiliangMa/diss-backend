package models

import (
	"flag"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
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
