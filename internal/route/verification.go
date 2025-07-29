// @Author  xuanshuiyuan 2021/12/28 17:18
package route

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"go_project_template/internal/utils"
	"net"
	"os"
	"strconv"
	"time"
)

const RequestHeaderParamsDevelop = "develop"
const RequestFailTip = "Im Invalid Request!"

//RequestHeaderParams B端的接口请求头部字段
type RequestHeaderParams struct {
	Token     string `json:"token"`
	Sign      string `json:"sign"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
	Userid    int64  `json:"userid"`
	Develop   string `json:"develop"`
	Random    string `json:"random"`
	Edition   string `json:"edition"`
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ip_net, ok := address.(*net.IPNet); ok && !ip_net.IP.IsLoopback() {
			if ip_net.IP.To4() != nil {
				return ip_net.IP.String()
			}
		}
	}
	return ""
}

// @Title ApiTimerVerification
// @Description 计划任务接口验证,内网才能访问
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) ApiTimerVerification(c context.Context) {
	ip := getLocalIP()
	next := false
	for _, v := range conf.Config.Conf.ApiTimer.Ip {
		if ip == v {
			next = true
		}
	}
	//fmt.Println(ip)
	if next == true {
		c.Next()
	} else {
		result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
		return
	}
}

// @Title WebLoginVerification
// @Description 登陆接口验证
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) WebLoginVerification(c context.Context) {
	if c.GetHeader("source") != "1010" && c.GetHeader("source") != "0110" {
		result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
		return
	}
	if err := s.verification(c); err != nil {
		result.Error(c, goxy.StatusValidationFailed, err.Error())
		return
	}
	c.Next()
}

// @Title AdminLoginVerification
// @Description 登陆接口验证
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) AdminLoginVerification(c context.Context) {
	if c.GetHeader("source") != "0101" {
		result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
		return
	}
	if err := s.verification(c); err != nil {
		result.Error(c, goxy.StatusValidationFailed, err.Error())
		return
	}
	c.Next()
}

// @Title WebVerification
// @Description 接口验证
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) WebVerification(c context.Context) {
	if c.GetHeader("source") != "1010" && c.GetHeader("source") != "0110" {
		result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
		return
	}
	if err := s.verification(c); err != nil {
		result.Error(c, goxy.StatusValidationFailed, err.Error())
		return
	}
	requestHeaderParamsKey := [...]string{"token", "userid"}
	for _, v := range requestHeaderParamsKey {
		if c.GetHeader(v) == "" {
			result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
			return
		}
	}
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	requestHeaderParams := RequestHeaderParams{
		Token:  c.GetHeader("token"),
		Userid: user_id,
	}
	key := utils.GetTokenKey(conf.RedisWebTokenKey, requestHeaderParams.Userid, c.GetHeader("source"))
	//redis验证token
	redis := utils.NewRedis()
	err := redis.RedisVerification(key, requestHeaderParams.Token)
	if err != nil { //登录失效，需重新登陆获取token
		result.Error(c, goxy.StatusTokenExpired, "登陆失效")
		return
	}
	c.Next()
}

// @Title AdminVerification
// @Description 接口验证
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) AdminVerification(c context.Context) {
	if c.GetHeader("source") != "0101" {
		result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
		return
	}
	if err := s.verification(c); err != nil {
		result.Error(c, goxy.StatusValidationFailed, err.Error())
		return
	}
	requestHeaderParamsKey := [...]string{"token", "userid"}
	for _, v := range requestHeaderParamsKey {
		if c.GetHeader(v) == "" {
			result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
			return
		}
	}
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	requestHeaderParams := RequestHeaderParams{
		Token:  c.GetHeader("token"),
		Userid: user_id,
	}
	key := utils.GetTokenKey(conf.RedisAdminTokenKey, requestHeaderParams.Userid, c.GetHeader("source"))
	//redis验证token
	redis := utils.NewRedis()
	err := redis.RedisVerification(key, requestHeaderParams.Token)
	if err != nil { ////登录失效，需重新登陆获取token
		result.Error(c, goxy.StatusTokenExpired, "登陆失效")
		return
	}
	if err = s.AuthApiVerification(c); err != nil {
		result.Error(c, goxy.StatusParameterError, err.Error())
		return
	}
	c.Next()
}

// @Title AuthApiVerification
// @Description 权限Api验证
// @Author xuanshuiyuan 2022-01-21 17:12
// @Param context.Context
// @Return
func (s *Service) AuthApiVerification(c context.Context) (err error) {
	admin_info, err := utils.GetAdminInfo(c)
	if err != nil {
		return
	}
	if admin_info["api"] == nil {
		return errors.New("无权操作")
	}
	isdo := false
	path := c.Path()
	if goxy.StringsInSearch(path, conf.AuthApiNoVerification) {
		isdo = true
	} else {
		for _, v := range admin_info["api"].([]interface{}) {
			if path == v {
				isdo = true
				break
			}
		}
	}
	if isdo == false {
		return errors.New("无权操作")
	}
	return
}

// @Title Verification
// @Description 接口验证
// @Author xuanshuiyuan 2021-10-31 16:01
// @Param context.Context
// @Return int64, string
func (s *Service) Verification(c context.Context) {
	if err := s.verification(c); err != nil {
		result.Error(c, goxy.StatusValidationFailed, err.Error())
		return
	}
	requestHeaderParamsKey := [...]string{"token"}
	for _, v := range requestHeaderParamsKey {
		if c.GetHeader(v) == "" {
			result.Error(c, goxy.StatusValidationFailed, RequestFailTip)
			return
		}
	}
	requestHeaderParams := RequestHeaderParams{
		Token: c.GetHeader("token"),
	}
	//redis验证token
	redis := utils.NewRedis()
	err := redis.RedisVerification(requestHeaderParams.Token, c.GetHeader("source"))
	if err != nil { ////登录失效，需重新登陆获取token
		result.Error(c, goxy.StatusTokenExpired, "登陆失效")
		return
	}
	c.Next()
}

func (s *Service) verification(c context.Context) error {
	requestHeaderParamsKey := [...]string{"timestamp", "sign", "source", "random", "edition"}
	for _, v := range requestHeaderParamsKey {
		if c.GetHeader(v) == "" {
			return errors.New(RequestFailTip)
		}
	}
	timestamp, _ := strconv.ParseInt(c.GetHeader("timestamp"), 10, 64)
	requestHeaderParams := RequestHeaderParams{
		Timestamp: timestamp,
		Sign:      c.GetHeader("sign"),
		Source:    c.GetHeader("source"),
		Develop:   c.GetHeader("develop"),
		Random:    c.GetHeader("random"),
		Edition:   c.GetHeader("edition"),
	}
	//验证版本
	edition := utils.GetEdition(c)
	if edition == 0 {
		return errors.New(RequestFailTip)
	}
	if requestHeaderParams.Develop == RequestHeaderParamsDevelop && os.Getenv("ENV") != "production" {
		return nil
	}
	//请求有效时间5分钟
	if (time.Now().Unix() - requestHeaderParams.Timestamp/1000) > 5*60 {
		return errors.New("您的手机时间与北京时间相差超过5分钟")
	}
	//请求来源
	if contain := goxy.Contain(requestHeaderParams.Source, conf.Config.Conf.Verification.SourceList); contain == false {
		return errors.New(RequestFailTip)
	}
	//通过source获取key
	key := conf.Config.Conf.Verification.KeyList[requestHeaderParams.Source]
	signTemp := fmt.Sprintf("Source=%s&TimeStamp=%d&Key=%s&Weight=%s", requestHeaderParams.Source, requestHeaderParams.Timestamp, key, requestHeaderParams.Random)
	localSign := fmt.Sprintf("%X", md5.Sum([]byte(signTemp)))
	fmt.Println(signTemp, localSign)
	//验证sign
	if requestHeaderParams.Sign != localSign {
		return errors.New(RequestFailTip)
	}
	redis := &utils.RedisService{}
	//sing放入redis 5分钟 防止重复访问
	if err := redis.SetKey(localSign).RedisGetSign(); err == nil {
		return errors.New("请勿重复访问")
	} else {
		redis.SetKey(localSign).SetValue("1").SetExp(conf.RedisSignExp).RedisSetAndEx()
	}
	return nil
}
