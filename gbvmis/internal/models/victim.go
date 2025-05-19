package models

import "time"

type Victim struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Sex               string    `json:"sex"`
	Occupation        string    `json:"occupation"`
	MaritalStatus     string    `json:"marital_status"`
	Residence         string    `json:"residence"`
	CaseNumber        string    `json:"case_number"`
	PoliceUnit        string    `json:"police_unit"`
	DateOfIncident    time.Time `json:"date_of_incident"`
	Narrator          string    `json:"narrator"`
	Relationship      string    `json:"relationship"`
	ApparentAge       string    `json:"apparent_age"`
	AgeEstimationNote string    `json:"age_estimation_note"`
	IncidentHistory   string    `json:"incident_history"`
	GeneralCondition  string    `json:"general_condition"`
	MentalStatus      string    `json:"mental_status"`
	HeadNeck          string    `json:"head_neck"`
	ChestBreast       string    `json:"chest_breast"`
	AbdomenBack       string    `json:"abdomen_back"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
