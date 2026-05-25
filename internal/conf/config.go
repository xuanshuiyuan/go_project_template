// @Author  xuanshuiyuan
// 全局常量配置：Redis Key 前缀、过期时间、上传限制、通用枚举等
package conf

const ParamsSignkey = "go_project_template"           // 接口签名密钥
const RedisTokenKey = "go_project_template_token_"     // 通用 Token 键前缀
const RedisAdminTokenKey = "go_project_template_admin_token_" // 管理端 Token 键前缀
const RedisWebTokenKey = "go_project_template_web_token_"     // Web 端 Token 键前缀
const RedisTokenExp = 7200                                // 登陆 token 有效时间，单位秒，默认 2 小时
const RedisLockExp = 30                                   // Redis 分布式锁超时时间，单位秒
const RedisWebRegisterKey = "1010KYXDFGCKV02LQPV9RG76"   // Web 端注册 Redis 锁键

const RedisPayOrderKey = "XBOF717ONVAP0GQC58NG_"     // Redis 支付锁键前缀
const RedisPayNotifyUrlKey = "1N6JW3FYHCDCYFBW7Z3A_" // Redis 支付回调锁键前缀

const RedisSignExp = 300 // 接口签名有效时间，单位秒，默认 5 分钟

const ErrorTips = "操作失败，系统异常" // 系统异常默认提示

// 允许上传的图片格式及大小限制
const AllowUploadImageFormat = ".png,.jpg,.jpeg,.gif"
const AllowUploadImageMaxSize = 20 << 20 // 20MB

// 允许上传的视频格式及大小限制
const AllowUploadVideoFormat = ".mp4,.wmv,.3gp,.mp4,.mov,.avi,.flv,.rmvb"
const AllowUploadVideoMaxSize = 20 << 20 // 20MB
const AllowUploadVideoMaxSizeTips = "请上传不超过20MB的视频哦"

// UploadCatalogueType OSS 上传目录映射，key 为业务类型编号
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
