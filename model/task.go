package model

import (
	"fmt"
	"gcrontab/constant"
	"time"

	"github.com/google/uuid"
)

// DBTask 是任务的数据库模型。
type DBTask struct {
	Base
	Name             string `gorm:"index;not null;type:varchar(1024)"`
	IntervalDuration int    `gorm:"not null"`                  // 频率
	UnitOfInterval   string `gorm:"not null;type:varchar(32)"` // 频率单位 例如 天 周等
	Protocol         string `gorm:"type:varchar(32)"`          // HTTP / COMMAND
	Command          string `gorm:"type:varchar(1024)"`        // HTTP请求的URL 或者执行的脚本/命令
	HTTPMethod       string `gorm:"type:varchar(20)"`
	Expired_time     int    `gorm:"type:integer"` // 请求或者执行任务的超时时间 默认60s 单位ms
	RetryTimes       int    `gorm:"type:integer"` // 重试次数 默认为0 不重试
	RetryInterval    int    `gorm:"type:integer"` // 重试间隔 单位ms
	Remark           string `gorm:"type:varchar(1024)"`
	Status           string `gorm:"type:varchar(32);index"`
	Lock             string `gorm:"type:varchar(256)"`
	NextRuntime      *time.Time
	LastRuntime      *time.Time
	Level            int    `gorm:"type:integer"` // 重要程度
	UpdateID         string `gorm:"index;type:varchar(36)"`
	UpdateFlag       int8   `gorm:"type:smallint"` // 是否手动修改过(针对修改时间等)
	Param            string // 默认text好了
	PostType         string `gorm:"type:varchar(24)"` // json or form
	CreaterName      string `gorm:"type:varchar(255);index"`
	FailNotify       string `gorm:"type:varchar(10)"`   // on or off
	Email            string `gorm:"type:varchar(255)"`  // 通知邮箱
	Headers          string `gorm:"type:varchar(2048)"` // header
}

var taskTableName = new(DBTask).TableName()

// TableName 返回批次数据表名
func (DBTask) TableName() string {
	return "tbl_task"
}

// FindTaskByID 根据ID 查找Task 状态不为删除：del
func FindTaskByID(id uuid.UUID) (*DBTask, error) {
	db := DB()
	dbTask := &DBTask{}
	err := db.Model(dbTask).Where("id = ? and status != ?", id, constant.STATUSDEL).First(dbTask).Error
	// 这边不管什么状态都会去执行  待迭代
	if err != nil {
		return nil, err
	}
	return dbTask, nil
}

func FindTaskByName(name string) (*DBTask, error) {
	db := DB()
	t := new(DBTask)
	err := db.Model(t).Where("name = ? and status != ?", name, constant.STATUSDEL).First(t).Error
	return t, err

}

func FindTaskByNameWithOutStatus(name string) ([]uuid.UUID, error) {
	db := DB()

	var taskIDs []uuid.UUID
	err := db.Table(taskTableName).Where("name ilike ?", fmt.Sprintf("%%%s%%", name)).Pluck("id", &taskIDs).Error
	return taskIDs, err
}

// DeleteTaskByID 删除任务
func DeleteTaskByID(id uuid.UUID) error {
	db := DB()
	return db.Table(taskTableName).Where("id = ?", id).Update("status", constant.STATUSDEL).Error
}

// FindActiveTasks 查找待运行的任务
func FindActiveTasks(now time.Time) ([]*DBTask, error) {
	db := DB()
	var res []*DBTask
	err := db.Table(taskTableName).Where("next_runtime <= ? and Status = ?", now, constant.STATUSDEL).Find(&res).Error
	return res, err
}

// FindTaskByCode 根据code 查找Task 状态不为删除：del
func FindTaskByCode(code string) (*DBTask, error) {
	db := DB()
	dbTask := &DBTask{}
	err := db.Model(dbTask).Where("task_code = ? and status != ?", code, constant.STATUSDEL).First(dbTask).Error
	// 这边不管什么状态都会去执行  待迭代
	if err != nil {
		return nil, err
	}
	return dbTask, nil
}
