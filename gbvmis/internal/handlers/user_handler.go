package handlers

import (
	"gbvmis/internal/models"
	"gbvmis/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PoliceOfficerHandler struct {
	Service service.PoliceOfficerService
}

type PoliceStationHandler struct {
	Service service.PoliceStationService
}

// Police Officer Handlers
func (h *PoliceOfficerHandler) RegisterOfficer(c *fiber.Ctx) error {
	var officer models.PoliceOfficer
	if err := c.BodyParser(&officer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.Service.RegisterOfficer(&officer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register officer"})
	}
	return c.Status(fiber.StatusCreated).JSON(officer)
}

func (h *PoliceOfficerHandler) GetOfficerByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	officer, err := h.Service.GetOfficerProfile(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Officer not found"})
	}
	return c.JSON(officer)
}

func (h *PoliceOfficerHandler) ListOfficers(c *fiber.Ctx) error {
	officers, err := h.Service.ListOfficers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list officers"})
	}
	return c.JSON(officers)
}

func (h *PoliceOfficerHandler) UpdateOfficer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var officer models.PoliceOfficer
	if err := c.BodyParser(&officer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	officer.ID = uint(id)
	if err := h.Service.UpdateOfficer(&officer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update officer"})
	}
	return c.JSON(officer)
}

func (h *PoliceOfficerHandler) DeleteOfficer(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteOfficer(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete officer"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Police Station Handlers
func (h *PoliceStationHandler) RegisterStation(c *fiber.Ctx) error {
	var station models.PoliceStation
	if err := c.BodyParser(&station); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.Service.RegisterStation(&station); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register station"})
	}
	return c.Status(fiber.StatusCreated).JSON(station)
}

func (h *PoliceStationHandler) GetStationByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	station, err := h.Service.GetStation(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Station not found"})
	}
	return c.JSON(station)
}

func (h *PoliceStationHandler) ListStations(c *fiber.Ctx) error {
	stations, err := h.Service.ListStations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list stations"})
	}
	return c.JSON(stations)
}

func (h *PoliceStationHandler) UpdateStation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var station models.PoliceStation
	if err := c.BodyParser(&station); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	station.ID = uint(id)
	if err := h.Service.UpdateStation(&station); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update station"})
	}
	return c.JSON(station)
}

func (h *PoliceStationHandler) DeleteStation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteStation(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete station"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
