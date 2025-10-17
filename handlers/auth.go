package handlers

import (
	"Fitness-Tracker/models"
	"Fitness-Tracker/utils"
	"fmt"
	"net/http"
	"regexp"
	"sort"
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

	return re.MatchString(s)
}

func containsSpecialCharacters(s string) bool {
	specialCharacters := "!@#$%^&*()"

	return strings.ContainsAny(s, specialCharacters)
}

func (ac *AuthController) Home(c *gin.Context) {
	user := utils.GetCurrentUser(c, ac.DB)

	var workouts []models.Workout

	if err := ac.DB.
		Where("user_id = ?", user.ID).
		Order("created_at desc").
		Find(&workouts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "home", gin.H{
			"user":  user,
			"error": "Error fetching workouts",
		})
		return
	}

	weeklyWorkouts, weekKeys := groupWorkoutsWeekly(workouts)

	c.HTML(http.StatusOK, "home", gin.H{
		"user":           user,
		"workouts":       workouts,
		"weeklyWorkouts": weeklyWorkouts,
		"weekKeys":       weekKeys,
	})
}


func groupWorkoutsWeekly(workouts []models.Workout) (map[string][]models.Workout, []string) {
	weekly := make(map[string][]models.Workout)

	for _, workout := range workouts {
		year, week := workout.OccuredAt.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		weekly[weekKey] = append(weekly[weekKey], workout)
	}

	// Sort workouts within each week by date ascending
	for weekKey, weekWorkouts := range weekly {
		sort.Slice(weekWorkouts, func(i, j int) bool {
			return weekWorkouts[i].OccuredAt.Before(weekWorkouts[j].OccuredAt)
		})
		weekly[weekKey] = weekWorkouts
	}

	// Collect all week keys
	var weekKeys []string
	for k := range weekly {
		weekKeys = append(weekKeys, k)
	}

	// Sort weeks by newest first (descending)
	sort.Slice(weekKeys, func(i, j int) bool {
		// Parse back to year/week numbers for comparison
		var yi, wi, yj, wj int
		fmt.Sscanf(weekKeys[i], "%d-W%d", &yi, &wi)
		fmt.Sscanf(weekKeys[j], "%d-W%d", &yj, &wj)
		if yi == yj {
			return wi > wj // higher week = newer
		}
		return yi > yj // higher year = newer
	})

	return weekly, weekKeys
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
		Role:         "user",
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

	c.Redirect(http.StatusFound, "/home")
}
