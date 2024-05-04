package model

import (
	"gcrontab/constant"
	"time"

	"github.com/google/uuid"
)

const (
	StatusNormal = "normal"
)

var userTableName = new(DBUser).TableName()

type DBUser struct {
	ID         uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserName   string     `gorm:"unique_index;type:varchar(255)"`
	NickName   string     `gorm:"type:varchar(255)"`
	PassWord   string     `gorm:"type:varchar(64)"`
	Salt       string     `gorm:"type:varchar(32)"`
	Status     string     `gorm:"type:varchar(32)"`
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

// FindUserByUserName 根据用户名查询用户
func FindUserByUserName(userName string) (*DBUser, error) {
	db := DB()
	dbUser := &DBUser{}
	err := db.Model(dbUser).Where("user_name = ? and status = ?", userName, StatusNormal).First(dbUser).Error
	return dbUser, err
}

// FindUserByID 根据id 查询用户
func FindUserByID(id uuid.UUID) (*DBUser, error) {
	db := DB()
	dbUser := &DBUser{}
	err := db.Model(dbUser).Where("id = ? and status = ?", id, StatusNormal).First(dbUser).Error

	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

func InsertUser(u *DBUser) error {
	db := DB()
	return db.Create(u).Error
}

// DeleteUserByID 逻辑删除用户
func DeleteUserByID(id uuid.UUID, newUserName string) error {
	db := DB()
	m := map[string]string{"status": constant.STATUSON, "user_name": newUserName}
	return db.Table(userTableName).Where("id = ?", id).Update(m).Error
}

// FindEmails 寻找待发送的邮箱
func FindEmails() ([]string, error) {
	db := DB()
	var emails []string
	err := db.Table(userTableName).Where("fail_notify = ? and status = ?", constant.NOTIFYON, StatusNormal).Pluck("email", &emails).Error
	return emails, err
}
