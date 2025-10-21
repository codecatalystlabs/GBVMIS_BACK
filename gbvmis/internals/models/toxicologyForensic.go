package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Symptom struct {
	gorm.Model
	Name           string          `json:"name" gorm:"unique;not null"`
	PersonSymptoms []PersonSymptom `json:"person_symptoms,omitempty"`
}

type PersonSymptom struct {
	gorm.Model
	PersonID  uint
	SymptomID uint
	State     string `json:"state"` // e.g. "mild", "severe", "resolved"

	Person  Person
	Symptom Symptom
}

type PostMortemSummary struct {
	gorm.Model
	Description     string          `json:"description"`
	PeopleSummaries []PersonSummary `json:"people_summaries,omitempty"`
}

type PersonSummary struct {
	gorm.Model
	PersonID            uint
	PostMortemSummaryID uint
	State               string `json:"state"` // e.g. "mild", "severe", "resolved"

	Person            Person
	PostMortemSummary PostMortemSummary
}

type Person struct {
	gorm.Model
	Name                 string                      `json:"name"`
	Occupation           string                      `json:"occupation"`
	Habits               datatypes.JSONSlice[string] `json:"habits" gorm:"type:json"`
	ApproximateAge       int                         `json:"approximate_age"`
	Gender               string                      `json:"gender"`
	Category             string                      `json:"category"` // "sick" or "deceased"
	PersonSymptoms       []PersonSymptom             `json:"person_symptoms,omitempty"`
	Summaries            []PersonSummary             `json:"summaries,omitempty"`
	DateHourOfPostMortem time.Time                   `gorm:"type:date" json:"date_hour_of_post_mortem"`
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
