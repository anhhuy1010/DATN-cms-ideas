package routes

import (
	"net/http"
	"time"

	"github.com/anhhuy1010/DATN-cms-ideas/controllers"
	"github.com/gin-contrib/cors"

	docs "github.com/anhhuy1010/DATN-cms-ideas/docs"
	"github.com/anhhuy1010/DATN-cms-ideas/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RouteInit(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                  // Cho phép tất cả các domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},            // Các phương thức HTTP cho phép
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-token"}, // Các header cho phép
		AllowCredentials: true,                                                           // Cho phép chia sẻ cookie nếu cần
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           256 * time.Hour, // Thời gian tối đa mà trình duyệt sẽ cache thông tin CORS
	}))
	engine.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	userCtr := new(controllers.UserController)
	ideaCtr := new(controllers.IdeasController)

	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.Use(middleware.Recovery())
	docs.SwaggerInfo.BasePath = "/v1"

	apiV1 := engine.Group("/v1")
	apiV1.GET("/ideas", ideaCtr.List)
	apiV1.GET("/ideas/search", ideaCtr.SearchIdeasByTitle)

	apiV1.GET("/ideas/:uuid/price", ideaCtr.Price)
	apiV1.PUT("/ideas/:uuid/updatestatus", ideaCtr.UpdateStatus)
	//////////////////////////////////////////////////////////////////////////// admin
	apiV1.PUT("/ideas/update-status-admin/:uuid", ideaCtr.UpdateStatusAdmin)
	apiV1.GET("/ideas/admin-list", ideaCtr.ListForAdmin)
	apiV1.PUT("/ideas/:uuid/accept-status", ideaCtr.AcceptStatus)
	apiV1.PUT("/ideas/:uuid/reject-status", ideaCtr.RejectStatus)
	apiV1.GET("/ideas/:uuid/detail-watting", ideaCtr.DetailWatting)
	apiV1.GET("/ideas/total", ideaCtr.TotalSoldIdeasPrice)
	apiV1.GET("/ideas/countideas", ideaCtr.CountForAdmin)

	apiV1.Use(middleware.RequestLog())

	// ✅ Các route cần xác thực nằm trong group này
	protected := apiV1.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{

		protected.GET("/ideas/list", ideaCtr.ListID)
		protected.POST("/ideas", ideaCtr.PostIdea)
		protected.GET("/ideas/my-list", ideaCtr.GetMyIdeas)
		protected.GET("/ideas/:uuid/my-detail", ideaCtr.GetMyIdeaDetail)
		protected.GET("/ideas/:uuid/my-detail-for-update", ideaCtr.DetailForUpdate)
		protected.DELETE("/ideas/:uuid/delete-myidea", ideaCtr.DeleteMyIdea)
		protected.PUT("/ideas/:uuid/update-myidea", ideaCtr.UpdateMyIdea)
		protected.GET("/ideas/buy-idea", ideaCtr.GetPurchasedIdeas)
		////////////////////////////////////////////////////////////////////////
		protected.GET("/customer", userCtr.List)
		protected.GET("/customer/:uuid", userCtr.Detail)
		protected.POST("/customer", userCtr.Create)
		protected.PUT("/customer/:uuid", userCtr.Update)
		// protected.PUT("/customer/:uuid/update-status", userCtr.UpdateStatus)
		protected.DELETE("/customer/:uuid", userCtr.Delete)
		//////////////////////////////////////////////////////////////////////
		protected.POST("/ideas/add-favorite", ideaCtr.AddFavorite)
		protected.GET("/ideas/list-favorite", ideaCtr.ListFavorite)
		protected.DELETE("/favorite/:post_uuid/delete", ideaCtr.DeleteFavorite)

		protectedWithDate := protected.Group("/")
		protectedWithDate.Use(middleware.JWTDateCheckMiddleware())
		{
			protectedWithDate.GET("/ideas/:uuid", ideaCtr.Detail)
		}

	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
	})
}
