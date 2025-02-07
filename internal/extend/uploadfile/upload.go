// @Author xuanshuiyuan
package uploadfile

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/third_party/alioss"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var log *goxy.Logs

type UploadEr interface {
	Upload() ([]string, error)
	Delete() error
}

type OssEr interface {
	Upload(file string, catalogue string) (string, error)
	Delete(file []string) error
}

func DefaultOss() *AliOss {
	return NewAliOss()
}

//上传到本地
type Local struct {
}

func (l Local) Delete(file []string) error {
	panic("implement me")
}

type AliOss struct {
	upload *alioss.Alioss
}

func NewAliOss() *AliOss {
	return &AliOss{
		upload: alioss.NewAlioss(),
	}
}

type UploadService struct {
	UploadParams
	File File
}

type UploadParams struct {
	Ctx       context.Context //为空则上传本地文件
	Oss       OssEr           //文件需要上传的云空间
	FilePath  []string        //上传的本地文件路径
	FileName  string          //MultipartForm 字段名称
	Catalogue string          //oss的目录
}

type File struct {
	AllowUploadMaxSize          int64
	AllowUploadMaxSizeErrorTips string
	AllowUploadFormat           string
	AllowUploadFormatErrorTips  string
}

func LoadDefaultConfig() *UploadService {
	return &UploadService{}
}

type Image struct {
	UploadService
}

func LoadImageConfig() *UploadService {
	return &UploadService{
		File: File{
			AllowUploadMaxSize:          20 << 20, //20m
			AllowUploadMaxSizeErrorTips: "请上传不超过20MB的图片哦",
			AllowUploadFormat:           ".png,.jpg,.jpeg,.gif",
			AllowUploadFormatErrorTips:  "文件格式错误,允许上传的图片格式为：",
		},
	}
}

type Video struct {
	UploadService
}

func LoadVideoConfig() *UploadService {
	return &UploadService{
		File: File{
			AllowUploadMaxSize:          20 << 20, //20m
			AllowUploadMaxSizeErrorTips: "请上传不超过20MB的视频哦",
			AllowUploadFormat:           ".mp4,.wmv,.3gp,.mp4,.mov,.avi,.flv,.rmvb",
			AllowUploadFormatErrorTips:  "文件格式错误,允许上传的视频格式为：",
		},
	}
}

type Docs struct {
	UploadService
}

func LoadDocsConfig() *UploadService {
	return &UploadService{
		File: File{
			AllowUploadMaxSize:          20 << 20, //20m
			AllowUploadMaxSizeErrorTips: "请上传不超过20MB的文档哦",
			AllowUploadFormat:           ".docx,.pdf,.xlsx",
			AllowUploadFormatErrorTips:  "文件格式错误,允许上传的文档格式为：",
		},
	}
}

func (u *UploadService) InjectParams(s UploadParams) *UploadService {
	u.UploadParams = s
	if u.FileName == "" {
		u.FileName = "file"
	}
	if u.Oss == nil {
		u.Oss = DefaultOss()
	} else {
		_, ok := interface{}(u.Oss).(*Local)
		if ok {
			u.Oss = nil
		}
	}
	return u
}

func (u *UploadService) uploadToLocal() (result []string, err error) {
	file_path := fmt.Sprint(fmt.Sprintf("%s", goxy.StaticFileDirectory()))
	goxy.CheckDir(file_path)
	form := u.Ctx.Request().MultipartForm
	headers := form.File[u.FileName]
	for _, fh := range headers {
		if fh.Size > u.File.AllowUploadMaxSize {
			return result, errors.New(u.File.AllowUploadMaxSizeErrorTips)
		}
		suffix := path.Ext(fh.Filename)
		if !strings.Contains(u.File.AllowUploadFormat, suffix) {
			return result, errors.New(u.File.AllowUploadFormatErrorTips + u.File.AllowUploadFormat)
		}
		fh.Filename = fmt.Sprintf("%d-%s%s", time.Now().Unix(), goxy.RandChar(5), suffix)
		_, err = saveUploadedFile(fh, file_path)
		if err != nil {
			return result, err
		}
		result = append(result, file_path+"/"+fh.Filename)
	}
	return result, nil
}

func (u *UploadService) Upload() (result []string, err error) {
	filepath := []string{}
	if u.Ctx != nil {
		filepath, err = u.uploadToLocal()
		if err != nil {
			return
		}
	} else {
		filepath = u.FilePath
	}
	if u.Oss == nil {
		return filepath, nil
	} else {
		if len(filepath) == 0 {
			return filepath, nil
		}
		for _, v := range filepath {
			f, err := u.Oss.Upload(v, u.Catalogue)
			if err != nil {
				return result, err
			}
			result = append(result, f)
		}
	}
	return result, nil
}

func (u *UploadService) Delete() (err error) {
	if u.Oss == nil { //删除本地
		for _, v := range u.FilePath {
			os.Remove(v)
		}
	} else {
		err = u.Oss.Delete(u.FilePath)
	}
	return
}

func (l Local) Upload(file string, catalogue string) (string, error) {
	panic("implement me")
}

func Upload(u UploadEr) ([]string, error) {
	return u.Upload()
}

func Delete(u UploadEr) error {
	return u.Delete()
}

func (a *AliOss) Delete(file []string) (err error) {
	fmt.Println(file)
	_, err = a.upload.SimpleUpload.DeleteObjects(file)
	if err != nil {
		return
	}
	return
}

func (a *AliOss) Upload(file string, catalogue string) (result string, err error) {
	result, err = a.upload.SimpleUpload.SetFile(file).SetCatalogue(catalogue).UploadLocalFile()
	if err != nil {
		return
	}
	osErr := os.Remove(file)
	if osErr != nil {
		// 删除失败
	}
	return
}

// @Title saveUploadedFile
// @Description 保存文件
func saveUploadedFile(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, src)
}
