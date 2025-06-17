package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ExaminationRepository interface {
	CreateExamination(examination *models.Examination) error
	GetPaginatedExaminations(c *fiber.Ctx) (*utils.Pagination, []models.Examination, error)
	UpdateExamination(id string, updates map[string]interface{}) error
	GetExaminationByID(id string) (models.Examination, error)
	DeleteByID(id string) error
	SearchPaginatedExaminations(c *fiber.Ctx) (*utils.Pagination, []models.Examination, error)
}

type ExaminationRepositoryImpl struct {
	db *gorm.DB
}

func ExaminationDbService(db *gorm.DB) ExaminationRepository {
	return &ExaminationRepositoryImpl{db: db}
}

// =================================

func (r *ExaminationRepositoryImpl) CreateExamination(examination *models.Examination) error {
	return r.db.Create(examination).Error
}

func (r *ExaminationRepositoryImpl) GetPaginatedExaminations(c *fiber.Ctx) (*utils.Pagination, []models.Examination, error) {
	pagination, examinations, err := utils.Paginate(c, r.db.
		Preload("Victim").
		Preload("Case").
		Preload("Facility").
		Preload("Practitioner"), models.Examination{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, examinations, nil
}

func (r *ExaminationRepositoryImpl) GetExaminationByID(id string) (models.Examination, error) {
	var examination models.Examination
	err := r.db.
		Preload("Victim").
		Preload("Case").
		Preload("Facility").
		Preload("Practitioner").First(&examination, "id = ?", id).Error
	return examination, err
}

func (r *ExaminationRepositoryImpl) UpdateExamination(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Examination{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a Examination by ID
func (r *ExaminationRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Examination{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ExaminationRepositoryImpl) SearchPaginatedExaminations(c *fiber.Ctx) (*utils.Pagination, []models.Examination, error) {
	// Get query parameters from request
	FacilityID := c.Query("facility_id")
	PractitionerID := c.Query("practitioner_id")

	// Start building the query
	query := r.db.
		Preload("Victim").
		Preload("Case").
		Preload("Facility").
		Preload("Practitioner").Model(&models.Examination{})

	// Apply filters based on provided parameters
	if FacilityID != "" {
		if _, err := strconv.Atoi(FacilityID); err == nil {
			query = query.Where("facility_id = ?", FacilityID)
		}
	}
	if PractitionerID != "" {
		if _, err := strconv.Atoi(PractitionerID); err == nil {
			query = query.Where("practitioner_id = ?", PractitionerID)
		}
	}

	// Call the pagination helper
	pagination, examinations, err := utils.Paginate(c, query, models.Examination{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, examinations, nil
}
