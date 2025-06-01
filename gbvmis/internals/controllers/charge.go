package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChargeController struct {
	repo repository.ChargeRepository
}

func NewChargeController(repo repository.ChargeRepository) *ChargeController {
	return &ChargeController{repo: repo}
}

// ================================

// CreateCharge godoc
//
//	@Summary		Create a new charge record
//	@Description	Creates a new charge entry in the system and returns the created record.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Param			charge	body		models.Charge	true	"Charge data to create"
//	@Success		201		{object}	fiber.Map	"Successfully created charge record"
//	@Failure		400		{object}	fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map	"Server error when creating charge"
//	@Router			/charge [post]
func (h *ChargeController) CreateCharge(c *fiber.Ctx) error {
	// Initialize a new charge instance
	charge := new(models.Charge)

	// Parse the request body into the charge instance
	if err := c.BodyParser(charge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the charge record using the repository
	if err := h.repo.CreateCharge(charge); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create charge",
			"data":    err.Error(),
		})
	}

	// Return the newly created charge record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "charge created successfully",
		"data":    charge,
	})
}

// ===========

// GetAllCharges godoc
//
//	@Summary		Retrieve a paginated list of charges
//	@Description	Fetches all charge records with pagination support.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Charges retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve charges"
//	@Router			/charges [get]
func (h *ChargeController) GetAllCharges(c *fiber.Ctx) error {
	pagination, charges, err := h.repo.GetPaginatedCharges(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve charges",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "charges retrieved successfully",
		"data":    charges,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleCharge godoc
//
//	@Summary		Retrieve a single charge record by ID
//	@Description	Fetches a charge record based on the provided ID.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Charge ID"
//	@Success		200		{object}	fiber.Map	"Charge retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"Charge not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving charge"
//	@Router			/charge/{id} [get]
func (h *ChargeController) GetSingleCharge(c *fiber.Ctx) error {
	// Get the Charge ID from the route parameters
	id := c.Params("id")

	// Fetch the charge by ID
	charge, err := h.repo.GetChargeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Charge not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve charge",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Charge and associated data retrieved successfully",
		"data":    charge,
	})
}

// =======================

// Define the UpdateCharge struct
type UpdateChargePayload struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Dob         string `json:"dob"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Nationality string `json:"nationality"`
	Nin         string `json:"nin"`
	UpdatedBy   string `json:"updated_by"`
}

// UpdateCharge godoc
//
//	@Summary		Update an existing charge record by ID
//	@Description	Updates the details of a charge record based on the provided ID and request body.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Charge ID"
//	@Param			charge	body		UpdateChargePayload	true	"Charge data to update"
//	@Success		200		{object}	fiber.Map	"Charge updated successfully"
//	@Failure		400		{object}	fiber.Map	"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map	"Charge not found"
//	@Failure		500		{object}	fiber.Map	"Server error when updating charge"
//	@Router			/charge/{id} [put]
func (h *ChargeController) UpdateCharge(c *fiber.Ctx) error {
	// Get the charge ID from the route parameters
	id := c.Params("id")

	// Find the charge in the database
	_, err := h.repo.GetChargeByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Charge not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve charge",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateChargePayload struct
	var payload UpdateChargePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateChargePayload{} == payload) {
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
	if payload.Dob != "" {
		updates["dob"] = payload.Dob
	}
	if payload.PhoneNumber != "" {
		updates["phone_number"] = payload.PhoneNumber
	}
	if payload.Address != "" {
		updates["address"] = payload.Address
	}
	if payload.Nationality != "" {
		updates["nationality"] = payload.Nationality
	}
	if payload.Nin != "" {
		updates["nin"] = payload.Nin
	}
	if payload.UpdatedBy != "" {
		updates["updated_by"] = payload.UpdatedBy
	}

	// Update the Charge in the database
	if err := h.repo.UpdateCharge(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update charge",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Charge updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteChargeByID godoc
//
//	@Summary		Delete a charge record by ID
//	@Description	Deletes a charge record based on the provided ID.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Charge ID"
//	@Success		200		{object}	fiber.Map	"Charge deleted successfully"
//	@Failure		404		{object}	fiber.Map	"Charge not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting charge"
//	@Router			/charge/{id} [delete]
func (h *ChargeController) DeleteChargeByID(c *fiber.Ctx) error {
	// Get the Charge ID from the route parameters
	id := c.Params("id")

	// Find the Charge in the database
	charge, err := h.repo.GetChargeByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Charge not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Charge",
			"data":    err.Error(),
		})
	}

	// Delete the Charge
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Charge",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Charge deleted successfully",
		"data":    charge,
	})
}

// =================

// SearchCharges godoc
//
//	@Summary		Search for charges with pagination
//	@Description	Retrieves a paginated list of charges based on search criteria.
//	@Tags			Charges
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Charges retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve charges"
//	@Router			/charges/search [get]
func (h *ChargeController) SearchCharges(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, charges, err := h.repo.SearchPaginatedCharges(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve charges",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Charges retrieved successfully",
		"pagination": pagination,
		"data":       charges,
	})
}
