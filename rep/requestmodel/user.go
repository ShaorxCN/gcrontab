package requestmodel

import "time"

type UserReq struct {
	UserName      string    `json:"userName,omitempty"`
	NickName      string    `json:"nickName,omitempty"`
	PassWord      string    `json:"passWord,omitempty"` // md5
	Email         string    `json:"email,omitempty"`
	FailNotify    string    `json:"failNotify,omitempty"`
	NewPassWord   string    `json:"newPassWord,omitempty"`
	Salt          string    `json:"-"`
	UpdateTimeUse time.Time `json:"-"`
	Role          string    `json:"role,omitempty"`
}
