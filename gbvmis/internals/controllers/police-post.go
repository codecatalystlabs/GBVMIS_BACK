package controllers

import (
	"errors"
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PolicePostController struct {
	repo repository.PolicePostRepository
}

func NewPolicePostController(repo repository.PolicePostRepository) *PolicePostController {
	return &PolicePostController{repo: repo}
}

type CreatePolicePostPayload struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

type PolicePostResponse struct {
	ID        uint                        `json:"id"`
	Name      string                      `json:"name"`
	Location  string                      `json:"location"`
	Contact   string                      `json:"contact"`
	Officers  []PoliceOfficerPostResponse `json:"officers"`
	CreatedAt time.Time                   `json:"created_at"`
}

type PoliceOfficerPostResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Rank      string `json:"rank"`
	BadgeNo   string `json:"badge_no"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

func ConvertToPolicePostResponse(post models.PolicePost) PolicePostResponse {
	var officers []PoliceOfficerPostResponse
	for _, officer := range post.Officers {
		officers = append(officers, PoliceOfficerPostResponse{
			ID:        officer.ID,
			FirstName: officer.FirstName,
			LastName:  officer.LastName,
			Rank:      officer.Rank,
			BadgeNo:   officer.BadgeNo,
			Phone:     officer.Phone,
			Username:  officer.Username,
			Email:     officer.Email,
		})
	}

	return PolicePostResponse{
		ID:        post.ID,
		Name:      post.Name,
		Location:  post.Location,
		Contact:   post.Contact,
		Officers:  officers,
		CreatedAt: post.CreatedAt,
	}
}

// ================================

// CreatePolicePost godoc
//
//	@Summary		Create a new police post record
//	@Description	Creates a new police post entry in the system and returns the created record.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Param			Policepost	body		CreatePolicePostPayload	true	"PolicePost data to create"
//	@Success		201			{object}	fiber.Map				"Successfully created police post record"
//	@Failure		400			{object}	fiber.Map				"Bad request due to invalid input"
//	@Failure		500			{object}	fiber.Map				"Server error when creating police post"
//	@Router			/police-post [post]
func (h *PolicePostController) CreatePolicePost(c *fiber.Ctx) error {
	var payload CreatePolicePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input provided",
			"data":    err.Error(),
		})
	}

	post := &models.PolicePost{
		Name:     payload.Name,
		Location: payload.Location,
		Contact:  payload.Contact,
	}

	if err := h.repo.CreatePolicePost(post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create police post",
			"data":    err.Error(),
		})
	}

	response := ConvertToPolicePostResponse(*post)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Police post created successfully",
		"data":    response,
	})
}

// ===========

// GetAllPolicePosts godoc
//
//	@Summary		Retrieve a paginated list of policePosts
//	@Description	Fetches all policePost records with pagination support.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"policePosts retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve policePosts"
//	@Router			/police-posts [get]
func (h *PolicePostController) GetAllPolicePosts(c *fiber.Ctx) error {
	pagination, posts, err := h.repo.GetPaginatedPolicePosts(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve police posts",
			"data":    err.Error(),
		})
	}

	// Convert to response
	var responses []PolicePostResponse
	for _, post := range posts {
		responses = append(responses, ConvertToPolicePostResponse(post))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police posts retrieved successfully",
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

// GetSinglePolicePost godoc
//
//	@Summary		Retrieve a single policePost record by ID
//	@Description	Fetches a policePost record based on the provided ID.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"PolicePost ID"
//	@Success		200	{object}	fiber.Map	"PolicePost retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"PolicePost not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving PolicePost"
//	@Router			/police-post/{id} [get]
func (h *PolicePostController) GetSinglePolicePost(c *fiber.Ctx) error {
	// Get the policePost ID from the route parameters
	id := c.Params("id")

	// Fetch the policePost by ID
	policePost, err := h.repo.GetPolicePostByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "policePost not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve PolicePost",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "PolicePost and associated data retrieved successfully",
		"data":    ConvertToPolicePostResponse(policePost),
	})
}

// =======================

// Define the UpdatePolicePost struct
type UpdatePolicePostPayload struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

// UpdatePolicePost godoc
//
//	@Summary		Update an existing policePost record by ID
//	@Description	Updates the details of a policePost record based on the provided ID and request body.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"PolicePost ID"
//	@Param			policePost	body		UpdatePolicePostPayload	true	"PolicePost data to update"
//	@Success		200			{object}	fiber.Map				"PolicePost updated successfully"
//	@Failure		400			{object}	fiber.Map				"Invalid input or empty request body"
//	@Failure		404			{object}	fiber.Map				"PolicePost not found"
//	@Failure		500			{object}	fiber.Map				"Server error when updating policePost"
//	@Router			/police-post/{id} [put]
func (h *PolicePostController) UpdatePolicePost(c *fiber.Ctx) error {
	// Get the policePost ID from the route parameters
	id := c.Params("id")

	// Find the policePost in the database
	_, err := h.repo.GetPolicePostByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "PolicePost not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve PolicePost",
			"data":    err.Error(),
		})
	}

	// Parse the request body into the UpdatePolicePostPayload struct
	var payload UpdatePolicePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Check if the request body is empty
	if (UpdatePolicePostPayload{} == payload) {
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

	// Update the PolicePost in the database
	if err := h.repo.UpdatePolicePost(id, updates); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update policePost",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "PolicePost updated successfully",
		"data":    updates,
	})
}

// ==================

// DeletePolicePostByID godoc
//
//	@Summary		Delete a PolicePost record by ID
//	@Description	Deletes a PolicePost record based on the provided ID.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"PolicePost ID"
//	@Success		200	{object}	fiber.Map	"PolicePost deleted successfully"
//	@Failure		404	{object}	fiber.Map	"PolicePost not found"
//	@Failure		500	{object}	fiber.Map	"Server error when deleting policePost"
//	@Router			/police-post/{id} [delete]
func (h *PolicePostController) DeletePolicePostByID(c *fiber.Ctx) error {
	// Get the PolicePost ID from the route parameters
	id := c.Params("id")

	// Find the PolicePost in the database
	policePost, err := h.repo.GetPolicePostByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "PolicePost not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find PolicePost",
			"data":    err.Error(),
		})
	}

	// Delete the PolicePost
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete policePost",
			"data":    err.Error(),
		})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "PolicePost deleted successfully",
		"data":    policePost,
	})
}

// =================

// SearchPolicePosts godoc
//
//	@Summary		Search for policePosts with pagination
//	@Description	Retrieves a paginated list of policePosts based on search criteria.
//	@Tags			Police Posts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"PolicePosts retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve PolicePosts"
//	@Router			/police-posts/search [get]
func (h *PolicePostController) SearchPolicePosts(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, policePosts, err := h.repo.SearchPaginatedPolicePosts(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve policePosts",
			"data":    err.Error(),
		})
	}

	// Convert to response
	var responses []PolicePostResponse
	for _, post := range policePosts {
		responses = append(responses, ConvertToPolicePostResponse(post))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police posts retrieved successfully",
		"data":    responses,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}
