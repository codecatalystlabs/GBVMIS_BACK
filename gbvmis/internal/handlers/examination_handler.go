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

func (h *ExaminationHandler) ListExaminations(c *fiber.Ctx) error {
	exams, err := h.Service.ListExaminations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not list examinations"})
	}
	return c.JSON(exams)
}

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
