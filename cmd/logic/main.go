// @Author xuanshuiyuan
package main

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"go_project_template/internal/route"
	"runtime"
	"time"
)

var log *goxy.Logs

func main() {
	app := iris.New()

	app.Use(loggerHandler, recoverHandler)
	app.OnErrorCode(iris.StatusNotFound)
	app.OnErrorCode(iris.StatusInternalServerError)
	app.OnErrorCode(iris.StatusNoContent)
	mvc.Configure(app.Party("/"), route.New)
	app.Run(iris.Addr(fmt.Sprintf(":%d", conf.Config.Iris.Port)))
}

func views(app *iris.Application) {
	pugEngine := iris.HTML(conf.Config.Iris.HtmlTemplate, ".html").Delims("<$", "$>")
	pugEngine.Reload(true)
	app.RegisterView(pugEngine)
}

func loggerHandler(c context.Context) {
	start := time.Now()
	path := c.Path()
	raw, _ := json.Marshal(c.FormValues())
	method := c.Method()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.GetStatusCode()
	addrs, _ := conf.GetLocalIP()
	str := fmt.Sprintf("METHOD:%s | PATH:%s | PARAMS:%s | CODE:%d | IP:%s | TIME:%d", method, path, raw, statusCode, addrs, latency/time.Millisecond)
	log.Info(conf.Config.Base.LogFileName, "logger.log").Println(goxy.FmtLog(str))
}

func recoverHandler(c context.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s", time.Now().Format("2006-01-02 15:04:05"), err)
			log.Error(conf.Config.Base.LogFileName, "recover.log").Println(goxy.FmtLog(pnc))
		}
	}()
	c.Next()
}
