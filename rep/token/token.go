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

func (r *TokenRep) SaveToken(te *token.Token) error {
	db := r.db
	dt, err := te.ToDBTokenModel()
	if err != nil {
		return err
	}
	return db.Save(dt).Error
}

func (r *TokenRep) DeleteTokenByPK(id string) error {
	db := r.db
	return db.Delete(model.DBToken{}, "user_id = ? ", id).Error
}

func (r *TokenRep) DelTokenByUIDAndToken(id, token string) error {
	db := r.db
	return db.Delete(model.DBToken{}, "user_id = ? and token = ?", id, token).Error
}

func (r *TokenRep) FindTokenByUserID(id string) (*token.Token, error) {
	dbRet := new(model.DBToken)
	err := r.db.Where("user_id = ?", id).First(dbRet).Error
	if err != nil {
		return nil, err
	}
	return token.FromDBTokenModel(dbRet), nil
}

func (r *TokenRep) UpdateTokenByToken(id, old, newStr string) error {
	dbt := new(model.DBToken)
	return r.db.Model(dbt).Where("user_id = ? and token = ?", id, old).Update("token", newStr).Error
}

func (r *TokenRep) FindTokenByUserIDAndToken(id, tokenStr string) (*token.Token, error) {
	dbRet := new(model.DBToken)
	err := r.db.Where("user_id = ? and token = ?", id, tokenStr).First(dbRet).Error
	if err != nil {
		return nil, err
	}
	return token.FromDBTokenModel(dbRet), nil
}

// DeleteTokenByUserID 根据userID 删除token
func (r *TokenRep) DeleteTokenByUserID(id string) error {
	err := r.db.Where("user_id = ?", id).Delete(&model.DBToken{}).Error
	return err
}
