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

// @Summary Register a new police officer
// @Description Register a new police officer in the system
// @Tags officers
// @Accept json
// @Produce json
// @Param officer body models.PoliceOfficer true "Police officer information"
// @Success 201 {object} models.PoliceOfficer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /officers [post]
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

// @Summary Get officer by ID
// @Description Get police officer details by their ID
// @Tags officers
// @Accept json
// @Produce json
// @Param id path int true "Officer ID"
// @Success 200 {object} models.PoliceOfficer
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /officers/{id} [get]
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

// @Summary List all officers
// @Description Get a list of all registered police officers
// @Tags officers
// @Accept json
// @Produce json
// @Success 200 {array} models.PoliceOfficer
// @Failure 500 {object} map[string]string
// @Router /officers [get]
func (h *PoliceOfficerHandler) ListOfficers(c *fiber.Ctx) error {
	officers, err := h.Service.ListOfficers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list officers"})
	}
	return c.JSON(officers)
}

// @Summary Update officer
// @Description Update police officer information
// @Tags officers
// @Accept json
// @Produce json
// @Param id path int true "Officer ID"
// @Param officer body models.PoliceOfficer true "Updated officer information"
// @Success 200 {object} models.PoliceOfficer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /officers/{id} [put]
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

// @Summary Delete officer
// @Description Delete a police officer record
// @Tags officers
// @Accept json
// @Produce json
// @Param id path int true "Officer ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /officers/{id} [delete]
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

// @Summary Register a new police station
// @Description Register a new police station in the system
// @Tags stations
// @Accept json
// @Produce json
// @Param station body models.PoliceStation true "Police station information"
// @Success 201 {object} models.PoliceStation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /stations [post]
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

// @Summary Get station by ID
// @Description Get police station details by ID
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} models.PoliceStation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /stations/{id} [get]
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

// @Summary List all stations
// @Description Get a list of all police stations
// @Tags stations
// @Accept json
// @Produce json
// @Success 200 {array} models.PoliceStation
// @Failure 500 {object} map[string]string
// @Router /stations [get]
func (h *PoliceStationHandler) ListStations(c *fiber.Ctx) error {
	stations, err := h.Service.ListStations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list stations"})
	}
	return c.JSON(stations)
}

// @Summary Update station
// @Description Update police station information
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param station body models.PoliceStation true "Updated station information"
// @Success 200 {object} models.PoliceStation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /stations/{id} [put]
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

// @Summary Delete station
// @Description Delete a police station record
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /stations/{id} [delete]
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
