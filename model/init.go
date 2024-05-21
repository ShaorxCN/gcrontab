package model

import (
	"fmt"
	"gcrontab/config"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	// 导入 postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var done chan struct{}

func init() {
	initDone()
}

func initDone() {
	done = make(chan struct{})
}

// 数据库
type globalDB struct {
	mu sync.Mutex
	*gorm.DB
}

var (
	db              = new(globalDB)
	PageSizeLimit   = 100
	PageSizeDefault = 20
)

// DB 是全局数据库。
func DB() *gorm.DB {
	<-done
	return db.New()
}

// DbConfig 是数据库配置。
type DbConfig struct {
	Host            string `json:"db_host" env:"DB_HOST"`
	Port            string `json:"db_port" env:"DB_PORT"`
	User            string `json:"db_user" env:"DB_USER"`
	Name            string `json:"db_name" env:"DB_NAME"`
	Password        string `json:"db_password" env:"DB_PASSWORD"`
	SslMode         string `json:"db_ssl_mode" env:"DB_SSL_MODE"`
	PageSizeLimit   int    `json:"db_pagesize_limit" env:"DB_PAGESIZE_LIMIT"`
	PageSizeDefault int    `json:"db_pagesize_default" env:"DB_PAGESIZE_DEFAULT"`
}

// NewDBConfig 返回一个 web 数据库配置。
func NewDBConfig() *DbConfig {
	return &DbConfig{}
}

// Init 初始化数据库。
func (d *DbConfig) Init() error {
	if err := initHandler(d, db); err != nil {
		return err
	}

	if err := db.DB.Exec(`create extension IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		panic(err)
	}

	close(done)

	return nil
}

// Restart 重启服务。
func (d *DbConfig) Restart() error {
	initDone()
	return restartHandler(d, db)
}

// Stop 处理的时候db最后关
func (d *DbConfig) Stop() {
	logrus.Warn("close db Connection...")
}

func initHandler(cf interface{}, gdb *globalDB) error {
	var (
		gormDB *gorm.DB
		err    error
	)
	if gormDB, err = connectDB(cf); err != nil {
		logrus.Warnf("DB can't connect: %v", err)
		return err
	}

	gdb.DB = gormDB
	return nil
}

func restartHandler(cf config.Configer, gdb *globalDB) error {
	if err := gdb.Close(); err != nil {
		return err
	}
	gdb.mu.Lock()
	defer gdb.mu.Unlock()
	if err := cf.Init(); err != nil {
		return err
	}
	return nil
}

// Models 返回所有 model。
func Models() []interface{} {
	return []interface{}{
		&DBUser{},
		&DBTask{},
		&DBTaskLog{},
	}
}

// AutoMigrate 自动更新数据库表结构。
func AutoMigrate() error {
	err := DB().AutoMigrate(Models()...).Error
	if err != nil {
		return err
	}
	return nil
}

// connectDB 连接数据库。
func connectDB(cf interface{}) (*gorm.DB, error) {
	var host, port, user, name, password, sslmode string
	switch c := cf.(type) {
	case *DbConfig:
		host = c.Host
		port = c.Port
		user = c.User
		name = c.Name
		password = c.Password
		sslmode = c.SslMode
		if c.PageSizeLimit != 0 {
			PageSizeLimit = c.PageSizeLimit
		}
		if c.PageSizeDefault != 0 {
			PageSizeDefault = c.PageSizeDefault
		}
	default:
		return nil, fmt.Errorf("config 不支持")
	}
	args := buildDBArgs(host, port, user, name, password, sslmode)
	gormDB, err := gorm.Open("postgres", args)
	if err != nil {
		logrus.Errorf("DB open error: %v", err)
		return gormDB, err
	}
	logrus.WithFields(logrus.Fields{"db_name": name}).Info("DB connect successful.")
	gormDB.DB().SetMaxIdleConns(8)
	gormDB.DB().SetMaxOpenConns(16)
	return gormDB, err
}

// buildDBArgs 构建连接 DB 的参数。
func buildDBArgs(host, port, user, name, password, sslmode string) string {
	args := ""
	if host != "" {
		args = " host=" + host
	}
	if port != "" {
		args = args + " port=" + port
	}
	if user != "" {
		args = args + " user=" + user
	}
	if name != "" {
		args = args + " dbname=" + name
	}
	if password != "" {
		args = args + " password=" + password
	}
	if sslmode != "" {
		args = args + " sslmode=" + sslmode
	}
	return args
}
