//go:build smoke_test
// +build smoke_test

package tests

import (
	"account/http/server"
	"account/utils"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestSmoke(t *testing.T) {
	var wg sync.WaitGroup

	absPath, _ := filepath.Abs("../.env.test")

	if err := godotenv.Load(absPath); err != nil {
		log.Fatal(err.Error())
	}
	wg.Add(1)
	go server.SetupTestServer(&wg)
	wg.Wait()

	res, err := http.Get(utils.GetTestPingUrl())

	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, strings.Contains(string(body), "ok"))
}
