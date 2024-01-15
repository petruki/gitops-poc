package services

import (
	"testing"
)

const DEFAULT_JSON = "../../resources/default.json"

func TestReadJsonFileToObject(t *testing.T) {
	json := ReadJsonFromFile(DEFAULT_JSON)
	AssertNotNil(t, json)
	AssertContains(t, json, "Release 1")
}

func TestCreateSnapshotObjectFromJsonData(t *testing.T) {
	json := ReadJsonFromFile(DEFAULT_JSON)
	snapshot := NewSnapshotFromJson([]byte(json))
	AssertNotNil(t, snapshot)
	AssertEqual(t, len(snapshot.Domain.Group), 1)
}

func TestCheckSnapshotDiff(t *testing.T) {
	jsonLeft := ReadJsonFromFile(DEFAULT_JSON)
	jsonRight := ReadJsonFromFile("../../resources/merge1.json")
	snapshotLeft := NewSnapshotFromJson([]byte(jsonLeft))
	snapshotRight := NewSnapshotFromJson([]byte(jsonRight))

	diff := CheckSnapshotDiff(snapshotLeft, snapshotRight)
	AssertNotNil(t, diff)
	// println(NewJsonStringFromSnapshot(diff))

	diff2 := CheckSnapshotDiff(snapshotRight, snapshotLeft)
	AssertNotNil(t, diff2)
	// println(NewJsonStringFromSnapshot(diff2))
}
