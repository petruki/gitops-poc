package services

import (
	"testing"
)

func TestMain(m *testing.M) {
	InitEnv()
	m.Run()
}
