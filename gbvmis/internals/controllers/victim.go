package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VictimController struct {
	repo repository.VictimRepository
}

func NewVictimController(repo repository.VictimRepository) *VictimController {
	return &VictimController{repo: repo}
}

// ================================

// CreateVictim godoc
//
//	@Summary		Create a new victim record
//	@Description	Creates a new victim entry in the system and returns the created record.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Param			victim	body		models.Victim	true	"Victim data to create"
//	@Success		201		{object}	fiber.Map		"Successfully created victim record"
//	@Failure		400		{object}	fiber.Map		"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map		"Server error when creating victim"
//	@Router			/victims [post]
func (h *VictimController) CreateVictim(c *fiber.Ctx) error {
	// Initialize a new victim instance
	victim := new(models.Victim)

	// Parse the request body into the victim instance
	if err := c.BodyParser(victim); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// Attempt to create the victim record using the repository
	if err := h.repo.CreateVictim(victim); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create victim",
			"data":    err.Error(),
		})
	}

	// Return the newly created victim record
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Victim created successfully",
		"data":    victim,
	})
}

// ===========

// GetAllVictims godoc
//
//	@Summary		Retrieve a paginated list of victims
//	@Description	Fetches all victim records with pagination support.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Victims retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve victims"
//	@Router			/victims [get]
func (h *VictimController) GetAllVictims(c *fiber.Ctx) error {
	pagination, victims, err := h.repo.GetPaginatedVictims(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve victims",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "victims retrieved successfully",
		"data":    victims,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleVictim godoc
//
//	@Summary		Retrieve a single victim record by ID
//	@Description	Fetches a victim record based on the provided ID.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Victim ID"
//	@Success		200	{object}	fiber.Map	"Victim retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"Victim not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving victim"
//	@Router			/victim/{id} [get]
func (h *VictimController) GetSingleVictim(c *fiber.Ctx) error {
	// Get the Victim ID from the route parameters
	id := c.Params("id")

	// Fetch the victim by ID
	victim, err := h.repo.GetVictimByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Victim not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Victim",
			"data":    err.Error(),
		})
	}

	// // Fetch victim addresses associated with the victim
	// addresses, err := h.repo.GetVictimAddresses(victim.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Failed to retrieve victim addresses",
	// 		"data":    err.Error(),
	// 	})
	// }

	// // Fetch victim contacts associated with the victim
	// contacts, err := h.repo.GetVictimContacts(victim.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Failed to retrieve victim contacts",
	// 		"data":    err.Error(),
	// 	})
	// }

	// // Prepare the response
	// response := fiber.Map{
	// 	"victim": victim,
	// 	"address":  addresses,
	// 	"contact":  contacts,
	// }

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Victim and associated data retrieved successfully",
		"data":    victim,
	})
}

// =======================

// Define the UpdateVictim struct
type UpdateVictimPayload struct {
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

// UpdateVictim godoc
//
//	@Summary		Update an existing victim record by ID
//	@Description	Updates the details of a victim record based on the provided ID and request body.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Victim ID"
//	@Param			victim	body		UpdateVictimPayload	true	"Victim data to update"
//	@Success		200		{object}	fiber.Map			"Victim updated successfully"
//	@Failure		400		{object}	fiber.Map			"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map			"Victim not found"
//	@Failure		500		{object}	fiber.Map			"Server error when updating victim"
//	@Router			/victim/{id} [put]
func (h *VictimController) UpdateVictim(c *fiber.Ctx) error {
	// Get the victim ID from the route parameters
	id := c.Params("id")

	// Find the victim in the database
	_, err := h.repo.GetVictimByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Victim not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve victim",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateVictimPayload struct
	var payload UpdateVictimPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateVictimPayload{} == payload) {
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

	// Update the Victim in the database
	if err := h.repo.UpdateVictim(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update victim",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Victim updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteVictimByID godoc
//
//	@Summary		Delete a victim record by ID
//	@Description	Deletes a victim record based on the provided ID.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Victim ID"
//	@Success		200	{object}	fiber.Map	"Victim deleted successfully"
//	@Failure		404	{object}	fiber.Map	"Victim not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting victim"
//	@Router			/victim/{id} [delete]
func (h *VictimController) DeleteVictimByID(c *fiber.Ctx) error {
	// Get the Victim ID from the route parameters
	id := c.Params("id")

	// Find the Victim in the database
	victim, err := h.repo.GetVictimByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Victim not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find Victim",
			"data":    err.Error(),
		})
	}

	// Delete the Victim
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Victim",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Victim deleted successfully",
		"data":    victim,
	})
}

// =================

// SearchVictims godoc
//
//	@Summary		Search for victims with pagination
//	@Description	Retrieves a paginated list of victims based on search criteria.
//	@Tags			Victims
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Victims retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve victims"
//	@Router			/victims/search [get]
func (h *VictimController) SearchVictims(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, victims, err := h.repo.SearchPaginatedVictims(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve victims",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Victims retrieved successfully",
		"pagination": pagination,
		"data":       victims,
	})
}
