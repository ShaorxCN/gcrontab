package model

import (
	"time"

	"github.com/google/uuid"
)

// DBTaskLog 是任务记录的数据库模型。
type DBTaskLog struct {
	TaskID     uuid.UUID `gorm:"primary_key;not null;type:uuid;"`
	TimeStamp  int64     `gorm:"primary_key;index;type:bigint"`
	TaskName   string    `gorm:"not null;type:varchar(1024)"`
	ResultCode int       // httpStatusCode or exitCode
	Result     string    `gorm:"type:text"`          // 返回结果 200/timeOut不填写  其他读取Body
	Command    string    `gorm:"type:varchar(1024)"` // HTTP请求的URL 或者执行的脚本/命令
	StartTime  time.Time `gorm:"type:timestamp(3)"`
	EndTime    time.Time `gorm:"type:timestamp(3)"`
	TotalTime  int64     `gorm:"type:bigint"`                     // 任务耗时 单位ms
	Host       string    `gorm:"type:varchar(1024)"`              // 执行的hostname 如果是立即运行就是用户名 否则是host
	Status     string    `gorm:"type:smallint;index;varchar(32)"` // 执行状态  processing or  success or fail
	User       string    `gorm:"type:varchar(255)"`               // 如果是立即执行 这里存储的是用户名
}

var taskLogTableName = new(DBTaskLog).TableName()

// TableName 返回任务日志表名
func (DBTaskLog) TableName() string {
	return "tbl_task_log"
}
