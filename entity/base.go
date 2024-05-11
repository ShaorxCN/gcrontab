package entity

import (
	entitier "gcrontab/interface/entity"
)

// BaseEntity 基本结构体。
type BaseEntity struct {
	ID       entitier.Key `json:"-" url:"-"`
	CreateAt string       `json:"createAt,omitempty" url:"-"`
	UpdateAt string       `json:"updateAt,omitempty" url:"-"`
	Creater  string       `json:"creater,omitempty" url:"creater,omitempty"`
}

func (b *BaseEntity) EntityKey() entitier.Key {
	return b.ID
}
