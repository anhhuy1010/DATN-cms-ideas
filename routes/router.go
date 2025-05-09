package routes

import (
	"net/http"

	"github.com/anhhuy1010/DATN-cms-ideas/controllers"

	docs "github.com/anhhuy1010/DATN-cms-ideas/docs"
	"github.com/anhhuy1010/DATN-cms-ideas/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteInit(engine *gin.Engine) {
	userCtr := new(controllers.UserController)
	ideaCtr := new(controllers.IdeasController)

	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.Use(middleware.Recovery())
	docs.SwaggerInfo.BasePath = "/v1"

	apiV1 := engine.Group("/v1")
	apiV1.Use(middleware.RequestLog())

	// ✅ Các route cần xác thực nằm trong group này
	protected := apiV1.Group("/")
	protected.Use(controllers.RoleMiddleware())
	{
		protected.GET("/customer", userCtr.List)
		protected.GET("/customer/:uuid", userCtr.Detail)
		protected.POST("/customer", userCtr.Create)
		protected.PUT("/customer/:uuid", userCtr.Update)
		protected.PUT("/customer/:uuid/update-status", userCtr.UpdateStatus)
		protected.DELETE("/customer/:uuid", userCtr.Delete)
		//////////////////////////////////////////////////////////////////////
		protected.GET("/ideas", ideaCtr.List)
		protected.GET("/ideas/:uuid", ideaCtr.Detail)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
