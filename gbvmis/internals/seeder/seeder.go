package seeder

import (
	"gbvmis/internals/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SeedDatabase(db *gorm.DB) {

	role := []models.Role{
		{
			Name: "Admin",
		},
		{
			Name: "Station Manager",
		},
		{
			Name: "User",
		},
	}

	// Check if the companies table already has data
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		if err := db.Create(&role).Error; err != nil {
			log.Fatalf("Failed to seed roles: %v", err)
		} else {
			log.Println("Roles data seeded successfully")
		}
	} else {
		log.Println("Roles table already seeded, skipping...")
	}

	companies := []models.PolicePost{
		{
			Name:     "Kampala Central Police Head quarters",
			Location: "Kampala Main Street",
			Contact:  "+25641000789",
		},
		{
			Name:     "Wandegeya Police Post",
			Location: "Wandegeya",
			Contact:  "+25641000788",
		},
	}

	// Check if the companies table already has data
	var companyCount int64
	db.Model(&models.PolicePost{}).Count(&companyCount)
	if companyCount == 0 {
		if err := db.Create(&companies).Error; err != nil {
			log.Fatalf("Failed to seed police posts: %v", err)
		} else {
			log.Println("Police post data seeded successfully")
		}
	} else {
		log.Println("Police posts table already seeded, skipping...")
	}

	// Hashing password for users
	passwordHash, err := hashPassword("Admin123")
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	users := []models.PoliceOfficer{
		{
			FirstName: "John",
			LastName:  "Doe",
			Rank:      "ASP",
			BadgeNo:   "B00192",
			Username:  "Admin",
			Email:     "admin@example.com",
			Password:  passwordHash,
			Phone:     "07812663647",
			PostID:    1,
		},
	}

	var userCount int64
	db.Model(&models.PoliceOfficer{}).Count(&userCount)
	if userCount == 0 {
		if err := db.Create(&users).Error; err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		} else {
			log.Println("User data seeded successfully")
		}
	} else {
		log.Println("Users table already seeded, skipping...")
	}

}
