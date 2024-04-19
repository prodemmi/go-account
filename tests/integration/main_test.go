//go:build integration_tests
// +build integration_tests

package integration

import (
	"account/http/server"
	"account/utils"
	"sync"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	TestData.Init()
	var wg sync.WaitGroup
	pool, resource := utils.SetupDockerTestDB()
	wg.Add(1)
	go server.SetupTestServer(&wg)
	wg.Wait()
	time.Sleep(2 * time.Second)

	t.Run("app store app with true body", TestAppStoreAppWithTrueBody)
	t.Run("app get apps", TestAppGetApps)
	t.Run("app get one app", TestAppGetOneApp)

	t.Run("user store with true body", UserStoreWithTrueBody)
	t.Run("user get users", UserGetUsers)

	t.Run("auth sign up", AuthSignUp)
	t.Run("auth verify sign up", AuthVerifySignUp)

	t.Run("App delete one app", TestAppDeleteOneApp)

	utils.DockerTestPurge(pool, resource)
}
