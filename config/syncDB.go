package config

import (
	"sync"

	"github.com/mmycin/ndc14/libs"
	"github.com/mmycin/ndc14/models"
)

func SyncDatabase() {
	libs.TimeElapsed(func() {
		var wg sync.WaitGroup

		// List of models to migrate
		modelsToMigrate := []interface{}{
			&models.User{},
			&models.Notice{},
			&models.Contact{},
		}
		// Start a goroutine for each model migration
		for _, model := range modelsToMigrate {
			wg.Add(1)
			go func(m interface{}) {
				defer wg.Done()
				DB.AutoMigrate(m)
			}(model)
		}
		// Wait for all goroutines to finish
		wg.Wait()
		libs.Success("Migrated models")
	})
}
