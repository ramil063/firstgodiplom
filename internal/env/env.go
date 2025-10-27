package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppEnv = ""

func InitEnvironmentVariables() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Получение переменных
	AppEnv = os.Getenv("APP_ENV")

	log.Printf("env: %s", AppEnv)
}
