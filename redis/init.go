package redis

import (
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

const (
	// 参考服务端 应当小于服务端
	defaultIdleTimeOut = 60 * 5
	defaultMaxIdle     = 10
	// defaultMaxActive = 20
)

var (
	redisPool *redisClient
	done      chan struct{}
)

type redisClient struct {
	mu   sync.Mutex
	pool *redis.Pool
}

// Config redis 初始化配置
type Config struct {
	Host        string `json:"redis_host" env:"REDIS_HOST"`
	Port        string `json:"redis_port" env:"REDIS_PORT"`
	PassWord    string `json:"redis_password" env:"REDIS_PASSWORD"`
	MaxIdle     int    `json:"redis_maxidle" env:"REDIS_MAXIDLE"`
	IdleTimeOut int    `json:"redis_idle_timeout" env:"REDIS_IDLE_TIMEOUT"`
	Index       int    `json:"redis_index" env:"REDIS_INDEX"`
}

func (c *Config) buildArgs() string {
	return fmt.Sprintf("redis://:%s@%s:%s/%d", c.PassWord, c.Host, c.Port, c.Index)
}

func initDone() {
	done = make(chan struct{})
}

// Init 初始化
func (c *Config) Init() error {
	logrus.Info("redis init...")
	initDone()
	url := c.buildArgs()
	maxIdle, idleTimeOut := defaultMaxIdle, defaultIdleTimeOut
	if c.IdleTimeOut != 0 {
		idleTimeOut = c.IdleTimeOut
	}
	if c.MaxIdle != 0 {
		maxIdle = c.MaxIdle
	}

	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Second * time.Duration(idleTimeOut),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(url)
			return conn, err
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	redisPool = new(redisClient)
	redisPool.pool = pool

	close(done)

	conn := GetConn()
	if conn.Err() != nil {
		return fmt.Errorf("redis init failed:%v", conn.Err())
	}

	defer conn.Close()
	logrus.Info("redis init end...")
	return nil
}

// Restart 重新初始化redis pool
func (c *Config) Restart() error {
	initDone()
	if err := redisPool.pool.Close(); err != nil {
		return err
	}
	redisPool.mu.Lock()
	defer redisPool.mu.Unlock()

	return c.Init()
}

// Stop 释放pool
func (c *Config) Stop() {
	logrus.Warn("close db Connections...")
	initDone()
	redisPool.pool.Close()
	logrus.Warn("close redis Connetctions end...")
}

// GetConn 获取一个连接
func GetConn() redis.Conn {
	<-done
	return redisPool.pool.Get()
}
