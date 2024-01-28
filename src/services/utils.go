package services

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	godotenv.Load("../../.env")
	return os.Getenv(key)
}

func FormatJSON(jsonString string) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, []byte(jsonString), "", "  ")
	if error != nil {
		return jsonString
	}
	return string(prettyJSON.String())
}
