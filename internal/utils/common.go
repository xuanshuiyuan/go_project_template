// @Author
package utils

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"math/rand"
	"net"
	"strconv"
	"time"
)

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

func GetRedisSpinKeyStr(key string, str string) string {
	return fmt.Sprintf("%s%s", key, str)
}

func GetRedisSpinKey(key string, id int64) string {
	return fmt.Sprintf("%s%d", key, id)
}

func GetRedisSpinStrKey(key string, id string) string {
	return fmt.Sprintf("%s%s", key, id)
}

func CreateLoginCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000))
}

func GetExcelFileName(filename string) string {
	var year = time.Now().Year()
	var month = time.Now().Format("01")
	var day = time.Now().Format("02")
	file_path := fmt.Sprint(fmt.Sprintf("%s", goxy.StaticFileDirectory()))
	goxy.CheckDir(file_path)
	name := fmt.Sprintf("%s-%d-%s-%s-%s.xlsx", filename, year, month, day, goxy.RandChar(5))
	return fmt.Sprintf("%s/%s", file_path, name)
}

func GetTokenKey(key string, id int64, source string) string {
	//return key + strconv.FormatInt(id, 10) + "_" + source + "_" + conf.Config.Conf.Verification.SourceEngExplainList[source]
	return key + strconv.FormatInt(id, 10) + "_" + source
}

// @Title GetCacheUserInfoId
// @Description 得到当前用户的ID
// @Author xuanshuiyuan 2021-1-11 10:19
// @Param
// @Return int64
func GetCacheUserInfoId(c context.Context) int64 {
	//return int64((data["userinfo"].(map[string]interface{})["user_id"]).(float64))
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	return user_id
}

// @Title GetEdition
// @Description 得到当前请求的版本号
// @Author xuanshuiyuan 2022-04-07 16:30
// @Param
// @Return int64
func GetEdition(c context.Context) int64 {
	return conf.Config.Conf.Verification.EditionList[c.GetHeader("edition")]
}

// @Title GetAdminInfo
// @Description 得到当前用户的信息
// @Author xuanshuiyuan 2021/12/30 09:52
// @Param context.Context
// @Return map[string]interface{}, string
func GetAdminInfo(c context.Context) (result map[string]interface{}, err error) {
	redis := NewRedis()
	user_id, _ := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	result, err = redis.GetInfoByKey(GetTokenKey(conf.RedisAdminTokenKey, user_id, c.GetHeader("source")))
	if err != nil {
		return
	}
	return
}

//去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// @Title GetRedisAdminTokenKey
// @Description 得到当前管理员的redis key
// @Author xuanshuiyuan 2021/12/30 09:52
// @Param user_id
// @Return string
func GetRedisAdminTokenKey(user_id int64) string {
	return fmt.Sprintf("%s%d", conf.RedisAdminTokenKey, user_id)
}

// @Title GetCacheAdminId
// @Description 得到当前登陆管理员的ID
// @Author xuanshuiyuan 2021/12/29 14:44
// @Param
// @Return int64
func GetCacheAdminId(data map[string]interface{}) int64 {
	return int64((data["admin_info"].(map[string]interface{})["id"]).(float64))
}

// @Title GetCacheAdminName
// @Description 得到当前登陆管理员的姓名
// @Author xuanshuiyuan 2021/12/29 14:44
// @Param
// @Return string
func GetCacheAdminName(data map[string]interface{}) string {
	return data["admin_info"].(map[string]interface{})["username"].(string)
}

// @Title GetCacheAdminSource
// @Description 得到当前登陆管理员的登陆来源
// @Author xuanshuiyuan 2021/12/29 14:44
// @Param
// @Return string
func GetCacheAdminSource(data map[string]interface{}) string {
	return conf.Config.Conf.Verification.SourceExplainList[data["source"].(string)]
}

func GetByAdminLogStruct(key string) string {
	return conf.AdminLog[key]
}

func GetVerifyCode() string {
	rand.Seed(time.Now().UnixNano())         // 初始化随机种子
	randomNumber := rand.Intn(10000)         // 生成0到9999的随机数
	return fmt.Sprintf("%04d", randomNumber) // 格式化输出为4位数
}

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
