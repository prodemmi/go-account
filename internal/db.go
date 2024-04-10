package internal

import (
	"account/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=root password=root dbname=account-db port=5434 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	path := filepath.Join("sql",  "prepareDB.sql")
	c, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		// handle error.
	}
	sql := string(c)
	_, err := d.Exec(sql)
	if err != nil {
		// handle error.
	}

	err = db.AutoMigrate(&models.App{}, &models.User{}, &models.Password{}, models.Profile{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
