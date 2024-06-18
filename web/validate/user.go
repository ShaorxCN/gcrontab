package validate

import (
	"gcrontab/constant"
	"gcrontab/custom"
	"gcrontab/rep/requestmodel"
	"gcrontab/utils"
	"regexp"
)

var roleSlice = []string{constant.ADMIN, constant.TASKADMIN, constant.USER, ""}

func IsValidMD5(s string) bool {
	if len(s) != 32 {
		return false
	}
	match, _ := regexp.MatchString("^[a-fA-F0-9]{32}$", s)
	return match
}

// CheckCreateUserRequest 创建用户request的基础校验
func CheckCreateUserRequest(in *requestmodel.UserReq) (err error) {
	switch {
	case in.UserName == "" || len(in.UserName) > 255:
		err = custom.ParamErrorReturn(userName)
	case in.PassWord == "" || !IsValidMD5(in.PassWord):
		err = custom.ParamErrorReturn(passWord)
	case in.NickName == "" || len(in.NickName) > 255:
		err = custom.ParamErrorReturn(nickName)
	case in.Email == "" || len(in.Email) > 126:
		err = custom.ParamErrorReturn(email)
	case !utils.StrInSlice(in.FailNotify, constant.NotifySlice):
		err = custom.ParamErrorReturn(failNotify)
	case !utils.StrInSlice(in.Role, constant.RoleSlice):
		err = custom.ParamErrorReturn(role)
	}

	if in.Role == "" {
		in.Role = constant.USER
	}

	return
}
