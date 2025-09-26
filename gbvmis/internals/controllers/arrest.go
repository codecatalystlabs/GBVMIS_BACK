package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ArrestController struct {
	repo repository.ArrestRepository
}

func NewArrestController(repo repository.ArrestRepository) *ArrestController {
	return &ArrestController{repo: repo}
}

type CreateArrestPayload struct {
	ArrestDate  time.Time `json:"arrest_date"`
	Location    string    `json:"location"`
	OfficerName string    `json:"officer_name"`
	SuspectID   uint      `json:"suspect_id"`
	Notes       string    `json:"notes"`
}

type ArrestResponse struct {
	ID          uint      `json:"id"`
	ArrestDate  time.Time `json:"arrest_date"`
	Location    string    `json:"location"`
	OfficerName string    `json:"officer_name"`
	SuspectID   uint      `json:"suspect_id"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

func ConvertToArrestResponse(a models.Arrest) ArrestResponse {

	return ArrestResponse{
		ID:          a.ID,
		ArrestDate:  a.ArrestDate,
		Location:    a.Location,
		OfficerName: a.OfficerName,
		SuspectID:   a.SuspectID,
		Notes:       a.Notes,
		CreatedAt:   a.CreatedAt,
	}
}

// ================================

// CreateArrest godoc
//
//	@Summary		Create a new arrest record
//	@Description	Creates a new arrest entry in the system and returns the created record.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Param			Arrest	body		CreateArrestPayload	true	"Arrest data to create"
//	@Success		201		{object}	fiber.Map			"Successfully created arrest record"
//	@Failure		400		{object}	fiber.Map			"Bad request due to invalid input"
//	@Failure		500		{object}	fiber.Map			"Server error when creating arrest"
//	@Router			/arrest [post]
func (h *ArrestController) CreateArrest(c *fiber.Ctx) error {
	var payload CreateArrestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	a := &models.Arrest{
		ArrestDate:  payload.ArrestDate,
		Location:    payload.Location,
		OfficerName: payload.OfficerName,
		SuspectID:   payload.SuspectID,
		Notes:       payload.Notes,
	}

	if err := h.repo.CreateArrest(a); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create arrest",
			"data":    err.Error(),
		})
	}

	response := ConvertToArrestResponse(*a)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Arrest created successfully",
		"data":    response,
	})
}

// ===========

// GetAllArrests godoc
//
//	@Summary		Retrieve a paginated list of arrests
//	@Description	Fetches all Arrest records with pagination support.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Arrests retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve arrests"
//	@Router			/arrests [get]
func (h *ArrestController) GetAllArrests(c *fiber.Ctx) error {
	pagination, arrests, err := h.repo.GetPaginatedArrests(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve arrests",
			"data":    err.Error(),
		})
	}

	var responses []ArrestResponse
	for _, a := range arrests {
		responses = append(responses, ConvertToArrestResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Arrests retrieved successfully",
		"data":    responses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleArrest godoc
//
//	@Summary		Retrieve a single Arrest record by ID
//	@Description	Fetches a Arrest record based on the provided ID.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Arrest ID"
//	@Success		200	{object}	fiber.Map	"Arrest retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"Arrest not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving arrest"
//	@Router			/arrest/{id} [get]
func (h *ArrestController) GetSingleArrest(c *fiber.Ctx) error {
	// Get the Arrest ID from the route parameters
	id := c.Params("id")

	// Fetch the Arrest by ID
	arrest, err := h.repo.GetArrestByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Arrest not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve arrest",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "arrest and associated data retrieved successfully",
		"data":    ConvertToArrestResponse(arrest),
	})
}

// =======================

// Define the UpdateArrest struct
type UpdateArrestPayload struct {
	ArrestDate  time.Time `json:"arrest_date"`
	Location    string    `json:"location"`
	OfficerName string    `json:"officer_name"`
	SuspectID   uint      `json:"suspect_id"`
	Notes       string    `json:"notes"`
}

// UpdateArrest godoc
//
//	@Summary		Update an existing arrest record by ID
//	@Description	Updates the details of a Arrest record based on the provided ID and request body.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Arrest ID"
//	@Param			Arrest	body		UpdateArrestPayload	true	"Arrest data to update"
//	@Success		200		{object}	fiber.Map			"Arrest updated successfully"
//	@Failure		400		{object}	fiber.Map			"Invalid input or empty request body"
//	@Failure		404		{object}	fiber.Map			"Arrest not found"
//	@Failure		500		{object}	fiber.Map			"Server error when updating Arrest"
//	@Router			/arrest/{id} [put]
func (h *ArrestController) UpdateArrest(c *fiber.Ctx) error {
	// Get the Arrest ID from the route parameters
	id := c.Params("id")

	// Find the Arrest in the database
	_, err := h.repo.GetArrestByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Arrest not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Arrest",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdateArrestPayload struct
	var payload UpdateArrestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdateArrestPayload{} == payload) {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Empty request body",
		})
	}

	// Convert payload to a map for partial update
	updates := make(map[string]interface{})

	if !payload.ArrestDate.IsZero() {
		updates["arrest_date"] = payload.ArrestDate
	}

	// Only include Location if it’s non-empty
	if payload.Location != "" {
		updates["location"] = payload.Location
	}

	// Only include OfficerName if it’s non-empty
	if payload.OfficerName != "" {
		updates["officer_name"] = payload.OfficerName
	}

	// Only include SuspectID if it’s non-zero
	if payload.SuspectID != 0 {
		updates["suspect_id"] = payload.SuspectID
	}

	// Only include Notes if it’s non-empty
	if payload.Notes != "" {
		updates["notes"] = payload.Notes
	}

	// Update the Arrest in the database
	if err := h.repo.UpdateArrest(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update arrest",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Arrest updated successfully",
		"data":    updates,
	})
}

// ==================

// DeleteArrestByID godoc
//
//	@Summary		Delete a Arrest record by ID
//	@Description	Deletes a Arrest record based on the provided ID.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Arrest ID"
//	@Success		200	{object}	fiber.Map	"Arrest deleted successfully"
//	@Failure		404	{object}	fiber.Map	"Arrest not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting Arrest"
//	@Router			/arrest/{id} [delete]
func (h *ArrestController) DeleteArrestByID(c *fiber.Ctx) error {
	// Get the Arrest ID from the route parameters
	id := c.Params("id")

	// Find the Arrest in the database
	arrest, err := h.repo.GetArrestByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Arrest not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find arrest",
			"data":    err.Error(),
		})
	}

	// Delete the Arrest
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete arrest",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Arrest deleted successfully",
		"data":    arrest,
	})
}

// =================

// SearchArrests godoc
//
//	@Summary		Search for arrests with pagination
//	@Description	Retrieves a paginated list of arrests based on search criteria.
//	@Tags			Arrests
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Arrests retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve Arrests"
//	@Router			/arrests/search [get]
func (h *ArrestController) SearchArrests(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, arrests, err := h.repo.SearchPaginatedArrests(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve arrests",
			"data":    err.Error(),
		})
	}

	var responses []ArrestResponse
	for _, a := range arrests {
		responses = append(responses, ConvertToArrestResponse(a))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Arrests retrieved successfully",
		"data":    responses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
