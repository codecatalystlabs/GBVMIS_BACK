package handlers

import (
	"gbvmis/internal/models"
	"gbvmis/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ExaminationHandler struct {
	Service service.MedicalExaminationService
}

// @Summary Register a new medical examination
// @Description Register a new medical examination record
// @Tags examinations
// @Accept json
// @Produce json
// @Param examination body models.MedicalExamination true "Medical examination information"
// @Success 201 {object} models.MedicalExamination
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /examinations [post]
func (h *ExaminationHandler) RegisterExamination(c *fiber.Ctx) error {
	var exam models.MedicalExamination
	if err := c.BodyParser(&exam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.Service.RegisterExamination(&exam); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register examination"})
	}
	return c.Status(fiber.StatusCreated).JSON(exam)
}

// @Summary Get examination by ID
// @Description Get medical examination details by ID
// @Tags examinations
// @Accept json
// @Produce json
// @Param id path int true "Examination ID"
// @Success 200 {object} models.MedicalExamination
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /examinations/{id} [get]
func (h *ExaminationHandler) GetExaminationByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	exam, err := h.Service.GetExamination(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Examination not found"})
	}
	return c.JSON(exam)
}

// @Summary List all examinations
// @Description Get a list of all medical examinations
// @Tags examinations
// @Accept json
// @Produce json
// @Success 200 {array} models.MedicalExamination
// @Failure 500 {object} map[string]string
// @Router /examinations [get]
func (h *ExaminationHandler) ListExaminations(c *fiber.Ctx) error {
	exams, err := h.Service.ListExaminations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list examinations"})
	}
	return c.JSON(exams)
}

// @Summary Update examination
// @Description Update medical examination information
// @Tags examinations
// @Accept json
// @Produce json
// @Param id path int true "Examination ID"
// @Param examination body models.MedicalExamination true "Updated examination information"
// @Success 200 {object} models.MedicalExamination
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /examinations/{id} [put]
func (h *ExaminationHandler) UpdateExamination(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var exam models.MedicalExamination
	if err := c.BodyParser(&exam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	exam.ID = uint(id)
	if err := h.Service.UpdateExamination(&exam); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update examination"})
	}
	return c.JSON(exam)
}

// @Summary Delete examination
// @Description Delete a medical examination record
// @Tags examinations
// @Accept json
// @Produce json
// @Param id path int true "Examination ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /examinations/{id} [delete]
func (h *ExaminationHandler) DeleteExamination(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := h.Service.DeleteExamination(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete examination"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
