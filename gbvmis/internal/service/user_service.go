package service

import (
	"gbvmis/internal/models"
	"gbvmis/internal/repository"
)

type PoliceOfficerService interface {
	RegisterOfficer(officer *models.PoliceOfficer) error
	GetOfficerProfile(id uint) (*models.PoliceOfficer, error)
	GetOfficerByUsername(username string) (*models.PoliceOfficer, error)
	ListOfficers() ([]models.PoliceOfficer, error)
	UpdateOfficer(officer *models.PoliceOfficer) error
	DeleteOfficer(id uint) error
}

type PoliceOfficerServiceImpl struct {
	Repo repository.PoliceOfficerRepository
}

func (s *PoliceOfficerServiceImpl) RegisterOfficer(officer *models.PoliceOfficer) error {
	return s.Repo.CreateOfficer(officer)
}

func (s *PoliceOfficerServiceImpl) GetOfficerProfile(id uint) (*models.PoliceOfficer, error) {
	return s.Repo.GetOfficerByID(id)
}

func (s *PoliceOfficerServiceImpl) GetOfficerByUsername(username string) (*models.PoliceOfficer, error) {
	return s.Repo.GetOfficerByUsername(username)
}

func (s *PoliceOfficerServiceImpl) ListOfficers() ([]models.PoliceOfficer, error) {
	return s.Repo.ListOfficers()
}

func (s *PoliceOfficerServiceImpl) UpdateOfficer(officer *models.PoliceOfficer) error {
	return s.Repo.UpdateOfficer(officer)
}

func (s *PoliceOfficerServiceImpl) DeleteOfficer(id uint) error {
	return s.Repo.DeleteOfficer(id)
}

type PoliceStationService interface {
	RegisterStation(station *models.PoliceStation) error
	GetStation(id uint) (*models.PoliceStation, error)
	ListStations() ([]models.PoliceStation, error)
	UpdateStation(station *models.PoliceStation) error
	DeleteStation(id uint) error
}

type PoliceStationServiceImpl struct {
	Repo repository.PoliceStationRepository
}

func (s *PoliceStationServiceImpl) RegisterStation(station *models.PoliceStation) error {
	return s.Repo.CreateStation(station)
}

func (s *PoliceStationServiceImpl) GetStation(id uint) (*models.PoliceStation, error) {
	return s.Repo.GetStationByID(id)
}

func (s *PoliceStationServiceImpl) ListStations() ([]models.PoliceStation, error) {
	return s.Repo.ListStations()
}

func (s *PoliceStationServiceImpl) UpdateStation(station *models.PoliceStation) error {
	return s.Repo.UpdateStation(station)
}

func (s *PoliceStationServiceImpl) DeleteStation(id uint) error {
	return s.Repo.DeleteStation(id)
}
