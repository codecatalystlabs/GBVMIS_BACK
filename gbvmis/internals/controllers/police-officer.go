package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"gbvmis/internals/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PoliceOfficerController struct {
	repo repository.PoliceOfficerRepository
}

func NewPoliceOfficerController(repo repository.PoliceOfficerRepository) *PoliceOfficerController {
	return &PoliceOfficerController{repo: repo}
}

// ================================

// CreatePoliceOfficer godoc
//
//	@Summary		Create a new police Officer record
//	@Description	Creates a new police Officer entry in the system and returns the created record.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			PoliceOfficer	body		models.PoliceOfficer	true	"PoliceOfficer data to create"
//	@Success		201		{object}	fiber.Map	"Successfully created police Officer record"
//	@Failure		400		{object}	fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map	"Server error when creating police Officer"
//	@Router			/police-officer [post]
func (h *PoliceOfficerController) CreatePoliceOfficer(c *fiber.Ctx) error {
	// Initialize a new police officer instance
	policeOfficer := new(models.PoliceOfficer)

	// Parse the request body into the PoliceOfficer instance
	if err := c.BodyParser(policeOfficer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	if policeOfficer.Username == "" || policeOfficer.Email == "" || policeOfficer.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "username, email, and password are required",
		})
	}

	hashed, err := utils.HashPassword(policeOfficer.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
			"data":    err.Error(),
		})
	}
	policeOfficer.Password = hashed

	// Attempt to create the policeOfficer record using the repository
	if err := h.repo.CreatePoliceOfficer(policeOfficer); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Username or email already exists",
				"data":    err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create police officer",
			"data":    err.Error(),
		})
	}

	// Return the newly created policeOfficer record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "police officer created successfully",
		"data":    policeOfficer,
	})
}

// ===========

// GetAllPoliceOfficers godoc
//
//	@Summary		Retrieve a paginated list of policeOfficers
//	@Description	Fetches all policeOfficer records with pagination support.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"policeOfficers retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve policeOfficers"
//	@Router			/police-officers [get]
func (h *PoliceOfficerController) GetAllPoliceOfficers(c *fiber.Ctx) error {
	pagination, policeOfficers, err := h.repo.GetPaginatedPoliceOfficers(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve police officers",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officers retrieved successfully",
		"data":    policeOfficers,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSinglePoliceOfficer godoc
//
//	@Summary		Retrieve a single policeOfficer record by ID
//	@Description	Fetches a policeOfficer record based on the provided ID.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"PoliceOfficer ID"
//	@Success		200		{object}	fiber.Map	"PoliceOfficer retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"PoliceOfficer not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving PoliceOfficer"
//	@Router			/police-officer/{id} [get]
func (h *PoliceOfficerController) GetSinglePoliceOfficer(c *fiber.Ctx) error {
	// Get the policeOfficer ID from the route parameters
	id := c.Params("id")

	// Fetch the policeOfficer by ID
	policeOfficer, err := h.repo.GetPoliceOfficerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "policeOfficer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve PoliceOfficer",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "PoliceOfficer and associated data retrieved successfully",
		"data":    policeOfficer,
	})
}

// =======================

// Define the UpdatePoliceOfficer struct
type UpdatePoliceOfficerPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rank      string `json:"rank"`
	BadgeNo   string `json:"badge_no"`
	Phone     string `json:"phone"`
	PostID    uint   `json:"post_id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

// UpdatePoliceOfficer godoc
//
//	@Summary		Update an existing policeOfficer record by ID
//	@Description	Updates the details of a policeOfficer record based on the provided ID and request body.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"PoliceOfficer ID"
//	@Param			policeOfficer	body		UpdatePoliceOfficerPayload	true	"PoliceOfficer data to update"
//	@Success		200		{object}	fiber.Map	"PoliceOfficer updated successfully"
//	@Failure		400		{object}	fiber.Map	"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map	"PoliceOfficer not found"
//	@Failure		500		{object}	fiber.Map	"Server error when updating policeOfficer"
//	@Router			/police-officer/{id} [put]
func (h *PoliceOfficerController) UpdatePoliceOfficer(c *fiber.Ctx) error {
	// Get the policeOfficer ID from the route parameters
	id := c.Params("id")

	// Find the policeOfficer in the database
	_, err := h.repo.GetPoliceOfficerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "PoliceOfficer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve PoliceOfficer",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdatePoliceOfficerPayload struct
	var payload UpdatePoliceOfficerPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdatePoliceOfficerPayload{} == payload) {
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
	if payload.Rank != "" {
		updates["rank"] = payload.Rank
	}
	if payload.BadgeNo != "" {
		updates["badge_no"] = payload.BadgeNo
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}
	if payload.PostID != 0 {
		updates["post_id"] = payload.PostID
	}
	// if payload.UpdatedBy != "" {
	// 	updates["updated_by"] = payload.UpdatedBy
	// }

	// Update the PoliceOfficer in the database
	if err := h.repo.UpdatePoliceOfficer(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update policeOfficer",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "PoliceOfficer updated successfully",
		"data":    updates,
	})
}

// ==================

// DeletePoliceOfficerByID godoc
//
//	@Summary		Delete a PoliceOfficer record by ID
//	@Description	Deletes a PoliceOfficer record based on the provided ID.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"PoliceOfficer ID"
//	@Success		200		{object}	fiber.Map	"PoliceOfficer deleted successfully"
//	@Failure		404		{object}	fiber.Map	"PoliceOfficer not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting policeOfficer"
//	@Router			/police-officer/{id} [delete]
func (h *PoliceOfficerController) DeletePoliceOfficerByID(c *fiber.Ctx) error {
	// Get the PoliceOfficer ID from the route parameters
	id := c.Params("id")

	// Find the PoliceOfficer in the database
	policeOfficer, err := h.repo.GetPoliceOfficerByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "PoliceOfficer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find PoliceOfficer",
			"data":    err.Error(),
		})
	}

	// Delete the PoliceOfficer
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete policeOfficer",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "PoliceOfficer deleted successfully",
		"data":    policeOfficer,
	})
}

// =================

// SearchPoliceOfficers godoc
//
//	@Summary		Search for policeOfficers with pagination
//	@Description	Retrieves a paginated list of policeOfficers based on search criteria.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"PoliceOfficers retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve PoliceOfficers"
//	@Router			/police-officers/search [get]
func (h *PoliceOfficerController) SearchPoliceOfficers(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, policeOfficers, err := h.repo.SearchPaginatedPoliceOfficers(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve policeOfficers",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "PoliceOfficers retrieved successfully",
		"pagination": pagination,
		"data":       policeOfficers,
	})
}
