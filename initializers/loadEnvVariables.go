package initializers

import "github.com/joho/godotenv"
import "log"

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
}