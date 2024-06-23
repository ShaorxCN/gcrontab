package service

import (
	"gcrontab/cache"
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/entity/token"
	"gcrontab/entity/user"
	"gcrontab/interface/entity"
	"gcrontab/model"
	"gcrontab/rep/requestmodel"
	trep "gcrontab/rep/token"
	urep "gcrontab/rep/user"
	"gcrontab/security"
	"gcrontab/utils"
	"gcrontab/web/response"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	ctx   *utils.ServiceContext
	db    *gorm.DB
	users []*user.User
}

func dealUser(in *requestmodel.UserReq, id string) (*user.User, error) {
	userEntity := new(user.User)
	now := utils.Now()
	userEntity.NickName = in.NickName
	userEntity.UserName = in.UserName

	salt, err := utils.Nonce(10)
	if err != nil {
		return nil, err
	}

	userEntity.Salt = salt
	// 前端传递的应该就是md5 然后这边加盐再次摘要防止脱库
	userEntity.PassWord = security.HashSha256(in.PassWord + userEntity.Salt)
	userEntity.Status = constant.STATUSNORMAL
	userEntity.CreateAt = now.Format(constant.TIMELAYOUT)
	userEntity.ID = entity.NewEntityKey(uuid.New().String(), user.UserEntityType)
	userEntity.Email = in.Email
	userEntity.FailNotify = in.FailNotify
	if in.FailNotify == "" {
		userEntity.FailNotify = constant.NOTIFYOFF
	}
	userEntity.Creater = id
	userEntity.Role = in.Role

	return userEntity, nil
}

func NewUserService(ctx *utils.ServiceContext, db *gorm.DB, users []*user.User) *UserService {
	if db == nil {
		db = model.DB()
	}
	if users == nil {
		return &UserService{ctx, db, make([]*user.User, 0, 10)}
	}
	return &UserService{ctx, db, users}
}

func (us *UserService) FindUserByUserName(userName string) (*user.User, error) {
	userRep := urep.NewUserRep(us.db.New())
	return userRep.FindUserByUserName(userName)
}

func (us *UserService) FindUserByID(id string) (*user.User, error) {
	userRep := urep.NewUserRep(us.db.New())
	return userRep.FindUserByID(id)
}

func (us *UserService) CreateUser(in *requestmodel.UserReq) error {
	userRep := urep.NewUserRep(us.db.New())
	if _, err := userRep.FindUserByUserName(in.UserName); err == nil {
		return custom.ErrorRecordExist
	}

	userEntity, err := dealUser(in, us.ctx.Operator)
	if err != nil {
		logrus.Errorf("dealUser failed :%v", err)
		return custom.ErrorInternalServerError
	}

	err = userRep.InsertUser(userEntity)

	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "uix_tbl_user_user_name"` {
		return custom.ErrorRecordExist
	}

	return err
}

func (us *UserService) Login(in *requestmodel.UserReq) (string, *response.BaseResponse) {
	userRep := urep.NewUserRep(us.db.New())
	var ue *user.User
	var err error
	if ue, err = userRep.FindUserByUserName(in.UserName); err != nil {
		logrus.WithError(err).Error("invalide user")
		return "", response.NewBusinessFailedBaseResponse(custom.StatusNonAuthoritativeInfo, custom.ErrorLoginFailed.Error())
	}

	dealPwd := in.PassWord + ue.Salt

	pwdCheck := security.HashSha256(dealPwd)

	if pwdCheck != ue.PassWord {
		logrus.WithError(err).Error("login failed")
		return "", response.NewBusinessFailedBaseResponse(custom.StatusNonAuthoritativeInfo, custom.ErrorLoginFailed.Error())
	}

	c := new(utils.Claims)

	c.UID = ue.ID.GetIDValue()
	c.NickName = ue.NickName
	var secret string
	secret, err = utils.Nonce(10)
	if err != nil {
		logrus.Errorf("generate nonce failed:%v", err)
		return "", response.NewBusinessFailedBaseResponse(custom.InternalServerError, custom.ErrorInternalServerError.Error())
	}

	c.Secret = secret
	c.Role = ue.Role

	now := utils.Now()
	exp := now.Add(time.Duration(constant.TokenTTL) * time.Second).Format(constant.TIMELAYOUT)
	dead := now.Add(time.Duration(10*constant.TokenTTL) * time.Second).Format(constant.TIMELAYOUT)
	c.Exp = exp
	c.DeadLine = dead
	tokenStr, err := utils.GenToken(c)
	if err != nil {
		logrus.Errorf("generate  token failed:%v", err)
		return "", response.NewBusinessFailedBaseResponse(custom.InternalServerError, custom.ErrorInternalServerError.Error())
	}

	tokenE := &token.Token{UserID: ue.ID.GetIDValue(), CreateTime: now.Format(constant.TIMELAYOUT), Token: tokenStr, Salt: c.Secret}

	tokenRep := trep.NewTokenRep(us.db.New())
	err = tokenRep.SaveToken(tokenE)
	if err != nil {
		logrus.Errorf("save token failed:%v", err)
		return "", response.NewBusinessFailedBaseResponse(custom.StatusFailedDependency, custom.ErrorSaveToDBFailed.Error())
	}
	cache.RemoveUIDSalt(c.UID)
	cache.SetSalt(c.UID, tokenStr, c.Secret)

	return tokenStr, nil
}

func InitAdmin(in *requestmodel.UserReq) error {
	userRep := urep.NewUserRep(model.DB().New())

	if _, err := userRep.FindUserByUserName(in.UserName); err == nil || err != gorm.ErrRecordNotFound {
		return err
	}
	userEntity, err := dealUser(in, uuid.Nil.String())
	if err != nil {
		logrus.Errorf("initadmin failed :%v", err)
		return err
	}

	return userRep.InsertUser(userEntity)

}
