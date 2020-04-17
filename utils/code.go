package utils

var (
	//diss-backend code 0-200
	Fail                            = 1
	GetHostListErr                  = 2
	AddHostErr                      = 3
	DeleteHostErr                   = 4
	GetHostZero                     = 5
	HostExistError                  = 6
	GetHostMetricError              = 7
	EditHostErr                     = 8
	SiginErr                        = 17
	AuthorizeErr                    = 18
	GetUserInfoErr                  = 19
	ElasticConnErr                  = 21
	ElasticSearchErr                = 22
	AddSYSTemplateErr               = 23
	GetSYSTemplateErr               = 24
	AddBenchMarkLogErr              = 25
	GetBenchMarkLogErr              = 26
	GetHostPsErr                    = 27
	AddHostPsErr                    = 28
	EditHostInfoErr                 = 29
	GetHostInfoErr                  = 30
	AddHostConfigErr                = 31
	EditHostConfigErr               = 32
	GetHostConfigErr                = 33
	EditContainerConfigErr          = 34
	AddImageConfigErr               = 35
	GetImageConfigErr               = 36
	AddContainerConfigErr           = 37
	GetContainerConfigErr           = 38
	EditImageConfigErr              = 39
	EditHostPsErr                   = 40
	GetContainerInfoErr             = 41
	AddContainerInfoErr             = 42
	EditContainerInfoErr            = 43
	DeleteHostPsErr                 = 44
	DeleteContainerInfoErr          = 45
	GetContainerPsErr               = 46
	AddContainerPsErr               = 47
	EditContainerPsErr              = 48
	DeleteContainerPsErr            = 49
	AddImageInfoErr                 = 50
	GetImageInfoErr                 = 51
	GetCmdHistoryErr                = 52
	AddCmdHistoryErr                = 53
	DeleteCmdHistoryErr             = 54
	DeleteImageInfoErr              = 55
	DeleteImageConfigErr            = 56
	AddGroupErr                     = 57
	GetGroupErr                     = 58
	EditGroupErr                    = 59
	DeleteGroupErr                  = 60
	GetAccountClusterErr            = 61
	NoccountClusterErr              = 62
	AddAccountClusterErr            = 63
	DeleteContainerConfigErr        = 64
	ContainerConfigNotFoundErr      = 65
	EmptyDirtyDataContinerConfigErr = 66
	EmptyDirtyDataContinerInfoErr   = 67
	EmptyDirtyDataNameSpaceErr      = 68
	EmptyDirtyDataPodErr            = 69
	BindContainerErr                = 70
	UnBindContainerErr              = 71
	DeleteSYSTemplateErr            = 72
	EditSYSTemplateErr              = 73

	//k8s 1001-1100
	AddNameSpaceErr    = 1001
	AddPodErr          = 1002
	AddClusterErr      = 1003
	GetClusterErr      = 1004
	GetPodErr          = 1005
	EditNameSpaceErr   = 1006
	EditPodErr         = 1007
	GetNameSpaceErr    = 1008
	IsBindErr          = 1009
	NoNameSpacedErr    = 1010
	NameSpaceExistErr  = 1011
	BindNameSpaceErr   = 1012
	UnBindNameSpaceErr = 1013
	BindPodErr         = 1014
	UnBindPodErr       = 1015

	//task 1100-1200
	AddTaskErr    = 1100
	GetTaskErr    = 1101
	DeleteTaskErr = 1102
	EditTaskErr   = 1103

	//system 1200-1300
	CheckK8sFilePostErr         = 1200
	CheckK8sFileCreateClientErr = 1201
	CheckK8sFileIsExistErr      = 1202
	UploadK8sFileErr            = 1203
	CheckK8sFileTestErr         = 1204
	UploadLogoErr               = 1205
	ChecLogoPostErr             = 1206
	CheckLogoIsNotExistErr      = 1207

	// timescaledb 安全日志1300-1400
	GetIntrudeDetectLogErr = 1300

	// diss_api 1400-1500
	GetAccountsErr     = 1400
	GetAccountUsersErr = 1401
	NoAccountUsersErr  = 1402
)
