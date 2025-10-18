package routes

import (
	"Fitness-Tracker/config"
	"Fitness-Tracker/handlers"
	"Fitness-Tracker/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMetricRoutes(router *gin.Engine) {
	db := config.GetDB()
	metrics := handlers.NewMetricsController(db)

	r := router.Group("/")
	r.Use(middleware.AuthRequired())
	{
		r.GET("/display-metrics", metrics.DisplayMetrics)
		r.GET("/create-metrics", metrics.CreateMetrics)
		r.POST("/create-metrics", metrics.CreateMetrics)
	}
}