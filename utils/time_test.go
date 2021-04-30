package utils

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"testing"
	"time"
)

func Test_SubTime(t *testing.T) {
	now := time.Now()
	createTime, _ := time.Parse(time.RFC3339Nano, "2019-12-13T01:45:33.000Z0700")
	subM := now.Sub(createTime)
	t.Log(int(subM.Minutes()), "Hours")
}

func Test_GetTimeFromNow(t *testing.T) {
	now := time.Now().Format("2006-01-02T15:04:05Z")
	//nowStr, _ := time.Parse(time.RFC3339, now.String())
	timepoint := time.Now().Add(time.Hour * -1).Format("2006-01-02T15:04:05Z")

	t.Log("\n Now:", now, "\n", "Timepoint:", timepoint)
}

func Test_Viper(t *testing.T) {
	config := viper.New()
	config.AddConfigPath(".")             //设置读取的文件路径
	config.SetConfigName("diss-kubescan") //设置读取的文件名
	config.SetConfigType("yaml")          //设置文件的类型
	//尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	//打印文件读取出来的内容:
	var lis []corev1.EnvVar
	var s1 corev1.EnvVar
	s1.Name = "KUBEHUNTER_HTTP_DISPATCH_URL"
	s1.Value = "https://10.28.147.11:10443/api/report/kubevulnscan"
	lis = append(lis, s1)
	var s2 corev1.EnvVar
	s2.Name = "TASK_ID"
	s2.Value = "task_id"
	lis = append(lis, s2)
	var s3 corev1.EnvVar
	s3.Name = "CLUSTER_ID"
	s3.Value = "cluster_id"
	lis = append(lis, s3)
	//config.Set("items.spec.template.spec.containers.env", lis)
	//fmt.Println(config.Get("items.spec.template.spec.containers.env"))

	bs, err := yaml.Marshal(config.AllSettings())
	logs.Error("===========", string(bs))
	list := new(corev1.List)
	err = json.Unmarshal(bs, &list)
	//logs.Error("------", job.Spec.Template.Spec.Containers[0].Env)
	logs.Error("------", err)
	//job.Spec.Template.Spec.Containers[0].Env = lis
	//fmt.Println("=====", job)

}
