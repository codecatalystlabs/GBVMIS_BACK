package models

import "gorm.io/gorm"

// Forward declaration for Role to fix undefined error
// type Role struct defined in role.go

type PoliceOfficer struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rank      string `json:"rank"`
	BadgeNo   string `json:"badge_no"`
	Phone     string `json:"phone"`
	PostID    uint   `json:"post_id"`

	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`

	Roles []*Role `gorm:"many2many:officer_roles;" json:"roles"`
	Cases []Case  `gorm:"foreignKey:OfficerID" json:"cases"`
}
