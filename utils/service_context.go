package utils

import (
	"context"
	"gcrontab/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceContext struct {
	ID           string
	Cancel       context.CancelFunc
	Ctx          context.Context
	Operator     string // 操作者的ID UserID or Host
	OperatorName string // 用户的话是用户昵称
}

func NewServiceContext(ctx *gin.Context) *ServiceContext {
	operator := ctx.GetHeader(constant.HEADEROPERATOR)
	nickName := ctx.GetHeader(constant.HEADEROPERATORNAME)

	sctx := &ServiceContext{
		Operator:     operator,
		OperatorName: nickName,
		ID:           uuid.New().String(),
	}

	sctx.Ctx, sctx.Cancel = context.WithCancel(context.Background())
	return sctx
}
