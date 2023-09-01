package router

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	api2 "videoGo/api"
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
		v1.POST("/user/register", api2.UserRegister)
		v1.POST("/user/login", api2.UserLogin)
		authorized := v1.Group("/", jwtMiddleware())
		{
			authorized.POST("/user/changePwd", api2.UserChangePwd)
			authorized.GET("/user/showInfo", api2.UserShowInfo)
			authorized.POST("/user/changeAvatar", api2.UserChangeAvatar)
			authorized.GET("/user/list", api2.UserList)

			authorized.POST("/video/createOne", api2.VideoCreateOne)
			authorized.GET("/video/list", api2.VideoList)
			authorized.GET("/video/One", api2.VideoOne)

			authorized.POST("/video/action", api2.VideoAction)
			authorized.POST("/video/comment", api2.VideoComment)
			authorized.GET("/video/comment/list", api2.VideoCommentList)

			authorized.POST("/video/timeComment", api2.VideoTimeCommentCreateOne)
			authorized.GET("/video/timeComment/list", api2.VideoTimeCommentList)

			authorized.GET("/video/top", api2.VideoTop)

			auditor := authorized.Group("/auditor", auditorAuthenticMiddleware())
			{
				auditor.POST("/video_audit", api2.AuditorVideo)
				auditor.POST("/user_audit", api2.AuditorUser)
				auditor.POST("/comment_audit", api2.AuditorComment)
			}
		}
	}

	return engine
}
