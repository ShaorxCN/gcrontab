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

func (ts *TokenService) FindTokenByUID(id string) (*token.Token, error) {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.FindTokenByUserID(id)
}

func (ts *TokenService) FindTokenByUIDAndToken(id, token string) (*token.Token, error) {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.FindTokenByUserIDAndToken(id, token)
}

func (ts *TokenService) DelTokenByUIDAndToken(id, token string) error {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.DelTokenByUIDAndToken(id, token)

}

func (ts *TokenService) UpdateTokenByToken(uid, old, newStr string) error {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.UpdateTokenByToken(uid, old, newStr)
}

func (ts *TokenService) UpdateToken(newToken *token.Token) error {
	tokenRep := rep.NewTokenRep(ts.db.New())
	return tokenRep.SaveToken(newToken)
}
