package model

import (
	"time"
)

// DBToken 是登陆token数据库模型。
type DBToken struct {
	UserID     string     `gorm:"primary_key;type:varchar(1024)"`
	CreateTime *time.Time `gorm:"index"`
	Token      string     `gorm:"type:varchar(1024)"`
	Salt       string     `gorm:"type:varchar(32)"`
}

var tokenTableName = new(DBToken).TableName()

// TableName 返回批次数据表名
func (DBToken) TableName() string {
	return "tbl_token"
}

func GetTokenTableName() string {
	return tokenTableName
}
