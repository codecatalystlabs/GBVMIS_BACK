package models

import (
	"time"

	"gorm.io/gorm"
)

type Victim struct {
	gorm.Model
	FirstName   string    `gorm:"size:50" json:"first_name"`
	LastName    string    `gorm:"size:50" json:"last_name"`
	Gender      string    `gorm:"size:10" json:"gender"`
	Dob         time.Time `gorm:"type:date" json:"dob"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Nationality string    `json:"nationality"`
	Nin         string    `json:"nin"`
	CreatedBy   string    `gorm:"size:50;not null" json:"created_by"`
	UpdatedBy   string    `gorm:"size:50" json:"updated_by"`

	// Relationships
	Cases []Case `gorm:"many2many:case_victims;" json:"cases"`
}
