package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CaseRepository interface {
	CreateCase(casee *models.Case) error
	GetPaginatedCases(c *fiber.Ctx) (*utils.Pagination, []models.Case, error)
	UpdateCase(id string, updates map[string]interface{}) error
	GetCaseByID(id string) (models.Case, error)
	DeleteByID(id string) error
	SearchPaginatedCases(c *fiber.Ctx) (*utils.Pagination, []models.Case, error)
	FindVictimsByIDs(ids []uint, victims *[]models.Victim) error
	BeginTransaction() *gorm.DB
}

type CaseRepositoryImpl struct {
	db *gorm.DB
}

func CaseDbService(db *gorm.DB) CaseRepository {
	return &CaseRepositoryImpl{db: db}
}

// =================================

func (r *CaseRepositoryImpl) CreateCase(casee *models.Case) error {
	return r.db.Create(casee).Error
}

func (r *CaseRepositoryImpl) GetPaginatedCases(c *fiber.Ctx) (*utils.Pagination, []models.Case, error) {
	pagination, cases, err := utils.Paginate(c, r.db, models.Case{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, cases, nil
}

func (r *CaseRepositoryImpl) GetCaseByID(id string) (models.Case, error) {
	var casee models.Case
	err := r.db.First(&casee, "id = ?", id).Error
	return casee, err
}

func (r *CaseRepositoryImpl) UpdateCase(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Case{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a case by ID
func (r *CaseRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Case{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *CaseRepositoryImpl) SearchPaginatedCases(c *fiber.Ctx) (*utils.Pagination, []models.Case, error) {
	// Get query parameters from request
	CaseNumber := c.Query("case_number")
	Title := c.Query("title")
	Status := c.Query("status")
	PolicePostID := c.Query("police_post_id")

	// Start building the query
	query := r.db.Model(&models.Charge{})

	// Apply filters based on provided parameters
	if CaseNumber != "" {
		query = query.Where("charge_title ILIKE ?", "%"+CaseNumber+"%")
	}
	if Title != "" {
		query = query.Where("title ILIKE ?", "%"+Title+"%")
	}
	if Status != "" {
		query = query.Where("charge_title ILIKE ?", "%"+Status+"%")
	}
	if PolicePostID != "" {
		if _, err := strconv.Atoi(PolicePostID); err == nil {
			query = query.Where("police_post_id = ?", PolicePostID)
		}
	}

	// Call the pagination helper
	pagination, cases, err := utils.Paginate(c, query, models.Case{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, cases, nil
}

func (r *CaseRepositoryImpl) FindVictimsByIDs(ids []uint, victims *[]models.Victim) error {
	return r.db.Where("id IN ?", ids).Find(victims).Error
}

func (r *CaseRepositoryImpl) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
