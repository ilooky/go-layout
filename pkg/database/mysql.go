package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilooky/go-layout/pkg/config"
	"github.com/ilooky/go-layout/pkg/guava"
	"github.com/ilooky/logger"
	"reflect"
	"strings"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"
)

var Db *xorm.Engine

func InitOrm(c config.Mysql) (db *xorm.Engine, err error) {
	mysqlUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
	logger.Infof("connect mysql url = %s", mysqlUrl)
	if Db, err = xorm.NewEngine("mysql", mysqlUrl); err != nil {
		return nil, err
	} else {
		l := log{level: levelMap[logger.GetLevel()], showSQL: c.ShowSql}
		Db.SetLogger(&l)
		Db.SetMaxIdleConns(10)
		Db.SetMaxOpenConns(10)
		Db.SetConnMaxLifetime(time.Minute * 60)
		loc, _ := time.LoadLocation("Local")
		Db.TZLocation = loc
		Db.DatabaseTZ = loc
		snakeMapper := names.SnakeMapper{}
		tbMapper := names.NewPrefixMapper(snakeMapper, "us_")
		Db.SetTableMapper(tbMapper)
		Db.SetColumnMapper(snakeMapper)
	}
	return Db, nil
}

type Base struct {
	Id        int64    `json:"id"`
	Created   JsonTime `json:"created"    xorm:"created"`
	Updated   JsonTime `json:"updated"    xorm:"updated"`
	DeletedAt JsonTime `json:"-"          xorm:"deleted"`
}
type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format(guava.DateLayout(guava.YYYY_MM_DD_HH_MM_SS)) + `"`), nil
}

func (j *JsonTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	loc, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	sv := string(b)
	if len(sv) == 10 {
		sv += " 00:00:00"
	} else if len(sv) == 16 {
		sv += ":00"
	}
	now, err := time.ParseInLocation(guava.DateLayout(guava.YYYY_MM_DD_HH_MM_SS), string(b), loc)
	if err != nil {
		if now, err = time.ParseInLocation(guava.DateLayout(guava.YYYY_MM_DD_HH_MM), string(b), loc); err != nil {
			now, err = time.ParseInLocation(guava.DateLayout(guava.YYYY_MM_DD), string(b), loc)
		}
	}
	*j = JsonTime(now)
	return
}

type log struct {
	level   xlog.LogLevel
	showSQL bool
}

var levelMap = map[string]xlog.LogLevel{
	"debug": xlog.LOG_DEBUG,
	"info":  xlog.LOG_INFO,
	"warn":  xlog.LOG_WARNING,
	"error": xlog.LOG_ERR,
	"panic": xlog.LOG_ERR,
}

func (l *log) Debug(v ...interface{}) {
	logger.Debug(v...)
}

func (l *log) Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func (l *log) Error(v ...interface{}) {
	logger.Error(v...)
}

func (l *log) Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (l *log) Info(v ...interface{}) {
	logger.Info(v...)
}

func (l *log) Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (l *log) Warn(v ...interface{}) {
	logger.Warn(v...)
}

func (l *log) Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func (l *log) Level() xlog.LogLevel {
	return l.level
}

func (l *log) SetLevel(level xlog.LogLevel) {
	l.level = level
}

func (l *log) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *log) IsShowSQL() bool {
	return l.showSQL
}

func Save(entity interface{}) (err error) {
	_, err = Db.InsertOne(entity)
	return
}

func SaveAll(entity ...interface{}) (err error) {
	_, err = Db.Insert(entity...)
	return
}

func Delete(entity interface{}) (err error) {
	_, err = Db.Delete(entity)
	return
}

func Update(entity interface{}) (err error) {
	_, err = Db.Update(entity)
	return
}

func FindOneById(id interface{}, dest interface{}) error {
	get, err := Db.ID(id).Get(dest)
	if !get || err != nil {
		return fmt.Errorf("not find entity where id = %s", id)
	}
	return nil
}

// FindOneByField dest is ptr
func FindOneByField(field string, value interface{}, dest interface{}) error {
	get, err := Db.Where(field+"=?", value).Get(dest)
	if get {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("not find entity ,where %s = %s", field, value)
}

// FindCols dest is *struct or *[]struct
func FindCols(field string, value interface{}, dest interface{}, cols ...string) error {
	session := Db.Where(field+"=?", value).Cols(cols...)
	var err error
	if reflect.ValueOf(dest).Elem().Kind() == reflect.Slice {
		err = session.Find(dest)
	} else {
		get, err := session.Get(dest)
		if !get {
			return fmt.Errorf("not find entity,where %s = %s", field, value)
		}
		return err
	}
	return err
}

func FindListByField(field string, value interface{}, dest interface{}) error {
	err := Db.Where(field+"=?", value).Find(dest)
	if err != nil {
		return err
	}
	return nil
}

func FindOneByCondition(conditions map[string]interface{}, dest interface{}) error {
	var condition []string
	var values []interface{}
	for k, v := range conditions {
		condition = append(condition, k+"=?")
		values = append(values, v)
	}
	if get, _ := Db.Where(strings.Join(condition, " and "), values...).Get(dest); get {
		return nil
	}
	return fmt.Errorf("not find entity ,where %+v ", conditions)
}

func FindListByCondition(conditions map[string]interface{}, dest interface{}) error {
	var condition []string
	var values []interface{}
	for k, v := range conditions {
		condition = append(condition, k+"=?")
		values = append(values, v)
	}
	return Db.Where(strings.Join(condition, " and "), values...).Find(dest)
}

func FindByValues(dest interface{}, conditions map[string]interface{}, field string, inValues ...interface{}) error {
	var condition []string
	var values []interface{}
	for k, v := range conditions {
		condition = append(condition, k+"=?")
		values = append(values, v)
	}
	return Db.Where(strings.Join(condition, " and "), values...).And(builder.In(field, inValues...)).Find(dest)
}

// Execute ("delete from us_diagram where STATION_CODE = ?", stationCode)
func Execute(sqlOrArgs ...interface{}) {
	session := Db.NewSession()
	defer session.Close()
	if _, err := session.Exec(sqlOrArgs...); err != nil {
		_ = session.Rollback()
	}
	_ = session.Commit()
}

func Executes(sql ...string) {
	session := Db.NewSession()
	defer session.Close()
	for _, s := range sql {
		if _, err := session.Exec(s); err != nil {
			_ = session.Rollback()
		}
	}
	_ = session.Commit()
}

func CreateTable(beans ...interface{}) {
	_ = Db.CreateTables(beans...)
}

func SyncTable(beans ...interface{}) {
	_ = Db.Sync2(beans...)
}
