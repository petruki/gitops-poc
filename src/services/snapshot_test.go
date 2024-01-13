package services

import (
	"testing"
)

func TestReadJsonFileToObject(t *testing.T) {
	json := ReadJsonFromFile("../../resources/default.json")
	AssertNotNil(t, json)
	AssertContains(t, json, "Playground")
}

func TestCreateSnapshotObjectFromJsonData(t *testing.T) {
	json := ReadJsonFromFile("../../resources/default.json")
	snapshot := NewSnapshotFromJson([]byte(json))
	AssertNotNil(t, snapshot)
	AssertContains(t, snapshot.Domain.Name, "Playground")
}
