// @Author  xuanshuiyuan
package conf

const ParamsSignkey = "go_project_template"
const RedisTokenKey = "go_project_template_"
const RedisAdminTokenKey = "go_project_template_"
const RedisWebTokenKey = "go_project_template_"
const RedisTokenExp = "7200"                           //登陆token的有效时间 2小时
const RedisLockExp = "30"                              //redis分布式锁 30秒
const RedisWebRegisterKey = "1010KYXDFGCKV02LQPV9RG76" //web端注册redis锁key

const RedisPayOrderKey = "XBOF717ONVAP0GQC58NG_"     //redis 支付锁
const RedisPayNotifyUrlKey = "1N6JW3FYHCDCYFBW7Z3A_" //redis 支付回调锁

const RedisSignExp = "300" //接口sign的有效时间 5分钟

const ErrorTips = "操作失败，系统异常"

//var Mutex sync.Mutex
//var RWMutex sync.RWMutex

//允许上传文件格式
const AllowUploadImageFormat = ".png,.jpg,.jpeg,.gif"
const AllowUploadImageMaxSize = 20 << 20 //20m

const AllowUploadVdeioFormat = ".mp4,.wmv,.3gp,.mp4,.mov,.avi,.flv,.rmvb"
const AllowUploadVedioMaxSize = 20 << 20 //20m
const AllowUploadVedioMaxSizeTips = "请上传不超过20MB的视频哦"

//oss上传目录
var UploadCatalogueType = map[int8]string{
	1: "default", //默认
}

type OptionFormat struct {
	Key   int8   `json:"key"`
	Value string `json:"value"`
}

var CommonStatus = map[int8]string{
	1: "启用",
	2: "禁用",
}

var CommonIsStatus = map[int8]string{
	1: "是",
	2: "否",
}
