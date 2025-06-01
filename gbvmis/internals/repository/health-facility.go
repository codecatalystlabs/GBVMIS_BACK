package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthFacilityRepository interface {
	CreateHealthFacility(healthFacility *models.HealthFacility) error
	GetPaginatedHealthFacilities(c *fiber.Ctx) (*utils.Pagination, []models.HealthFacility, error)
	UpdateHealthFacility(id string, updates map[string]interface{}) error
	GetHealthFacilityByID(id string) (models.HealthFacility, error)
	DeleteByID(id string) error
	SearchPaginatedHealthFacilities(c *fiber.Ctx) (*utils.Pagination, []models.HealthFacility, error)
}

type HealthFacilityRepositoryImpl struct {
	db *gorm.DB
}

func HealthFacilityDbService(db *gorm.DB) HealthFacilityRepository {
	return &HealthFacilityRepositoryImpl{db: db}
}

// =================================

func (r *HealthFacilityRepositoryImpl) CreateHealthFacility(healthFacility *models.HealthFacility) error {
	return r.db.Create(healthFacility).Error
}

func (r *HealthFacilityRepositoryImpl) GetPaginatedHealthFacilities(c *fiber.Ctx) (*utils.Pagination, []models.HealthFacility, error) {
	pagination, healthFacilities, err := utils.Paginate(c, r.db, models.HealthFacility{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, healthFacilities, nil
}

func (r *HealthFacilityRepositoryImpl) GetHealthFacilityByID(id string) (models.HealthFacility, error) {
	var healthFacility models.HealthFacility
	err := r.db.First(&healthFacility, "id = ?", id).Error
	return healthFacility, err
}

func (r *HealthFacilityRepositoryImpl) UpdateHealthFacility(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.HealthFacility{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a healthFacility by ID
func (r *HealthFacilityRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.HealthFacility{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *HealthFacilityRepositoryImpl) SearchPaginatedHealthFacilities(c *fiber.Ctx) (*utils.Pagination, []models.HealthFacility, error) {
	// Get query parameters from request
	Name := c.Query("chargetitle")
	Location := c.Query("severity")

	// Start building the query
	query := r.db.Model(&models.HealthFacility{})

	// Apply filters based on provided parameters
	if Name != "" {
		query = query.Where("name ILIKE ?", "%"+Name+"%")
	}
	if Location != "" {
		query = query.Where("location ILIKE ?", "%"+Location+"%")
	}

	// Call the pagination helper
	pagination, healthFacilities, err := utils.Paginate(c, query, models.HealthFacility{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, healthFacilities, nil
}
