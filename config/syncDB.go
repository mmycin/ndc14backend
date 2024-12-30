package config

import (
	"github.com/mmycin/ndc14/libs"
	"github.com/mmycin/ndc14/models"
)

func SyncDatabase() {
	libs.TimeElapsed(func() {
		// Define models in order of dependencies
		modelsToMigrate := []interface{}{
			&models.User{},
			&models.Notice{},
			&models.File{},
			&models.Contact{},
		}

		// Migrate sequentially to ensure proper foreign key creation
		for _, model := range modelsToMigrate {
			if err := DB.AutoMigrate(model); err != nil {
				libs.Error("Failed to migrate model: " + err.Error())
			}
		}

		libs.Success("Migrated models")
	})
}
