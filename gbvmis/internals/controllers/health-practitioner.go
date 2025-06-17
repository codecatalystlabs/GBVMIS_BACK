package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthPractitionerController struct {
	repo repository.HealthPractitionerRepository
}

func NewHealthPractitionerController(repo repository.HealthPractitionerRepository) *HealthPractitionerController {
	return &HealthPractitionerController{repo: repo}
}

type CreateHealthPractitionerPayload struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	Profession string `json:"profession" validate:"required"`
	FacilityID uint   `json:"facility_id" validate:"required"`
}

type ExaminationResponse struct {
	ID           uint   `json:"id"`
	VictimID     uint   `json:"victim_id"`
	CaseID       uint   `json:"case_id"`
	FacilityID   uint   `json:"facility_id"`
	ExamDate     string `json:"exam_date"`
	Findings     string `json:"findings"`
	Treatment    string `json:"treatment"`
	Referral     string `json:"referral"`
	ConsentGiven bool   `json:"consent_given"`
}

type HealthPractitionerExamResponse struct {
	ID           uint                  `json:"id"`
	FirstName    string                `json:"first_name"`
	LastName     string                `json:"last_name"`
	Gender       string                `json:"gender"`
	Phone        string                `json:"phone"`
	Profession   string                `json:"profession"`
	FacilityID   uint                  `json:"facility_id"`
	CreatedAt    time.Time             `json:"created_at"`
	Examinations []ExaminationResponse `json:"examinations"`
}

func ConvertToHealthPractitionerResponse(hp models.HealthPractitioner) HealthPractitionerExamResponse {
	examinations := make([]ExaminationResponse, len(hp.Examinations))
	for i, examination := range hp.Examinations {
		examinations[i] = ExaminationResponse{
			ID:           examination.ID,
			VictimID:     examination.VictimID,
			CaseID:       examination.CaseID,
			FacilityID:   examination.FacilityID,
			ExamDate:     examination.ExamDate,
			Findings:     examination.Findings,
			Treatment:    examination.Treatment,
			Referral:     examination.Findings,
			ConsentGiven: examination.ConsentGiven,
		}
	}

	return HealthPractitionerExamResponse{
		ID:           hp.ID,
		FirstName:    hp.FirstName,
		LastName:     hp.LastName,
		Gender:       hp.Gender,
		Phone:        hp.Phone,
		Profession:   hp.Profession,
		FacilityID:   hp.FacilityID,
		CreatedAt:    hp.CreatedAt,
		Examinations: examinations,
	}
}

// ================================

// CreateHealthPractitioner godoc
//
//	@Summary		Create a new health practitioner record
//	@Description	Creates a new health practitioner entry in the system and returns the created record.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Param			Health	practitioner	body		CreateHealthPractitionerPayload	true	"Health practitioner data to create"
//	@Success		201		{object}		fiber.Map	"Successfully created health practitioner record"
//	@Failure		400		{object}		fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}		fiber.Map	"Server error when creating health practitioner"
//	@Router			/health-practitioner [post]
func (h *HealthPractitionerController) CreateHealthPractitioner(c *fiber.Ctx) error {
	var payload CreateHealthPractitionerPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	hp := &models.HealthPractitioner{
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Gender:     payload.Gender,
		Phone:      payload.Phone,
		Profession: payload.Profession,
		FacilityID: payload.FacilityID,
	}

	if err := h.repo.CreateHealthPractitioner(hp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create health practitioner",
			"data":    err.Error(),
		})
	}

	response := ConvertToHealthPractitionerResponse(*hp)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Health practitioner created successfully",
		"data":    response,
	})
}

// ===========

// GetAllHealthPractitioners godoc
//
//	@Summary		Retrieve a paginated list of health practitioners
//	@Description	Fetches all healthPractitioner records with pagination support.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Health practitioners retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve health practitioners"
//	@Router			/health-practitioners [get]
func (h *HealthPractitionerController) GetAllHealthPractitioners(c *fiber.Ctx) error {
	pagination, healthPractitioners, err := h.repo.GetPaginatedHealthPractitioners(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve health practitioners",
			"data":    err.Error(),
		})
	}

	hpResponses := make([]HealthPractitionerExamResponse, len(healthPractitioners))
	for i, hp := range healthPractitioners {
		hpResponses[i] = ConvertToHealthPractitionerResponse(hp)
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Health practitioners retrieved successfully",
		"data":    hpResponses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleHealthPractitioner godoc
//
//	@Summary		Retrieve a single HealthPractitioner record by ID
//	@Description	Fetches a HealthPractitioner record based on the provided ID.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"HealthPractitioner ID"
//	@Success		200	{object}	fiber.Map	"Health Practitioner retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"Health Practitioner not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving Health Practitioner"
//	@Router			/health-practitioner/{id} [get]
func (h *HealthPractitionerController) GetSingleHealthPractitioner(c *fiber.Ctx) error {
	// Get the healthPractitioner ID from the route parameters
	id := c.Params("id")

	// Fetch the healthPractitioner by ID
	healthPractitioner, err := h.repo.GetHealthPractitionerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "healthPractitioner not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve healthPractitioner",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "healthPractitioner and associated data retrieved successfully",
		"data":    healthPractitioner,
	})
}

