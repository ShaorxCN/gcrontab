package middleware

import (
	"gcrontab/constant"
	"gcrontab/utils"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	salt, err := utils.Nonce(10)
	if err != nil {
		t.Errorf("generate nonce failed:%v", err)
	}

	newClaims := &utils.Claims{
		UID:      uuid.New().String(),
		Exp:      time.Now().Format(constant.TIMELAYOUT),
		NickName: "testnickname",
		DeadLine: time.Now().Format(constant.TIMELAYOUT),
		Secret:   salt,
		Role:     "test",
	}
	newToken, err := utils.GenToken(newClaims)

	if err != nil {
		t.Errorf("generate token failed error:%v", err)
	}
	t.Log(newToken)
}
