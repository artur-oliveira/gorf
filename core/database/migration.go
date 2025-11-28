package database

import (
	"grf/core/config"
	"log"

	"gorm.io/gorm"
)

type MigrationOptions struct {
	DB     *gorm.DB
	Config *config.Config
	Models []interface{}
}

func RegisterMigrations(options *MigrationOptions) {
	log.Println("starting database migrations")

	if err := PerformMigration(options.DB, options.Config, options.Models...); err != nil {
		log.Fatalf("Failed to perform automigrations: %v", err)
	}

	log.Println("database migrations complete.")
}
