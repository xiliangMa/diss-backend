package k8s

import (
	"encoding/json"
	"flag"
	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"math/rand"
	"time"
)

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

type PodMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp  time.Time `json:"timestamp"`
		Window     string    `json:"window"`
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	} `json:"items"`
}

type PodMetrics struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Timestamp  time.Time `json:"timestamp"`
	Window     string    `json:"window"`
	Containers []struct {
		Name  string `json:"name"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"containers"`
}

type NodeMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
		} `json:"metadata"`
		Timestamp time.Time `json:"timestamp"`
		Window    string    `json:"window"`
		Usage     struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"items"`
}

type NodeMetrics struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Timestamp time.Time `json:"timestamp"`
	Window    string    `json:"window"`
	Usage     struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"usage"`
}

func CreateK8sClient(params *ApiParams) ClientGo {
	var config *restclient.Config
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

func (clientgo *ClientGo) GetPodsByNameSpace(namespace string) (*v1.PodList, error) {
	return clientgo.ClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
}

func (clientgo *ClientGo) GetNodes() (*v1.NodeList, error) {
	return clientgo.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
}

func (clientgo *ClientGo) GetNameSpaces() (*v1.NamespaceList, error) {
	return clientgo.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
}

func (clientgo *ClientGo) GetPodLogs(namespace, pod string) *restclient.Request {
	return clientgo.ClientSet.CoreV1().Pods(namespace).GetLogs(pod, &v1.PodLogOptions{})
}

func (clientgo *ClientGo) GetJob(namespace, job string) (*batchv1.Job, error) {
	return clientgo.ClientSet.BatchV1().Jobs(namespace).Get(job, metav1.GetOptions{})
}

func (clientgo *ClientGo) CreateJobByYml(kubeBenchJobCommand []string, file, namespace string) (*batchv1.Job, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logs.Error("Read file: %s,  error:  %s", file, err)
	}
	var job *batchv1.Job
	err = yaml.Unmarshal(bytes, &job)
	// 更新 kube-bench  脚本
	if kubeBenchJobCommand != nil {
		job.Spec.Template.Spec.Containers[0].Command = kubeBenchJobCommand
	}
	if err != nil {
		logs.Error("Job Unmarshal error:  %s", err)
	}
	return clientgo.ClientSet.BatchV1().Jobs(namespace).Create(job)
}

func (clientgo *ClientGo) DeleteJob(namespace, jobName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return clientgo.ClientSet.BatchV1().Jobs(namespace).Delete(jobName, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}

func (clientgo *ClientGo) GetPodMetrics(namespace, podName string) (*PodMetrics, error) {
	url := "apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods/" + podName
	data, err := clientgo.ClientSet.RESTClient().Get().AbsPath(url).DoRaw()
	var podMetrics *PodMetrics
	if err != nil {
		logs.Error("Get pod: %s metrics err: %s", podName, err.Error())
	} else {
		err = json.Unmarshal(data, &podMetrics)
	}
	return podMetrics, err
}

func (clientgo *ClientGo) GetPodMetricsList(namespace string) (*PodMetricsList, error) {
	url := "apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods"
	data, err := clientgo.ClientSet.RESTClient().Get().AbsPath(url).DoRaw()
	var podMetricsList *PodMetricsList
	if err != nil {
		logs.Error("Get Pod metrics list err: %s", err.Error())
	} else {
		err = json.Unmarshal(data, &podMetricsList)
	}
	return podMetricsList, err
}

func (clientgo *ClientGo) GetNodeMetricsList() (*NodeMetricsList, error) {
	url := "apis/metrics.k8s.io/v1beta1/nodes"
	data, err := clientgo.ClientSet.RESTClient().Get().AbsPath(url).DoRaw()
	var nodeMetricsList *NodeMetricsList
	if err != nil {
		logs.Error("Get List metrics err: %s", err.Error())
	} else {
		err = json.Unmarshal(data, &nodeMetricsList)
	}
	return nodeMetricsList, err
}

func (clientgo *ClientGo) GetNodeMetrics(nodeName string) (*NodeMetrics, error) {
	url := "apis/metrics.k8s.io/v1beta1/nodes/" + nodeName
	data, err := clientgo.ClientSet.RESTClient().Get().AbsPath(url).DoRaw()
	var nodeMetrics *NodeMetrics
	if err != nil {
		logs.Error("Get List metrics err: %s", err.Error())
	} else {
		err = json.Unmarshal(data, &nodeMetrics)
	}
	return nodeMetrics, err
}

func (clientgo *ClientGo) CreateDaemonSetByYml(kubeBenchJobCommand []string, file, namespace string) (*appsv1.DaemonSet, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logs.Error("Read file: %s,  error:  %s", file, err)
	}
	var daemonSet *appsv1.DaemonSet
	err = yaml.Unmarshal(bytes, &daemonSet)

	if err != nil {
		logs.Error("DaemonSet Unmarshal error:  %s", err)
	}
	return clientgo.ClientSet.AppsV1().DaemonSets(namespace).Create(daemonSet)
}
