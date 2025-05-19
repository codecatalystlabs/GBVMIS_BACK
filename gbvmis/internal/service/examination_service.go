package service

import (
	"gbvmis/internal/models"
	"gbvmis/internal/repository"
)

type MedicalExaminationService interface {
	RegisterExamination(exam *models.MedicalExamination) error
	GetExamination(id uint) (*models.MedicalExamination, error)
	ListExaminations() ([]models.MedicalExamination, error)
	UpdateExamination(exam *models.MedicalExamination) error
	DeleteExamination(id uint) error
}

type MedicalExaminationServiceImpl struct {
	Repo repository.MedicalExaminationRepository
}

func (s *MedicalExaminationServiceImpl) RegisterExamination(exam *models.MedicalExamination) error {
	return s.Repo.CreateExamination(exam)
}

func (s *MedicalExaminationServiceImpl) GetExamination(id uint) (*models.MedicalExamination, error) {
	return s.Repo.GetExaminationByID(id)
}

func (s *MedicalExaminationServiceImpl) ListExaminations() ([]models.MedicalExamination, error) {
	return s.Repo.ListExaminations()
}

func (s *MedicalExaminationServiceImpl) UpdateExamination(exam *models.MedicalExamination) error {
	return s.Repo.UpdateExamination(exam)
}

func (s *MedicalExaminationServiceImpl) DeleteExamination(id uint) error {
	return s.Repo.DeleteExamination(id)
}
