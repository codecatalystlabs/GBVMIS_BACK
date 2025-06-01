package models

import "gorm.io/gorm"

type PolicePost struct {
	gorm.Model
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`

	// Relationships
	Officers []PoliceOfficer `gorm:"foreignKey:PostID" json:"officers"`
}
