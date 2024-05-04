package web

import (
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

// Init 初始化 rest 服务。
func (g *GinConfig) Init() (err error) {
	r := gin.New()
	logrus.WithFields(logrus.Fields{"host": g.Host, "port": g.Port}).Info("REST Server 启动")
	if g.Mode != "" {
		gin.SetMode(g.Mode)
	}

	return r.Run(g.Host + ":" + g.Port)
}

// Restart 重启服务。
func (g *GinConfig) Restart() error {
	return g.Init()
}

// Stop 停止服务
func (g *GinConfig) Stop() {
}
