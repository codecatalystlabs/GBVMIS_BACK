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

// @Summary Register a new victim
// @Description Register a new victim in the system
// @Tags victims
// @Accept json
// @Produce json
// @Param victim body models.Victim true "Victim information"
// @Success 201 {object} models.Victim
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /victims [post]
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

// @Summary Get victim by ID
// @Description Get victim details by their ID
// @Tags victims
// @Accept json
// @Produce json
// @Param id path int true "Victim ID"
// @Success 200 {object} models.Victim
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /victims/{id} [get]
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

// @Summary List all victims
// @Description Get a list of all registered victims
// @Tags victims
// @Accept json
// @Produce json
// @Success 200 {array} models.Victim
// @Failure 500 {object} map[string]string
// @Router /victims [get]
func (h *VictimHandler) ListVictims(c *fiber.Ctx) error {
	victims, err := h.Service.ListVictims()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list victims"})
	}
	return c.JSON(victims)
}

// @Summary Update victim
// @Description Update victim information
// @Tags victims
// @Accept json
// @Produce json
// @Param id path int true "Victim ID"
// @Param victim body models.Victim true "Updated victim information"
// @Success 200 {object} models.Victim
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /victims/{id} [put]
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

// @Summary Delete victim
// @Description Delete a victim record
// @Tags victims
// @Accept json
// @Produce json
// @Param id path int true "Victim ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /victims/{id} [delete]
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
