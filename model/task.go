package model

import (
	"time"
)

// DBTask 是任务的数据库模型。
type DBTask struct {
	Base
	Name             string `gorm:"unique_index;not null;type:varchar(1024)"`
	IntervalDuration int    `gorm:"not null"`                  // 频率
	UnitOfInterval   string `gorm:"not null;type:varchar(32)"` // 频率单位 例如 天 周等
	Protocol         string `gorm:"type:varchar(32)"`          // HTTP / COMMAND
	Command          string `gorm:"type:varchar(1024)"`        // HTTP请求的URL 或者执行的脚本/命令
	HTTPMethod       string `gorm:"type:varchar(20)"`
	Expired_time     int    `gorm:"type:integer"` // 请求或者执行任务的超时时间 默认60s 单位ms
	RetryTimes       int    `gorm:"type:integer"` // 重试次数 默认为0 不重试
	RetryInterval    int    `gorm:"type:integer"` // 重试间隔 单位ms
	Remark           string `gorm:"type:varchar(1024)"`
	Status           int    `gorm:"type:smallint;index"`
	NextRuntime      *time.Time
	LastRuntime      *time.Time
	Level            int    `gorm:"type:integer"` // 重要程度
	UpdateID         string `gorm:"index;type:varchar(36)"`
	UpdateFlag       int8   `gorm:"type:smallint"` // 是否手动修改过(针对修改时间等)
	Param            string // 默认text好了
	PostType         string `gorm:"type:varchar(24)"` // json or form
	CreaterName      string `gorm:"type:varchar(255);index"`
	FailNotify       string `gorm:"type:varchar(10)"` // on or off
	// Email            string `gorm:"type:varchar(255)"`  // 通知邮箱
	Headers string `gorm:"type:varchar(2048)"` // header
}

var taskTableName = new(DBTask).TableName()

func GetTaskTableName() string {
	return taskTableName
}

// TableName 返回批次数据表名
func (DBTask) TableName() string {
	return "tbl_task"
}
