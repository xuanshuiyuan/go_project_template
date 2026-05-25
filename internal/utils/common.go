// @Author xuanshuiyuan
// 通用工具函数包：IP获取、Redis Key生成、验证码、数组去重、缓存读取等
package utils

import (
	crypto_rand "crypto/rand"
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// GetLocalIP 获取本机非回环、全局单播的 IPv4 地址
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// GetRedisSpinKeyStr 拼接 Redis 键名（字符串后缀）
// 使用方式: GetRedisSpinKeyStr("lock_", "order_123") → "lock_order_123"
func GetRedisSpinKeyStr(key string, str string) string {
	return fmt.Sprintf("%s%s", key, str)
}

// GetRedisSpinKey 拼接 Redis 键名（int64 后缀）
// 使用方式: GetRedisSpinKey("user_lock_", 100) → "user_lock_100"
func GetRedisSpinKey(key string, id int64) string {
	return fmt.Sprintf("%s%d", key, id)
}

// GetRedisSpinStrKey 拼接 Redis 键名（字符串 ID 后缀）
// 使用方式: GetRedisSpinStrKey("token_", "abc123") → "token_abc123"
func GetRedisSpinStrKey(key string, id string) string {
	return fmt.Sprintf("%s%s", key, id)
}

// CreateLoginCode 生成 4 位随机登录验证码
func CreateLoginCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000))
}

// GetExcelFileName 生成带日期和随机字符的 Excel 文件名
// 使用方式: GetExcelFileName("用户报表") → "/static/files/用户报表-2026-01-15-A3B5C.xlsx"
func GetExcelFileName(filename string) string {
	now := time.Now()
	var year = now.Year()
	var month = now.Format("01")
	var day = now.Format("02")
	file_path := fmt.Sprint(fmt.Sprintf("%s", goxy.StaticFileDirectory()))
	goxy.CheckDir(file_path)
	name := fmt.Sprintf("%s-%d-%s-%s-%s.xlsx", filename, year, month, day, goxy.RandChar(5))
	return fmt.Sprintf("%s/%s", file_path, name)
}

// GetTokenKey 拼接 Token 的 Redis 存储键
// 使用方式: GetTokenKey(conf.RedisAdminTokenKey, 100, "0101") → "go_project_template_admin_token_100_0101"
func GetTokenKey(key string, id int64, source string) string {
	return key + strconv.FormatInt(id, 10) + "_" + source
}

// GetCacheUserInfoId 从请求头获取当前登录用户 ID
func GetCacheUserInfoId(c context.Context) int64 {
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	return user_id
}

// GetEdition 从请求头获取当前请求的版本号，版本不存在时返回 0
func GetEdition(c context.Context) int64 {
	edition := conf.Config.Conf.Verification.EditionList[c.GetHeader("edition")]
	return edition
}

// GetAdminInfo 从 Redis 获取当前管理员信息
func GetAdminInfo(c context.Context) (result map[string]interface{}, err error) {
	redis := NewRedis()
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	result, err = redis.GetInfoByKey(GetTokenKey(conf.RedisAdminTokenKey, user_id, c.GetHeader("source")))
	return
}

// RemoveRepeatedElement 字符串数组去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	seen := make(map[string]bool, len(arr))
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			newArr = append(newArr, v)
		}
	}
	return
}

// GetRedisAdminTokenKey 获取当前管理员的 Redis Token 键名
func GetRedisAdminTokenKey(user_id int64) string {
	return fmt.Sprintf("%s%d", conf.RedisAdminTokenKey, user_id)
}

// GetCacheAdminId 从缓存数据中安全获取管理员 ID
func GetCacheAdminId(data map[string]interface{}) (int64, error) {
	adminInfo, ok := data["admin_info"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("admin_info 类型断言失败")
	}
	id, ok := adminInfo["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("admin_info.id 类型断言失败")
	}
	return int64(id), nil
}

// GetCacheAdminName 从缓存数据中安全获取管理员姓名
func GetCacheAdminName(data map[string]interface{}) (string, error) {
	adminInfo, ok := data["admin_info"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("admin_info 类型断言失败")
	}
	name, ok := adminInfo["username"].(string)
	if !ok {
		return "", fmt.Errorf("admin_info.username 类型断言失败")
	}
	return name, nil
}

// GetCacheAdminSource 从缓存数据中安全获取管理员登录来源名称
func GetCacheAdminSource(data map[string]interface{}) (string, error) {
	source, ok := data["source"].(string)
	if !ok {
		return "", fmt.Errorf("source 类型断言失败")
	}
	return conf.Config.Conf.Verification.SourceExplainList[source], nil
}

// GetByAdminLogStruct 根据操作类型获取日志模板
func GetByAdminLogStruct(key string) string {
	return conf.AdminLog[key]
}

// GetVerifyCode 使用加密随机数生成 4 位数字验证码
func GetVerifyCode() string {
	n := make([]byte, 2)
	_, _ = crypto_rand.Read(n)
	randomNumber := int(n[0]) % 10000
	return fmt.Sprintf("%04d", randomNumber)
}

// SplitSliceBySize 将字符串切片按指定大小分割为多个子切片
// 使用方式: SplitSliceBySize(["a","b","c","d","e"], 2) → [["a","b"],["c","d"],["e"]]
func SplitSliceBySize(slice []string, n int) [][]string {
	if n <= 0 {
		return nil
	}
	var divided [][]string
	length := len(slice)
	for i := 0; i < length; i += n {
		end := i + n
		if end > length {
			end = length
		}
		divided = append(divided, slice[i:end])
	}
	return divided
}
