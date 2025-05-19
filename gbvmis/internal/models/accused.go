package models

import "time"

type Accused struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Sex               string    `json:"sex"`
	Occupation        string    `json:"occupation"`
	Residence         string    `json:"residence"`
	WorkPlace         string    `json:"work_place"`
	Telephone         string    `json:"telephone"`
	CaseNumber        string    `json:"case_number"`
	PoliceUnit        string    `json:"police_unit"`
	DateOfIncident    time.Time `json:"date_of_incident"`
	ApparentAge       string    `json:"apparent_age"`
	AgeEstimationNote string    `json:"age_estimation_note"`
	HIVTestResults    string    `json:"hiv_test_results"`
	GeneralCondition  string    `json:"general_condition"`
	MentalStatus      string    `json:"mental_status"`
	HeadNeck          string    `json:"head_neck"`
	ChestBreast       string    `json:"chest_breast"`
	AbdomenBack       string    `json:"abdomen_back"`
	UpperLowerLimbs   string    `json:"upper_lower_limbs"`
	AnoGenital        string    `json:"ano_genital"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
