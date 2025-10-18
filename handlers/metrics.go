package handlers

import (
	"Fitness-Tracker/models"
	"Fitness-Tracker/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MetricsController struct {
	DB *gorm.DB
}

func NewMetricsController(db *gorm.DB) *MetricsController {
	return &MetricsController{db}
}

func (m *MetricsController) DisplayMetrics(ctx *gin.Context) {
	user := utils.GetCurrentUser(ctx, m.DB)

	var metrics []models.Metrics
	if err := m.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&metrics).Error; err != nil {
		ctx.HTML(http.StatusInternalServerError, "display-metrics", gin.H{"user": user, "error": "Failed to fetch metrics"})
		return
	}

	ctx.HTML(http.StatusOK, "display-metrics", gin.H{"user": user, "metrics": metrics})
}

func (m *MetricsController) CreateMetrics(ctx *gin.Context) {
	user := utils.GetCurrentUser(ctx, m.DB)

	if ctx.Request.Method == "GET" {
		ctx.HTML(http.StatusOK, "create-metrics", gin.H{"user":user})
		return
	}

	weightStr := ctx.PostForm("weight")
	heightStr := ctx.PostForm("height")

	var weight, height float32

	if weightStr != "" {
		val, err := strconv.ParseFloat(weightStr, 32)
		if err != nil {
			return
		}
		weight = float32(val)
	}

	if heightStr != "" {
		val, err := strconv.ParseFloat(heightStr, 32)
		if err != nil {
			return
		}
		height = float32(val)
	}

	newRecord := models.Metrics {
		UserID: user.ID,
		Date: time.Now(),
		WeightKg: weight,
		HeightCm: height,
	}

	if err := m.DB.Create(&newRecord).Error; err != nil {
		ctx.HTML(http.StatusInternalServerError, "create-metrics", gin.H{"error": "Failed to create a new record"})
		return
	}

	ctx.Redirect(http.StatusFound, "/display-metrics")
}
