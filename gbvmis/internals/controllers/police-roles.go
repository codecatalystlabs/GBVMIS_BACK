package controllers

import (
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PoliceRolesController struct {
	repo repository.RoleRepository
}

func NewPoliceRolesController(repo repository.RoleRepository) *PoliceRolesController {
	return &PoliceRolesController{repo: repo}
}

type PoliceRolesResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ConvertToPoliceRolesResponse(role models.Role) PoliceRolesResponse {
	return PoliceRolesResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func (h *PoliceRolesController) CreatePoliceRoles(c *fiber.Ctx) error {
	payload := new(PoliceRolesResponse)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}

	role := models.Role{
		Name: payload.Name,
	}

	if err := h.repo.CreateRole(&role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create police role",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Police role created successfully",
		"data":    ConvertToPoliceRolesResponse(role),
	})
}

func (h *PoliceRolesController) GetAllPoliceRoles(c *fiber.Ctx) error {
	_, roles, err := h.repo.GetPaginatedRoles(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve police roles",
			"data":    err.Error(),
		})
	}
	roleResponses := make([]PoliceRolesResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = ConvertToPoliceRolesResponse(role)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police roles retrieved successfully",
		"data":    roleResponses,
	})
}

func (h *PoliceRolesController) GetSinglePoliceRole(c *fiber.Ctx) error {
	id := c.Params("id")
	role, err := h.repo.GetRoleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Police role not found",
			"data":    err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police role retrieved successfully",
		"data":    ConvertToPoliceRolesResponse(role),
	})
}

func (h *PoliceRolesController) UpdatePoliceRole(c *fiber.Ctx) error {
	id := c.Params("id")
	payload := new(PoliceRolesResponse)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err.Error(),
		})
	}
	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if err := h.repo.UpdateRole(id, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update police role",
			"data":    err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police role updated successfully",
	})
}

func (h *PoliceRolesController) DeletePoliceRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.DeleteByID(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete police role",
			"data":    err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police role deleted successfully",
	})
}

func (h *PoliceRolesController) SearchPoliceRoles(c *fiber.Ctx) error {
	_, roles, err := h.repo.SearchPaginatedRoles(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to search police roles",
			"data":    err.Error(),
		})
	}
	roleResponses := make([]PoliceRolesResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = ConvertToPoliceRolesResponse(role)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Police roles searched successfully",
		"data":    roleResponses,
	})
}
