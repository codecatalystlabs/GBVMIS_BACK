package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"gbvmis/internals/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CaseController struct {
	repo repository.CaseRepository
}

func NewCaseController(repo repository.CaseRepository) *CaseController {
	return &CaseController{repo: repo}
}

type CreateCasePayload struct {
	CaseNumber   string    `json:"case_number" validate:"required"`
	Title        string    `json:"title" validate:"required"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	DateOpened   time.Time `json:"date_opened"`
	SuspectID    uint      `json:"suspect_id"`
	OfficerID    uint      `json:"officer_id"`
	PolicePostID uint      `json:"police_post_id"`

	Charges   []ChargePayload `json:"charges"`    // Optional inline
	VictimIDs []uint          `json:"victim_ids"` // For existing victims
}

type ChargePayload struct {
	ChargeTitle string `json:"charge_title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// =======

type CaseResponse struct {
	ID          uint      `json:"id"`
	CaseNumber  string    `json:"case_number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DateOpened  time.Time `json:"date_opened"`

	SuspectID    uint `json:"suspect_id"`
	OfficerID    uint `json:"officer_id"`
	PolicePostID uint `json:"police_post_id"`

	Charges []ChargeResponse `json:"charges"`
	Victims []VictimResponse `json:"victims"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChargeResponse struct {
	ID          uint   `json:"id"`
	ChargeTitle string `json:"charge_title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type VictimResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
}

func ConvertToCaseResponse(casee models.Case) CaseResponse {
	var charges []ChargeResponse
	for _, ch := range casee.Charges {
		charges = append(charges, ChargeResponse{
			ID: ch.ID, ChargeTitle: ch.ChargeTitle, Description: ch.Description, Severity: ch.Severity,
		})
	}

	var victims []VictimResponse
	for _, v := range casee.Victims {
		victims = append(victims, VictimResponse{
			ID: v.ID, FirstName: v.FirstName, LastName: v.LastName,
			Gender: v.Gender, PhoneNumber: v.PhoneNumber,
		})
	}

	return CaseResponse{
		ID: casee.ID, CaseNumber: casee.CaseNumber, Title: casee.Title,
		Description: casee.Description, Status: casee.Status, DateOpened: casee.DateOpened,
		SuspectID: casee.SuspectID, OfficerID: casee.OfficerID, PolicePostID: casee.PolicePostID,
		Charges: charges, Victims: victims, CreatedAt: casee.CreatedAt, UpdatedAt: casee.UpdatedAt,
	}
}

// ================================

// CreateCase godoc
//
//	@Summary		Create a new case record
//	@Description	Creates a new case entry in the system and returns the created record.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			case	body		CreateCasePayload	true	"Case data to create"
//	@Success		201		{object}	fiber.Map			"Successfully created case record"
//	@Failure		400		{object}	fiber.Map			"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map			"Server error when creating case"
//	@Router			/case [post]
func (h *CaseController) CreateCase(c *fiber.Ctx) error {
	var payload CreateCasePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid input", err))
	}

	casee := &models.Case{
		CaseNumber:   payload.CaseNumber,
		Title:        payload.Title,
		Description:  payload.Description,
		Status:       payload.Status,
		DateOpened:   payload.DateOpened,
		SuspectID:    payload.SuspectID,
		OfficerID:    payload.OfficerID,
		PolicePostID: payload.PolicePostID,
	}

	// Create and attach charges
	for _, ch := range payload.Charges {
		casee.Charges = append(casee.Charges, models.Charge{
			ChargeTitle: ch.ChargeTitle,
			Description: ch.Description,
			Severity:    ch.Severity,
		})
	}

	// Attach existing victims by IDs
	if len(payload.VictimIDs) > 0 {
		var victims []models.Victim
		if err := h.repo.FindVictimsByIDs(payload.VictimIDs, &victims); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid victim IDs", err))
		}
		casee.Victims = victims
	}

	// Save case with associations
	if err := h.repo.CreateCase(casee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create case", err))
	}

	// Convert to response
	resp := ConvertToCaseResponse(*casee)
	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse("Case created successfully", resp))
}

// ===========

// GetAllCases godoc
//
//	@Summary		Retrieve a paginated list of cases
//	@Description	Fetches all case records with pagination support.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Cases retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve cases"
//	@Router			/cases [get]
func (h *CaseController) GetAllCases(c *fiber.Ctx) error {
	pagination, cases, err := h.repo.GetPaginatedCases(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to retrieve cases", err))
	}

	caseResponses := make([]CaseResponse, len(cases))
	for i, c := range cases {
		caseResponses[i] = ConvertToCaseResponse(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cases retrieved successfully",
		"data":    caseResponses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleCase godoc
//
//	@Summary		Retrieve a single case record by ID
//	@Description	Fetches a case record based on the provided ID.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Case ID"
//	@Success		200	{object}	fiber.Map	"Case retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"Case not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving case"
//	@Router			/case/{id} [get]
func (h *CaseController) GetSingleCase(c *fiber.Ctx) error {
	id := c.Params("id")

	// Fetch case with preloaded relations
	casee, err := h.repo.GetCaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Case not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve case",
			"data":    err.Error(),
		})
	}

	// Convert to safe response structure
	caseResponse := ConvertToCaseResponse(casee)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Case and associated data retrieved successfully",
		"data":    caseResponse,
	})
}

// =======================

// Define the UpdateCase struct
type ChargeUpdatePayload struct {
	ChargeTitle string `json:"charge_title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type UpdateCasePayload struct {
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	Status       string                `json:"status"`
	DateOpened   time.Time             `json:"date_opened"`
	SuspectID    uint                  `json:"suspect_id"`
	OfficerID    uint                  `json:"officer_id"`
	PolicePostID uint                  `json:"police_post_id"`
	Charges      []ChargeUpdatePayload `json:"charges"`    // <== new field
	VictimIDs    []uint                `json:"victim_ids"` // NEW
}

func (p UpdateCasePayload) IsEmpty() bool {
	return p.Title == "" &&
		p.Description == "" &&
		p.Status == "" &&
		p.SuspectID == 0 &&
		p.OfficerID == 0 &&
		p.PolicePostID == 0 &&
		p.DateOpened.IsZero() &&
		len(p.Charges) == 0 &&
		len(p.VictimIDs) == 0
}

// UpdateCase godoc
//
//	@Summary		Update an existing case record by ID
//	@Description	Updates the details of a case record based on the provided ID and request body.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Case ID"
//	@Param			case	body		UpdateCasePayload	true	"Case data to update"
//	@Success		200		{object}	fiber.Map			"Case updated successfully"
//	@Failure		400		{object}	fiber.Map			"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map			"Case not found"
//	@Failure		500		{object}	fiber.Map			"Server error when updating case"
//	@Router			/case/{id} [put]
func (h *CaseController) UpdateCase(c *fiber.Ctx) error {
	id := c.Params("id")

	caseRecord, err := h.repo.GetCaseByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Case not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve case",
			"data":    err.Error(),
		})
	}

	var payload UpdateCasePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	if payload.IsEmpty() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Begin DB transaction
	tx := h.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{}

	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Description != "" {
		updates["description"] = payload.Description
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}
	if !payload.DateOpened.IsZero() {
		updates["date_opened"] = payload.DateOpened
	}
	if payload.SuspectID != 0 {
		updates["suspect_id"] = payload.SuspectID
	}
	if payload.OfficerID != 0 {
		updates["officer_id"] = payload.OfficerID
	}
	if payload.PolicePostID != 0 {
		updates["police_post_id"] = payload.PolicePostID
	}

	if err := tx.Model(caseRecord).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update case",
			"data":    err.Error(),
		})
	}

	// Update charges if provided
	if len(payload.Charges) > 0 {
		// Delete existing charges
		if err := tx.Where("case_id = ?", caseRecord.ID).Delete(&models.Charge{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to delete existing charges",
				"data":    err.Error(),
			})
		}

		// Insert new charges
		for _, ch := range payload.Charges {
			newCharge := models.Charge{
				CaseID:      caseRecord.ID,
				ChargeTitle: ch.ChargeTitle,
				Description: ch.Description,
				Severity:    ch.Severity,
			}
			if err := tx.Create(&newCharge).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to create charge",
					"data":    err.Error(),
				})
			}
		}
	}

	if len(payload.VictimIDs) > 0 {
		var victims []models.Victim
		if err := tx.Where("id IN ?", payload.VictimIDs).Find(&victims).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid victim IDs",
				"data":    err.Error(),
			})
		}

		// Replace victims association
		if err := tx.Model(caseRecord).Association("Victims").Replace(victims); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update case victims",
				"data":    err.Error(),
			})
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to finalize update",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Case and charges updated successfully",
	})
}

