package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthFacilityController struct {
	repo repository.HealthFacilityRepository
}

func NewHealthFacilityController(repo repository.HealthFacilityRepository) *HealthFacilityController {
	return &HealthFacilityController{repo: repo}
}

// ================================

// CreateHealthFacility godoc
//
//	@Summary		Create a new health facility record
//	@Description	Creates a new health facility entry in the system and returns the created record.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Param			HealthFacility	body		models.HealthFacility	true	"HealthFacility data to create"
//	@Success		201		{object}	fiber.Map	"Successfully created health facility record"
//	@Failure		400		{object}	fiber.Map	"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map	"Server error when creating health facility"
//	@Router			/health-facility [post]
func (h *HealthFacilityController) CreateHealthFacility(c *fiber.Ctx) error {
	// Initialize a new health facility instance
	healthFacility := new(models.HealthFacility)

	// Parse the request body into the HealthFacility instance
	if err := c.BodyParser(healthFacility); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the healthFacility record using the repository
	if err := h.repo.CreateHealthFacility(healthFacility); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create health facility",
			"data":    err.Error(),
		})
	}

	// Return the newly created healthFacility record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "health facility created successfully",
		"data":    healthFacility,
	})
}

// ===========

// GetAllHealthFacilities godoc
//
//	@Summary		Retrieve a paginated list of health facilities
//	@Description	Fetches all healthFacility records with pagination support.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Health facilities retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve health facilities"
//	@Router			/health-facilities [get]
func (h *HealthFacilityController) GetAllHealthFacilities(c *fiber.Ctx) error {
	pagination, healthFacilities, err := h.repo.GetPaginatedHealthFacilities(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve health facilities",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Health facilities retrieved successfully",
		"data":    healthFacilities,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleHealthFacility godoc
//
//	@Summary		Retrieve a single HealthFacility record by ID
//	@Description	Fetches a HealthFacility record based on the provided ID.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"HealthFacility ID"
//	@Success		200		{object}	fiber.Map	"Health facility retrieved successfully"
//	@Failure		404		{object}	fiber.Map	"Health facility not found"
//	@Failure		500		{object}	fiber.Map	"Server error when retrieving Health facility"
//	@Router			/health-facility/{id} [get]
func (h *HealthFacilityController) GetSingleHealthFacility(c *fiber.Ctx) error {
	// Get the healthFacility ID from the route parameters
	id := c.Params("id")

	// Fetch the healthFacility by ID
	healthFacility, err := h.repo.GetHealthFacilityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "healthFacility not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve healthFacility",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "healthFacility and associated data retrieved successfully",
		"data":    healthFacility,
	})
}

// =======================

// Define the UpdateHealthFacility struct
type UpdateHealthFacilityPayload struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

// UpdateHealthFacility godoc
//
//	@Summary		Update an existing Health facility record by ID
//	@Description	Updates the details of a HealthFacility record based on the provided ID and request body.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Health facility ID"
//	@Param			HealthFacility	body		UpdateHealthFacilityPayload	true	"Health facility data to update"
//	@Success		200		{object}	fiber.Map	"Health facility updated successfully"
//	@Failure		400		{object}	fiber.Map	"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map	"HealthFacility not found"
//	@Failure		500		{object}	fiber.Map	"Server error when updating Health facility"
//	@Router			/health-facility/{id} [put]
func (h *HealthFacilityController) UpdateHealthFacility(c *fiber.Ctx) error {
	// Get the HealthFacility ID from the route parameters
	id := c.Params("id")

	// Find the HealthFacility in the database
	_, err := h.repo.GetHealthFacilityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "HealthFacility not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve HealthFacility",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateHealthFacilityPayload struct
	var payload UpdateHealthFacilityPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateHealthFacilityPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Location != "" {
		updates["location"] = payload.Location
	}
	if payload.Contact != "" {
		updates["contact"] = payload.Contact
	}

	// Update the HealthFacility in the database
	if err := h.repo.UpdateHealthFacility(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update HealthFacility",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "HealthFacility updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteHealthFacilityByID godoc
//
//	@Summary		Delete a HealthFacility record by ID
//	@Description	Deletes a Health facility record based on the provided ID.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"HealthFacility ID"
//	@Success		200		{object}	fiber.Map	"Health facility deleted successfully"
//	@Failure		404		{object}	fiber.Map	"HealthFacility not found"
//	@Failure		500		{object}	fiber.Map	"Server error when deleting HealthFacility"
//	@Router			/health-facility/{id} [delete]
func (h *HealthFacilityController) DeleteHealthFacilityByID(c *fiber.Ctx) error {
	// Get the HealthFacility ID from the route parameters
	id := c.Params("id")

	// Find the HealthFacility in the database
	healthFacility, err := h.repo.GetHealthFacilityByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "HealthFacility not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find HealthFacility",
			"data":    err.Error(),
		})
	}

	// Delete the HealthFacility
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete HealthFacility",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "HealthFacility deleted successfully",
		"data":    healthFacility,
	})
}

// =================

// SearchHealthFacilities godoc
//
//	@Summary		Search for health facilities with pagination
//	@Description	Retrieves a paginated list of health facilities based on search criteria.
//	@Tags			Health facilities
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	fiber.Map	"Health facilities retrieved successfully"
//	@Failure		500		{object}	fiber.Map	"Failed to retrieve HealthFacilities"
//	@Router			/health-facilities/search [get]
func (h *HealthFacilityController) SearchHealthFacilities(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, healthFacilities, err := h.repo.SearchPaginatedHealthFacilities(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve HealthFacilities",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "HealthFacilities retrieved successfully",
		"pagination": pagination,
		"data":       healthFacilities,
	})
}
