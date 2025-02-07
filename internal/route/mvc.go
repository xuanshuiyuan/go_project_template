// @Author  xuanshuiyuan
package route

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"go_project_template/internal/controllers/web"
)

var result *goxy.IrisHttpResult
var log *goxy.Logs

// Service 模块对象 定义了不同客户端的接口
type Service struct {
	web *web.WebService
}

func Init() *Service {
	service := &Service{
		web: web.NewWeb(),
	}
	return service
}

func New(app *mvc.Application) {
	conf.ConfInit() //初始化配置文件
	//engine.NewEngine() //初始化数据库
	app.Handle(Init())
}

func (s *Service) BeforeActivation(b mvc.BeforeActivation) {
	//业务逻辑接口-需要验证签名
	s.initApiRouter(b)
	//B端逻辑接口-需要验证签名
	s.initAdminRouter(b)
	//C端逻辑接口-需要验证签名
	s.initWebRouter(b)
	//C端逻辑接口-不需要验证token
	s.initNotVerificationWebRouter(b)
	//B端逻辑接口-不需要验证token
	s.initNotVerificationAdminRouter(b)
	//计划任务
	s.initApiTimerRouter(b)
	//DEBUG 调试接口
	s.initDebugApiRouter(b)

}

//API
func (s *Service) initApiRouter(b mvc.BeforeActivation) {

}

//ADMIN
func (s *Service) initAdminRouter(b mvc.BeforeActivation) {

}

//WEB
func (s *Service) initWebRouter(b mvc.BeforeActivation) {
}

func (s *Service) initNotVerificationAdminRouter(b mvc.BeforeActivation) {
}

func (s *Service) initNotVerificationWebRouter(b mvc.BeforeActivation) {
	b.Handle("GET", "/web/utils/test", "WebUtilsTest") //
}

func (s *Service) initDebugApiRouter(b mvc.BeforeActivation) {
}

func (s *Service) initApiTimerRouter(b mvc.BeforeActivation) {
}
