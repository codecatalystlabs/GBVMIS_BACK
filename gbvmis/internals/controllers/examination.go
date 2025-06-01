package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExaminationController struct {
	repo repository.ExaminationRepository
}

func NewExaminationController(repo repository.ExaminationRepository) *ExaminationController {
	return &ExaminationController{repo: repo}
}

// ================================

// CreateExamination godoc
//
//	@Summary		Create a new examination record
//	@Description	Creates a new examination entry in the system and returns the created record.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Param			examination	body		models.Examination	true	"Examination data to create"
//	@Success		201		{object}	fiber.Map	"Successfully created examination record"
//	@Failure		400		{object}	fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map	"Server error when creating examination"
//	@Router			/examination [post]
func (h *ExaminationController) CreateExamination(c *fiber.Ctx) error {
	// Initialize a new examination instance
	examination := new(models.Examination)

	// Parse the request body into the examination instance
	if err := c.BodyParser(examination); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the examination record using the repository
	if err := h.repo.CreateExamination(examination); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create examination",
			"data":    err.Error(),
		})
	}

	// Return the newly created examination record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "examination created successfully",
		"data":    examination,
	})
}

// ===========

// GetAllExaminations godoc
//
//	@Summary		Retrieve a paginated list of examinations
//	@Description	Fetches all examination records with pagination support.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Examinations retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve examinations"
//	@Router			/examinations [get]
func (h *ExaminationController) GetAllExaminations(c *fiber.Ctx) error {
	pagination, examinations, err := h.repo.GetPaginatedExaminations(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve examinations",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "examinations retrieved successfully",
		"data":    examinations,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleExamination godoc
//
//	@Summary		Retrieve a single examination record by ID
//	@Description	Fetches a examination record based on the provided ID.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Examination ID"
//	@Success		200		{object}	fiber.Map	"Examination retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"Examination not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving examination"
//	@Router			/examination/{id} [get]
func (h *ExaminationController) GetSingleExamination(c *fiber.Ctx) error {
	// Get the Examination ID from the route parameters
	id := c.Params("id")

	// Fetch the examination by ID
	examination, err := h.repo.GetExaminationByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Examination not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Examination",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Examination and associated data retrieved successfully",
		"data":    examination,
	})
}

// =======================

// Define the UpdateExamination struct
type UpdateExaminationPayload struct {
	VictimID       uint   `json:"victim_id"`
	CaseID         uint   `json:"case_id"`
	FacilityID     uint   `json:"facility_id"`
	PractitionerID uint   `json:"practitioner_id"`
	ExamDate       string `json:"exam_date"`
	Findings       string `json:"findings"`
	Treatment      string `json:"treatment"`
	Referral       string `json:"referral"` // Optional referral info
	ConsentGiven   bool   `json:"consent_given"`
}

// UpdateExamination godoc
//
//	@Summary		Update an existing examination record by ID
//	@Description	Updates the details of a examination record based on the provided ID and request body.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Examination ID"
//	@Param			examination	body		UpdateExaminationPayload	true	"Examination data to update"
//	@Success		200		{object}	fiber.Map	"Examination updated successfully"
//	@Failure		400		{object}	fiber.Map	"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map	"Examination not found"
//	@Failure		500		{object}	fiber.Map	"Server error when updating examination"
//	@Router			/examination/{id} [put]
func (h *ExaminationController) UpdateExamination(c *fiber.Ctx) error {
	// Get the examination ID from the route parameters
	id := c.Params("id")

	// Find the examination in the database
	_, err := h.repo.GetExaminationByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Examination not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve examination",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateExaminationPayload struct
	var payload UpdateExaminationPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateExaminationPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if payload.VictimID != 0 {
		updates["victim_id"] = payload.VictimID
	}
	if payload.CaseID != 0 {
		updates["case_id"] = payload.CaseID
	}
	if payload.FacilityID != 0 {
		updates["facility_id"] = payload.FacilityID
	}
	if payload.PractitionerID != 0 {
		updates["practitioner_id"] = payload.PractitionerID
	}
	if payload.ExamDate != "" {
		updates["exam_date"] = payload.ExamDate
	}
	if payload.Findings != "" {
		updates["findings"] = payload.Findings
	}
	if payload.Treatment != "" {
		updates["treatment"] = payload.Treatment
	}
	if payload.Referral != "" {
		updates["referral"] = payload.Referral
	}
	if payload.ConsentGiven != false {
		updates["consent_given"] = payload.ConsentGiven
	}

	// Update the Examination in the database
	if err := h.repo.UpdateExamination(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update examination",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Examination updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteExaminationByID godoc
//
//	@Summary		Delete a examination record by ID
//	@Description	Deletes a examination record based on the provided ID.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Examination ID"
//	@Success		200		{object}	fiber.Map	"Examination deleted successfully"
//	@Failure		404		{object}	fiber.Map	"Examination not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting examination"
//	@Router			/examination/{id} [delete]
func (h *ExaminationController) DeleteExaminationByID(c *fiber.Ctx) error {
	// Get the Examination ID from the route parameters
	id := c.Params("id")

	// Find the Examination in the database
	examination, err := h.repo.GetExaminationByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Examination not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Examination",
			"data":    err.Error(),
		})
	}

	// Delete the Examination
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Examination",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Examination deleted successfully",
		"data":    examination,
	})
}

// =================

// SearchExaminations godoc
//
//	@Summary		Search for examinations with pagination
//	@Description	Retrieves a paginated list of examinations based on search criteria.
//	@Tags			Examinations
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Examinations retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve examinations"
//	@Router			/examinations/search [get]
func (h *ExaminationController) SearchExaminations(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, examinations, err := h.repo.SearchPaginatedExaminations(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve examinations",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Examinations retrieved successfully",
		"pagination": pagination,
		"data":       examinations,
	})
}
