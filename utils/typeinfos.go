package utils

type endpoint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EndpointPattern struct {
	Endpoints []*endpoint `json:"endpoints"`
}

func ACRRegion() *EndpointPattern {
	ep := &EndpointPattern{
		Endpoints: []*endpoint{
			{Key: "cn-hangzhou", Value: "https://registrys.cn-hangzhou.aliyuncs.com"},
			{Key: "cn-shanghai", Value: "https://registrys.cn-shanghai.aliyuncs.com"},
			{Key: "cn-qingdao", Value: "https://registrys.cn-qingdao.aliyuncs.com"},
			{Key: "cn-beijing", Value: "https://registrys.cn-beijing.aliyuncs.com"},
			{Key: "cn-zhangjiakou", Value: "https://registrys.cn-zhangjiakou.aliyuncs.com"},
			{Key: "cn-huhehaote", Value: "https://registrys.cn-huhehaote.aliyuncs.com"},
			{Key: "cn-shenzhen", Value: "https://registrys.cn-shenzhen.aliyuncs.com"},
			{Key: "cn-chengdu", Value: "https://registrys.cn-chengdu.aliyuncs.com"},
			{Key: "cn-hongkong", Value: "https://registrys.cn-hongkong.aliyuncs.com"},
			{Key: "ap-southeast-1", Value: "https://registrys.ap-southeast-1.aliyuncs.com"},
			{Key: "ap-southeast-2", Value: "https://registrys.ap-southeast-2.aliyuncs.com"},
			{Key: "ap-southeast-3", Value: "https://registrys.ap-southeast-3.aliyuncs.com"},
			{Key: "ap-southeast-5", Value: "https://registrys.ap-southeast-5.aliyuncs.com"},
			{Key: "ap-northeast-1", Value: "https://registrys.ap-northeast-1.aliyuncs.com"},
			{Key: "ap-south-1", Value: "https://registrys.ap-south-1.aliyuncs.com"},
			{Key: "eu-central-1", Value: "https://registrys.eu-central-1.aliyuncs.com"},
			{Key: "eu-west-1", Value: "https://registrys.eu-west-1.aliyuncs.com"},
			{Key: "us-west-1", Value: "https://registrys.us-west-1.aliyuncs.com"},
			{Key: "us-east-1", Value: "https://registrys.us-east-1.aliyuncs.com"},
			{Key: "me-east-1", Value: "https://registrys.me-east-1.aliyuncs.com"},
		},
	}
	return ep
}

func SWRRegion() *EndpointPattern {
	ep := &EndpointPattern{
		Endpoints: []*endpoint{
			{Key: "af-south-1", Value: "af-south-1.myhuaweicloud.com"},
			{Key: "ap-southeast-1", Value: "ap-southeast-1.myhuaweicloud.com"},
			{Key: "ap-southeast-2", Value: "ap-southeast-2.myhuaweicloud.com"},
			{Key: "ap-southeast-3", Value: "ap-southeast-3.myhuaweicloud.com"},
			{Key: "cn-east-2", Value: "cn-east-2.myhuaweicloud.com"},
			{Key: "cn-east-3", Value: "cn-east-3.myhuaweicloud.com"},
			{Key: "cn-north-1", Value: "cn-north-1.myhuaweicloud.com"},
			{Key: "cn-north-2", Value: "cn-north-2.myhuaweicloud.com"},
			{Key: "cn-north-4", Value: "cn-north-4.myhuaweicloud.com"},
			{Key: "cn-south-1", Value: "cn-south-1.myhuaweicloud.com"},
			{Key: "cn-south-2", Value: "cn-south-2.myhuaweicloud.com"},
			{Key: "cn-southwest-2", Value: "cn-southwest-2.myhuaweicloud.com"},
			{Key: "ru-northwest-2", Value: "ru-northwest-2.myhuaweicloud.com"},
		},
	}
	return ep
}

func AWSRegion() *EndpointPattern {
	ep := &EndpointPattern{
		Endpoints: []*endpoint{
			{
				Key:   "ap-northeast-1",
				Value: "https://api.ecr.ap-northeast-1.amazonaws.com",
			},
			{
				Key:   "us-east-1",
				Value: "https://api.ecr.us-east-1.amazonaws.com",
			},
			{
				Key:   "us-east-2",
				Value: "https://api.ecr.us-east-2.amazonaws.com",
			},
			{
				Key:   "us-west-1",
				Value: "https://api.ecr.us-west-1.amazonaws.com",
			},
			{
				Key:   "us-west-2",
				Value: "https://api.ecr.us-west-2.amazonaws.com",
			},
			{
				Key:   "ap-east-1",
				Value: "https://api.ecr.ap-east-1.amazonaws.com",
			},
			{
				Key:   "ap-south-1",
				Value: "https://api.ecr.ap-south-1.amazonaws.com",
			},
			{
				Key:   "ap-northeast-2",
				Value: "https://api.ecr.ap-northeast-2.amazonaws.com",
			},
			{
				Key:   "ap-southeast-1",
				Value: "https://api.ecr.ap-southeast-1.amazonaws.com",
			},
			{
				Key:   "ap-southeast-2",
				Value: "https://api.ecr.ap-southeast-2.amazonaws.com",
			},
			{
				Key:   "ca-central-1",
				Value: "https://api.ecr.ca-central-1.amazonaws.com",
			},
			{
				Key:   "eu-central-1",
				Value: "https://api.ecr.eu-central-1.amazonaws.com",
			},
			{
				Key:   "eu-west-1",
				Value: "https://api.ecr.eu-west-1.amazonaws.com",
			},
			{
				Key:   "eu-west-2",
				Value: "https://api.ecr.eu-west-2.amazonaws.com",
			},
			{
				Key:   "eu-west-3",
				Value: "https://api.ecr.eu-west-3.amazonaws.com",
			},
			{
				Key:   "eu-north-1",
				Value: "https://api.ecr.eu-north-1.amazonaws.com",
			},
			{
				Key:   "sa-east-1",
				Value: "https://api.ecr.sa-east-1.amazonaws.com",
			},
			{
				Key:   "cn-north-1",
				Value: "https://api.ecr.cn-north-1.amazonaws.com.cn",
			},
			{
				Key:   "cn-northwest-1",
				Value: "https://api.ecr.cn-northwest-1.amazonaws.com.cn",
			},
		},
	}
	return ep
}

func DockerHubRegion() *EndpointPattern {
	ep := &EndpointPattern{
		Endpoints: []*endpoint{
			{
				Key:   "hub.docker.com",
				Value: "https://hub.docker.com",
			},
		},
	}
	return ep
}

func GoogleRegion() *EndpointPattern {
	ep := &EndpointPattern{
		Endpoints: []*endpoint{
			{
				Key:   "gcr.io",
				Value: "https://gcr.io",
			},
			{
				Key:   "us.gcr.io",
				Value: "https://us.gcr.io",
			},
			{
				Key:   "eu.gcr.io",
				Value: "https://eu.gcr.io",
			},
			{
				Key:   "asia.gcr.io",
				Value: "https://asia.gcr.io",
			},
		},
	}
	return ep
}
