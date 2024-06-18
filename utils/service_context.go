package utils

import (
	"context"
	"gcrontab/constant"
	"gcrontab/rep/requestmodel"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceContext struct {
	ID           string
	Cancel       context.CancelFunc
	Ctx          context.Context
	Operator     string // 操作者的ID UserID or Host
	OperatorName string // 用户的话是用户昵称
	Query        *requestmodel.Params
}

func NewServiceContext(ctx *gin.Context, query *requestmodel.Params) *ServiceContext {
	operator := ctx.GetHeader(constant.HEADEROPERATORID)
	nickName := ctx.GetHeader(constant.HEADEROPERATORNAME)
	sctx := &ServiceContext{
		Operator:     operator,
		OperatorName: nickName,
		ID:           uuid.New().String(),
		Query:        query,
	}

	sctx.Ctx, sctx.Cancel = context.WithCancel(context.Background())
	return sctx
}
