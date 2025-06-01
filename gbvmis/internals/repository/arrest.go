package repository

import (
	"fmt"
	"gbvmis/internals/models"
	"gbvmis/internals/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ArrestRepository interface {
	CreateArrest(arrest *models.Arrest) error
	GetPaginatedArrests(c *fiber.Ctx) (*utils.Pagination, []models.Arrest, error)
	UpdateArrest(id string, updates map[string]interface{}) error
	GetArrestByID(id string) (models.Arrest, error)
	DeleteByID(id string) error
	SearchPaginatedArrests(c *fiber.Ctx) (*utils.Pagination, []models.Arrest, error)
}

type ArrestRepositoryImpl struct {
	db *gorm.DB
}

func ArrestDbService(db *gorm.DB) ArrestRepository {
	return &ArrestRepositoryImpl{db: db}
}

// =================================

func (r *ArrestRepositoryImpl) CreateArrest(arrest *models.Arrest) error {
	return r.db.Create(arrest).Error
}

func (r *ArrestRepositoryImpl) GetPaginatedArrests(c *fiber.Ctx) (*utils.Pagination, []models.Arrest, error) {
	pagination, arrests, err := utils.Paginate(c, r.db, models.Arrest{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, arrests, nil
}

func (r *ArrestRepositoryImpl) GetArrestByID(id string) (models.Arrest, error) {
	var arrest models.Arrest
	err := r.db.First(&arrest, "id = ?", id).Error
	return arrest, err
}

func (r *ArrestRepositoryImpl) UpdateArrest(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.Arrest{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID deletes a arrest by ID
func (r *ArrestRepositoryImpl) DeleteByID(id string) error {
	if err := r.db.Delete(&models.Arrest{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ArrestRepositoryImpl) SearchPaginatedArrests(c *fiber.Ctx) (*utils.Pagination, []models.Arrest, error) {
	// Get query parameters from request
	loaction := c.Query("location")
	officerName := c.Query("officer_name")
	minArrestDate := c.Query("min_arrest_date")
	maxArrestDate := c.Query("max_arrest_date")

	// Start building the query
	query := r.db.Model(&models.Arrest{})

	// Apply filters based on provided parameters
	if loaction != "" {
		query = query.Where("location ILIKE ?", "%"+loaction+"%")
	}
	if officerName != "" {
		query = query.Where("officer_name ILIKE ?", "%"+officerName+"%")
	}

	if minArrestDate != "" {
		parsedMinArrestDate, err := time.Parse("2006-01-02", minArrestDate)
		if err == nil {
			query = query.Where("arrest_date >= ?", parsedMinArrestDate)
		} else {
			fmt.Println("Error parsing minArrestDate:", err)
		}
	}

	if maxArrestDate != "" {
		parsedMaxArrestDate, err := time.Parse("2006-01-02", maxArrestDate)
		if err == nil {
			query = query.Where("arrest_date <= ?", parsedMaxArrestDate)
		} else {
			fmt.Println("Error parsing maxArrestDate:", err)
		}
	}

	// Call the pagination helper
	pagination, arrests, err := utils.Paginate(c, query, models.Arrest{})
	if err != nil {
		return nil, nil, err
	}

	return &pagination, arrests, nil
}
