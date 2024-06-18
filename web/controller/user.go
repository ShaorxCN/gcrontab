package controller

import (
	"gcrontab/rep/requestmodel"

	"github.com/gin-gonic/gin"
)

const (
	defaultUserName = "admin"
	defaultPassWord = "admin"
)

// User 用来实现用户的 rest 接口。
type User struct{}

// AddUserRouter 注册用户 router。
func AddUserRouter(e *gin.Engine) {
	// e.POST("/users/login", User{}.Login)
	e.POST("/users", User{}.CreateUser)
	// e.GET("/users/:userID", User{}.GetUser)
	// e.PUT("/users/:userID", User{}.ModifyUser)
	// e.POST("/users/:userID/logout", User{}.LoginOut)
	// e.DELETE("/users/:userID", User{}.DeleteUserByID)
}

func (u User) CreateUser(ctx *gin.Context) {
	in := new(requestmodel.UserReq)
}
