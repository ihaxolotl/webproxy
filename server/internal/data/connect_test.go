package data

import (
	"testing"
)

func TestDatabaseSetup(t *testing.T) {
	if err := New().Setup(); err != nil {
		t.Fatal(err)
	}
}
