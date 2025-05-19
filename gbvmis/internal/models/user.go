package models

import "time"

type PoliceOfficer struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Username     string    `json:"username" gorm:"unique"`
	PasswordHash string    `json:"password_hash"`
	Rank         string    `json:"rank"`
	ForceNumber  string    `json:"force_number"`
	Telephone    string    `json:"telephone"`
	Email        string    `json:"email"`
	StationID    uint      `json:"station_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PoliceStation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
