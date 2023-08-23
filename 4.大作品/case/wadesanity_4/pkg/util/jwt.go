package util

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	Key = []byte("1234abcd")
	AlgList = []string{"HS256"}
)

func NewJwt(userID uint) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":  "my-auth-server-video",
			"sub":  "userLogin",
			"foo":  2,
			"exp":  jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			"id": userID,
		})
	s, err := t.SignedString(Key)
	if err != nil {
		Logger.Errorf("jwt new error:%v", err)
	}
	return s
}
