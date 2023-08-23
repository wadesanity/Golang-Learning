package router

import (
	"github.com/gin-gonic/gin"
	user "videoGo/case/wadesanity_4/api/user"
)

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	engine.POST("/user/register", user.UserRegister)
	engine.POST("/user/login", user.UserLogin)

	authorized := engine.Group("/", jwtMiddleware())
	{
		authorized.POST("/user/changPwd", user.UserChangePwd)
		//authorized.POST("/task", taskPost)
		//authorized.PUT("/task", taskPut)
		//authorized.DELETE("/task", taskDelete)
	}

	return engine
}