// =======================

// Define the UpdateHealthPractitioner struct
type UpdateHealthPractitionerPayload struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Phone      string `json:"phone"`
	Profession string `json:"profession"` // e.g., Doctor, Nurse, Counselor
	FacilityID uint   `json:"facility_id"`
}

// UpdateHealthPractitioner godoc
//
//	@Summary		Update an existing Health practitioner record by ID
//	@Description	Updates the details of a HealthPractitioner record based on the provided ID and request body.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string							true	"Health practitioner ID"
//	@Param			HealthPractitioner	body		UpdateHealthPractitionerPayload	true	"Health practitioner data to update"
//	@Success		200					{object}	fiber.Map						"Health practitioner updated successfully"
//	@Failure		400					{object}	fiber.Map						"Invalid input or empty request body"
//	@Failure		404					{object}	fiber.Map						"HealthPractitioner not found"
//	@Failure		500					{object}	fiber.Map						"Server error when updating Health practitioner"
//	@Router			/health-practitioner/{id} [put]
func (h *HealthPractitionerController) UpdateHealthPractitioner(c *fiber.Ctx) error {
	// Get the HealthPractitioner ID from the route parameters
	id := c.Params("id")

	// Find the HealthPractitioner in the database
	_, err := h.repo.GetHealthPractitionerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "HealthPractitioner not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve HealthPractitioner",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateHealthPractitionerPayload struct
	var payload UpdateHealthPractitionerPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateHealthPractitionerPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if payload.FirstName != "" {
		updates["first_name"] = payload.FirstName
	}
	if payload.LastName != "" {
		updates["last_name"] = payload.LastName
	}
	if payload.Gender != "" {
		updates["gender"] = payload.Gender
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}
	if payload.Profession != "" {
		updates["profession"] = payload.Profession
	}
	if payload.FacilityID != 0 {
		updates["facility_id"] = payload.FacilityID
	}

	// Update the HealthPractitioner in the database
	if err := h.repo.UpdateHealthPractitioner(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update HealthPractitioner",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "HealthPractitioner updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteHealthPractitionerByID godoc
//
//	@Summary		Delete a HealthPractitioner record by ID
//	@Description	Deletes a Health Practitioner record based on the provided ID.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"HealthPractitioner ID"
//	@Success		200	{object}	fiber.Map	"Health Practitioner deleted successfully"
//	@Failure		404	{object}	fiber.Map	"HealthPractitioner not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting HealthPractitioner"
//	@Router			/health-practitioner/{id} [delete]
func (h *HealthPractitionerController) DeleteHealthPractitionerByID(c *fiber.Ctx) error {
	// Get the HealthPractitioner ID from the route parameters
	id := c.Params("id")

	// Find the HealthPractitioner in the database
	healthPractitioner, err := h.repo.GetHealthPractitionerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "HealthPractitioner not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find HealthPractitioner",
			"data":    err.Error(),
		})
	}

	// Delete the HealthPractitioner
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete HealthPractitioner",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "HealthPractitioner deleted successfully",
		"data":    healthPractitioner,
	})
}

// =================

// SearchHealthPractitioners godoc
//
//	@Summary		Search for health practitioners with pagination
//	@Description	Retrieves a paginated list of health practitioners based on search criteria.
//	@Tags			Health practitioners
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Health practitioners retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve HealthPractitioners"
//	@Router			/health-practitioners/search [get]
func (h *HealthPractitionerController) SearchHealthPractitioners(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, healthPractitioners, err := h.repo.SearchPaginatedHealthPractitioners(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve HealthPractitioners",
			"data":    err.Error(),
		})
	}

	hpResponses := make([]HealthPractitionerExamResponse, len(healthPractitioners))
	for i, hp := range healthPractitioners {
		hpResponses[i] = ConvertToHealthPractitionerResponse(hp)
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Health practitioners retrieved successfully",
		"data":    hpResponses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
