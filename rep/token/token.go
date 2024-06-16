package token

import (
	"gcrontab/entity/token"
	"gcrontab/model"

	"github.com/jinzhu/gorm"
)

type TokenRep struct {
	db *gorm.DB
}

func NewTokenRep(db *gorm.DB) *TokenRep {
	if db == nil {
		db = model.DB()
	}

	return &TokenRep{db}
}

// FindTokenByPK 根据主键查找
func (r *TokenRep) FindTokenByPK(ts string) (*token.Token, error) {
	t := new(model.DBToken)
	t.Token = ts
	err := r.db.Model(t).First(t).Error
	if err != nil {
		return nil, err
	}

	return token.FromDBTokenModel(t), nil
}

func (r *TokenRep) InsertToken(te *token.Token) error {
	db := r.db
	dt, err := te.ToDBTokenModel()
	if err != nil {
		return err
	}
	return db.Create(dt).Error
}

func (r *TokenRep) DeleteTokenByPK(tokenStr string) error {
	db := r.db
	return db.Delete(model.DBToken{}, "token = ? ", tokenStr).Error
}
