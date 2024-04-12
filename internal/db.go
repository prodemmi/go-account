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

	err1 := db.AutoMigrate(&models.App{})
	if err1 != nil {
		return err1
	}
	err2 := db.AutoMigrate(&models.User{})
	if err2 != nil {
		return err2
	}
	err3 := db.AutoMigrate(&models.Profile{})
	if err3 != nil {
		return err3
	}

	DB = db

	return nil
}
