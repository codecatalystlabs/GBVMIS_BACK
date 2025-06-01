package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
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

// ================================

// CreateCase godoc
//
//	@Summary		Create a new case record
//	@Description	Creates a new case entry in the system and returns the created record.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			case	body		models.Case	true	"Case data to create"
//	@Success		201		{object}	fiber.Map	"Successfully created case record"
//	@Failure		400		{object}	fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map	"Server error when creating case"
//	@Router			/case [post]
func (h *CaseController) CreateCase(c *fiber.Ctx) error {
	// Initialize a new case instance
	casee := new(models.Case)

	// Parse the request body into the Case instance
	if err := c.BodyParser(casee); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the case record using the repository
	if err := h.repo.CreateCase(casee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create case",
			"data":    err.Error(),
		})
	}

	// Return the newly created case record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "case created successfully",
		"data":    casee,
	})
}

// ===========

// GetAllCases godoc
//
//	@Summary		Retrieve a paginated list of cases
//	@Description	Fetches all case records with pagination support.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Cases retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve cases"
//	@Router			/cases [get]
func (h *CaseController) GetAllCases(c *fiber.Ctx) error {
	pagination, cases, err := h.repo.GetPaginatedCases(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve cases",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "cases retrieved successfully",
		"data":    cases,
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
//	@Param			id		path		string	true	"Case ID"
//	@Success		200		{object}	fiber.Map	"Case retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"Case not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving case"
//	@Router			/case/{id} [get]
func (h *CaseController) GetSingleCase(c *fiber.Ctx) error {
	// Get the Case ID from the route parameters
	id := c.Params("id")

	// Fetch the case by ID
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

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Case and associated data retrieved successfully",
		"data":    casee,
	})
}

// =======================

// Define the UpdateCase struct
type UpdateCasePayload struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	DateOpened   time.Time `json:"date_opened"`
	SuspectID    uint      `json:"suspect_id"`
	OfficerID    uint      `json:"officer_id"`
	PolicePostID uint      `json:"police_post_id"`
}

// UpdateCase godoc
//
//	@Summary		Update an existing case record by ID
//	@Description	Updates the details of a case record based on the provided ID and request body.
//	@Tags			Cases
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Case ID"
//	@Param			case	body		UpdateCasePayload	true	"Case data to update"
//	@Success		200		{object}	fiber.Map	"Case updated successfully"
//	@Failure		400		{object}	fiber.Map	"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map	"Case not found"
//	@Failure		500		{object}	fiber.Map	"Server error when updating case"
//	@Router			/case/{id} [put]
func (h *CaseController) UpdateCase(c *fiber.Ctx) error {
	// Get the case ID from the route parameters
	id := c.Params("id")

	// Find the case in the database
	_, err := h.repo.GetCaseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Case not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve case",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateCasePayload struct
	var payload UpdateCasePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateCasePayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Description != "" {
		updates["description"] = payload.Description
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}
	// if payload.DateOpened != "" {
	// 	updates["date_opened"] = payload.DateOpened
	// }
	if payload.SuspectID != 0 {
		updates["phone_number"] = payload.SuspectID
	}
	if payload.OfficerID != 0 {
		updates["officer_id"] = payload.OfficerID
	}
	if payload.PolicePostID != 0 {
		updates["police_post_id"] = payload.PolicePostID
	}

	// Update the Case in the database
	if err := h.repo.UpdateCase(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update case",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Case updated successfully",
		"data":    updates,
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
//	@Param			id		path		string	true	"Case ID"
//	@Success		200		{object}	fiber.Map	"Case deleted successfully"
//	@Failure		404		{object}	fiber.Map	"Case not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting case"
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
//	@Success		200		{object}	fiber.Map	"Cases retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve Cases"
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

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Cases retrieved successfully",
		"pagination": pagination,
		"data":       cases,
	})
}
