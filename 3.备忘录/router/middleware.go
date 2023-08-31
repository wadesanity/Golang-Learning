package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"todolistGo/common"
	"todolistGo/logger"
)

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
			//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			//}
			return common.Key, nil
		}, jwt.WithValidMethods(common.AlgList))
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
					logger.Logger.Errorf("token claims id get error, claims:%v", claims)
					c.Abort()
				}

			} else {
				logger.Logger.Errorf("token mapclaims pase error, token:%v", token)
				c.Abort()
			}
		} else {
			logger.Logger.Errorf("token pase error:%v", err)
			c.JSON(http.StatusBadRequest, common.ResponseError{
				Status: http.StatusBadRequest,
				Error:  err.Error(),
			})
			c.Abort()
		}

	}
}
