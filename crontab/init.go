package crontab

import (
	"gcrontab/entity/task"
	"gcrontab/utils"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxGoroutine    = 10000
	DefaultDBScanInterval  = 60000
	DefaultMemScanInterval = 1000
	DefaultRunSize         = 1000
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
	// db扫表间隔
	DBInterval int `json:"crontab_db_interval" env:"CRONTAB_DB_INTERVAL"`
	// 内存数据扫描间隔
	MemInterval  int    `json:"crontab_mem_interval" env:"CRONTAB_MEM_INTERVAL"`
	TimeLocation string `json:"crontab_timeLocation" env:"CRONTAB_TIMELOCATION"`
	// 任务队列阻塞警告时间线  单位秒
	BlockAlertTime int `json:"block_alert_time" env:"BLOCK_ALERT_TIME"`
}

// Init 开启定时任务。
func (c *CrontabConfig) Init() error {
	if c.MaxGoroutine == 0 {
		c.MaxGoroutine = DefaultMaxGoroutine
	}

	if c.DBInterval == 0 {
		c.DBInterval = DefaultDBScanInterval
	}

	if c.MemInterval == 0 {
		c.MemInterval = DefaultMemScanInterval
	}

	if c.RunSize == 0 {
		c.RunSize = DefaultRunSize
	}

	ts = new(taskScheduler)
	ts.MaxGoroutine = make(chan int, c.MaxGoroutine)
	ts.DBScanInterval = c.DBInterval
	ts.MemScanInterval = c.MemInterval
	ts.updateTaskChan = make(chan *task.Task, 30)
	ts.updateExecTask = make([]*task.Task, 0, 100)
	ts.taskUUIDMap = make(map[string]int)
	imme_tasks = make(chan *task.Task, c.RunSize)
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
