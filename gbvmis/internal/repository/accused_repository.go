package repository

import (
	"gbvmis/internal/models"

	"gorm.io/gorm"
)

type AccusedRepository interface {
	CreateAccused(accused *models.Accused) error
	GetAccusedByID(id uint) (*models.Accused, error)
	ListAccused() ([]models.Accused, error)
	UpdateAccused(accused *models.Accused) error
	DeleteAccused(id uint) error
}

type GormAccusedRepository struct {
	DB *gorm.DB
}

func (r *GormAccusedRepository) CreateAccused(accused *models.Accused) error {
	return r.DB.Create(accused).Error
}

func (r *GormAccusedRepository) GetAccusedByID(id uint) (*models.Accused, error) {
	var accused models.Accused
	if err := r.DB.First(&accused, id).Error; err != nil {
		return nil, err
	}
	return &accused, nil
}

func (r *GormAccusedRepository) ListAccused() ([]models.Accused, error) {
	var accused []models.Accused
	if err := r.DB.Find(&accused).Error; err != nil {
		return nil, err
	}
	return accused, nil
}

func (r *GormAccusedRepository) UpdateAccused(accused *models.Accused) error {
	return r.DB.Save(accused).Error
}

func (r *GormAccusedRepository) DeleteAccused(id uint) error {
	return r.DB.Delete(&models.Accused{}, id).Error
}
