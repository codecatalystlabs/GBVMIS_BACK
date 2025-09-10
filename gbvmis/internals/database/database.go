package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gbvmis/internals/config"
	"gbvmis/internals/models"

	"gbvmis/internals/seeder"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database is an interface for database operations
type Database interface {
	Connect()
	GetDB() *gorm.DB
	Migrate()
	Seed()
	Close()
}

// DBInstance implements the Database interface
type DBInstance struct {
	Db *gorm.DB
}

var dbInstance Database = &DBInstance{}

// NewDatabase returns a new instance of the database interface
func NewDatabase() Database {
	return dbInstance
}

// Connect establishes the database connection
func (d *DBInstance) Connect() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Config("DB_HOST"),
		config.Config("DB_USER"),
		config.Config("DB_PASS"),
		config.Config("DB_NAME"),
		port,
		config.Config("DB_SSLMODE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected to database")
	d.Db = db
}

// Migrate applies schema migrations
func (d *DBInstance) Migrate() {
	log.Println("Running migrations...")
	d.Db.AutoMigrate(
		&models.Suspect{},
		&models.Case{},
		&models.Arrest{},
		&models.Charge{},
		&models.Victim{},
		&models.HealthFacility{},
		&models.HealthPractitioner{},
		&models.Examination{},
		&models.Role{},
		&models.PolicePost{},
		&models.PoliceOfficer{},
		&models.Witness{},
		&models.Symptom{},
		&models.PostMortemSummary{},
		&models.Person{},
		&models.PoliceReport{},
		&models.ToxicologyForensicReport{},
	)
}

// Seed populates the database with initial data
func (d *DBInstance) Seed() {
	log.Println("Seeding database...")
	seeder.SeedDatabase(d.Db)
}

// GetDB returns the database instance
func (d *DBInstance) GetDB() *gorm.DB {
	return d.Db
}

// Close closes the database connection
func (d *DBInstance) Close() {
	sqlDB, err := d.Db.DB()
	if err != nil {
		log.Println("Error closing database:", err)
		return
	}
	sqlDB.Close()
}
