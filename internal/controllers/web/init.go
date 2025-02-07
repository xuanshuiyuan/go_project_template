// @Author xuanshuiyuan
package web

import (
	"github.com/xuanshuiyuan/goxy"
)

var result *goxy.IrisHttpResult
var log *goxy.Logs

type WebService struct {
	Utils *Utils
}

//初始化
func NewWeb() *WebService {
	web := &WebService{
		Utils: NewUtils(),
	}
	return web
}
