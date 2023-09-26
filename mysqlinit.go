package gormutils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/xiaotushaoxia/errx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var InvalidMysqlAddr = fmt.Errorf("invalid mysql addr config")

type Address struct {
	Address  string `json:"address" yaml:"address"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	DBName   string `json:"dbname" yaml:"dbname"`

	// 部分mysql dsn配置 不全  具体看 https://github.com/Go-SQL-Driver/MySQL
	DialTimeout  int `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  int `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int `json:"write_timeout" yaml:"write_timeout"`

	// sql.Conn的配置
	MaxIdleConns    int `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int `json:"max_open_conns" yaml:"max_open_conns"`
	MaxLifeTime     int `json:"max_life_time" yaml:"max_life_time"`
	ConnMaxIdleTime int `json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}

// NewDB 创建一个gorm.DB， 并且会自动创建一个database Address.DBName
func NewDB(addr *Address, opts ...gorm.Option) (*gorm.DB, error) {
	pool, err := openSqlDBAndCreateDatabase(addr) // 先创建数据库
	if err != nil {
		return nil, err
	}
	return NewDBWithConnPool(pool, opts...)
}

func NewDBWithConnPool(pool gorm.ConnPool, opts ...gorm.Option) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool}), opts...)
	if err != nil {
		return nil, failedTo(err, "gorm.Open")
	}
	return db, nil
}

func createDatabase(upa, dbname string, ops string) error {
	// upa: user password address.  addr.User + ":" + addr.Password + "@tcp(" + addr.Address + ")/"
	dsn := upa + "?" + ops
	db0, err := sql.Open("mysql", dsn)
	if nil != err {
		return failedTo(err, "sql.Open")
	}
	// db0是连接池，create database以后需要每个连接都设置一下use database. 太麻烦 直接关了重连
	defer db0.Close()

	_, err = db0.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;", dbname))
	if err != nil {
		return failedTo(err, "create database "+dbname)
	}
	return nil
}

// 初始化mysql的数据库——没有配置文件中的数据库就创建，有就跳过
func openSqlDBAndCreateDatabase(addr *Address) (*sql.DB, error) {
	upa, ops, err := addr.toDsn2()
	if err != nil {
		return nil, err
	}

	err = createDatabase(upa, addr.DBName, ops)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", upa+addr.DBName+"?"+ops)
	if nil != err {
		return nil, failedTo(err, "sql.Open")
	}
	if addr.MaxIdleConns != 0 {
		db.SetMaxIdleConns(addr.MaxIdleConns)
	}
	if addr.MaxOpenConns != 0 {
		db.SetMaxOpenConns(addr.MaxOpenConns)
	}
	if addr.MaxLifeTime != 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(addr.MaxLifeTime))
	}
	if addr.MaxLifeTime != 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(addr.MaxLifeTime))
	}
	return db, nil
}

func (addr *Address) toDsn2() (upa, ops string, err error) {
	if addr == nil {
		err = errx.WithMessage(InvalidMysqlAddr, "nil Address")
		return
	}
	if addr.Address == "" {
		err = errx.WithMessage(InvalidMysqlAddr, " empty Address")
		return
	}
	if addr.Password == "" {
		err = errx.WithMessage(InvalidMysqlAddr, " empty Password")
		return
	}
	if addr.User == "" {
		err = errx.WithMessage(InvalidMysqlAddr, " empty User")
		return
	}
	if addr.DBName == "" {
		err = errx.WithMessage(InvalidMysqlAddr, " empty DBName")
		return
	}

	upa = addr.User + ":" + addr.Password + "@tcp(" + addr.Address + ")/"
	// 默认dsn配置
	var ss = []string{"multiStatements=true", "charset=utf8mb4", "parseTime=True", "loc=Local"}

	if addr.DialTimeout != 0 {
		ss = append(ss, fmt.Sprintf("timeout=%s", time.Duration(addr.DialTimeout)*time.Second))
	}
	if addr.ReadTimeout != 0 {
		ss = append(ss, fmt.Sprintf("readTimeout=%s", time.Duration(addr.ReadTimeout)*time.Second))
	}
	if addr.WriteTimeout != 0 {
		ss = append(ss, fmt.Sprintf("writeTimeout=%s", time.Duration(addr.WriteTimeout)*time.Second))
	}
	ops = strings.Join(ss, "&")
	return
}
