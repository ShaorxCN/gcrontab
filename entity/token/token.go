package token

import (
	"gcrontab/constant"
	"gcrontab/model"
	"gcrontab/utils"
	"time"
)

const (
	// TokenEntityType entity type
	TokenEntityType = "token"
)

// Token 实体
type Token struct {
	UserID     string `json:"userID,omitempty" url:"userID,omitempty"`
	CreateTime string `json:"-" url:"-"`
	Token      string `json:"token,omitempty" url:"token,omitempty"`
	Salt       string `json:"-" url:"-"`
}

// ToDBTokenModel 将token实体转换为DBToken数据库模型。
func (t *Token) ToDBTokenModel() (*model.DBToken, error) {
	d := &model.DBToken{}
	d.UserID = t.UserID
	d.Token = t.Token
	d.Salt = t.Salt

	create, err := time.ParseInLocation(constant.TIMELAYOUT, t.CreateTime, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}

	d.CreateTime = &create
	return d, nil
}

// FromDBTokenModel 将 token Model转为实体。
func FromDBTokenModel(d *model.DBToken) *Token {
	t := new(Token)
	t.CreateTime = d.CreateTime.In(utils.DefaultLocation).Format(constant.TIMELAYOUT)
	t.UserID = d.UserID
	t.Token = d.Token
	t.Salt = d.Salt
	return t
}
