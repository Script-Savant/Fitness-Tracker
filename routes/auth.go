package routes

import (
	"Fitness-Tracker/config"
	"Fitness-Tracker/handlers"
	"Fitness-Tracker/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine) {
	db := config.DB

	authController := handlers.NewAuthController(db)

	

	// Auth routes

	r.GET("/register", authController.Register)
	r.POST("/register", authController.Register)
	r.GET("/login", authController.Login)
	r.POST("/login", authController.Login)

	// private routes
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// home
		auth.GET("/", authController.Home)
		auth.GET("/home", authController.Home)

		auth.GET("/logout", handlers.Logout)
	}
}
