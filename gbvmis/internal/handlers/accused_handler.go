package handlers

import (
	"gbvmis/internal/models"
	"gbvmis/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AccusedHandler struct {
	Service service.AccusedService
}

// @Summary Register a new accused
// @Description Register a new accused person in the system
// @Tags accused
// @Accept json
// @Produce json
// @Param accused body models.Accused true "Accused information"
// @Success 201 {object} models.Accused
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accused [post]
func (h *AccusedHandler) RegisterAccused(c *fiber.Ctx) error {
	var accused models.Accused
	if err := c.BodyParser(&accused); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.Service.RegisterAccused(&accused); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register accused"})
	}
	return c.Status(fiber.StatusCreated).JSON(accused)
}

// @Summary Get accused by ID
// @Description Get accused details by their ID
// @Tags accused
// @Accept json
// @Produce json
// @Param id path int true "Accused ID"
// @Success 200 {object} models.Accused
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /accused/{id} [get]
func (h *AccusedHandler) GetAccusedByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	accused, err := h.Service.GetAccusedProfile(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Accused not found"})
	}
	return c.JSON(accused)
}

// @Summary List all accused
// @Description Get a list of all registered accused persons
// @Tags accused
// @Accept json
// @Produce json
// @Success 200 {array} models.Accused
// @Failure 500 {object} map[string]string
// @Router /accused [get]
func (h *AccusedHandler) ListAccused(c *fiber.Ctx) error {
	accused, err := h.Service.ListAccused()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list accused"})
	}
	return c.JSON(accused)
}

// @Summary Update accused
// @Description Update accused information
// @Tags accused
// @Accept json
// @Produce json
// @Param id path int true "Accused ID"
// @Param accused body models.Accused true "Updated accused information"
// @Success 200 {object} models.Accused
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accused/{id} [put]
func (h *AccusedHandler) UpdateAccused(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var accused models.Accused
	if err := c.BodyParser(&accused); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	accused.ID = uint(id)
	if err := h.Service.UpdateAccused(&accused); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update accused"})
	}
	return c.JSON(accused)
}

// @Summary Delete accused
// @Description Delete an accused record
// @Tags accused
// @Accept json
// @Produce json
// @Param id path int true "Accused ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accused/{id} [delete]
func (h *AccusedHandler) DeleteAccused(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteAccused(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete accused"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
