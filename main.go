package main

import (
	"Fitness-Tracker/config"
	"Fitness-Tracker/routes"
	"Fitness-Tracker/utils"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// database connection
	config.ConnectDatabase()

	// router
	router := gin.Default()

	// sessions
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// templates
	router.HTMLRender = utils.SetupTemplates()

	// static files
	router.Static("/static", "./static")

	// routes
	routes.SetupAuthRoutes(router)
	routes.SetupWorkoutRoutes(router)

	// start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
