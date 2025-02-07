// @Author  xuanshuiyuan
package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"os"
	"path"
	"strings"
	"time"
)

type SimpleUploadService struct {
	bucket     *oss.Bucket
	file       string //单个文件
	catalogue  string //上传到oss的目录
	objectName string
}

func NewSimpleUpload() *SimpleUploadService {
	var s = &SimpleUploadService{
		bucket: NewBucket(),
	}
	return s
}

func (s *SimpleUploadService) SetFile(file string) *SimpleUploadService {
	s.file = file
	return s
}

func (s *SimpleUploadService) SetCatalogue(catalogue string) *SimpleUploadService {
	s.catalogue = catalogue
	return s
}

func (s *SimpleUploadService) SetObjectName(objectName string) *SimpleUploadService {
	s.objectName = objectName
	return s
}

// @Title UploadLocalFile
// @Description 上传文件到阿里云OSS
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return string, error
func (s *SimpleUploadService) UploadLocalFile() (string, error) {
	// 读取本地文件。
	fd, err := os.Open(s.file)
	defer fd.Close()
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return "", err
	}
	// 上传文件流。
	suffix := path.Ext(s.file)
	var filename = fmt.Sprintf("%s/%s/%s/%d-%s%s", conf.Config.Conf.Oss.RootPath, s.catalogue, goxy.YmdStr(), time.Now().Unix(), goxy.RandChar(5), suffix)
	err = s.bucket.PutObject(filename, fd)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return "", err
	}
	return filename, nil
}

// @Title UploadLocalFileNoName
// @Description 上传文件到阿里云OSS
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return string, error
func (s *SimpleUploadService) UploadLocalFileNoName() (string, error) {
	// 读取本地文件。
	fd, err := os.Open(s.file)
	defer fd.Close()
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return "", err
	}
	// 上传文件流。
	d := strings.Split(s.file, "/")
	var filename = fmt.Sprintf("%s/%s/%s/%s", conf.Config.Conf.Oss.RootPath, s.catalogue, goxy.YmdStr(), d[len(d)-1])
	err = s.bucket.PutObject(filename, fd)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return "", err
	}
	return filename, nil
}

// @Title IsExist
// @Description 检查文件是够存在
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param
// @Return bool, error
func (s *SimpleUploadService) IsExist() (bool, error) {
	isExist, err := s.bucket.IsObjectExist(s.objectName)
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return false, err
	}
	return isExist, nil
}

// @Title DeleteObjects
// @Description 删除文件
// @Author xuanshuiyuan 2021-10-22 17:14:47
// @Param objectNames
// @Return bool, error
func (s *SimpleUploadService) DeleteObjects(objectNames []string) (bool, error) {
	// 不返回删除的结果。
	_, err := s.bucket.DeleteObjects(objectNames, oss.DeleteObjectsQuiet(true))
	if err != nil {
		log.Error(conf.Config.Base.LogFileName, "alioss.log").Println(goxy.FmtLog(err.Error()))
		return false, err
	}
	return true, nil
}
