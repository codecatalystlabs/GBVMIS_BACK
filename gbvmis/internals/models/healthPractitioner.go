package models

import "gorm.io/gorm"

type HealthPractitioner struct {
	gorm.Model
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	Profession string `json:"profession"`  // e.g., Doctor, Nurse, Counselor
	FacilityID uint   `json:"facility_id"` // Foreign key to HealthFacility

	// Relationships
	Examinations []Examination `gorm:"foreignKey:PractitionerID" json:"examinations"`
}

type Examination struct {
	gorm.Model
	VictimID       uint   `json:"victim_id"`
	CaseID         uint   `json:"case_id"`
	FacilityID     uint   `json:"facility_id"`
	PractitionerID uint   `json:"practitioner_id"`
	ExamDate       string `gorm:"type:date" json:"exam_date"`
	Findings       string `gorm:"type:text" json:"findings"`
	Treatment      string `gorm:"type:text" json:"treatment"`
	Referral       string `json:"referral"` // Optional referral info
	ConsentGiven   bool   `json:"consent_given"`

	Victim       Victim             `gorm:"foreignKey:VictimID"`
	Case         Case               `gorm:"foreignKey:CaseID"`
	Facility     HealthFacility     `gorm:"foreignKey:FacilityID"`
	Practitioner HealthPractitioner `gorm:"foreignKey:PractitionerID"`
}
