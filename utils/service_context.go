package utils

import (
	"context"
	"gcrontab/interface/entity"
	"gcrontab/redis"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type ServiceContext struct {
	ID string
	*gorm.DB
	Context      context.Context
	entitiers    []entity.Key // 只保存实体 key
	Operator     string       // 操作者的ID UserID or Host
	OperatorName string       // 用户的话是用户昵称
}

func NewServiceContext(ctx context.Context, db *gorm.DB, operator, nickName string) *ServiceContext {
	context := &ServiceContext{
		DB: db, Context: ctx,
		Operator:     operator,
		OperatorName: nickName,
		ID:           uuid.New().String(),
	}

	go func() {
		<-context.Context.Done()
		context.UnregisterAllEntityInRedis()
	}()
	return context
}

// UnregisterAllEntityInRedis 从redis 中注销实体锁
func (ctx *ServiceContext) UnregisterAllEntityInRedis() {
	for _, ek := range ctx.entitiers {
		err := redis.UnLock(ek.JSON())
		logrus.Errorf("unlock keys in redis failed:%v,key:[%v]", err, ek)
	}
}
