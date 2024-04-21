package utils

import (
	"account/database/models"
	"account/internal"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
)

func SetupDockerTestDB() (*dockertest.Pool, *dockertest.Resource) {
	absPath, _ := filepath.Abs("../.env.test")

	if err := godotenv.Load(absPath); err != nil {
		log.Fatal(err.Error())
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Name:       "account-test-db",
		Env: []string{
			fmt.Sprintf("POSTGRES_HOST=%v", os.Getenv("POSTGRES_HOST")),
			fmt.Sprintf("POSTGRES_PASSWORD=%v", os.Getenv("POSTGRES_PASSWORD")),
			fmt.Sprintf("POSTGRES_USER=%v", os.Getenv("POSTGRES_USER")),
			fmt.Sprintf("POSTGRES_DB=%v", os.Getenv("POSTGRES_DB")),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%v:%v@%s/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		hostAndPort,
		os.Getenv("POSTGRES_DB"),
	)

	log.Println("Connecting to database on url: ", databaseUrl)

	errExpire := resource.Expire(10)
	if errExpire != nil {
		log.Fatal(errExpire.Error())
	}

	var testDB *gorm.DB
	pool.MaxWait = 120 * time.Second
	time.Sleep(2 * time.Second)
	if err = pool.Retry(func() error {
		testDB, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
		d, _ := testDB.DB()
		if err != nil {
			return err
		}
		return d.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testDB.Exec("CREATE DATABASE" + " " + os.Getenv("POSTGRES_DB") + ";")
	testDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	errMigrate := testDB.AutoMigrate(&models.App{}, models.User{}, models.Otp{})
	if errMigrate != nil {
		log.Fatal(errMigrate.Error())
	}
	internal.DB = testDB

	return pool, resource
}

func DockerTestPurge(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
