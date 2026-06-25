package v1

import (
	"minjust-website/internal/transport/http/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, secretKey string, reqHandler *RequestHandler, empHandler *EmployeeHandler, authHandler *AuthHandler) {
	//r.Use(cors.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authHandler.LoginH)
		v1.POST("/register", empHandler.CreateAccountH)
		v1.GET("/handbook", empHandler.GetHandbookH)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(secretKey))
		{
			protected.GET("/employee/profile", empHandler.GetProfileH)

			protected.POST("/password-requests", reqHandler.CreateRequestH)
			protected.GET("/password-requests/my", reqHandler.GetMyRequestsH)

			adminRoutes := protected.Group("/admin")
			adminRoutes.Use(middleware.RoleBlockMiddleware("admin"))
			{
				adminRoutes.GET("/password-requests/all", reqHandler.GetAllRequestsH)
				adminRoutes.PUT("/password-requests/:id/review", reqHandler.ReviewRequestH)
			}
		}
	}
}
