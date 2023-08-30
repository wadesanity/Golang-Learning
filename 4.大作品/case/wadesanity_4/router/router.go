package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"videoGo/case/wadesanity_4/api"
)

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	//engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 开启swag
	//store := cookie.NewStore([]byte("something-very-secret"))
	//sessionNames := []string{"a", "b"}
	//engine.Use(sessions.SessionsMany(sessionNames, store))
	//engine.Use(sessions.Sessions("mysession", store))

	store, _ := redis.NewStore(10, "tcp", "192.168.50.133:6379", "", []byte("secret"))
	engine.Use(sessions.Sessions("mysession", store))

	engine.Use(Cors()) //允许跨域
	v1 := engine.Group("/api/v1")
	{
		v1.GET("/cookie", func(c *gin.Context) {

			cookie, err := c.Cookie("gin_cookie")

			if err != nil {
				cookie = "NotSet"
				c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
			}

			fmt.Printf("Cookie value: %s \n", cookie)
		})
		v1.GET("/hello", func(c *gin.Context) {
			session := sessions.Default(c)

			if session.Get("hello") != "world" {
				fmt.Println("no hello world")
				session.Set("hello", "world")
				session.Save()
			}

			c.JSON(200, gin.H{"hello": session.Get("hello")})
		})
		v1.GET("/hello1", func(c *gin.Context) {
			sessionA := sessions.DefaultMany(c, "a")
			sessionB := sessions.DefaultMany(c, "b")

			if sessionA.Get("hello") != "world!" {
				sessionA.Set("hello", "world!")
				sessionA.Save()
			}

			if sessionB.Get("hello") != "world?" {
				sessionB.Set("hello", "world?")
				sessionB.Save()
			}

			c.JSON(200, gin.H{
				"a": sessionA.Get("hello"),
				"b": sessionB.Get("hello"),
			})
		})
		v1.GET("/incr", func(c *gin.Context) {
			session := sessions.Default(c)
			session.Options(sessions.Options{
				Path:     "/",
				Domain:   "",
				MaxAge:   60,
				Secure:   false,
				HttpOnly: true,
				SameSite: 1,
			})
			var count int
			v := session.Get("count")
			if v == nil {
				count = 0
			} else {
				count = v.(int)
				count++
			}
			session.Set("count", count)
			session.Save()
			c.JSON(200, gin.H{"count": count})
		})

		v1.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		v1.POST("/user/register", api.UserRegister)
		v1.POST("/user/login", api.UserLogin)
		authorized := v1.Group("/", jwtMiddleware())
		{
			authorized.POST("/user/changePwd", api.UserChangePwd)
			authorized.GET("/user/showInfo", api.UserShowInfo)
			authorized.POST("/user/changeAvatar", api.UserChangeAvatar)

			authorized.POST("/video/createOne", api.VideoCreateOne)
			authorized.GET("/video/list", api.VideoList)
			authorized.GET("/video/One", api.VideoOne)

			authorized.POST("/video/action", api.VideoAction)
			authorized.POST("/video/comment", api.VideoComment)
			authorized.GET("/video/comment/list", api.VideoCommentList)

			authorized.POST("/video/timeComment", api.VideoTimeCommentCreateOne)
			authorized.GET("/video/timeComment/list", api.VideoTimeCommentList)

			authorized.GET("/video/top", api.VideoTop)

			auditor := authorized.Group("/auditor", auditorAuthenticMiddleware())
			{
				auditor.POST("/video_audit", api.AuditorVideo)
				auditor.POST("/user_audit", api.AuditorUser)
				auditor.POST("/comment_audit", api.AuditorComment)
			}
		}
	}

	return engine
}
