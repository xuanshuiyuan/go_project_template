// @Author xuanshuiyuan
package logs

import "github.com/kataras/iris/v12/context"

var LogOperationConfig = map[string]string{
	"UploadFile": "上传文件：%s",
}

var LogAdminOperationConfig = map[string]string{
	"UploadFile": "上传文件：%s",
}

type LogsParams struct {
	Ctx       context.Context //
	Means     MeansEr
	UserId    int64
	UserName  string
	Mobile    string
	Action    string
	Operation string
	RelateId  int64
	Source    string
	Remark    string
}

type OperationLog struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	Mobile     string `json:"mobile"`
	Action     string `json:"action"`
	Operation  string `json:"operation"`
	RelateId   int64  `json:"relate_id"`
	Source     string `json:"source"`
	CreateTime int64  `json:"create_time"`
	Remark     string `json:"remark"`
}

type AdminLog struct {
	Id         int64  `json:"id"`
	AdminId    int64  `json:"admin_id"`
	AdminName  string `json:"admin_name"`
	Action     string `json:"action"`
	Operation  string `json:"operation"`
	RelateId   int64  `json:"relate_id"`
	Source     string `json:"source"`
	CreateTime int64  `json:"create_time"`
	Remark     string `json:"remark"`
}
