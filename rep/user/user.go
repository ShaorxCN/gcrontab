package user

import (
	"gcrontab/constant"
	"gcrontab/entity/user"
	"gcrontab/model"

	"github.com/jinzhu/gorm"
)

type UserRep struct {
	db *gorm.DB
}

func NewUserRep(db *gorm.DB) *UserRep {
	if db == nil {
		db = model.DB()
	}

	return &UserRep{db}
}

// FindTokenByPK 根据userName查找
func (r *UserRep) FindUserByUserName(userName string) (*user.User, error) {
	u := new(model.DBUser)
	err := r.db.Model(u).Where("user_name = ? and status = ?", userName, constant.STATUSNORMALDB).First(u).Error

	if err != nil {
		return nil, err
	}
	return user.FromDBUserModel(u)
}
