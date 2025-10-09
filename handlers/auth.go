package handlers

import (
	"Fitness-Tracker/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Login(c *gin.Context) {

	if c.Request.Method == "GET" {
		c.HTML(http.StatusFound, "login", gin.H{})
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "login", gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		c.HTML(http.StatusNotFound, "login", gin.H{"error": "Invalid email or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_email", email)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusFound, "/home")
}
