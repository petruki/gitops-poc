package services

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	godotenv.Load("../../.env")
	return os.Getenv(key)
}
