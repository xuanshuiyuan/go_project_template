// @Author  xuanshuiyuan
package logic

//ConfigConf go格式的配置文件
type ConfigConf struct {
	Env           string
	Oss           *OssConfig
	OssPath       string
	Verification  *Verification
	ApiTimer      *ApiTimer
	SnowflakeList *SnowflakeList
	LocalIP       string
	Base          *Base
	Adminer       *[]Adminer
	AdminMobile   *[]AdminMobile
}

type ApiTimer struct {
	Ip []string
}

//Verification 登陆验证的字段
type Verification struct {
	KeyList              map[string]string
	SourceList           []string
	SourceExplainList    map[string]string
	SourceRedisList      map[string]string
	SourceEngExplainList map[string]string
	EditionList          map[string]int64
}

//OssConfig 阿里OSS的配置属性
type OssConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	RoleArn         string
	BucketName      string
	TokenExpireTime int64
	RegionId        string
	RootPath        string
}

type SnowflakeList struct {
	Workerid     map[string]int64
	Datacenterid map[string]int64
}

type Base struct {
	RootUrl       string
	DockerBaseUrl string
}

type Adminer struct {
	UserId int64
}

type AdminMobile struct {
	Mobile string
}
