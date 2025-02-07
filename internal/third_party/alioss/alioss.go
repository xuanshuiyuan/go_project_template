// @Author  xuanshuiyuan
package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
)

type Alioss struct {
	SimpleUpload *SimpleUploadService
}

var (
	//Client *base
	log       *goxy.Logs
	ossBucket *oss.Bucket
)

func NewAlioss() *Alioss {
	return &Alioss{
		SimpleUpload: NewSimpleUpload(),
	}
}

// @Title NewBucket
// @Description 初始化阿里云OSS
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return *oss.Bucket
func NewBucket() *oss.Bucket {
	if ossBucket != nil {
		return ossBucket
	}
	// 创建OSSClient实例。
	client, err := oss.New(conf.Config.Conf.Oss.Endpoint, conf.Config.Conf.Oss.AccessKeyId, conf.Config.Conf.Oss.AccessKeySecret)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "error.log").Println(goxy.FmtLog(err.Error()))
		panic(err.Error())
	}
	// 获取存储空间。
	ossBucket, err = client.Bucket(conf.Config.Conf.Oss.BucketName)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "error.log").Println(goxy.FmtLog(err.Error()))
		panic(err.Error())
	}
	return ossBucket
}
