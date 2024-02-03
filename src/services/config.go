package services

import (
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	goEnv := os.Getenv("GO_ENV")

	if goEnv != "" {
		godotenv.Load("../../.env." + goEnv)
	} else {
		godotenv.Load("../../.env")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
