package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"videoGo/pkg/e"
	util2 "videoGo/pkg/util"
	"videoGo/repository/cache"
	typesRes "videoGo/types/res"
)

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token, err := jwt.Parse(c.GetHeader("Authorization"), func(token *jwt.Token) (interface{}, error) {
		//	return util.Key, nil
		//}, jwt.WithValidMethods(util.AlgList))
		token, err := jwt.ParseWithClaims(c.GetHeader("Authorization"), &util2.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return util2.Key, nil
		}, jwt.WithValidMethods(util2.AlgList))
		if token.Valid {
			if claims, ok := token.Claims.(*util2.MyClaims); ok {
				res, _ := cache.GetBlackBitUserID(c.Request.Context(), claims.UserID)
				if res == 1 {
					c.JSON(http.StatusForbidden, typesRes.NewResError(http.StatusForbidden, e.AuthorizeError))
					c.Abort()
					return
				}
				c.Set("userID", claims.UserID)
				c.Set("isAuditor", claims.IsAuditor)
				// c.Request = c.Request.WithContext(
				// 	ctl.NewContextWithUserIDKey(c.Request.Context(),
				// 	&ctl.UserInfo{Id: userID}))
				c.Next()
			} else {
				util2.Logger.Errorf("token myclaims pase get error, token:%v", token)
				c.JSON(http.StatusUnauthorized, typesRes.NewResError(http.StatusUnauthorized, e.AuthorizeError))
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
			c.JSON(http.StatusUnauthorized, typesRes.NewResError(http.StatusUnauthorized, e.AuthorizeError))
			c.Abort()
			return
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			// 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,"+
				"session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, "+
				"X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, "+
				"Content-Type, Pragma")
			// 允许跨域设置
			// 可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			// 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")
			// 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")
			//  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")
			// 设置返回格式是json
		}
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func auditorAuthenticMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("isAuditor") == true {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, typesRes.NewResError(http.StatusUnauthorized, e.AuthorizeError))
			c.Abort()
			return
		}
	}
}
