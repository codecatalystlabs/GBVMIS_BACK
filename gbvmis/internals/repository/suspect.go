package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SuspectRepository interface {
	CreateSuspect(suspect *models.Suspect) error
	GetPaginatedSuspects(c *fiber.Ctx) (*utils.Pagination, []models.Suspect, error)
	UpdateSuspect(id string, updates map[string]interface{}) error
	GetSuspectByID(id string) (models.Suspect, error)
	DeleteByID(id string) error
	SearchPaginatedSuspects(c *fiber.Ctx) (*utils.Pagination, []models.Suspect, error)
}

type SuspectRepositoryImpl struct {
	db *gorm.DB
}

func SuspectDbService(db *gorm.DB) SuspectRepository {
	return &SuspectRepositoryImpl{db: db}
}

// =====================
func (r *SuspectRepositoryImpl) CreateSuspect(suspect *models.Suspect) error {
	// Create the suspect record in the database
	result := r.db.Create(suspect)
	return result.Error
}

func (r *SuspectRepositoryImpl) GetPaginatedSuspects(c *fiber.Ctx) (*utils.Pagination, []models.Suspect, error) {
	pagination, suspects, err := utils.Paginate(c, r.db, models.Suspect{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, suspects, nil
}

func (r *SuspectRepositoryImpl) GetSuspectByID(id string) (models.Suspect, error) {
	var suspect models.Suspect
	err := r.db.First(&suspect, "id = ?", id).Error
	return suspect, err
}

func (r *SuspectRepositoryImpl) UpdateSuspect(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Suspect{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a suspect by ID
func (r *SuspectRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Suspect{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *SuspectRepositoryImpl) SearchPaginatedSuspects(c *fiber.Ctx) (*utils.Pagination, []models.Suspect, error) {
	// Get query parameters from request
	FirstName := c.Query("first_name")
	MiddleName := c.Query("middle_name")
	LastName := c.Query("last_name")
	Gender := c.Query("gender")
	PhoneNumber := c.Query("phone_number")
	Nin := c.Query("nin")
	Nationality := c.Query("nationality")
	Occupation := c.Query("occupation")
	Status := c.Query("status")

	// Start building the query
	query := r.db.Model(&models.Charge{})

	// Apply filters based on provided parameters
	if FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+FirstName+"%")
	}
	if MiddleName != "" {
		query = query.Where("middle_name ILIKE ?", "%"+MiddleName+"%")
	}
	if LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+LastName+"%")
	}
	if Gender != "" {
		query = query.Where("gender ILIKE ?", "%"+Gender+"%")
	}
	if PhoneNumber != "" {
		query = query.Where("phone_number ILIKE ?", "%"+PhoneNumber+"%")
	}
	if Nin != "" {
		query = query.Where("nin ILIKE ?", "%"+Nin+"%")
	}
	if Nationality != "" {
		query = query.Where("nationality ILIKE ?", "%"+Nationality+"%")
	}
	if Occupation != "" {
		query = query.Where("occupation ILIKE ?", "%"+Occupation+"%")
	}
	if Status != "" {
		query = query.Where("status ILIKE ?", "%"+Status+"%")
	}

	// Call the pagination helper
	pagination, suspects, err := utils.Paginate(c, query, models.Suspect{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, suspects, nil
}
