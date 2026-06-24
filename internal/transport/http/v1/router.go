package v1

import (
	"minjust-website/internal/transport/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, secretKey string, reqHandler *RequestHandler, empHandler *EmployeeHandler, authHandler *AuthHandler) {
	r.Use(cors.Default())
	v1 := r.Group("/api/v1/")

	{
		v1.POST("/auth/login", authHandler.LoginH)
		v1.POST("/register", empHandler.CreateAccountH)
		v1.POST("/password-reset", reqHandler.RequestResetPasswordH)
		v1.GET("/password-reset/:id", reqHandler.CheckStatusH)
		v1.GET("/handbook", empHandler.GetHandbookH)

		adminRoutes := v1.Group("/admin", middleware.AuthMiddleware(secretKey, "admin"))
		{
			adminRoutes.POST("/password-reset/:id/process", reqHandler.ProcessRequestH)
			adminRoutes.GET("/employees", empHandler.GetAllAccountsH)
		}
	}
}
