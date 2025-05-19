package service

import (
	"gbvmis/internal/models"
	"gbvmis/internal/repository"
)

type VictimService interface {
	RegisterVictim(victim *models.Victim) error
	GetVictimProfile(id uint) (*models.Victim, error)
	ListVictims() ([]models.Victim, error)
	UpdateVictim(victim *models.Victim) error
	DeleteVictim(id uint) error
}

type VictimServiceImpl struct {
	Repo repository.VictimRepository
}

func (s *VictimServiceImpl) RegisterVictim(victim *models.Victim) error {
	return s.Repo.CreateVictim(victim)
}

func (s *VictimServiceImpl) GetVictimProfile(id uint) (*models.Victim, error) {
	return s.Repo.GetVictimByID(id)
}

func (s *VictimServiceImpl) ListVictims() ([]models.Victim, error) {
	return s.Repo.ListVictims()
}

func (s *VictimServiceImpl) UpdateVictim(victim *models.Victim) error {
	return s.Repo.UpdateVictim(victim)
}

func (s *VictimServiceImpl) DeleteVictim(id uint) error {
	return s.Repo.DeleteVictim(id)
}
