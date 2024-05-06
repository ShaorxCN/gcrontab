package model

import (
	"time"

	"github.com/google/uuid"
)

var (
	// TokenTTL 过期时间  单位s
	TokenTTL int
	// TimeOut 任务的默认超时时间
	TimeOut int
)

// Base 是数据库基本结构
type Base struct {
	// ID 是 UUID 主键
	ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// CreateAt 创建时间
	CreateAt *time.Time `gorm:"index;default:now()"`
	// UpdateAt 更新时间
	UpdateAt *time.Time
	// 任务或者用户的创建者(task)
	Creater string `gorm:"index;type:varchar(36)"`
}
