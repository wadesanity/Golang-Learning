package common

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"todolistGo/logger"
)

var (
	Key     = []byte("1234abcd")
	AlgList = []string{"HS256"}
)

func NewJwt(userID uint) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "my-auth-server-case",
			"sub": "userLogin",
			"foo": 2,
			"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			"id":  userID,
		})
	s, err := t.SignedString(Key)
	if err != nil {
		logger.Logger.Errorf("jwt new error:%v", err)
	}
	return s
}
