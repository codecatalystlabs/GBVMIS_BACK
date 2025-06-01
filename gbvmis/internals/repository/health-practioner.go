package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthPractitionerRepository interface {
	CreateHealthPractitioner(healthPractitioner *models.HealthPractitioner) error
	GetPaginatedHealthPractitioners(c *fiber.Ctx) (*utils.Pagination, []models.HealthPractitioner, error)
	UpdateHealthPractitioner(id string, updates map[string]interface{}) error
	GetHealthPractitionerByID(id string) (models.HealthPractitioner, error)
	DeleteByID(id string) error
	SearchPaginatedHealthPractitioners(c *fiber.Ctx) (*utils.Pagination, []models.HealthPractitioner, error)
}

type HealthPractitionerRepositoryImpl struct {
	db *gorm.DB
}

func HealthPractitionerDbService(db *gorm.DB) HealthPractitionerRepository {
	return &HealthPractitionerRepositoryImpl{db: db}
}

// =================================

func (r *HealthPractitionerRepositoryImpl) CreateHealthPractitioner(healthPractitioner *models.HealthPractitioner) error {
	return r.db.Create(healthPractitioner).Error
}

func (r *HealthPractitionerRepositoryImpl) GetPaginatedHealthPractitioners(c *fiber.Ctx) (*utils.Pagination, []models.HealthPractitioner, error) {
	pagination, healthPractitioners, err := utils.Paginate(c, r.db, models.HealthPractitioner{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, healthPractitioners, nil
}

func (r *HealthPractitionerRepositoryImpl) GetHealthPractitionerByID(id string) (models.HealthPractitioner, error) {
	var healthPractitioner models.HealthPractitioner
	err := r.db.First(&healthPractitioner, "id = ?", id).Error
	return healthPractitioner, err
}

func (r *HealthPractitionerRepositoryImpl) UpdateHealthPractitioner(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.HealthPractitioner{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a healthPractitioner by ID
func (r *HealthPractitionerRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.HealthPractitioner{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *HealthPractitionerRepositoryImpl) SearchPaginatedHealthPractitioners(c *fiber.Ctx) (*utils.Pagination, []models.HealthPractitioner, error) {

	// Get query parameters from request
	FirstName := c.Query("first_name")
	LastName := c.Query("last_name")
	Profession := c.Query("profession")
	Gender := c.Query("gender")
	FacilityID := c.Query("facility_id")

	// Start building the query
	query := r.db.Model(&models.HealthPractitioner{})

	// Apply filters based on provided parameters
	if FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+FirstName+"%")
	}
	if LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+LastName+"%")
	}
	if Profession != "" {
		query = query.Where("profession ILIKE ?", "%"+Profession+"%")
	}
	if Gender != "" {
		query = query.Where("gender ILIKE ?", "%"+Gender+"%")
	}

	if FacilityID != "" {
		if _, err := strconv.Atoi(FacilityID); err == nil {
			query = query.Where("post_id = ?", FacilityID)
		}
	}

	// Call the pagination helper
	pagination, healthPractitioners, err := utils.Paginate(c, query, models.HealthPractitioner{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, healthPractitioners, nil
}
