package dto

import "time"

type PersonSymptomDTO struct {
	SymptomID uint   `json:"symptom_id"`
	State     string `json:"state"`
}

type PersonSummaryDTO struct {
	PostMortemSummaryID uint   `json:"post_mortem_summary_id"`
	State               string `json:"state"`
}

type PersonDTO struct {
	Name                 string             `json:"name"`
	Occupation           string             `json:"occupation"`
	Habits               []string           `json:"habits"`
	ApproximateAge       int                `json:"approximate_age"`
	Gender               string             `json:"gender"`
	Category             string             `json:"category"`
	DateHourOfPostMortem time.Time          `json:"date_hour_of_post_mortem"`
	PersonSymptoms       []PersonSymptomDTO `json:"person_symptoms"`
	PersonSummaries      []PersonSummaryDTO `json:"person_summaries"`
}

type PoliceReportDTO struct {
	OfficerID             uint   `json:"officer_id"`
	IsPersonPoisoned      bool   `json:"is_person_poisoned"`
	IsSuicideOrAccident   bool   `json:"is_suicide_or_accident"`
	IsDeceasedOnTreatment bool   `json:"is_deceased_on_treatment"`
	TreatmentDetails      string `json:"treatment_details"`
}

type ToxicologyReportCreateDTO struct {
	Person               PersonDTO       `json:"person"`
	WitnessID            uint            `json:"witness_id"`
	PractitionerID       uint            `json:"practitioner_id"`
	PoliceReport         PoliceReportDTO `json:"police_report"`
	DateOnset            time.Time       `json:"date_onset"`
	DateHourOfDeath      time.Time       `json:"date_hour_of_death"`
	DateHourOfBurial     time.Time       `json:"date_hour_of_burial"`
	DateHourOfExhumation time.Time       `json:"date_hour_of_exhumation"`
	SpecimenSealedBy     string          `json:"specimen_sealed_by"`
	WitnessedBy          string          `json:"witnessed_by"`
	HandedOverTo         string          `json:"handed_over_to"`
}
