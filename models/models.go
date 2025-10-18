package models

import (
	"math"
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

type Metrics struct {
	gorm.Model
	UserID   uint      `gorm:"index; not null"`
	Date     time.Time `gorm:"index"`
	WeightKg float32   `gorm:"not null"`
	HeightCm float32   `gorm:"not null"`
	BMI      float32
	Status   string
	User     User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (m *Metrics) BeforeSave(tx *gorm.DB) (err error) {
	if m.HeightCm <= 0 {
		m.BMI = 0
		m.Status = "Invalid height"
		return nil
	}

	divisor := math.Pow(float64(m.HeightCm)/100, 2)
	m.BMI = float32(float64(m.WeightKg) / divisor)

	switch {
	case m.BMI < 18.50:
		m.Status = "Underweight"
	case m.BMI >= 18.50 && m.BMI < 25.00:
		m.Status = "Normal"
	case m.BMI >= 25.00 && m.BMI < 30.00:
		m.Status = "Overweight"
	default:
		m.Status = "Obese"
	}

	return nil
}
