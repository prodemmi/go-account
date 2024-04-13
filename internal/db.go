package internal

import (
	"account/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func SetupDatabase() error {
	dsn := "host=localhost user=root password=root dbname=account-db port=5434 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	path := filepath.Join("sql", "prepareDB.sql")
	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		// handle error.
	}
	sql := string(c)
	db.Exec(sql)

	errMigrate := db.AutoMigrate(&models.App{}, models.User{}, models.Profile{})
	if errMigrate != nil {
		return errMigrate
	}

	DB = db

	return nil
}
