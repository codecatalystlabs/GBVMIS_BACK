package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *models.Role) error
	GetPaginatedRoles(c *fiber.Ctx) (*utils.Pagination, []models.Role, error)
	UpdateRole(id string, updates map[string]interface{}) error
	GetRoleByID(id string) (models.Role, error)
	DeleteByID(id string) error
	SearchPaginatedRoles(c *fiber.Ctx) (*utils.Pagination, []models.Role, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func RoleDbService(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// =====================
func (r *RoleRepositoryImpl) CreateRole(role *models.Role) error {
	// Create the Role record in the database
	result := r.db.Create(role)
	return result.Error
}

func (r *RoleRepositoryImpl) GetPaginatedRoles(c *fiber.Ctx) (*utils.Pagination, []models.Role, error) {
	pagination, roles, err := utils.Paginate(c, r.db, models.Role{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, roles, nil
}

func (r *RoleRepositoryImpl) GetRoleByID(id string) (models.Role, error) {
	var role models.Role
	err := r.db.First(&role, "id = ?", id).Error
	return role, err
}

func (r *RoleRepositoryImpl) UpdateRole(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Role{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a role by ID
func (r *RoleRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Role{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleRepositoryImpl) SearchPaginatedRoles(c *fiber.Ctx) (*utils.Pagination, []models.Role, error) {
	// Get query parameters from request
	Name := c.Query("name")

	// Start building the query
	query := r.db.Model(&models.Charge{})

	// Apply filters based on provided parameters
	if Name != "" {
		query = query.Where("name ILIKE ?", "%"+Name+"%")
	}

	// Call the pagination helper
	pagination, roles, err := utils.Paginate(c, query, models.Role{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, roles, nil
}
