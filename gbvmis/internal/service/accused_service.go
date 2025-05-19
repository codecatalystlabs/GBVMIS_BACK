package service

import (
	"gbvmis/internal/models"
	"gbvmis/internal/repository"
)

type AccusedService interface {
	RegisterAccused(accused *models.Accused) error
	GetAccusedProfile(id uint) (*models.Accused, error)
	ListAccused() ([]models.Accused, error)
	UpdateAccused(accused *models.Accused) error
	DeleteAccused(id uint) error
}

type AccusedServiceImpl struct {
	Repo repository.AccusedRepository
}

func (s *AccusedServiceImpl) RegisterAccused(accused *models.Accused) error {
	return s.Repo.CreateAccused(accused)
}

func (s *AccusedServiceImpl) GetAccusedProfile(id uint) (*models.Accused, error) {
	return s.Repo.GetAccusedByID(id)
}

func (s *AccusedServiceImpl) ListAccused() ([]models.Accused, error) {
	return s.Repo.ListAccused()
}

func (s *AccusedServiceImpl) UpdateAccused(accused *models.Accused) error {
	return s.Repo.UpdateAccused(accused)
}

func (s *AccusedServiceImpl) DeleteAccused(id uint) error {
	return s.Repo.DeleteAccused(id)
}
