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
			{Key: "cn-hangzhou", Value: "https://registry.cn-hangzhou.aliyuncs.com"},
			{Key: "cn-shanghai", Value: "https://registry.cn-shanghai.aliyuncs.com"},
			{Key: "cn-qingdao", Value: "https://registry.cn-qingdao.aliyuncs.com"},
			{Key: "cn-beijing", Value: "https://registry.cn-beijing.aliyuncs.com"},
			{Key: "cn-zhangjiakou", Value: "https://registry.cn-zhangjiakou.aliyuncs.com"},
			{Key: "cn-huhehaote", Value: "https://registry.cn-huhehaote.aliyuncs.com"},
			{Key: "cn-shenzhen", Value: "https://registry.cn-shenzhen.aliyuncs.com"},
			{Key: "cn-chengdu", Value: "https://registry.cn-chengdu.aliyuncs.com"},
			{Key: "cn-hongkong", Value: "https://registry.cn-hongkong.aliyuncs.com"},
			{Key: "ap-southeast-1", Value: "https://registry.ap-southeast-1.aliyuncs.com"},
			{Key: "ap-southeast-2", Value: "https://registry.ap-southeast-2.aliyuncs.com"},
			{Key: "ap-southeast-3", Value: "https://registry.ap-southeast-3.aliyuncs.com"},
			{Key: "ap-southeast-5", Value: "https://registry.ap-southeast-5.aliyuncs.com"},
			{Key: "ap-northeast-1", Value: "https://registry.ap-northeast-1.aliyuncs.com"},
			{Key: "ap-south-1", Value: "https://registry.ap-south-1.aliyuncs.com"},
			{Key: "eu-central-1", Value: "https://registry.eu-central-1.aliyuncs.com"},
			{Key: "eu-west-1", Value: "https://registry.eu-west-1.aliyuncs.com"},
			{Key: "us-west-1", Value: "https://registry.us-west-1.aliyuncs.com"},
			{Key: "us-east-1", Value: "https://registry.us-east-1.aliyuncs.com"},
			{Key: "me-east-1", Value: "https://registry.me-east-1.aliyuncs.com"},
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
				Key:   "us-east-1",
				Value: "dkr.ecr.us-east-1.amazonaws.com",
			},
			{
				Key:   "us-east-2",
				Value: "dkr.ecr.us-east-2.amazonaws.com",
			},
			{
				Key:   "us-west-1",
				Value: "dkr.ecr.us-west-1.amazonaws.com",
			},
			{
				Key:   "us-west-2",
				Value: "dkr.ecr.us-west-2.amazonaws.com",
			},
			{
				Key:   "af-south-1",
				Value: "dkr.ecr.af-south-1.amazonaws.com",
			},
			{
				Key:   "ap-east-1",
				Value: "dkr.ecr.ap-east-1.amazonaws.com",
			},
			{
				Key:   "ap-south-1",
				Value: "dkr.ecr.ap-south-1.amazonaws.com",
			},
			{
				Key:   "ap-northeast-3",
				Value: "dkr.ecr.ap-northeast-3.amazonaws.com",
			},
			{
				Key:   "ap-northeast-2",
				Value: "dkr.ecr.ap-northeast-2.amazonaws.com",
			},
			{
				Key:   "ap-southeast-1",
				Value: "dkr.ecr.ap-southeast-1.amazonaws.com",
			},
			{
				Key:   "ap-southeast-2",
				Value: "dkr.ecr.ap-southeast-2.amazonaws.com",
			},
			{
				Key:   "ap-northeast-1",
				Value: "dkr.ecr.ap-northeast-1.amazonaws.com",
			},
			{
				Key:   "ca-central-1",
				Value: "dkr.ecr.ca-central-1.amazonaws.com",
			},
			{
				Key:   "eu-central-1",
				Value: "dkr.ecr.eu-central-1.amazonaws.com",
			},
			{
				Key:   "eu-west-1",
				Value: "dkr.ecr.eu-west-1.amazonaws.com",
			},
			{
				Key:   "eu-west-2",
				Value: "dkr.ecr.eu-west-2.amazonaws.com",
			},
			{
				Key:   "eu-south-1",
				Value: "dkr.ecr.eu-south-1.amazonaws.com",
			},
			{
				Key:   "eu-west-3",
				Value: "dkr.ecr.eu-west-3.amazonaws.com",
			},
			{
				Key:   "eu-north-1",
				Value: "dkr.ecr.eu-north-1.amazonaws.com",
			},
			{
				Key:   "me-south-1",
				Value: "dkr.ecr.me-south-1.amazonaws.com",
			},
			{
				Key:   "sa-east-1",
				Value: "dkr.ecr.sa-east-1.amazonaws.com",
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
