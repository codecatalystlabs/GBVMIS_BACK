package handlers

import (
	"gbvmis/internal/models"
	"gbvmis/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VictimHandler struct {
	Service service.VictimService
}

func (h *VictimHandler) RegisterVictim(c *fiber.Ctx) error {
	var victim models.Victim
	if err := c.BodyParser(&victim); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.Service.RegisterVictim(&victim); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register victim"})
	}
	return c.Status(fiber.StatusCreated).JSON(victim)
}

func (h *VictimHandler) GetVictimByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	victim, err := h.Service.GetVictimProfile(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Victim not found"})
	}
	return c.JSON(victim)
}

func (h *VictimHandler) ListVictims(c *fiber.Ctx) error {
	victims, err := h.Service.ListVictims()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list victims"})
	}
	return c.JSON(victims)
}

func (h *VictimHandler) UpdateVictim(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var victim models.Victim
	if err := c.BodyParser(&victim); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	victim.ID = uint(id)
	if err := h.Service.UpdateVictim(&victim); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update victim"})
	}
	return c.JSON(victim)
}

func (h *VictimHandler) DeleteVictim(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteVictim(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete victim"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
