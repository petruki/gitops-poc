package services

import (
	"encoding/json"
	"os"

	"github.com/petruki/gitops-poc/src/model"
)

func ReadJsonFromFile(path string) string {
	file, _ := os.Open(path)
	defer file.Close()

	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	file.Read(bs)
	return string(bs)
}

func NewSnapshotFromJson(jsonData []byte) model.Snapshot {
	var snapshot model.Snapshot
	json.Unmarshal(jsonData, &snapshot)
	return snapshot
}
