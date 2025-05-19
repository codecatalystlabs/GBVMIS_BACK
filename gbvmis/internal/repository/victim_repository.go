package repository

import (
	"gbvmis/internal/models"

	"gorm.io/gorm"
)

type VictimRepository interface {
	CreateVictim(victim *models.Victim) error
	GetVictimByID(id uint) (*models.Victim, error)
	ListVictims() ([]models.Victim, error)
	UpdateVictim(victim *models.Victim) error
	DeleteVictim(id uint) error
}

type GormVictimRepository struct {
	DB *gorm.DB
}

func (r *GormVictimRepository) CreateVictim(victim *models.Victim) error {
	return r.DB.Create(victim).Error
}

func (r *GormVictimRepository) GetVictimByID(id uint) (*models.Victim, error) {
	var victim models.Victim
	if err := r.DB.First(&victim, id).Error; err != nil {
		return nil, err
	}
	return &victim, nil
}

func (r *GormVictimRepository) ListVictims() ([]models.Victim, error) {
	var victims []models.Victim
	if err := r.DB.Find(&victims).Error; err != nil {
		return nil, err
	}
	return victims, nil
}

func (r *GormVictimRepository) UpdateVictim(victim *models.Victim) error {
	return r.DB.Save(victim).Error
}

func (r *GormVictimRepository) DeleteVictim(id uint) error {
	return r.DB.Delete(&models.Victim{}, id).Error
}
