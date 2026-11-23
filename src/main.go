package main

import (
	"os"

	"github.com/joho/godotenv"

	"aman/internal/adapters/inbound/http"
	"aman/internal/adapters/outbound/notification"
	"aman/internal/adapters/outbound/otpgen"
	"aman/internal/adapters/outbound/postgres"
	"aman/internal/adapters/outbound/security"
	"aman/internal/adapters/outbound/webhook"
	"aman/internal/domains/services"

	"go.uber.org/fx"
)

func main() {
	if os.Getenv("APP_ENV") == "staging" {
		_ = godotenv.Load(".staging.env")
	} else {
		_ = godotenv.Load()
	}

	app := fx.New(
		postgres.Module(),
		security.Module(),
		otpgen.Module(),
		notification.Module(),
		webhook.Module(),
		services.Module(),
		httpadapter.Module(),
	)

	app.Run()
}
