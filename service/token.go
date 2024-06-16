package service

import (
	"gcrontab/entity/token"
	"gcrontab/model"
	rep "gcrontab/rep/token"
	"gcrontab/utils"

	"github.com/jinzhu/gorm"
)

type TokenService struct {
	ctx    *utils.ServiceContext
	db     *gorm.DB
	tokens []*token.Token
}

func NewTokenService(ctx *utils.ServiceContext, db *gorm.DB, tokens []*token.Token) *TokenService {
	if db == nil {
		db = model.DB()
	}
	if tokens == nil {
		return &TokenService{ctx, db, make([]*token.Token, 0, 10)}
	}
	return &TokenService{ctx, db, tokens}
}

func (ts *TokenService) FindSaltByToken(atoken string) (*token.Token, error) {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.FindTokenByPK(atoken)
}

func (ts *TokenService) UpdateToken(del string, newToken *token.Token) error {
	tokenRep := rep.NewTokenRep(ts.db.New())

	var err error
	err = tokenRep.DeleteTokenByPK(del)
	if err != nil {
		return err
	}

	return tokenRep.InsertToken(newToken)
}
