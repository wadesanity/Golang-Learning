package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	engine.POST("/user/register", userRegister)
	engine.POST("/user/login", userLogin)

	authorized := engine.Group("/", jwtMiddleware())
	{
		authorized.GET("/task", taskGet)
		authorized.POST("/task", taskPost)
		authorized.PUT("/task", taskPut)
		authorized.DELETE("/task", taskDelete)
	}

	return engine
}
