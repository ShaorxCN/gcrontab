package service

import (
	"gcrontab/entity/token"
	"gcrontab/entity/user"
	"gcrontab/model"
	rep "gcrontab/rep/user"
	"gcrontab/utils"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	ctx    *utils.ServiceContext
	db     *gorm.DB
	tokens []*token.Token
}

func NewUserService(ctx *utils.ServiceContext, db *gorm.DB, tokens []*token.Token) *UserService {
	if db == nil {
		db = model.DB()
	}
	if tokens == nil {
		return &UserService{ctx, db, make([]*token.Token, 0, 10)}
	}
	return &UserService{ctx, db, tokens}
}

func (us *UserService) FindUserByUserName(userName string) (*user.User, error) {
	userRep := rep.NewUserRep(us.db.New())
	return userRep.FindUserByUserName(userName)
}
