package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VictimRepository interface {
	CreateVictim(victim *models.Victim) error
	GetPaginatedVictims(c *fiber.Ctx) (*utils.Pagination, []models.Victim, error)
	UpdateVictim(id string, updates map[string]interface{}) error
	GetVictimByID(id string) (models.Victim, error)
	DeleteByID(id string) error
	SearchPaginatedVictims(c *fiber.Ctx) (*utils.Pagination, []models.Victim, error)
}

type VictimRepositoryImpl struct {
	db *gorm.DB
}

func VictimDbService(db *gorm.DB) VictimRepository {
	return &VictimRepositoryImpl{db: db}
}

// =================================

func (r *VictimRepositoryImpl) CreateVictim(victim *models.Victim) error {
	return r.db.Create(victim).Error
}

func (r *VictimRepositoryImpl) GetPaginatedVictims(c *fiber.Ctx) (*utils.Pagination, []models.Victim, error) {
	pagination, victims, err := utils.Paginate(c, r.db, models.Victim{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, victims, nil
}

func (r *VictimRepositoryImpl) GetVictimByID(id string) (models.Victim, error) {
	var victim models.Victim
	err := r.db.First(&victim, "id = ?", id).Error
	return victim, err
}

func (r *VictimRepositoryImpl) UpdateVictim(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Victim{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a victim by ID
func (r *VictimRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Victim{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *VictimRepositoryImpl) SearchPaginatedVictims(c *fiber.Ctx) (*utils.Pagination, []models.Victim, error) {
	// Get query parameters from request
	lastname := c.Query("lastname")
	firstname := c.Query("firstname")
	gender := c.Query("gender")
	nationality := c.Query("nationality")
	nin := c.Query("nin")

	// Start building the query
	query := r.db.Model(&models.Victim{})

	// Apply filters based on provided parameters
	if lastname != "" {
		query = query.Where("last_name ILIKE ?", "%"+lastname+"%")
	}
	if firstname != "" {
		query = query.Where("first_name ILIKE ?", "%"+firstname+"%")
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if nationality != "" {
		query = query.Where("nationality = ?", nationality)
	}
	if nin != "" {
		query = query.Where("nin = ?", nin)
	}

	// Call the pagination helper
	pagination, victims, err := utils.Paginate(c, query, models.Victim{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, victims, nil
}