// ==================

// DeleteCaseByID godoc
//
//	@Summary		Delete a Case record by ID
//	@Description	Deletes a Case record based on the provided ID.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Case ID"
//	@Success		200	{object}	fiber.Map	"Case deleted successfully"
//	@Failure		404	{object}	fiber.Map	"Case not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting case"
//	@Router			/case/{id} [delete]
func (h *CaseController) DeleteCaseByID(c *fiber.Ctx) error {
	// Get the Case ID from the route parameters
	id := c.Params("id")

	// Find the Case in the database
	casee, err := h.repo.GetCaseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Case not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Case",
			"data":    err.Error(),
		})
	}

	// Delete the Case
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Case",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Case deleted successfully",
		"data":    casee,
	})
}

// =================

// SearchCases godoc
//
//	@Summary		Search for Cases with pagination
//	@Description	Retrieves a paginated list of cases based on search criteria.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Cases retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve Cases"
//	@Router			/cases/search [get]
func (h *CaseController) SearchCases(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, cases, err := h.repo.SearchPaginatedCases(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cases",
			"data":    err.Error(),
		})
	}

	caseResponses := make([]CaseResponse, len(cases))
	for i, c := range cases {
		caseResponses[i] = ConvertToCaseResponse(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Cases retrieved successfully",
		"data":    caseResponses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
