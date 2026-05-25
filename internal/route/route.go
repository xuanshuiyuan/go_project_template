// @Author  xuanshuiyuan
package route

import "github.com/kataras/iris/v12/context"

func (s *Service) WebUtilsTest(c context.Context) {
	s.web.Utils.Test(c)
}
