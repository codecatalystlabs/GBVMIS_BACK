package repository

import (
	"gbvmis/internal/models"

	"gorm.io/gorm"
)

type PoliceOfficerRepository interface {
	CreateOfficer(officer *models.PoliceOfficer) error
	GetOfficerByID(id uint) (*models.PoliceOfficer, error)
	GetOfficerByUsername(username string) (*models.PoliceOfficer, error)
	ListOfficers() ([]models.PoliceOfficer, error)
	UpdateOfficer(officer *models.PoliceOfficer) error
	DeleteOfficer(id uint) error
}

type GormPoliceOfficerRepository struct {
	DB *gorm.DB
}

func (r *GormPoliceOfficerRepository) CreateOfficer(officer *models.PoliceOfficer) error {
	return r.DB.Create(officer).Error
}

func (r *GormPoliceOfficerRepository) GetOfficerByID(id uint) (*models.PoliceOfficer, error) {
	var officer models.PoliceOfficer
	if err := r.DB.First(&officer, id).Error; err != nil {
		return nil, err
	}
	return &officer, nil
}

func (r *GormPoliceOfficerRepository) GetOfficerByUsername(username string) (*models.PoliceOfficer, error) {
	var officer models.PoliceOfficer
	if err := r.DB.Where("username = ?", username).First(&officer).Error; err != nil {
		return nil, err
	}
	return &officer, nil
}

func (r *GormPoliceOfficerRepository) ListOfficers() ([]models.PoliceOfficer, error) {
	var officers []models.PoliceOfficer
	if err := r.DB.Find(&officers).Error; err != nil {
		return nil, err
	}
	return officers, nil
}

func (r *GormPoliceOfficerRepository) UpdateOfficer(officer *models.PoliceOfficer) error {
	return r.DB.Save(officer).Error
}

func (r *GormPoliceOfficerRepository) DeleteOfficer(id uint) error {
	return r.DB.Delete(&models.PoliceOfficer{}, id).Error
}

type PoliceStationRepository interface {
	CreateStation(station *models.PoliceStation) error
	GetStationByID(id uint) (*models.PoliceStation, error)
	ListStations() ([]models.PoliceStation, error)
	UpdateStation(station *models.PoliceStation) error
	DeleteStation(id uint) error
}

type GormPoliceStationRepository struct {
	DB *gorm.DB
}

func (r *GormPoliceStationRepository) CreateStation(station *models.PoliceStation) error {
	return r.DB.Create(station).Error
}

func (r *GormPoliceStationRepository) GetStationByID(id uint) (*models.PoliceStation, error) {
	var station models.PoliceStation
	if err := r.DB.First(&station, id).Error; err != nil {
		return nil, err
	}
	return &station, nil
}

func (r *GormPoliceStationRepository) ListStations() ([]models.PoliceStation, error) {
	var stations []models.PoliceStation
	if err := r.DB.Find(&stations).Error; err != nil {
		return nil, err
	}
	return stations, nil
}

func (r *GormPoliceStationRepository) UpdateStation(station *models.PoliceStation) error {
	return r.DB.Save(station).Error
}

func (r *GormPoliceStationRepository) DeleteStation(id uint) error {
	return r.DB.Delete(&models.PoliceStation{}, id).Error
}
