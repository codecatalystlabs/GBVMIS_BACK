package service

import (
	"gbvmis/internals/models"
	"gbvmis/internals/repository"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
)

type ToxicologyService interface {
	Create(report *models.ToxicologyForensicReport) error
	GetAll() ([]models.ToxicologyForensicReport, error)
	GetByID(id uint) (*models.ToxicologyForensicReport, error)
	Update(report *models.ToxicologyForensicReport) error
	Delete(id uint) error

	GetPaginatedReports(c *fiber.Ctx) (*utils.Pagination, []models.ToxicologyForensicReport, error)
}

type toxicologyService struct {
	repo repository.ToxicologyReportRepository
}

func NewToxicologyService(repo repository.ToxicologyReportRepository) ToxicologyService {
	return &toxicologyService{repo: repo}
}

func (s *toxicologyService) Create(report *models.ToxicologyForensicReport) error {
	return s.repo.Create(report)
}

func (s *toxicologyService) GetAll() ([]models.ToxicologyForensicReport, error) {
	return s.repo.GetAll()
}

func (s *toxicologyService) GetPaginatedReports(c *fiber.Ctx) (*utils.Pagination, []models.ToxicologyForensicReport, error) {
	pagination, reports, _ := s.repo.GetPaginatedReports(c)
	return pagination, reports, nil
}

func (s *toxicologyService) GetByID(id uint) (*models.ToxicologyForensicReport, error) {
	return s.repo.GetByID(id)
}

func (s *toxicologyService) Update(report *models.ToxicologyForensicReport) error {
	return s.repo.Update(report)
}

func (s *toxicologyService) Delete(id uint) error {
	return s.repo.Delete(id)
}
