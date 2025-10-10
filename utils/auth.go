package utils

import (
	"Fitness-Tracker/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCurrentUser(c *gin.Context, db *gorm.DB) *models.User{
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		return nil
	}

	var user models.User
	if err := db.Find(&user, userID).Error; err != nil {
		return nil
	}

	return &user
}