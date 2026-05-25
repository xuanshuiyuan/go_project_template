// @Author xuanshuiyuan
package logs

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"go_project_template/internal/engine"
	"strings"
	"time"
)

type LogsDbService struct {
	LogsParams
}

func LoadDbConfig(c context.Context) *LogsDbService {
	info, _ := GetUserInfo(c)
	return &LogsDbService{
		LogsParams{
			Ctx:      c,
			UserId:   GetUserId(c),
			UserName: GetFieldVal(info, "username"),
			Mobile:   GetFieldVal(info, "mobile"),
			Source:   GetSource(info),
		},
	}
}

func (l *LogsDbService) Add() {
	l.Means.Add(&l.LogsParams)
}

func (l *LogsDbService) InjectParams(params *LogsParams) *LogsDbService {
	l.Action = params.Action
	l.RelateId = params.RelateId
	l.Remark = params.Remark
	l.Means = params.Means
	if l.Means == nil {
		path := strings.Split(l.Ctx.Path(), "/")
		if path[1] == "admin" {
			l.Means = NewDbAdminMeans()
		} else {
			l.Means = DefaultMeans()
		}
	}
	return l
}

func (l *LogsDbService) SetOperation(args ...interface{}) *LogsDbService {
	l.Operation = l.Operation + l.Means.SetOperation(l.Action, args...)
	return l
}

//数据库日志
type DbMeans struct {
}

func DefaultMeans() *DbMeans {
	return NewDbMeans()
}

func NewDbMeans() *DbMeans {
	return &DbMeans{}
}

func (d DbMeans) SetOperation(action string, args ...interface{}) string {
	return fmt.Sprintf(LogOperationConfig[action], FmtOperation(args...)...)
}

func (d DbMeans) Add(params *LogsParams) (err error) {
	create := &OperationLog{
		UserId:     params.UserId,
		UserName:   params.UserName,
		Mobile:     params.Mobile,
		Action:     params.Action,
		Operation:  params.Operation,
		RelateId:   params.RelateId,
		Source:     params.Source,
		Remark:     params.Remark,
		CreateTime: time.Now().Unix(),
	}
	err = engine.DB.Mysql.Model(&OperationLog{}).Create(&create).Error
	if err != nil {
		return
	}
	return
}

type DbAdminMeans struct {
}

func NewDbAdminMeans() *DbAdminMeans {
	return &DbAdminMeans{}
}

func (d DbAdminMeans) SetOperation(action string, args ...interface{}) string {
	return fmt.Sprintf(LogAdminOperationConfig[action], FmtOperation(args...)...)
}

func (d DbAdminMeans) Add(params *LogsParams) (err error) {
	create := &AdminLog{
		AdminId:    params.UserId,
		AdminName:  params.UserName,
		Action:     params.Action,
		Operation:  params.Operation,
		RelateId:   params.RelateId,
		Source:     params.Source,
		Remark:     params.Remark,
		CreateTime: time.Now().Unix(),
	}
	err = engine.DB.Mysql.Model(&AdminLog{}).Create(&create).Error
	if err != nil {
		return
	}
	return
}
