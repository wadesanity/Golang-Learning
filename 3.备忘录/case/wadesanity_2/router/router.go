package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	engine.POST("/user", userRegister)
	engine.GET("/user", userLogin)


	engine.Use(jwtMiddleware)
	engine.GET("/todoList")
	engine.POST("/todoList")
	engine.

	return engine
}
