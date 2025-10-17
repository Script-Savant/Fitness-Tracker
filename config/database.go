package config

import (
	"Fitness-Tracker/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("fitness_tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to make a database connection")
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Workout{},
		&models.WeeklyMetric{},
	); err != nil {
		log.Fatal("Failed to migrate the database")
	}

	fmt.Println("Successfully connected to the database")

	DB = db
}

func GetDB() * gorm.DB {
	return DB
}