package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PoliceOfficerRepository interface {
	CreatePoliceOfficer(policeOfficer *models.PoliceOfficer) error
	GetPaginatedPoliceOfficers(c *fiber.Ctx) (*utils.Pagination, []models.PoliceOfficer, error)
	UpdatePoliceOfficer(id string, updates map[string]interface{}) error
	GetPoliceOfficerByID(id string) (models.PoliceOfficer, error)
	DeleteByID(id string) error
	SearchPaginatedPoliceOfficers(c *fiber.Ctx) (*utils.Pagination, []models.PoliceOfficer, error)
	FindRolesByIDs(ids []uint, out *[]*models.Role) error
	UpdateOfficerRoles(officer *models.PoliceOfficer, roles []*models.Role) error
}

type PoliceOfficerRepositoryImpl struct {
	db *gorm.DB
}

func PoliceOfficerDbService(db *gorm.DB) PoliceOfficerRepository {
	return &PoliceOfficerRepositoryImpl{db: db}
}

// =================================

func (r *PoliceOfficerRepositoryImpl) CreatePoliceOfficer(policeOfficer *models.PoliceOfficer) error {
	return r.db.Create(policeOfficer).Error
}

func (r *PoliceOfficerRepositoryImpl) GetPaginatedPoliceOfficers(c *fiber.Ctx) (*utils.Pagination, []models.PoliceOfficer, error) {
	pagination, policeOfficers, err := utils.Paginate(c, r.db.Preload("Roles"), models.PoliceOfficer{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, policeOfficers, nil
}

func (r *PoliceOfficerRepositoryImpl) GetPoliceOfficerByID(id string) (models.PoliceOfficer, error) {
	var policeOfficer models.PoliceOfficer
	err := r.db.Preload("Roles").First(&policeOfficer, "id = ?", id).Error
	return policeOfficer, err
}

func (r *PoliceOfficerRepositoryImpl) UpdatePoliceOfficer(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.PoliceOfficer{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a policeOfficer by ID
func (r *PoliceOfficerRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.PoliceOfficer{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PoliceOfficerRepositoryImpl) SearchPaginatedPoliceOfficers(c *fiber.Ctx) (*utils.Pagination, []models.PoliceOfficer, error) {
	// Get query parameters from request
	FirstName := c.Query("first_name")
	LastName := c.Query("last_name")
	BadgeNo := c.Query("badge_no")
	Username := c.Query("username")
	PostID := c.Query("post_id")

	// Start building the query
	query := r.db.Preload("Roles").Model(&models.PoliceOfficer{})

	// Apply filters based on provided parameters
	if FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+FirstName+"%")
	}
	if LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+LastName+"%")
	}
	if BadgeNo != "" {
		query = query.Where("badge_no ILIKE ?", "%"+BadgeNo+"%")
	}
	if Username != "" {
		query = query.Where("username ILIKE ?", "%"+Username+"%")
	}

	if PostID != "" {
		if _, err := strconv.Atoi(PostID); err == nil {
			query = query.Where("post_id = ?", PostID)
		}
	}

	// Call the pagination helper
	pagination, policeOfficers, err := utils.Paginate(c, query, models.PoliceOfficer{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, policeOfficers, nil
}

func (r *PoliceOfficerRepositoryImpl) FindRolesByIDs(ids []uint, out *[]*models.Role) error {
	return r.db.Where("id IN ?", ids).Find(out).Error
}

func (r *PoliceOfficerRepositoryImpl) UpdateOfficerRoles(officer *models.PoliceOfficer, roles []*models.Role) error {
	return r.db.Model(officer).Association("Roles").Replace(roles)
}
