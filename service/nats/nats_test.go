package nats

import (
	"github.com/nats-io/stan.go"
	"testing"
)

var (
	CLuster_Id = "diss-cluster"
	Client_Id  = "diss-agent22222"
	Nats_Url   = "nats://diss:dissP@ssw0rd@82.156.252.193:4222"
)

func Test_Nats_Pub(t *testing.T) {
	msg := `{
    "code":200,
    "msg":"Success",
    "type":"Metric",
    "tag":"HostImageVulnScan",
    "rc_type":"Put",
    "data":[
        {
            "HostId":"89e38324-e2bd-4c4b-abf3-d132f268187c",
            "ImageId":"22",
            "Target":"derekhjray/metrics:latest (ubuntu 18.04)",
            "Type":"ubuntu",
            "Vulnerabilities":[
                {
                    "VulnerabilityID":"CVE-2019-18276",
                    "PkgName":"bash",
                    "InstalledVersion":"4.4.18-2ubuntu1.2",
                    "Layer":{
                        "DiffID":"sha256:837d6facb613e572926fbfe8cd7171ddf5919c1454cf4d5b4e78f3d2a7729000"
                    },
                    "SeveritySource":"ubuntu",
                    "PrimaryURL":"https://avd.aquasec.com/nvd/cve-2019-18276",
                    "Title":"bash: when effective UID is not equal to its real UID the saved UID is not dropped",
                    "Description":"An issue was discovered in disable_priv_mode in shell.c in GNU Bash through 5.0 patch 11. ",
                    "Severity":"LOW",
                    "CweIDs":[
                        "CWE-273"
                    ],
                    "CVSS":{
                        "nvd":{
                            "V2Score":7.2,
                            "V3Score":7.8
                        },
                        "redhat":{
                            "V3Score":7.8
                        }
                    },
                    "References":[
                        "https://www.youtube.com/watch?v=-wGtxJ8opa8"
                    ],
                    "PublishedDate": 1574903700000000000,
                    "LastModifiedDate": 1614891840000000000
                }
            ]
        }
    ],
    "config":""
}`
	task := `{
  "code": 0,
  "msg": "",
  "type": "Control",
  "tag": "Task",
  "rc_type": "Put",
  "data": {
    "Id": "92c4e201-9ee1-4d59-863a-079855e2bcf2",
        "Account": "admin",
        "Name": "系统任务-92c4e201-9ee1-4d59-863a-079855e2bcf2",
        "Description": "系统任务-HostImageVulnScan",
        "Spec": "",
        "Type": "Once",
        "Status": "Running",
        "Batch": 1617952210669429500,
        "SystemTemplate": {
          "Id": "9b4320a4-57f5-4559-a91b-d84b0d0a0e92",
          "Account": "admin",
          "Name": "主机镜像漏洞扫描",
          "Description": "",
          "Type": "HostImageVulnScan",
          "Version": "1.0",
          "Commands": "",
          "Status": "Enable",
          "IsDefault": true,
          "SystemTemplateGroup": null,
          "Job": null,
          "Task": null,
          "ConfigMode": "",
          "DefaultTargets": "",
          "CheckMasterJson": "",
          "CheckNodeJson": "",
          "CheckControlPlaneJson": "",
          "CheckEtcdJson": "",
          "CheckPoliciesJson": "",
          "CheckManagedServicesJson": "",
          "CheckIdsMaster": "",
          "CheckIdsNode": "",
          "CheckIdsControlPlane": "",
          "CheckIdsEtcd": "",
          "CheckIdsPolicies": "",
          "CheckIdsManagedServices": "",
          "CheckIdsDocker": "",
          "CheckIdsDockerCheck": ""
        },
        "Host": null,
        "Container": null,
        "Image": {
          "Id": "c67483ff-34ad-4af4-95da-aca48ee5b1c7---sha256:24ac5e25ee2afe880e573f68494c9893d2b07ccb86158fb97289d71fd59651d2---falcosecurity/falco-driver-loader:latest",
          "ImageId": "sha256:24ac5e25ee2afe880e573f68494c9893d2b07ccb86158fb97289d71fd59651d2",
          "HostId": "c67483ff-34ad-4af4-95da-aca48ee5b1c7",
          "HostName": "node-179.com",
          "Name": "falcosecurity/falco-driver-loader:latest",
          "Size": "684.0MB",
          "OS": "",
          "DissStatus": 0,
          "Age": "2 月前",
          "CreateTime": 1610989290000000000,
          "DBType": "",
          "GetLatestTask": false,
          "TaskList": null,
          "Registry": null
        },
        "CreateTime": 1617952210673011500,
        "UpdateTime": 1617952210673011500,
        "Job": {
          "Id": "9b4320a4-57f5-4559-a91b-d84b0d0a0e05",
          "Account": "admin",
          "Name": "主机镜像漏洞扫描",
          "Description": "系统任务-主机镜像漏洞扫描",
          "Spec": "",
          "Type": "Once",
          "SystemTemplateType": "HostImageVulnScan",
          "Status": "Enable",
          "CreateTime": 1589795318000,
          "UpdateTime": 1589795318000,
          "SystemTemplate": {
            "Id": "9b4320a4-57f5-4559-a91b-d84b0d0a0e92",
            "Account": "admin",
            "Name": "主机镜像漏洞扫描",
            "Description": "",
            "Type": "HostImageVulnScan",
            "Version": "1.0",
            "Commands": "",
            "Status": "Enable",
            "IsDefault": true,
            "SystemTemplateGroup": null,
            "Job": null,
            "Task": null,
            "ConfigMode": "",
            "DefaultTargets": "",
            "CheckMasterJson": "",
            "CheckNodeJson": "",
            "CheckControlPlaneJson": "",
            "CheckEtcdJson": "",
            "CheckPoliciesJson": "",
            "CheckManagedServicesJson": "",
            "CheckIdsMaster": "",
            "CheckIdsNode": "",
            "CheckIdsControlPlane": "",
            "CheckIdsEtcd": "",
            "CheckIdsPolicies": "",
            "CheckIdsManagedServices": "",
            "CheckIdsDocker": "",
            "CheckIdsDockerCheck": ""
          },
          "SystemTemplateGroup": null,
          "Task": null,
          "HostConfig": null,
          "ContainerConfig": null,
          "JobLevel": "System",
          "IsUpdateHost": false
        },
        "ClusterId": "",
        "IsOne": false,
        "RunCount": 0,
        "Action": "",
        "ContainerHostId": ""
  },
  "config": ""
}`
	print(msg)
	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
	err = nc.Publish("Common-c67483ff-34ad-4af4-95da-aca48ee5b1c7", []byte(task))

	if err != nil {
		t.Errorf("Pub......... fail, err: %s", err)
	} else {
		t.Log("Pub......... success")
	}
}

//func Test_Nats_Sub(t *testing.T) {
//	subject := "89e38324-e2bd-4c4b-abf3-d132f268187c"
//	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
//	if err != nil {
//		t.Errorf("Client1 connet Nats server fail. err: %s", err)
//	}
//	_, err = nc.Subscribe(subject, func(m *stan.Msg) {
//		fmt.Println("22222222", string(m.Data))
//		fmt.Println("Client1 get a message:", string(m.Data))
//	}, stan.DurableName(subject))
//	if err != nil {
//		t.Errorf("Client1 get mesage fail. err: %s", err)
//	}
//}

//
//func Test_Nats_Sub2(t *testing.T) {
//	subject := "bogon_Task"
//	nc, err := stan.Connect(CLuster_Id, Client_Id, stan.NatsURL(Nats_Url))
//	if err != nil {
//		t.Errorf("Client2 connet Nats server fail. err: %s", err)
//	}
//	_, err = nc.Subscribe(subject, func(m *stan.Msg) {
//		fmt.Println("Client2 get a message:", string(m.Data))
//	}, stan.DurableName(subject))
//	if err != nil {
//		t.Errorf("Client2 get mesage fail. err: %s", err)
//	}
//}
