package main

import (
	"Fitness-Tracker/models"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  admin create <email> <password>")
		fmt.Println("  admin update <email> <password>")
		fmt.Println("  admin delete <email>")
		os.Exit(1)
	}

	command := os.Args[1]
	email := os.Args[2]

	db, err := gorm.Open(sqlite.Open("fitness_tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate user model: ", err)
	}

	switch command {
	case "create":
		if len(os.Args) != 4 {
			fmt.Println("Usage: admin create <email> <password>")
			os.Exit(1)
		}
		password := os.Args[3]

		var existingUser models.User
		if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
			fmt.Printf("User with email: '%s' already exists.", email)
			os.Exit(1)
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("failed to hash password: ", err)
		}

		adminUser := models.User{
			Email:        email,
			PasswordHash: string(hashedPass),
			Role:         "admin",
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatal("Failed to create admin user:", err)
		}

		fmt.Printf("Admin user '%s' created successfully!\n", email)

	case "update":
		if len(os.Args) != 4 {
			fmt.Println("Usage: admin update <email> <password>")
			os.Exit(1)
		}
		password := os.Args[3]

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			fmt.Printf("User with email '%s' not found.\n", email)
			os.Exit(1)
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		user.PasswordHash = string(hashedPass)
		user.Role = "admin"
		if err := db.Save(&user).Error; err != nil {
			log.Fatal("Failed to update user: ", err)
		}

		fmt.Printf("User '%s' updated with new password and promoted to admin.\n", email)

	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: admin delete <email>")
			os.Exit(1)
		}

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			fmt.Printf("User with email '%s' not found.\n", email)
			os.Exit(1)
		}

		if err := db.Delete(&user).Error; err != nil {
			log.Fatal("Failed to delete user:", err)
		}

		fmt.Printf("User '%s' deleted successfully.\n", email)

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Usage:")
		fmt.Println("  admin create <email> <password>")
		fmt.Println("  admin update <email> <password>")
		fmt.Println("  admin delete <email>")
		os.Exit(1)
	}
}
