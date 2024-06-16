package model

import (
	"gcrontab/constant"
	"time"

	"github.com/google/uuid"
)

var userTableName = new(DBUser).TableName()

type DBUser struct {
	ID         uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserName   string     `gorm:"unique_index;type:varchar(255)"`
	NickName   string     `gorm:"type:varchar(255)"`
	PassWord   string     `gorm:"type:varchar(64)"`
	Salt       string     `gorm:"type:varchar(32)"`
	Status     int        `gorm:"type:smallint"`
	Email      string     `gorm:"type:varchar(255)"`
	FailNotify string     `gorm:"index;type:varchar(8)"` // 任务失败是否通知
	CreateAt   *time.Time `gorm:"index;default:now()"`
	UpdateAt   *time.Time
	Creater    string `gorm:"type:varchar(36)"`
	Role       string `gorm:"type:varchar(32)"` // 角色  admin taskAdmin  user
}

// TableName 返回表名
func (DBUser) TableName() string {
	return "tbl_user"
}

func GetUserTableName() string {
	return userTableName
}

// FindEmails 寻找待发送的邮箱
func FindEmails() ([]string, error) {
	db := DB()
	var emails []string
	err := db.Table(userTableName).Where("fail_notify = ? and status = ?", constant.NOTIFYON, constant.STATUSNORMALDB).Pluck("email", &emails).Error
	return emails, err
}
