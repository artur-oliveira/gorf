package controller_test

import (
	"grf/core/bootstrap"
	"grf/core/config"
	"grf/core/server"
	"grf/domain/auth"
	"log"
	"os"
	"testing"
)

var testApp *server.App
var err error

func TestMain(m *testing.M) {
	testApp, err = bootstrap.NewApp(config.Config{
		DBName:               "file:memdb1?mode=memory&cache=shared",
		DBVendor:             "sqlite",
		DBMigrate:            true,
		DBLogLevel:           "info",
		DBMaxIdle:            10,
		DBMaxOpened:          30,
		DBMaxLifeTimeSeconds: 30,

		AppName: "TestGRF",

		Env: "development",

		ServerPort:         "1234",
		ServerIdleTimeout:  3,
		ServerReadTimeout:  3,
		ServerWriteTimeout: 3,

		JWTSecret:               "test_super_secret_jwt_secret_do_not_verify",
		JWTExpiresInMinutes:     30,
		JWTRefreshExpiresInDays: 1,
	}, auth.GetModels())
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	log.Println("Suíte de testes 'auth' concluída.")
	os.Exit(code)
}
