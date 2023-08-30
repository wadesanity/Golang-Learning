package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	Key     = []byte("1234abcd")
	AlgList = []string{"HS256"}
)

type MyClaims struct {
	UserID    uint
	IsAuditor bool
	jwt.RegisteredClaims
}

func NewJwt(userID uint, isAuditor bool) string {
	//t := jwt.NewWithClaims(jwt.SigningMethodHS256,
	//	jwt.MapClaims{
	//		"iss": "my-auth-server-video",
	//		"sub": "userLogin",
	//		"foo": 2,
	//		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	//		"id":  userID,
	//	})
	claims := MyClaims{
		userID,
		isAuditor,
		jwt.RegisteredClaims{ // Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "test",
		}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(Key)
	if err != nil {
		Logger.Errorf("jwt new error:%v", err)
	}
	return s
}
