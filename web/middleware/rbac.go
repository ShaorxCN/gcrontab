package middleware

import (
	"fmt"
	"gcrontab/cache"
	"gcrontab/casbin"
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/entity/token"
	"gcrontab/service"
	"gcrontab/utils"
	"gcrontab/web/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func checkAndUpdateToken(tokenStr string, c *gin.Context) error {
	now := utils.Now()

	if tokenStr == "" {
		return custom.ErrorInvalideAccessToken
	}

	uid, goon := utils.ParseClaimWithoutValidate(tokenStr)

	if !goon {
		return custom.ErrorInvalideAccessToken
	}

	tokenService := service.NewTokenService(utils.NewServiceContext(c, nil), nil, nil)
	// cache?
	var saltCache interface{}
	var ok bool
	saltCache, ok = cache.GetSaltByUID(uid)
	if !ok {
		te, err := tokenService.FindTokenByUID(uid)
		if err != nil {
			logrus.Errorf("find token by uid:[%s] error:%v", uid, err)
			return custom.ErrorInvalideAccessToken
		}

		saltCache = te.Salt
	}

	cm, err := utils.ValideToken(tokenStr, saltCache.(string))
	if err != nil {
		logrus.Errorf("valide token:[%s] error:%v", tokenStr, err)
		return custom.ErrorInvalideAccessToken
	}

	expTime, deadTime, err := parseExpAndDeadLine(cm.Exp, cm.DeadLine)
	if err != nil {
		logrus.Errorf("valide token:[%s] error:%v", tokenStr, err)
		return custom.ErrorInvalideAccessToken
	}

	if deadTime.Before(now) {
		logrus.Errorf("token is expired,userId:[%s],nickName:[%s]", cm.UID, cm.NickName)
		if err := tokenService.DelTokenByUID(uid); err != nil {
			logrus.Errorf("token[%s] del failed:%v", tokenStr, err)
		}

		cache.RemoveSalt(uid)
		return custom.ErrorInvalideAccessToken
	}

	userService := service.NewUserService(utils.NewServiceContext(c, nil), nil, nil)

	user, err := userService.FindUserByID(cm.UID)
	if err != nil {
		logrus.Errorf("find user by name[%s] failed:%v", cm.UID, err)
		return custom.ErrorInvalideAccessToken
	}

	role := utils.If(user.Role == "", constant.ANONYMOUS, user.Role).(string)
	c.Request.Header.Set(constant.HEADEROPERATORNAME, user.NickName)
	c.Request.Header.Set(constant.HEADEROPERATACCT, user.UserName)
	c.Request.Header.Set(constant.HEADEROPERATORROLE, role)
	c.Request.Header.Set(constant.HEADEROPERATORID, user.ID.GetIDValue())

	if now.Before(expTime) {
		return nil
	}

	salt, err := utils.Nonce(10)
	if err != nil {
		logrus.Errorf("generate nonce failed:%v", err)
		return err
	}

	// 此处新老token deadline不变
	// TODO: 这边jwt已经存库校验了 deadline是否也存库算了？或者不存 只根据createat 计算？
	newClaims := &utils.Claims{
		UID:      cm.UID,
		Exp:      now.Add(time.Duration(constant.TokenTTL) * time.Second).Format(constant.TIMELAYOUT),
		NickName: user.NickName,
		DeadLine: cm.DeadLine,
		Secret:   salt,
		Role:     user.Role,
	}
	newToken, err := utils.GenToken(newClaims)

	if err != nil {
		logrus.Errorf("generate token failed error:%v", err)
		return custom.ErrorInvalideAccessToken
	}

	c.Header(constant.HEADERTOKEN, newToken)

	newTokenE := &token.Token{
		UserID:     cm.UID,
		CreateTime: now.Format(constant.TIMELAYOUT),
		Token:      newToken,
		Salt:       salt,
	}

	err = tokenService.UpdateToken(newTokenE)
	if err != nil {
		logrus.Errorf("update token failed:%v", err)
		return custom.ErrorSaveToDBFailed
	}

	cache.SetSalt(cm.UID, salt)

	return nil
}

// 检查角色权限
func checkCasbinRBAC(ctx *gin.Context) error {
	if ctx.GetHeader(constant.HEADEROPERATACCT) == constant.ADMIN {
		return nil
	}
	role := ctx.GetHeader(constant.HEADEROPERATORROLE)
	obj := ctx.Request.RequestURI
	act := ctx.Request.Method
	ok := casbin.CasEnforecer.Enforce(role, obj, act)
	if ok {
		return nil
	}
	return fmt.Errorf("request failed: [sub:%s,obj:%s,act:%s]", role, obj, act)
}

// CheckPermission 返回一个中间件，检查用户登陆状态
func CheckPermission(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		uri := c.Request.URL.Path

		if method == "POST" && uri == "/users/login" {
			return
		}

		token := c.GetHeader(constant.HEADERTOKEN)
		if err := checkAndUpdateToken(token, c); err != nil {
			logger.Error("can't find valide access-token header")
			c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusNonAuthoritativeInfo, err.Error()))
			return
		}

		if err := checkCasbinRBAC(c); err != nil {
			logger.Errorf("check permission failed:[uid:%s,path:%s,method:%s]", c.GetHeader(constant.HEADEROPERATORID), c.Request.RequestURI, c.Request.Method)
			c.AbortWithStatusJSON(http.StatusOK, response.NewBusinessFailedBaseResponse(custom.StatusForbidden, custom.ErrorForbidden.Error()))
			return
		}
	}
}

func parseExpAndDeadLine(expr, deadline string) (exp, dead time.Time, err error) {
	exp, err = time.ParseInLocation(constant.TIMELAYOUT, expr, utils.DefaultLocation)
	if err != nil {
		return
	}

	dead, err = time.ParseInLocation(constant.TIMELAYOUT, deadline, utils.DefaultLocation)

	return
}
