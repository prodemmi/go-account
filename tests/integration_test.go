//go:build integration_tests
// +build integration_tests

package tests

import (
	"account/http/server"
	"account/tests/integration"
	"account/utils"
	"sync"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	integration.TestData.Init()
	var wg sync.WaitGroup
	pool, resource := utils.SetupDockerTestDB()
	wg.Add(1)
	go server.SetupTestServer(&wg)
	wg.Wait()
	time.Sleep(2 * time.Second)

	t.Run("app store app with true body", integration.TestAppStoreAppWithTrueBody)
	t.Run("app get apps", integration.TestAppGetApps)
	t.Run("app get one app", integration.TestAppGetOneApp)

	t.Run("user store with true body", integration.UserStoreWithTrueBody)
	t.Run("user get users", integration.UserGetUsers)

	t.Run("auth sign up", integration.AuthSignUp)
	t.Run("auth verify sign up", integration.AuthVerifyOtp)
	t.Run("auth logout after sign up", integration.AuthLogout)

	t.Run("auth sign in", integration.AuthSigIn)
	t.Run("auth verify sign in", integration.AuthVerifyOtp)
	t.Run("auth logout after sign in", integration.AuthLogout)

	t.Run("App delete one app", integration.TestAppDeleteOneApp)

	utils.DockerTestPurge(pool, resource)
}
