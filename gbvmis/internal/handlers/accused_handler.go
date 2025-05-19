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

func (h *AccusedHandler) ListAccused(c *fiber.Ctx) error {
	accused, err := h.Service.ListAccused()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list accused"})
	}
	return c.JSON(accused)
}

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
