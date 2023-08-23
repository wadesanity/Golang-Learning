package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"videoGo/case/wadesanity_4/pkg/util"
)

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
			return util.Key, nil
		}, jwt.WithValidMethods(util.AlgList))
		if token != nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIDAny, ok := claims["id"]; ok {
					var userID uint
					switch s := userIDAny.(type) {
					case float64:
						userID = uint(s)
					case int:
						userID = uint(s)
					case int64:
						userID = uint(s)
					case int32:
						userID = uint(s)
					case float32:
						userID = uint(s)
					}
					c.Set("userID", userID)
					c.Next()
				} else {
					util.Logger.Errorf("token claims id get error, claims:%v", claims)
					c.Abort()
				}

			} else {
				util.Logger.Errorf("token mapclaims pase error, token:%v", token)
				c.Abort()
			}
		} else {
			util.Logger.Errorf("token pase error:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			})
			c.Abort()
		}

	}
}