package routes

import (
	"Fitness-Tracker/config"
	"Fitness-Tracker/handlers"
	"Fitness-Tracker/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWorkoutRoutes(router *gin.Engine){
	db := config.GetDB()
	workouts := handlers.NewWorkoutController(db)

	r := router.Group("/")
	r.Use(middleware.AuthRequired())
	{
		r.GET("/create-workout", workouts.CreateWorkout)
		r.POST("/create-workout", workouts.CreateWorkout)
		r.GET("/update-workout/:id", workouts.UpdateWorkout)
		r.POST("/update-workout/:id", workouts.UpdateWorkout)
		r.POST("/delete-workout/:id", workouts.DeleteWorkout)
		r.POST("/workout/:id/done", workouts.MarkWorkoutDone)
	}
}