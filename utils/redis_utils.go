package utils

import (
	entitier "gcrontab/interface/entity"
	"gcrontab/redis"
)

// RegisterEntityInRedis 注册entity到redis。
func RegisterEntityInRedis(e entitier.Entitier, owner string, timeout int) bool {
	entityKey := e.EntityKey()
	err := redis.GetLock(entityKey.JSON(), owner, timeout)

	return err == nil
}

// UnregisterEntityInRedis 从redis解除entity
func UnregisterEntityInRedis(e entitier.Entitier) error {
	entityKey := e.EntityKey()
	err := redis.UnLock(entityKey.JSON())

	if err != nil {
		return err
	}

	return nil
}
