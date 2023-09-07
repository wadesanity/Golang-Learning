package router

import (
	"errors"
	"fmt"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/res"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError))
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &util.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return util.Key, nil
		}, jwt.WithValidMethods(util.AlgList))
		if token.Valid {
			if claims, ok := token.Claims.(*util.MyClaims); ok {
				util.Logger.Debugf("claims.UserID:%v", claims.UserID)
				c.Set("userID", uint(claims.UserID))
				c.Next()
			} else {
				util.Logger.Errorf("token myclaims pase get error, token:%v", token)
				c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError))
				c.Abort()
				return
			}
		} else {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				fmt.Println("That's not even a token")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				// Invalid signature
				fmt.Println("Invalid signature")
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}
			c.JSON(http.StatusUnauthorized, res.NewResError(http.StatusUnauthorized, e.AuthorizeError))
			c.Abort()
			return
		}
	}
}
