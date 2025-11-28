package main

import (
	"grf/core/bootstrap"
	"grf/core/config"
	"grf/domain/auth"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./", "app")
	if err != nil {
		log.Fatal(err)
	}

	var allModels []interface{}
	allModels = append(allModels, auth.GetModels()...)

	app, err := bootstrap.NewApp(cfg, allModels)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Start())
}
