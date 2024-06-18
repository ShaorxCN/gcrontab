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

// FindUserByID 根据id查找
func (r *UserRep) FindUserByID(id string) (*user.User, error) {
	u := new(model.DBUser)
	err := r.db.Model(u).Where(" id = ? and status = ?", id, constant.STATUSNORMALDB).First(u).Error

	if err != nil {
		return nil, err
	}
	return user.FromDBUserModel(u)
}

// InsertUser 创建user
func (r *UserRep) InsertUser(u *user.User) error {
	dbu, err := u.ToDBUserModel()
	if err != nil {
		return err
	}

	return r.db.Create(dbu).Error
}
