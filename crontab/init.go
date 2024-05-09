package gocron

import (
	"gcrontab/model"
	"gcrontab/utils"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxGoroutine = 10000
	DefaultScanInterval = 1000
	DefaultRunSize      = 1000
)

var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// CrontabConfig 定时任务的config
type CrontabConfig struct {
	// scheduler 本身最大go协程数量
	MaxGoroutine int `json:"crontab_maxgoroutine" env:"CRONTAB_MAXGOROUTINE"`
	// 立即运行的channel容量
	RunSize int `json:"crontab_run_size" env:"CRONTAB_RUN_SIZE"`
	// 扫表间隔
	Interval     int    `json:"crontab_interval" env:"CRONTAB_INTERVAL"`
	TimeLocation string `json:"crontab_timeLocation" env:"CRONTAB_TIMELOCATION"`
}

// Init 开启定时任务。
func (c *CrontabConfig) Init() error {
	if c.MaxGoroutine == 0 {
		c.MaxGoroutine = DefaultMaxGoroutine
	}

	if c.Interval == 0 {
		c.Interval = DefaultScanInterval
	}

	if c.RunSize == 0 {
		c.RunSize = DefaultRunSize
	}

	ts = new(taskScheduler)
	ts.MaxGoroutine = make(chan int, c.MaxGoroutine)
	ts.ScanInterval = c.Interval
	TaskChannel = make(chan *model.DBTask, c.RunSize)
	ts.exit = make(chan struct{})

	err := utils.InitTimeLocation(c.TimeLocation)
	if err != nil {
		return err
	}

	ts.Start()

	return nil

}

// Stop 停止调度
func (c *CrontabConfig) Stop() {
	logger.WithTime(utils.Now()).Warn("stop scheduler...")
	ts.Stop()
	logger.WithTime(utils.Now()).Warn("scheduler stoped...")
}

// Restart 重启服务。
func (c *CrontabConfig) Restart() error {
	c.Stop()
	return c.Init()

}
