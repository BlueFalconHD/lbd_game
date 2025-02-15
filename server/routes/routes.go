package routes

import (
	"time"

	"github.com/bluefalconhd/lbd_game/server/config"
	"github.com/bluefalconhd/lbd_game/server/controllers"
	"github.com/bluefalconhd/lbd_game/server/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()

	// CORS
	config := cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	// Public routes
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)
	router.GET("/phrase", controllers.GetCurrentPhrase)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/privelege", controllers.Privelege)
		protected.POST("/phrase", controllers.SubmitPhrase)
		protected.GET("/can_submit_phrase", controllers.CanSubmitPhrase)
		protected.POST("/verify", controllers.VerifyUser)
		protected.GET("/verifications", controllers.GetCurrentVerifications)
		protected.GET("/unverified_users", controllers.GetUnverifiedUsers)
	}

	// Admin routes
	admin := protected.Group("/admin")
	//FIXME: add back admin check
	// admin.Use(middleware.PrivilegeMiddleware(1)) // Requires at least Admin Level 1
	protected.Use(middleware.AuthMiddleware())
	{
		admin.GET("/stats/users", controllers.GetUserStatistics)
		// admin.PUT("/user/:id/resurrect", controllers.ResurrectUser)
		admin.PUT("/edit_phrase", controllers.EditPhrase)
		admin.PUT("/unsubmit_phrase", controllers.UnsubmitPhrase)
		admin.GET("/scheduled_windows", controllers.GetScheduledWindows)
		admin.DELETE("/scheduled_windows/:id", controllers.CancelScheduledWindow)
		admin.PUT("/manual_reset", controllers.ManualReset)
	}

	// Super Admin routes
	superAdmin := protected.Group("/superadmin")
	//FIXME: add back superadmin check
	// superAdmin.Use(middleware.PrivilegeMiddleware(2)) // Requires Admin Level 2
	protected.Use(middleware.AuthMiddleware())
	{
		superAdmin.PUT("/user/:id/promote", controllers.PromoteUser)
		superAdmin.PUT("/user/:id/demote", controllers.DemoteUser)
	}

	return router
}
