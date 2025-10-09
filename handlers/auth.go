package handlers

import (
	"Fitness-Tracker/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func containsUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func containsLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	re := regexp.MustCompile(`\d`)

	if re.MatchString(s) {
		return true
	}
	return false
}

func containsSpecialCharacters(s string) bool {
	specialCharacters := "!@#$%^&*()"

	if strings.ContainsAny(s, specialCharacters) {
		return true
	}
	return false
}

func (ac *AuthController) Register(c *gin.Context) {

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "register", gin.H{})
		return
	}

	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")
	gender := c.PostForm("gender")
	heightStr := c.PostForm("height")
	dobStr := c.PostForm("dob")

	if len(password) < 8 {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "Password must be at least 8 characters long"})
		return
	}
	if !containsNumber(password) {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "Password must contain at least one number"})
		return
	}
	if !containsUpperCase(password) {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "Password must have at least one upper case character"})
		return
	}
	if !containsLowerCase(password) {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "Password must have at least one lower case character"})
		return
	}
	if !containsSpecialCharacters(password) {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "Password must have at least one special character (!@#$%^&*())"})
		return
	}

	var user models.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register", gin.H{"error": "A user with that email already exists"})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register", gin.H{"error": "Failed to hash password"})
		return
	}

	var heightCm *uint
	if heightStr == "" {
		if h, err := strconv.ParseUint(heightStr, 10, 32); err == nil {
			val := uint(h)
			heightCm = &val
		}
	}

	var dob *time.Time
	if dobStr != "" {
		if converted, err := time.Parse("2006-01-01", dobStr); err == nil {
			dob = &converted
		}
	}

	newUser := models.User{
		Email:        email,
		PasswordHash: string(hashedPass),
		DisplayName:  name,
		HeightCm:     heightCm,
		Gender:       gender,
		DOB:          dob,
	}

	if err := ac.DB.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register", gin.H{"error": "Failed to create new user"})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) Login(c *gin.Context) {

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login", gin.H{})
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

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/login")
}
