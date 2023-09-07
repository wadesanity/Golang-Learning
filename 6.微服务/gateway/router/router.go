package router

import (
	"gateway/api"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	engine.Use(cors())
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/user_register", api.UserRegister)
		v1.POST("/user_login", api.UserLogin)
		authorized := v1.Group("/", jwtMiddleware())
		{
			authorized.POST("/user_changePwd", api.UserChangePwd)
			authorized.GET("/user_showInfo", api.UserShowInfo)
			authorized.POST("/user_changeAvatar", api.UserChangeAvatar)
			authorized.GET("/user_list", api.UserList)
		}

	}
}
