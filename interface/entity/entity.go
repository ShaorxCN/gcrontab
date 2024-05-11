package entity

import (
	"fmt"
)

// Entitier 实体需要实现的接口
type Entitier interface {
	EntityKey() Key
}

type Key struct {
	id string
	tp string
}

// NewEntityKey 生成key
func NewEntityKey(idValue, tpValue string) Key {
	return Key{
		id: idValue,
		tp: tpValue,
	}
}

// GetIDValue 返回id
func (e Key) GetIDValue() string {
	return e.id
}

// JSON 返回json 方便redis存储
func (e Key) JSON() string {
	return fmt.Sprintf("{\"id\":%s,\"tp\":%s}", e.id, e.tp)
}
