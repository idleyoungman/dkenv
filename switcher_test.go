package main

import (
	"testing"
)

func TestApiToVersion(t *testing.T) {
	v1, _ := apiToVersion("1.12")
	if v1 != "1.0.1" {
		t.Error("Expected 1.0.1, got ", v1)
	}

	v2, _ := apiToVersion("1.17")
	if v2 != "1.5.0" {
		t.Error("Expected 1.5.0, got ", v2)
	}
}

func TestlistDownloadedVersions(t *testing.T) {

}
