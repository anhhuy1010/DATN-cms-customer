package routes

import (
	"net/http"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/controllers"
	"github.com/gin-contrib/cors"

	docs "github.com/anhhuy1010/DATN-cms-customer/docs"
	"github.com/anhhuy1010/DATN-cms-customer/middleware"
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
		MaxAge:           24 * time.Hour, // Thời gian tối đa mà trình duyệt sẽ cache thông tin CORS
	}))
	engine.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})
	userCtr := new(controllers.UserController)

	engine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.Use(middleware.Recovery())
	docs.SwaggerInfo.BasePath = "/v1"

	apiV1 := engine.Group("/v1")
	// ❌ Không có RoleMiddleware ở đây
	// Các route không cần xác thực
	apiV1.POST("/customer/login", userCtr.Login)
	apiV1.POST("/customer/sign", userCtr.SignUp)
	apiV1.POST("/customer/verify-otp", userCtr.VerifyOTP)
	apiV1.GET("/customer", userCtr.List)

	apiV1.Use(middleware.RequestLog())

	// ✅ Các route cần xác thực nằm trong group này
	protected := apiV1.Group("/")
	protected.Use(controllers.RoleMiddleware())
	{

		protected.GET("/customer/my-profile", userCtr.MyProfile)
		protected.GET("/customer/:uuid", userCtr.Detail)
		protected.POST("/customer", userCtr.Create)
		protected.PUT("/customer/:uuid", userCtr.Update)
		protected.PUT("/customer/:uuid/update-status", userCtr.UpdateStatus)
		protected.DELETE("/customer/:uuid", userCtr.Delete)
		protected.POST("/customer/logout", userCtr.Logout)
		protected.POST("/customer/post-idea", userCtr.PostIdea)
		///////////////////////////////////////////////////////////////////////////////
		protected.POST("/customer/rating", userCtr.CreateRating)
		protected.GET("/customer/rating/:expert_uuid", userCtr.ListRating)

	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
	})
}
