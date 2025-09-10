package models

import (
	"time"

	"gorm.io/gorm"
)

type Suspect struct {
	gorm.Model
	FirstName    string    `gorm:"size:50" json:"first_name"`
	MiddleName   string    `gorm:"size:50" json:"middle_name"`
	LastName     string    `gorm:"size:50" json:"last_name"`
	Dob          time.Time `gorm:"type:date" json:"dob"`
	Gender       string    `gorm:"size:50" json:"gender"`
	PhoneNumber  string    `json:"phone_number"`
	Nin          string    `json:"nin"`
	Nationality  string    `json:"nationality"`
	Address      string    `json:"address"`
	Occupation   string    `json:"occupation"`
	Status       string    `gorm:"size:50" json:"status"`
	Fingerprints []byte    `gorm:"type:bytea" json:"fingerprints"`
	Photo        []byte    `gorm:"type:bytea" json:"photo"`
	CreatedBy    string    `gorm:"size:50;not null" json:"created_by"`
	UpdatedBy    string    `gorm:"size:50" json:"updated_by"`

	// Relationships
	Cases   []Case   `gorm:"many2many:case_suspects;" json:"cases"`
	Arrests []Arrest `gorm:"foreignKey:SuspectID" json:"arrests"`
}

// <img src="data:image/jpeg;base64,{{ .photo }}" alt="Suspect photo">
// <img src="data:image/png;base64,{{ .fingerprints }}" alt="Fingerprint">

type Case struct {
	gorm.Model
	CaseNumber   string        `gorm:"uniqueIndex" json:"case_number"`
	Title        string        `json:"title"`
	Description  string        `gorm:"type:text" json:"description"`
	Status       string        `json:"status"`
	DateOpened   time.Time     `gorm:"type:date" json:"date_opened"`
	OfficerID    uint          `json:"officer_id"`
	PolicePostID uint          `json:"police_post_id"`
	Suspects     []Suspect     `gorm:"many2many:case_suspects;" json:"suspects"`
	Charges      []Charge      `gorm:"many2many:case_charges;" json:"charges"`
	Victims      []Victim      `gorm:"many2many:case_victims;" json:"victims"`
	Witnesses    []Witness     `gorm:"many2many:case_witnesses;" json:"witnesses"`
	Officer      PoliceOfficer `gorm:"foreignKey:OfficerID"`
	PolicePost   PolicePost    `gorm:"foreignKey:PolicePostID"`
}

type Arrest struct {
	gorm.Model
	ArrestDate  time.Time `gorm:"type:date" json:"arrest_date"`
	Location    string    `json:"location"`
	OfficerName string    `json:"officer_name"`
	SuspectID   uint      `json:"suspect_id"`
	Notes       string    `gorm:"type:text" json:"notes"`
}

type Charge struct {
	gorm.Model
	ChargeTitle string `json:"charge_title"`
	Description string `gorm:"type:text" json:"description"`
	Severity    string `json:"severity"` // e.g., Felony, Misdemeanor
	Cases       []Case `gorm:"many2many:case_charges;" json:"cases"`
}
