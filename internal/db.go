package internal

import (
	"account/database/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

var DB *gorm.DB

func SetupDatabase() error {
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgUser := os.Getenv("POSTGRES_USER")
	pgName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	path, _ := filepath.Abs("../database/sql/prepareDB.sql")
	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		// handle error.
	}
	sql := string(c)
	db.Exec(sql)

	errMigrate := db.AutoMigrate(&models.App{}, models.User{}, models.Otp{})
	if errMigrate != nil {
		return errMigrate
	}

	if sqlDB, _ := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetMaxOpenConns(3)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	DB = db

	return nil
}
