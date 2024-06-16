package user

import (
	"gcrontab/constant"
	"gcrontab/entity"
	entitier "gcrontab/interface/entity"
	"gcrontab/model"
	"gcrontab/utils"
	"time"

	"github.com/google/uuid"
)

const (
	UserEntityType = "tbl_user"
)

// User 是任务实体。
type User struct {
	entity.BaseEntity
	UserName   string `json:"userName,omitempty"`
	NickName   string `json:"nickName,omitempty"`
	PassWord   string `json:"-"`
	Salt       string `json:"-"`
	Status     string `json:"status,omitempty"`
	FailNotify string `json:"failNotify,omitempty"`
	Email      string `json:"email,omitempty"`
	Creater    string `json:"-"`
	Role       string `json:"role,omitempty"`
}

// ToDBUserModel user entity -> model
func (u *User) ToDBUserModel() (*model.DBUser, error) {

	user := &model.DBUser{}

	uid, err := uuid.Parse(u.ID.GetIDValue())
	if err != nil {
		return nil, err
	}

	user.ID = uid
	create, err := time.ParseInLocation(constant.TIMELAYOUT, u.CreateAt, utils.DefaultLocation)
	if err != nil {
		return nil, err
	}

	if u.UpdateAt != "" {
		update, err := time.ParseInLocation(constant.TIMELAYOUT, u.UpdateAt, utils.DefaultLocation)
		if err != nil {
			return nil, err
		}
		user.UpdateAt = &update
	}

	user.CreateAt = &create

	user.UserName = u.UserName
	user.NickName = u.NickName
	user.PassWord = u.PassWord
	user.Salt = u.Salt
	user.Status = u.Status
	user.Email = u.Email
	user.FailNotify = u.FailNotify
	user.Creater = u.Creater
	user.Role = u.Role

	return user, nil
}

// FromDBUserModel model-entity
func FromDBUserModel(u *model.DBUser) *User {
	e := new(User)

	e.ID = entitier.NewEntityKey(u.ID.String(), UserEntityType)
	e.CreateAt = u.CreateAt.Format(constant.TIMELAYOUT)
	if e.UpdateAt != "" {
		e.UpdateAt = u.UpdateAt.Format(constant.TIMELAYOUT)
	}
	e.UserName = u.UserName
	e.NickName = u.NickName
	e.PassWord = u.PassWord
	e.Salt = u.Salt
	e.Status = u.Status
	e.Email = u.Email
	e.FailNotify = u.FailNotify
	e.Creater = u.Creater
	e.Role = u.Role
	return e
}
