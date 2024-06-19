package utils

import (
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	UID      string
	Exp      string //	过期时间
	NickName string //  昵称
	DeadLine string //  最后可以refresh token 时间  double exp
	Secret   string //  salt
	Role     string
}

// GenToken 生成token
func GenToken(c *Claims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["uid"] = c.UID
	claims["exp"] = c.Exp
	claims["nickName"] = c.NickName
	claims["deadLine"] = c.DeadLine
	claims["role"] = c.Role

	token.Claims = claims

	return token.SignedString([]byte(c.Secret))
}

// ValideToken 校验token
func ValideToken(tokenStr, secret string) (*Claims, error) {
	tv, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if m, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Errorf("parse token error")
			return nil, errors.New("method error")
		} else if m != jwt.SigningMethodHS256 {
			return nil, errors.New("method error")
		}
		return []byte(secret), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		logrus.Errorf("parse token error:%v", err)
		return nil, errors.New("parse and valide token error")
	}

	if claims, ok := tv.Claims.(jwt.MapClaims); ok && tv.Valid {
		return &Claims{
			UID:      claims["uid"].(string),
			Exp:      claims["exp"].(string),
			NickName: claims["nickName"].(string),
			DeadLine: claims["deadLine"].(string),
			Role:     claims["role"].(string),
		}, nil
	}

	return nil, errors.New("parse and valide token error")
}
