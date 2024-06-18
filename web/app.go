package web

import (
	"fmt"
	"gcrontab/model"
	"gcrontab/web/middleware"

	"gcrontab/web/controller"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinConfig 配置
type GinConfig struct {
	Host           string `json:"gin_host" env:"GIN_HOST"`
	Port           string `json:"gin_port" env:"GIN_PORT"`
	Mode           string `json:"gin_mode" env:"GIN_MODE"`
	TokenTTL       int    `json:"gin_token_ttl" env:"GIN_TOKEN_TTL"`
	TokenCacheSize int    `json:"gin_token_cache_size" env:"GIN_TOKEN_CACHE_SIZE"`
	AdminUserName  string `json:"admin_userName" env:"ADMIN_USERNAME"`
	AdminPassWord  string `json:"admin_passWord" env:"ADMIN_PASSWORD"`
	AdminEmail     string `json:"admin_email" env:"ADMIN_EMAIL"`
	APISecret      string `json:"api_secret" env:"API_SECRET"`
	TaskTimeOut    int    `json:"task_timeout" env:"TASK_TIMEOUT"`
	StatePort      int    `json:"state_port" env:"STATE_PORT"`
}

const (
	DefaultTokenTTL    = 120
	DefaultTaskExpired = 3600000
)

// Init 初始化 rest 服务。
func (g *GinConfig) Init() (err error) {
	r := gin.New()
	logrus.WithFields(logrus.Fields{"host": g.Host, "port": g.Port}).Info("REST Server 启动")
	if g.Mode != "" {
		gin.SetMode(g.Mode)
	}

	middleware.TokenTTL = g.TokenTTL
	if middleware.TokenTTL == 0 {
		middleware.TokenTTL = DefaultTokenTTL
	}

	model.TaskExpired = g.TaskTimeOut
	if model.TaskExpired == 0 {
		model.TaskExpired = DefaultTaskExpired
	}

	logger := logrus.StandardLogger()
	r.Use(
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.Cors(),
		middleware.CheckPermission(logger),
		middleware.QueryTrans(),
	)

	controller.AddTaskRouter(r)
	controller.AddTaskLogRouter(r)
	err = controller.InsertAdminUser(g.AdminUserName, g.AdminPassWord, g.AdminEmail)
	if err != nil {
		return fmt.Errorf("init admin user failed:%v", err)
	}
	return r.Run(g.Host + ":" + g.Port)
}

// Restart 重启服务。
func (g *GinConfig) Restart() error {
	return g.Init()
}

// Stop 停止服务
func (g *GinConfig) Stop() {
	logrus.Warn("close gin server...")
}
