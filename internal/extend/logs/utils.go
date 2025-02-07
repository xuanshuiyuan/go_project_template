// @Author xuanshuiyuan
package logs

import (
	"encoding/json"
	"github.com/kataras/iris/v12/context"
	"go_project_template/internal/conf"
	"go_project_template/internal/utils"
	"reflect"
	"sort"
	"strconv"
)

func FmtOperation(args ...interface{}) []interface{} {
	keys := make([]int, 0, len(args))
	for key, _ := range args {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for k, _ := range keys {
		switch reflect.TypeOf(args[k]).Kind() {
		case reflect.Int, reflect.Int8:
			args[k] = strconv.Itoa(args[k].(int))
		case reflect.Int64:
			args[k] = strconv.FormatInt(args[k].(int64), 10)
		case reflect.Float64:
			args[k] = strconv.FormatFloat(args[k].(float64), 'f', -1, 64)
		case reflect.Slice, reflect.Map, reflect.Struct, reflect.Ptr:
			str, _ := json.Marshal(args[k])
			args[k] = str
		}
	}
	return args
}

func GetUserId(c context.Context) int64 {
	user_id, err := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	if err != nil {
		return 0
	}
	return user_id
}

func GetFieldVal(data map[string]interface{}, key string) (username string) {
	if data != nil {
		return data[conf.Config.Conf.Verification.SourceRedisList[data["source"].(string)]].(map[string]interface{})[key].(string)
	}
	return
}

func GetUserInfo(c context.Context) (result map[string]interface{}, err error) {
	redis := utils.NewRedis()
	user_id, err := strconv.ParseInt(c.GetHeader("userid"), 10, 64)
	if err != nil {
		return nil, nil
	}
	result, err = redis.GetInfoByKey(GetTokenKey(conf.RedisTokenKey, user_id, c.GetHeader("source")))
	if err != nil {
		return nil, nil
	}
	return
}

func GetTokenKey(key string, id int64, source string) string {
	return key + strconv.FormatInt(id, 10) + "_" + source
}

func GetSource(data map[string]interface{}) string {
	if data != nil {
		return conf.Config.Conf.Verification.SourceExplainList[data["source"].(string)]
	}
	return ""
}
