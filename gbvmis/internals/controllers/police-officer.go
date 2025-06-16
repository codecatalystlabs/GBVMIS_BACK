package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"gbvmis/internals/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PoliceOfficerController struct {
	repo repository.PoliceOfficerRepository
}

func NewPoliceOfficerController(repo repository.PoliceOfficerRepository) *PoliceOfficerController {
	return &PoliceOfficerController{repo: repo}
}

type PoliceOfficerResponse struct {
	ID        uint           `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Rank      string         `json:"rank"`
	BadgeNo   string         `json:"badge_no"`
	Phone     string         `json:"phone"`
	PostID    uint           `json:"post_id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Roles     []RoleResponse `json:"roles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Nested response structs
type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateOfficerPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Rank      string `json:"rank" validate:"required"`
	BadgeNo   string `json:"badge_no" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	PostID    uint   `json:"post_id" validate:"required"`

	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`

	RoleIDs []uint `json:"role_ids"` // Optional: to assign roles
}

func ConvertToPoliceOfficerResponse(officer models.PoliceOfficer) PoliceOfficerResponse {
	roles := make([]RoleResponse, len(officer.Roles))
	for i, role := range officer.Roles {
		roles[i] = RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}

	return PoliceOfficerResponse{
		ID:        officer.ID,
		FirstName: officer.FirstName,
		LastName:  officer.LastName,
		Rank:      officer.Rank,
		BadgeNo:   officer.BadgeNo,
		Phone:     officer.Phone,
		PostID:    officer.PostID,
		Username:  officer.Username,
		Email:     officer.Email,
		Roles:     roles,
		CreatedAt: officer.CreatedAt,
		UpdatedAt: officer.UpdatedAt,
	}
}

// ================================

// CreatePoliceOfficer godoc
//
//	@Summary		Create a new police Officer record
//	@Description	Creates a new police Officer entry in the system and returns the created record.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			PoliceOfficer	body		CreateOfficerPayload	true	"PoliceOfficer data to create"
//	@Success		201				{object}	fiber.Map				"Successfully created police Officer record"
//	@Failure		400				{object}	fiber.Map				"Bad request due to invalid input"
//	@Failure		500				{object}	fiber.Map				"Server error when creating police Officer"
//	@Router			/police-officer [post]
func (h *PoliceOfficerController) CreatePoliceOfficer(c *fiber.Ctx) error {
	var payload CreateOfficerPayload

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	// You can add validation logic here if you're using a validator
	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "username, email, and password are required",
		})
	}

	hashed, err := utils.HashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not hash password",
			"data":    err.Error(),
		})
	}

	// Build PoliceOfficer from payload
	officer := &models.PoliceOfficer{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Rank:      payload.Rank,
		BadgeNo:   payload.BadgeNo,
		Phone:     payload.Phone,
		PostID:    payload.PostID,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashed,
	}

	// Optional: Attach roles
	if len(payload.RoleIDs) > 0 {
		var roles []*models.Role
		if err := h.repo.FindRolesByIDs(payload.RoleIDs, &roles); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch roles",
				"data":    err.Error(),
			})
		}
		officer.Roles = roles
	}

	// Save
	if err := h.repo.CreatePoliceOfficer(officer); err != nil {
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

	// Format response
	response := ConvertToPoliceOfficerResponse(*officer)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officer created successfully",
		"data":    response,
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
//	@Success		200	{object}	fiber.Map	"policeOfficers retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve policeOfficers"
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
	officerResponses := make([]PoliceOfficerResponse, len(policeOfficers))
	for i, officer := range policeOfficers {
		officerResponses[i] = ConvertToPoliceOfficerResponse(officer)
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officers retrieved successfully",
		"data":    officerResponses,
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
//	@Param			id	path		string		true	"PoliceOfficer ID"
//	@Success		200	{object}	fiber.Map	"PoliceOfficer retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"PoliceOfficer not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving PoliceOfficer"
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
				"message": "Police officer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve police officer",
			"data":    err.Error(),
		})
	}

	// Convert to response struct
	response := ConvertToPoliceOfficerResponse(policeOfficer)

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officer and associated data retrieved successfully",
		"data":    response,
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
	Password  string `json:"password"`
	RoleIDs   []uint `json:"role_ids"`
}

func (p UpdatePoliceOfficerPayload) IsEmpty() bool {
	return p.FirstName == "" &&
		p.LastName == "" &&
		p.Rank == "" &&
		p.BadgeNo == "" &&
		p.Phone == "" &&
		p.PostID == 0 &&
		p.Email == "" &&
		p.Password == "" &&
		len(p.RoleIDs) == 0
}

// UpdatePoliceOfficer godoc
//
//	@Summary		Update an existing policeOfficer record by ID
//	@Description	Updates the details of a policeOfficer record based on the provided ID and request body.
//	@Tags			Police Officers
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string						true	"PoliceOfficer ID"
//	@Param			policeOfficer	body		UpdatePoliceOfficerPayload	true	"PoliceOfficer data to update"
//	@Success		200				{object}	fiber.Map					"PoliceOfficer updated successfully"
//	@Failure		400				{object}	fiber.Map					"Invalid input or empty request body"
//	@Failure		404				{object}	fiber.Map					"PoliceOfficer not found"
//	@Failure		500				{object}	fiber.Map					"Server error when updating policeOfficer"
//	@Router			/police-officer/{id} [put]
func (h *PoliceOfficerController) UpdatePoliceOfficer(c *fiber.Ctx) error {
	id := c.Params("id")

	// Step 1: Check if the police officer exists
	officer, err := h.repo.GetPoliceOfficerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Police officer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve police officer",
			"data":    err.Error(),
		})
	}

	// Step 2: Parse input
	var payload UpdatePoliceOfficerPayload
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
			"message": "No update fields provided",
		})
	}

	// Step 3: Build update map
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
	if payload.Password != "" {
		hashed, err := utils.HashPassword(payload.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to hash password",
				"data":    err.Error(),
			})
		}
		updates["password"] = hashed
	}

	// Update roles if provided
	if payload.RoleIDs != nil {
		var roles []*models.Role
		if err := h.repo.FindRolesByIDs(payload.RoleIDs, &roles); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch roles",
				"data":    err.Error(),
			})
		}
		// Replace roles
		if err := h.repo.UpdateOfficerRoles(&officer, roles); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update roles",
				"data":    err.Error(),
			})
		}
	}

	// Step 4: Perform update
	if err := h.repo.UpdatePoliceOfficer(id, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update police officer",
			"data":    err.Error(),
		})
	}

	// Step 5: Load updated officer
	updatedOfficer, err := h.repo.GetPoliceOfficerByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to load updated police officer",
			"data":    err.Error(),
		})
	}

	// Step 6: Return clean response
	response := ConvertToPoliceOfficerResponse(updatedOfficer)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officer updated successfully",
		"data":    response,
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
//	@Param			id	path		string		true	"PoliceOfficer ID"
//	@Success		200	{object}	fiber.Map	"PoliceOfficer deleted successfully"
//	@Failure		404	{object}	fiber.Map	"PoliceOfficer not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting policeOfficer"
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
	response := ConvertToPoliceOfficerResponse(policeOfficer)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "PoliceOfficer deleted successfully",
		"data":    response,
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
//	@Success		200	{object}	fiber.Map	"PoliceOfficers retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve PoliceOfficers"
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

	officerResponses := make([]PoliceOfficerResponse, len(policeOfficers))
	for i, officer := range policeOfficers {
		officerResponses[i] = ConvertToPoliceOfficerResponse(officer)
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police officers retrieved successfully",
		"data":    officerResponses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
