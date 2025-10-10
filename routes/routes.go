package routes

import (
	"Fitness-Tracker/config"
	"Fitness-Tracker/handlers"
	"Fitness-Tracker/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	db := config.DB

	// home
	r.GET("/", handlers.Home)
	r.GET("/home", handlers.Home)

	// Auth routes
	authController := handlers.NewAuthController(db)

	r.GET("/register", authController.Register)
	r.POST("/register", authController.Register)
	r.GET("/login", authController.Login)
	r.POST("/login", authController.Login)

	// private routes
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/logout", handlers.Logout)
	}
}
