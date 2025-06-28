package controllers

import (
	"errors"
	"gbvmis/internals/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleController struct {
	repo repository.RoleRepository
}

func NewRoleController(repo repository.RoleRepository) *RoleController {
	return &RoleController{repo: repo}
}

// ================================

// GetAllRoles godoc
//
//	@Summary		Retrieve a paginated list of roles
//	@Description	Fetches all role records with pagination support.
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Roles retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve roles"
//	@Router			/roles [get]
func (h *RoleController) GetAllRoles(c *fiber.Ctx) error {
	pagination, roles, err := h.repo.GetPaginatedRoles(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve roles",
			"data":    err.Error(),
		})
	}

	// Return the paginated response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "roles retrieved successfully",
		"data":    roles,
		"pagination": fiber.Map{
			"total_items":  pagination.TotalItems,
			"total_pages":  pagination.TotalPages,
			"current_page": pagination.CurrentPage,
			"limit":        pagination.ItemsPerPage,
		},
	})
}

// =========

// GetSingleRole godoc
//
//	@Summary		Retrieve a single role record by ID
//	@Description	Fetches a role record based on the provided ID.
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Role ID"
//	@Success		200	{object}	fiber.Map	"Role retrieved successfully"
//	@Failure		404	{object}	fiber.Map	"Role not found"
//	@Failure		500	{object}	fiber.Map	"Server error when retrieving role"
//	@Router			/role/{id} [get]
func (h *RoleController) GetSingleRole(c *fiber.Ctx) error {
	// Get the Role ID from the route parameters
	id := c.Params("id")

	// Fetch the role by ID
	role, err := h.repo.GetRoleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Role not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve Role",
			"data":    err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Role and associated data retrieved successfully",
		"data":    role,
	})
}

// =================

// SearchRoles godoc
//
//	@Summary		Search for roles with pagination
//	@Description	Retrieves a paginated list of roles based on search criteria.
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	fiber.Map	"Roles retrieved successfully"
//	@Failure		500	{object}	fiber.Map	"Failed to retrieve roles"
//	@Router			/roles/search [get]
func (h *RoleController) SearchRoles(c *fiber.Ctx) error {
	// Call the repository function to get paginated search results
	pagination, roles, err := h.repo.SearchPaginatedRoles(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve roles",
			"data":    err.Error(),
		})
	}

	// Return the response with pagination details
	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "Roles retrieved successfully",
		"pagination": pagination,
		"data":       roles,
	})
}
