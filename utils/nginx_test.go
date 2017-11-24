package utils

import "testing"

func TestNginxReloadOk(t *testing.T) {
	err := NginxReload("/bin/true")

	if err != nil {
		t.Fatalf("Unexpected error: '%s'", err)
	}
}

func TestNginxReloadFail(t *testing.T) {
	err := NginxReload("/bin/false")

	if err == nil {
		t.Fatal("Expected to get error :(")
	}
}
