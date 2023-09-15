package router

import (
	"errors"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/res"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError.Error()))
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &util.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return util.Key, nil
		}, jwt.WithValidMethods(util.AlgList))
		if token.Valid {
			if claims, ok := token.Claims.(*util.MyClaims); ok {
				util.Logger.WithFields(logrus.Fields{
					"trace_id": c.Request.Context().Value(util.TraceIdKey),
					"user_id":  claims.UserID,
				}).Debugln("jwt check UserID.")
				c.Set("userID", uint(claims.UserID))
				c.Next()
			} else {
				util.Logger.WithFields(logrus.Fields{
					"token": token,
				}).Errorln("token myClaims parse get error.")
				c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError.Error()))
				c.Abort()
				return
			}
		} else {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				util.Logger.WithFields(logrus.Fields{
					"token":  token,
					"detail": err,
				}).Errorln("tokenInvalid: That's not even a token.")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				util.Logger.WithFields(logrus.Fields{
					"token":  token,
					"detail": err,
				}).Errorln("tokenInvalid: Invalid signature.")
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				util.Logger.WithFields(logrus.Fields{
					"token":  token,
					"detail": err,
				}).Errorln("tokenInvalid: Timing is everything.")
			} else {
				util.Logger.WithFields(logrus.Fields{
					"token":  token,
					"detail": err,
				}).Errorln("tokenInvalid: Couldn't handle this token.")
			}
			c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError.Error()))
			c.Abort()
			return
		}
	}
}
