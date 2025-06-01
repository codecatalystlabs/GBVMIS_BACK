package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChargeRepository interface {
	CreateCharge(charge *models.Charge) error
	GetPaginatedCharges(c *fiber.Ctx) (*utils.Pagination, []models.Charge, error)
	UpdateCharge(id string, updates map[string]interface{}) error
	GetChargeByID(id string) (models.Charge, error)
	DeleteByID(id string) error
	SearchPaginatedCharges(c *fiber.Ctx) (*utils.Pagination, []models.Charge, error)
}

type ChargeRepositoryImpl struct {
	db *gorm.DB
}

func ChargeDbService(db *gorm.DB) ChargeRepository {
	return &ChargeRepositoryImpl{db: db}
}

// =================================

func (r *ChargeRepositoryImpl) CreateCharge(charge *models.Charge) error {
	return r.db.Create(charge).Error
}

func (r *ChargeRepositoryImpl) GetPaginatedCharges(c *fiber.Ctx) (*utils.Pagination, []models.Charge, error) {
	pagination, charges, err := utils.Paginate(c, r.db, models.Charge{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, charges, nil
}

func (r *ChargeRepositoryImpl) GetChargeByID(id string) (models.Charge, error) {
	var charge models.Charge
	err := r.db.First(&charge, "id = ?", id).Error
	return charge, err
}

func (r *ChargeRepositoryImpl) UpdateCharge(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Charge{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a charge by ID
func (r *ChargeRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Charge{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChargeRepositoryImpl) SearchPaginatedCharges(c *fiber.Ctx) (*utils.Pagination, []models.Charge, error) {
	// Get query parameters from request
	chargetitle := c.Query("chargetitle")
	severity := c.Query("severity")

	// Start building the query
	query := r.db.Model(&models.Charge{})

	// Apply filters based on provided parameters
	if chargetitle != "" {
		query = query.Where("charge_title ILIKE ?", "%"+chargetitle+"%")
	}
	if severity != "" {
		query = query.Where("severity ILIKE ?", "%"+severity+"%")
	}

	// Call the pagination helper
	pagination, charges, err := utils.Paginate(c, query, models.Charge{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, charges, nil
}
