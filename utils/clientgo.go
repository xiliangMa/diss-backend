package utils

import (
	"flag"
	"github.com/astaxie/beego/logs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

type ClientGo struct {
	ClientSet *kubernetes.Clientset
	err error
}

func CreateK8sClient(kubeconfig, path, configName string) ClientGo {
	// 配置k8s集群 kubeconfig 配置文件
	if path != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(path, configName), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	logs.Info("path %s", filepath.Join(path, "config"))

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// 根据指定的 config 创建 clientset
	clientSet, err := kubernetes.NewForConfig(config)
	return ClientGo{clientSet, err}
}

func (clientgo *ClientGo) GetPodsByNameSpace(namespace string) (*v1.PodList, error) {
	return clientgo.ClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
}


func (clientgo *ClientGo) GetNodes() (*v1.NodeList, error) {
	return clientgo.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
}
