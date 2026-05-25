// @Author xuanshuiyuan 2025/4/17 11:26:00
package engine

import (
	"errors"
	"fmt"
	"go_project_template/internal/conf"

	"github.com/xuanshuiyuan/goxy"
	"gorm.io/gorm"
)

var (
	ErrDBNotInitialized = errors.New("database is not initialized")
	ErrTableNotSet      = errors.New("sql table is not set")
	ErrNoRowsAffected   = errors.New("no rows affected")
)

// SqlService GORM SQL 操作封装，提供链式调用的 CRUD 方法
// 使用方式:
//   sql := NewSql().SetTable(&User{})
//   ok, err := sql.Last(map[string]interface{}{"id": 1}, &user)
//   err = sql.Create(&user, "id,name", "id", "name")
type SqlService struct {
	DataBase  *gorm.DB    // 数据库连接实例
	Table     interface{} // 数据表对应的模型
	Db        *gorm.DB    // 实际执行的查询（真实执行）
	Statement *gorm.DB    // DryRun 模式，仅用于生成 SQL 文本供日志记录
}

// NewSql 创建 SqlService 实例，默认使用全局 MySQL 连接
func NewSql() *SqlService {
	return &SqlService{
		DataBase: DB.Mysql,
	}
}

// SetDb 设置自定义数据库连接（链式调用）
func (s *SqlService) SetDb(data_base *gorm.DB) *SqlService {
	s.DataBase = data_base
	return s
}

// SetTable 设置数据表模型（链式调用）
func (s *SqlService) SetTable(table interface{}) *SqlService {
	s.Table = table
	return s
}

// init 初始化 Db 和 Statement，Db 用于真实查询，Statement 用于日志 SQL 提取
func (s *SqlService) init() *SqlService {
	s.Db = s.DataBase.Model(s.Table)
	s.Statement = s.DataBase.Session(&gorm.Session{DryRun: true}).Model(s.Table)
	return s
}

// prepare 执行前检查：确保数据库和表已设置，然后初始化查询
func (s *SqlService) prepare() error {
	if s.DataBase == nil {
		return ErrDBNotInitialized
	}
	if s.Table == nil {
		return ErrTableNotSet
	}
	s.init()
	return nil
}

// @Title Last
// @Description 获取最新一条数据
// @Author xuanshuiyuan 2025-04-21 10:22
// @Param
// @Return
// Last 获取最新一条数据，ok=false 表示未找到记录
func (s *SqlService) Last(where interface{}, res interface{}) (ok bool, err error) {
	if err = s.prepare(); err != nil {
		return
	}
	err = s.Db.Where(where).Last(res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logs(s.Statement.Where(where), "Last", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

// @Title Lasts
// @Description 获取最新一条数据
// @Author xuanshuiyuan 2025-04-21 14:03
// @Param
// @Return
// Lasts 获取最新一条数据（支持额外查询条件），ok=false 表示未找到记录
func (s *SqlService) Lasts(where interface{}, res interface{}, query string, args ...interface{}) (ok bool, err error) {
	if err = s.prepare(); err != nil {
		return
	}
	var db *gorm.DB
	var statement *gorm.DB
	if query != "" {
		db = s.Db.Where(where).Where(query, args...)
		statement = s.Statement.Where(where).Where(query, args...)
	} else {
		db = s.Db.Where(where)
		statement = s.Statement.Where(where)
	}
	err = db.Last(res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logs(statement, "Lasts", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

// @Title Update
// @Description 编辑
// @Author xuanshuiyuan 2025-04-21 14:07
// @Param context.Context
// @Return
// Update 更新数据，无匹配行时返回 ErrNoRowsAffected
func (s *SqlService) Update(where interface{}, data map[string]interface{}) (err error) {
	if err = s.prepare(); err != nil {
		return
	}
	result := s.Db.Where(where).Updates(data)
	if result.Error != nil {
		s.Logs(s.Statement.Where(where).Updates(data), "Update", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return
}

// @Title Update
// @Description Scan
// @Author xuanshuiyuan 2025-04-21 14:16
// @Param context.Context
// @Return
// Scan 扫描查询结果到结构体，按 ID 倒序，ok=false 表示未找到记录
func (s *SqlService) Scan(where interface{}, res interface{}, query string, args ...interface{}) (ok bool, err error) {
	if err = s.prepare(); err != nil {
		return
	}
	var db *gorm.DB
	var statement *gorm.DB
	if query != "" {
		db = s.Db.Where(where).Where(query, args...)
		statement = s.Statement.Where(where).Where(query, args...)
	} else {
		db = s.Db.Where(where)
		statement = s.Statement.Where(where)
	}
	db = db.Order("id desc")
	err = db.Scan(res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logs(statement.Order("id desc"), "Scan", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}
	ok = true
	return
}

// @Title Create
// @Description 添加数据
// @Author xuanshuiyuan 2025-04-21 14:16
// @Param openid,unionid
// @Return
// Create 创建数据，query 指定插入的字段名
func (s *SqlService) Create(data interface{}, query interface{}, fields ...interface{}) (err error) {
	if err = s.prepare(); err != nil {
		return
	}
	err = s.Db.Select(query, fields...).Create(data).Error
	if err != nil {
		s.Logs(s.Statement.Select(query, fields...), "Create", err)
		return
	}
	return
}

// Logs 记录 SQL 执行错误日志，从 DryRun 的 Statement 中提取 SQL 文本
func (s *SqlService) Logs(statement *gorm.DB, title string, err error) {
	if err == nil {
		return
	}
	sqlText := "unknown"
	if statement != nil && statement.Statement != nil && statement.Statement.Dialector != nil {
		sqlStm := statement.Statement
		sqlText = sqlStm.Dialector.Explain(sqlStm.SQL.String(), sqlStm.Vars...)
	}
	msg := goxy.FmtLog("方法.title", title, "sql语句.title", sqlText, "错误.title", err.Error())
	if log != nil && conf.Config != nil && conf.Config.Base != nil {
		log.Error(conf.Config.Base.LogFileName, "mysql.log").Println(msg)
		return
	}
	fmt.Println(msg)
}
