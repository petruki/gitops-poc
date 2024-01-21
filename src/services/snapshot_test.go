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

	// Check changes
	diff := CheckSnapshotDiff(snapshotLeft, snapshotRight, CHANGED)
	AssertNotNil(t, diff)
	println("CHANGED")
	println(NewJsonStringFromSnapshot(diff))

	// Check new
	diff2 := CheckSnapshotDiff(snapshotRight, snapshotLeft, NEW)
	AssertNotNil(t, diff2)
	println("NEW")
	println(NewJsonStringFromSnapshot(diff2))

	// Check deleted
	diff3 := CheckSnapshotDiff(snapshotLeft, snapshotRight, DELETED)
	AssertNotNil(t, diff3)
	println("DELETED")
	println(NewJsonStringFromSnapshot(diff3))
}
