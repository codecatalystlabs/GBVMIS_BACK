package models

import "gorm.io/gorm"

type HealthFacility struct {
	gorm.Model
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`

	// Relationships
	Practitioners []HealthPractitioner `gorm:"foreignKey:FacilityID" json:"practitioners"`
}
