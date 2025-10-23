package controllers

import (
	"gbvmis/internals/dto"
	"gbvmis/internals/models"
	"gbvmis/internals/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ToxicologyHandler struct {
	service service.ToxicologyService
}

func NewToxicologyForensicController(service service.ToxicologyService) *ToxicologyHandler {
	return &ToxicologyHandler{service: service}
}

// POST /reports
func (h *ToxicologyHandler) Create(c *fiber.Ctx) error {
	var payload dto.ToxicologyReportCreateDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// 1️⃣ Map PersonDTO -> Person model
	if payload.Person.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Person name is required"})
	}
	if payload.PoliceReport.OfficerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "OfficerID is required in PoliceReport"})
	}
	if payload.PractitionerID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PractitionerID is required"})
	}

	person := models.Person{
		Name:                 payload.Person.Name,
		Occupation:           payload.Person.Occupation,
		Habits:               payload.Person.Habits,
		ApproximateAge:       payload.Person.ApproximateAge,
		Gender:               payload.Person.Gender,
		Category:             payload.Person.Category,
		DateHourOfPostMortem: payload.Person.DateHourOfPostMortem,
	}

	// 2️⃣ Attach Symptoms
	for _, s := range payload.Person.PersonSymptoms {
		if s.SymptomID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Each PersonSymptom must include a valid SymptomID",
			})
		}
		person.PersonSymptoms = append(person.PersonSymptoms, models.PersonSymptom{
			State:     s.State,
			SymptomID: s.SymptomID,
		})
	}

	// 3️⃣ Attach Summaries
	for _, ps := range payload.Person.PersonSummaries {
		if ps.PostMortemSummaryID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Each PersonSummary must include a valid PostMortemSummaryID",
			})
		}
		person.Summaries = append(person.Summaries, models.PersonSummary{
			State:               ps.State,
			PostMortemSummaryID: ps.PostMortemSummaryID,
		})
	}

	// 4️⃣ Create PoliceReport model
	policeReport := models.PoliceReport{
		OfficerID:             payload.PoliceReport.OfficerID,
		IsPersonPoisoned:      payload.PoliceReport.IsPersonPoisoned,
		IsSuicideOrAccident:   payload.PoliceReport.IsSuicideOrAccident,
		IsDeceasedOnTreatment: payload.PoliceReport.IsDeceasedOnTreatment,
		TreatmentDetails:      payload.PoliceReport.TreatmentDetails,
	}

	// 5️⃣ Create ToxicologyForensicReport model
	report := models.ToxicologyForensicReport{
		Person:               person,
		WitnessID:            payload.WitnessID,
		PractitionerID:       payload.PractitionerID,
		PoliceReport:         policeReport,
		DateOnset:            payload.DateOnset,
		DateHourOfDeath:      payload.DateHourOfDeath,
		DateHourOfBurial:     payload.DateHourOfBurial,
		DateHourOfExhumation: payload.DateHourOfExhumation,
		SpecimenSealedBy:     payload.SpecimenSealedBy,
		WitnessedBy:          payload.WitnessedBy,
		HandedOverTo:         payload.HandedOverTo,
	}

	// 6️⃣ Persist via service
	if err := h.service.Create(&report); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(report)
}

// GET /reports
func (h *ToxicologyHandler) GetAll(c *fiber.Ctx) error {
	reports, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reports)
}

func (h *ToxicologyHandler) GetPaginatedReports(c *fiber.Ctx) error {
	pagination, reports, err := h.service.GetPaginatedReports(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve reports",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Suspects retrieved successfully",
		"data":    reports,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// GET /reports/:id
func (h *ToxicologyHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	report, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(report)
}

// PUT /reports/:id
func (h *ToxicologyHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid report ID"})
	}

	var payload dto.ToxicologyReportCreateDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload: " + err.Error(),
		})
	}

	// ✅ Fetch existing report
	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Toxicology report not found"})
	}

	// ✅ Update report fields
	existing.WitnessID = payload.WitnessID
	existing.PractitionerID = payload.PractitionerID
	existing.DateOnset = payload.DateOnset
	existing.DateHourOfDeath = payload.DateHourOfDeath
	existing.DateHourOfBurial = payload.DateHourOfBurial
	existing.DateHourOfExhumation = payload.DateHourOfExhumation
	existing.SpecimenSealedBy = payload.SpecimenSealedBy
	existing.WitnessedBy = payload.WitnessedBy
	existing.HandedOverTo = payload.HandedOverTo

	// ✅ Update or replace associated PoliceReport
	if existing.PoliceReport.ID != 0 {
		existing.PoliceReport.OfficerID = payload.PoliceReport.OfficerID
		existing.PoliceReport.IsPersonPoisoned = payload.PoliceReport.IsPersonPoisoned
		existing.PoliceReport.IsSuicideOrAccident = payload.PoliceReport.IsSuicideOrAccident
		existing.PoliceReport.IsDeceasedOnTreatment = payload.PoliceReport.IsDeceasedOnTreatment
		existing.PoliceReport.TreatmentDetails = payload.PoliceReport.TreatmentDetails
	} else {
		existing.PoliceReport = models.PoliceReport{
			OfficerID:             payload.PoliceReport.OfficerID,
			IsPersonPoisoned:      payload.PoliceReport.IsPersonPoisoned,
			IsSuicideOrAccident:   payload.PoliceReport.IsSuicideOrAccident,
			IsDeceasedOnTreatment: payload.PoliceReport.IsDeceasedOnTreatment,
			TreatmentDetails:      payload.PoliceReport.TreatmentDetails,
		}
	}

	// ✅ Update Person details
	if existing.Person.ID != 0 {
		existing.Person.Name = payload.Person.Name
		existing.Person.Occupation = payload.Person.Occupation
		existing.Person.Habits = payload.Person.Habits
		existing.Person.ApproximateAge = payload.Person.ApproximateAge
		existing.Person.Gender = payload.Person.Gender
		existing.Person.Category = payload.Person.Category
		existing.Person.DateHourOfPostMortem = payload.Person.DateHourOfPostMortem
	}

	// ✅ Replace PersonSymptoms
	existing.Person.PersonSymptoms = nil
	for _, s := range payload.Person.PersonSymptoms {
		existing.Person.PersonSymptoms = append(existing.Person.PersonSymptoms, models.PersonSymptom{
			State:     s.State,
			SymptomID: s.SymptomID,
			PersonID:  existing.Person.ID,
		})
	}

	// ✅ Replace PersonSummaries
	existing.Person.Summaries = nil
	for _, ps := range payload.Person.PersonSummaries {
		existing.Person.Summaries = append(existing.Person.Summaries, models.PersonSummary{
			State:               ps.State,
			PostMortemSummaryID: ps.PostMortemSummaryID,
			PersonID:            existing.Person.ID,
		})
	}

	// ✅ Save updates
	if err := h.service.Update(existing); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update report: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Toxicology report updated successfully",
		"data":    existing,
	})
}

// DELETE /reports/:id
func (h *ToxicologyHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
