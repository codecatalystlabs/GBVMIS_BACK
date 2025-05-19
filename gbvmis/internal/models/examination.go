package models

import "time"

type MedicalExamination struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	VictimID         *uint     `json:"victim_id"`
	AccusedID        *uint     `json:"accused_id"`
	Place            string    `json:"place"`
	Date             time.Time `json:"date"`
	Findings         string    `json:"findings"`
	GeneralCondition string    `json:"general_condition"`
	MentalStatus     string    `json:"mental_status"`
	HeadNeck         string    `json:"head_neck"`
	ChestBreast      string    `json:"chest_breast"`
	AbdomenBack      string    `json:"abdomen_back"`
	UpperLowerLimbs  string    `json:"upper_lower_limbs"`
	AnoGenital       string    `json:"ano_genital"`
	PictogramFront   string    `json:"pictogram_front"`
	PictogramBack    string    `json:"pictogram_back"`
	PictogramMale    string    `json:"pictogram_male"`
	PictogramFemale  string    `json:"pictogram_female"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
