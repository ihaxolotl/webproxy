package data

import (
	"testing"
)

func TestConnect(t *testing.T) {
	if _, err := Connect(); err != nil {
		t.Fatal(err)
	}
}

func TestSetupDatabase(t *testing.T) {
	if _, err := SetupDatabase(); err != nil {
		t.Fatal(err)
	}
}
