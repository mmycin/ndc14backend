package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mmycin/ndc14/libs"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		cwd, _ := os.Getwd()
		libs.Fatal("Error loading .env file from " + cwd)
	}
	libs.Success("Loaded .env file")
}
