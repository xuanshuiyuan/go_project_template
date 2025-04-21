// @Author xuanshuiyuan 2025/4/17 11:26:00
package engine

import (
	"errors"
	"github.com/xuanshuiyuan/goxy"
	"go_project_template/internal/conf"
	"gorm.io/gorm"
)

type SqlService struct {
	DataBase  *gorm.DB    `json:"data_base"`
	Table     interface{} `json:"table"`
	Db        *gorm.DB    `json:"db"`
	Statement *gorm.DB    `json:"statement"`
}

func NewSql() *SqlService {
	return &SqlService{
		DataBase: DB.Mysql,
	}
}

func (s *SqlService) SetDb(data_base *gorm.DB) *SqlService {
	s.DataBase = data_base
	return s
}

func (s *SqlService) SetTable(table interface{}) *SqlService {
	s.Table = table
	return s
}

func (s *SqlService) init() *SqlService {
	s.Db = s.DataBase.Model(s.Table)
	s.Statement = s.DataBase.Session(&gorm.Session{DryRun: true}).Model(s.Table)
	return s
}

// @Title Last
// @Description 获取最新一条数据
// @Author xuanshuiyuan 2025-04-21 10:22
// @Param
// @Return
func (s *SqlService) Last(where interface{}, res interface{}) (ok bool, err error) {
	s.init()
	err = s.Db.Where(where).Last(res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logs(s.Statement.Where(where).Last(res), "Last", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) { //为空
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
func (s *SqlService) Lasts(where interface{}, res interface{}, query string, args ...interface{}) (ok bool, err error) {
	s.init()
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
		s.Logs(statement.Last(res), "Lasts", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) { //为空
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
func (s *SqlService) Update(where interface{}, data map[string]interface{}) (err error) {
	s.init()
	result := s.Db.Where(where).Updates(data)
	if result.Error != nil {
		s.Logs(s.Statement.Where(where).Updates(data), "Update", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("请至少修改一项数据")
	}
	return
}

// @Title Update
// @Description Scan
// @Author xuanshuiyuan 2025-04-21 14:16
// @Param context.Context
// @Return
func (s *SqlService) Scan(where interface{}, res interface{}, query string, args ...interface{}) (ok bool, err error) {
	s.init()
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
		s.Logs(statement.Order("id desc").Scan(&res), "Scan", err)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) { //为空
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
func (s *SqlService) Create(data interface{}, query interface{}, fields ...interface{}) (err error) {
	s.init()
	err = s.Db.Select(query, fields...).Create(data).Error
	if err != nil {
		s.Logs(s.Statement.Select(query, fields...).Create(data), "Create", err)
		return
	}
	return
}

func (s *SqlService) Logs(statement *gorm.DB, title string, err error) {
	sqlStm := statement.Statement
	log.Error(conf.Config.Base.LogFileName, "mysql.log").Println(goxy.FmtLog("方法.title", title, "sql语句.title", sqlStm.Dialector.Explain(sqlStm.SQL.String(), sqlStm.Vars...), "错误.title", err.Error()))
}
