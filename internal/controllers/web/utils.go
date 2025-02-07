// @Author xuanshuiyuan 2025/2/7 10:51:00
package web

import (
	"github.com/kataras/iris/v12/context"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/service/web_service"
)

type Utils struct {
	utils *web_service.UtilsService
}

func NewUtils() *Utils {
	return &Utils{
		utils: web_service.NewUtils(),
	}
}

func (u *Utils) Test(c context.Context) {
	result.Echo(c, goxy.StatusOK, "")
}
