package email

import (
	"crypto/tls"
	"net/smtp"

	"github.com/sirupsen/logrus"
)

// Config 邮件模块的config
type Config struct {
	Addr      string `json:"email_addr" env:"EMAIL_ADDR"`
	Host      string `json:"email_host" env:"EMAIL_HOST"`
	PassWord  string `json:"email_passWord" env:"EMAIL_PASSWORD"`
	User      string `json:"email_user" env:"EMAIL_USER"`
	TestEmail string `json:"test_email" env:"TEST_EMAIL"`
	TimeOut   int    `json:"timeout" env:"TIMEOUT"`
	ViewURL   string `json:"view_url" env:"VIEW_URL"`
}

var (
	auth        smtp.Auth
	ok          chan struct{}
	user        string
	addr        string
	tlsConfig   *tls.Config
	testAddress string
	viewURL     string // log详情url 模板
	// pool        *email.Pool
	// sendTimeOut int // 发送超时时间 单位s
)

func initDone() {
	ok = make(chan struct{})
}

// Init 初始化
func (c *Config) Init() error {
	initDone()
	logrus.Info("module email init...")
	var err error
	auth = smtp.PlainAuth("", c.User, c.PassWord, c.Host)
	user = c.User
	addr = c.Addr
	tlsConfig = &tls.Config{ServerName: c.Host}
	testAddress = c.TestEmail
	viewURL = c.ViewURL

	err = sendTest()
	if err != nil {
		return err
	}
	close(ok)
	logrus.Info("module email init end...")
	return nil
}

// Restart 重启
func (c *Config) Restart() error {
	c.Stop()
	return c.Init()
}

// Stop 停止
func (c *Config) Stop() {
	initDone()
}
