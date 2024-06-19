package controller

import (
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/rep/requestmodel"
	"gcrontab/security"
	"gcrontab/service"
	"gcrontab/utils"
	"gcrontab/web/response"
	"gcrontab/web/validate"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	defaultUserName = "admin"
	defaultPassWord = "123456"
	defaultEmail    = "test@test.com"
)

// User 用来实现用户的 rest 接口。
type User struct{}

// AddUserRouter 注册用户 router。
func AddUserRouter(e *gin.Engine) {
	e.POST("/users/login", User{}.Login)
	e.POST("/users", User{}.CreateUser)
	// e.GET("/users/:userID", User{}.GetUser)
	// e.PUT("/users/:userID", User{}.ModifyUser)
	// e.POST("/users/:userID/logout", User{}.LoginOut)
	// e.DELETE("/users/:userID", User{}.DeleteUserByID)
}

// Login 用户登陆.
func (u User) Login(ctx *gin.Context) {
	in := new(requestmodel.UserReq)
	err := ctx.BindJSON(in)
	if err != nil {
		logrus.WithError(err).Error()
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ErrorInvalideRequest.Error()))
		return
	}

	err = validate.CheckUserLoginRequest(in)

	if err != nil {
		logrus.Errorf("check request failed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.ParamError, err.Error()))
		return
	}

	userService := service.NewUserService(utils.NewServiceContext(ctx, nil), nil, nil)

	tokenStr, res := userService.Login(in)
	if res != nil {
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.Header(constant.HEADERTOKEN, tokenStr)
	ctx.JSON(http.StatusOK, response.NewSuccessBaseResponse())
}

func (u User) CreateUser(ctx *gin.Context) {
	in := new(requestmodel.UserReq)

	err := ctx.BindJSON(in)
	if err != nil {
		logrus.WithError(err).Error("input invalide")
		ctx.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusBadRequest, custom.ErrorInvalideRequest.Error()))
		return
	}

	err = validate.CheckCreateUserRequest(in)

	if err != nil {
		logrus.Errorf("check request failed:%v", err)
		ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.ParamError, err.Error()))
		return
	}

	userService := service.NewUserService(utils.NewServiceContext(ctx, nil), nil, nil)

	err = userService.CreateUser(in)

	if err != nil {
		if err == custom.ErrorRecordExist {
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.RecordExist, custom.ErrorRecordExist.Error()))
			return
		} else {
			logrus.Errorf("insert user to db failed:%v", err)
			ctx.JSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusFailedDependency, custom.ErrorSaveToDBFailed.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, response.NewSuccessBaseResponse())

}

func InsertAdminUser(username, password, email string) error {
	username = utils.If(username == "", defaultUserName, username).(string)
	password = utils.If(password == "", security.HashMD5(defaultPassWord), password).(string)
	email = utils.If(email == "", defaultEmail, email).(string)

	in := &requestmodel.UserReq{UserName: username, PassWord: password, Email: email, NickName: username, Role: constant.ADMIN}

	err := validate.CheckCreateUserRequest(in)
	if err != nil {
		logrus.Errorf("check failed :%v", err)
		return err
	}

	return service.InitAdmin(in)
}
