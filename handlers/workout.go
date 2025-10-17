package handlers

import (
	"Fitness-Tracker/models"
	"Fitness-Tracker/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkoutController struct {
	DB *gorm.DB
}

func NewWorkoutController(db *gorm.DB) *WorkoutController {
	return &WorkoutController{db}
}

func (wc *WorkoutController) getWorkout(c *gin.Context) *models.Workout {
	user := utils.GetCurrentUser(c, wc.DB)
	workoutID := c.Param("id")

	var workout models.Workout
	if err := wc.DB.First(&workout, workoutID).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Error finding workout"})
		return nil
	}

	if workout.UserID != user.ID {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "No rights to access that workout"})
		return nil
	}

	return &workout
}

func (wc *WorkoutController) CreateWorkout(c *gin.Context) {
	user := utils.GetCurrentUser(c, wc.DB)

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "create-workout", gin.H{"user": user})
		fmt.Println(user)
		return
	}

	customErr := ""
	data := gin.H{
		"user":  user,
		"error": customErr,
	}

	workoutType := c.PostForm("Type")
	durationStr := c.PostForm("Duration")
	distanceStr := c.PostForm("Distance")
	occuredAt := c.PostForm("OccuredAt")
	notes := c.PostForm("Notes")

	duration, distance, occuredAtTime, err := validate(durationStr, distanceStr, occuredAt)
	if err != nil {
		data["error"] = err.Error()
		c.HTML(http.StatusBadRequest, "create-workout", data)
		return
	}

	newWorkout := models.Workout{
		UserID:    user.ID,
		Type:      workoutType,
		Duration:  duration,
		DistanceM: distance,
		Notes:     notes,
		OccuredAt: occuredAtTime,
	}

	// create workout
	if err := wc.DB.Create(&newWorkout).Error; err != nil {
		customErr = "Error creating  new workout"
		c.HTML(http.StatusInternalServerError, "create-workout", data)
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

func (wc *WorkoutController) UpdateWorkout(c *gin.Context) {
	user := utils.GetCurrentUser(c, wc.DB)
	workout := wc.getWorkout(c)

	fmt.Println(workout)

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "update-workout", gin.H{"user": user, "workout": workout})
		return
	}

	customErr := ""
	data := gin.H{
		"user":  user,
		"error": customErr,
	}

	workoutType := c.PostForm("Type")
	durationStr := c.PostForm("Duration")
	distanceStr := c.PostForm("Distance")
	occuredAt := c.PostForm("OccuredAt")
	notes := c.PostForm("Notes")

	duration, distance, occuredAtTime, err := validate(durationStr, distanceStr, occuredAt)
	if err != nil {
		data["error"] = err.Error()
		c.HTML(http.StatusBadRequest, "create-workout", data)
		return
	}

	workout.Type = workoutType
	workout.Duration = duration
	workout.DistanceM = distance
	workout.Notes = notes
	workout.OccuredAt = occuredAtTime

	if err := wc.DB.Save(&workout).Error; err != nil {
		customErr = "Failed to update workout"
		c.HTML(http.StatusInternalServerError, "update-workout", data)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/update-workout/%d", workout.ID))
}

func (wc *WorkoutController) DeleteWorkout(c *gin.Context) {
	workout := wc.getWorkout(c)

	if err := wc.DB.Delete(&workout).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": "Error deleting workout"})
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

func (wc *WorkoutController) MarkWorkoutDone(c *gin.Context) {
	workout := wc.getWorkout(c)

	workout.Done = !workout.Done
	if err := wc.DB.Save(&workout).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "home", gin.H{"error": "Error marking workout as done"})
		return
	}

	c.Redirect(http.StatusFound, "/home")
}

func validate(durationStr, distanceStr, occuredAt string) (int, int, time.Time, error) {
	var duration, distance int

	if durationStr != "" {
		val, err := strconv.Atoi(durationStr)
		if err != nil {
			return 0, 0, time.Time{}, fmt.Errorf("invalid duration value")
		}
		duration = val
	}

	if distanceStr != "" {
		val, err := strconv.Atoi(distanceStr)
		if err != nil {
			return 0, 0, time.Time{}, fmt.Errorf("invalid distance value")
		}
		distance = val
	}

	layout := "2006-01-02T15:04"
	parsedTime, err := time.Parse(layout, occuredAt)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("invalid date/time format value")
	}

	return duration, distance, parsedTime, nil
}
