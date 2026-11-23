package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	pgmodels "aman/internal/adapters/outbound/postgres/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB opens a connection to Postgres, runs auto-migration and returns the
// *gorm.DB. Environment variables are expected to be loaded by the caller
// (main loads .env before building the fx graph).
func NewDB() (*gorm.DB, error) {
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgUser := os.Getenv("POSTGRES_USER")
	pgName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgName,
	)

	// SQL logging is a development aid only: production keeps gorm silent.
	logLevel := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&pgmodels.App{}, &pgmodels.Session{}, &pgmodels.Otp{}, &pgmodels.OTT{}); err != nil {
		return nil, err
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetMaxOpenConns(3)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("postgres: connected and migrated")
	return db, nil
}
