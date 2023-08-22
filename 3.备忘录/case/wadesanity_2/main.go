package main

import (
	"github.com/swaggo/files"           // swagger embed files
	"github.com/swaggo/gin-swagger"     // gin-swagger middleware
	"todolistGo/case/wadesanity_2/docs" // docs is generated by Swag CLI, you have to import it.
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/router"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	e:=router.SetupRouter()

	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := e.Run()
	if err != nil {
		logger.Logger.Fatalf("gin run error:%v",err)
		return 
	}
}