package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex; not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string
	DOB          *time.Time
	HeightCm     *uint
	Gender       string
	Role         string `gorm:"default:user"`
}

type Workout struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null"`
	Type      string `gorm:"not null"`
	Duration  int
	DistanceM int
	Notes     string
	OccuredAt time.Time `gorm:"index"`
	Done      bool      `gorm:"default:false"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type WeeklyMetric struct {
	gorm.Model
	UserID   uint      `gorm:"index; not null"`
	Date     time.Time `gorm:"index"`
	WeightKg float32
	User     User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
