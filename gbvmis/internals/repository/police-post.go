package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PolicePostRepository interface {
	CreatePolicePost(policePost *models.PolicePost) error
	GetPaginatedPolicePosts(c *fiber.Ctx) (*utils.Pagination, []models.PolicePost, error)
	UpdatePolicePost(id string, updates map[string]interface{}) error
	GetPolicePostByID(id string) (models.PolicePost, error)
	DeleteByID(id string) error
	SearchPaginatedPolicePosts(c *fiber.Ctx) (*utils.Pagination, []models.PolicePost, error)
}

type PolicePostRepositoryImpl struct {
	db *gorm.DB
}

func PolicePostDbService(db *gorm.DB) PolicePostRepository {
	return &PolicePostRepositoryImpl{db: db}
}

// =================================

func (r *PolicePostRepositoryImpl) CreatePolicePost(policePost *models.PolicePost) error {
	return r.db.Create(policePost).Error
}

func (r *PolicePostRepositoryImpl) GetPaginatedPolicePosts(c *fiber.Ctx) (*utils.Pagination, []models.PolicePost, error) {
	pagination, policePosts, err := utils.Paginate(c, r.db, models.PolicePost{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, policePosts, nil
}

func (r *PolicePostRepositoryImpl) GetPolicePostByID(id string) (models.PolicePost, error) {
	var policePost models.PolicePost
	err := r.db.First(&policePost, "id = ?", id).Error
	return policePost, err
}

func (r *PolicePostRepositoryImpl) UpdatePolicePost(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.PolicePost{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a policePost by ID
func (r *PolicePostRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.PolicePost{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PolicePostRepositoryImpl) SearchPaginatedPolicePosts(c *fiber.Ctx) (*utils.Pagination, []models.PolicePost, error) {
	// Get query parameters from request
	Name := c.Query("name")
	Location := c.Query("location")

	// Start building the query
	query := r.db.Model(&models.PolicePost{})

	// Apply filters based on provided parameters
	if Name != "" {
		query = query.Where("name ILIKE ?", "%"+Name+"%")
	}
	if Location != "" {
		query = query.Where("location ILIKE ?", "%"+Location+"%")
	}

	// Call the pagination helper
	pagination, policePosts, err := utils.Paginate(c, query, models.PolicePost{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, policePosts, nil
}
