package repository

import (
	"gbvmis/internals/models"
	"gbvmis/internals/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ToxicologyReportRepository interface {
	Create(report *models.ToxicologyForensicReport) error
	GetAll() ([]models.ToxicologyForensicReport, error)
	GetByID(id uint) (*models.ToxicologyForensicReport, error)
	Update(report *models.ToxicologyForensicReport) error
	Delete(id uint) error

	GetPaginatedReports(c *fiber.Ctx) (*utils.Pagination, []models.ToxicologyForensicReport, error)
}

type toxicologyReportRepository struct {
	db *gorm.DB
}

func NewToxicologyReportRepository(db *gorm.DB) ToxicologyReportRepository {
	return &toxicologyReportRepository{db: db}
}

func (r *toxicologyReportRepository) GetPaginatedReports(c *fiber.Ctx) (*utils.Pagination, []models.ToxicologyForensicReport, error) {
	pagination, reports, err := utils.Paginate(c, r.db.
		Preload("Person").
		Preload("Witness").
		Preload("Practitioner").
		Preload("PoliceReport"), models.ToxicologyForensicReport{})
	if err != nil {
		return nil, nil, err
	}
	return &pagination, reports, nil
}

func (r *toxicologyReportRepository) Create(report *models.ToxicologyForensicReport) error {
	return r.db.Create(report).Error
}

func (r *toxicologyReportRepository) GetAll() ([]models.ToxicologyForensicReport, error) {
	var reports []models.ToxicologyForensicReport
	err := r.db.
		Preload("Person").
		Preload("Witness").
		Preload("Practitioner").
		Preload("PoliceReport").
		Find(&reports).Error
	return reports, err
}

func (r *toxicologyReportRepository) GetByID(id uint) (*models.ToxicologyForensicReport, error) {
	var report models.ToxicologyForensicReport
	err := r.db.
		Preload("Person").
		Preload("Witness").
		Preload("Practitioner").
		Preload("PoliceReport").
		First(&report, id).Error
	return &report, err
}

func (r *toxicologyReportRepository) Update(report *models.ToxicologyForensicReport) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update main report
		if err := tx.Save(&report).Error; err != nil {
			return err
		}

		// Update or create related entities
		if err := tx.Save(&report.Person).Error; err != nil {
			return err
		}
		if err := tx.Save(&report.PoliceReport).Error; err != nil {
			return err
		}

		// Replace nested slices (symptoms and summaries)
		if err := tx.Model(&report.Person).Association("PersonSymptoms").Replace(report.Person.PersonSymptoms); err != nil {
			return err
		}
		if err := tx.Model(&report.Person).Association("Summaries").Replace(report.Person.Summaries); err != nil {
			return err
		}

		return nil
	})
}

func (r *toxicologyReportRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var report models.ToxicologyForensicReport

		// 1️⃣ Load the report with associations
		if err := tx.Preload("Person.PersonSymptoms").
			Preload("Person.Summaries").
			Preload("PoliceReport").
			First(&report, id).Error; err != nil {
			return err
		}

		// 2️⃣ Delete related PersonSymptoms
		if len(report.Person.PersonSymptoms) > 0 {
			if err := tx.Where("person_id = ?", report.Person.ID).Delete(&models.PersonSymptom{}).Error; err != nil {
				return err
			}
		}

		// 3️⃣ Delete related PersonSummaries
		if len(report.Person.Summaries) > 0 {
			if err := tx.Where("person_id = ?", report.Person.ID).Delete(&models.PersonSummary{}).Error; err != nil {
				return err
			}
		}

		// 4️⃣ Delete related Person
		if report.Person.ID != 0 {
			if err := tx.Delete(&models.Person{}, report.Person.ID).Error; err != nil {
				return err
			}
		}

		// 5️⃣ Delete related PoliceReport
		if report.PoliceReport.ID != 0 {
			if err := tx.Delete(&models.PoliceReport{}, report.PoliceReport.ID).Error; err != nil {
				return err
			}
		}

		// 6️⃣ Finally, delete the ToxicologyForensicReport
		if err := tx.Delete(&models.ToxicologyForensicReport{}, report.ID).Error; err != nil {
			return err
		}

		return nil
	})
}
