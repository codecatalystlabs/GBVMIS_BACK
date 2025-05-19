package repository

import (
	"gbvmis/internal/models"

	"gorm.io/gorm"
)

type MedicalExaminationRepository interface {
	CreateExamination(exam *models.MedicalExamination) error
	GetExaminationByID(id uint) (*models.MedicalExamination, error)
	ListExaminations() ([]models.MedicalExamination, error)
	UpdateExamination(exam *models.MedicalExamination) error
	DeleteExamination(id uint) error
}

type GormMedicalExaminationRepository struct {
	DB *gorm.DB
}

func (r *GormMedicalExaminationRepository) CreateExamination(exam *models.MedicalExamination) error {
	return r.DB.Create(exam).Error
}

func (r *GormMedicalExaminationRepository) GetExaminationByID(id uint) (*models.MedicalExamination, error) {
	var exam models.MedicalExamination
	if err := r.DB.First(&exam, id).Error; err != nil {
		return nil, err
	}
	return &exam, nil
}

func (r *GormMedicalExaminationRepository) ListExaminations() ([]models.MedicalExamination, error) {
	var exams []models.MedicalExamination
	if err := r.DB.Find(&exams).Error; err != nil {
		return nil, err
	}
	return exams, nil
}

func (r *GormMedicalExaminationRepository) UpdateExamination(exam *models.MedicalExamination) error {
	return r.DB.Save(exam).Error
}

func (r *GormMedicalExaminationRepository) DeleteExamination(id uint) error {
	return r.DB.Delete(&models.MedicalExamination{}, id).Error
}
