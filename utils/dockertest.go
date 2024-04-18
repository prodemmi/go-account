package utils

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"time"
)

var testDB *sql.DB

func DockerTest() (*dockertest.Pool, *dockertest.Resource) {
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
			fmt.Sprintf("POSTGRES_PORT=%v", os.Getenv("POSTGRES_PORT")),
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

	errExpire := resource.Expire(120)
	if errExpire != nil {
		log.Fatal(errExpire.Error())
	}

	pool.MaxWait = 10 * time.Second
	if err = pool.Retry(func() error {
		testDB, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return pool, resource
}
