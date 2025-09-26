package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Symptom struct {
	gorm.Model
	Name   string   `json:"name" gorm:"unique;not null"`
	People []Person `gorm:"many2many:person_symptoms;" json:"people,omitempty"`
}

type PostMortemSummary struct {
	gorm.Model
	Description string   `json:"description"`
	People      []Person `gorm:"many2many:postmortem_people;" json:"people,omitempty"`
}

type Person struct {
	gorm.Model
	Name           string                      `json:"name"`
	Occupation     string                      `json:"occupation"`
	Habits         datatypes.JSONSlice[string] `json:"habits" gorm:"type:json"`
	ApproximateAge int                         `json:"approximate_age"`
	Gender         string                      `json:"gender"`
	Category       string                      `json:"category"` // "sick" or "deceased"
	Symptoms       []Symptom                   `gorm:"many2many:person_symptoms;" json:"symptoms,omitempty"`
	Summaries      []PostMortemSummary         `gorm:"many2many:postmortem_people;" json:"summaries,omitempty"`
}

type PoliceReport struct {
	gorm.Model
	OfficerID uint          `json:"officer_id"`
	Officer   PoliceOfficer `gorm:"foreignKey:OfficerID"`
	Date      time.Time     `gorm:"type:date" json:"date"`
}

type ToxicologyForensicReport struct {
	gorm.Model
	PersonID             uint               `json:"sick_deceased_person_id"`
	Person               Person             `gorm:"foreignKey:PersonID"`
	WitnessID            uint               `json:"witness_id"`
	Witness              Witness            `gorm:"foreignKey:WitnessID"`
	DateOnset            time.Time          `gorm:"type:date" json:"date_onset"`
	DateHourOfDeath      time.Time          `gorm:"type:date" json:"date_hour_of_death"`
	DateHourOfBurial     time.Time          `gorm:"type:date" json:"date_hour_of_burial"`
	DateHourOfExhumation time.Time          `gorm:"type:date" json:"date_hour_of_exhumation"`
	SpecimenSealedBy     string             `json:"specimen_sealed_by"`
	WitnessedBy          string             `json:"witnessed_by"`
	HandedOverTo         string             `json:"handed_over_to"`
	PractitionerID       uint               `json:"practitioner_id"`
	Practitioner         HealthPractitioner `gorm:"foreignKey:PractitionerID"`
	PoliceReportID       uint               `json:"police_report_id"`
	PoliceReport         PoliceReport       `gorm:"foreignKey:PoliceReportID"`
}
