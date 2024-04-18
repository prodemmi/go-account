package controllers

import (
	http2 "account/http"
	"account/utils"
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	testServer, errServer := http2.SetupTestServer()
	if errServer != nil {
		log.Fatal(errServer.Error())
	}

	serverErr := testServer.Run(":" + os.Getenv("APP_PORT"))

	if serverErr != nil {
		log.Fatal(serverErr.Error())
	}

	pool, resource := utils.DockerTest()

	code := t.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAppController_Store(t *testing.T) {
	t.Log("testing")
	client := http.Client{}

	payload := CreateAppBody{}

	url := fmt.Sprintf("%v:%v/v1/users/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	myJson := bytes.NewBuffer([]byte(string(b)))

	res, errRes := client.Post(url, "application/json", myJson)
	if errRes != nil {
		_ = fmt.Errorf("error %s", errRes)
		return
	}

	t.Log("res   ", res)
}
