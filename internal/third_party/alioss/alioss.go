// @Author  xuanshuiyuan
package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"sync"
)

type Alioss struct {
	SimpleUpload *SimpleUploadService
}

var (
	log       *goxy.Logs
	ossBucket *oss.Bucket
	ossOnce   sync.Once
)

func NewAlioss() *Alioss {
	return &Alioss{
		SimpleUpload: NewSimpleUpload(),
	}
}

func NewBucket() (*oss.Bucket, error) {
	var initErr error
	ossOnce.Do(func() {
		client, err := oss.New(conf.Config.Conf.Oss.Endpoint, conf.Config.Conf.Oss.AccessKeyId, conf.Config.Conf.Oss.AccessKeySecret)
		if err != nil {
			log.Error(conf.Config.Base.LogFileName, "error.log").Println(goxy.FmtLog(err.Error()))
			initErr = fmt.Errorf("oss client init failed: %w", err)
			return
		}
		ossBucket, err = client.Bucket(conf.Config.Conf.Oss.BucketName)
		if err != nil {
			log.Error(conf.Config.Base.LogFileName, "error.log").Println(goxy.FmtLog(err.Error()))
			initErr = fmt.Errorf("oss bucket init failed: %w", err)
			return
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return ossBucket, nil
}
