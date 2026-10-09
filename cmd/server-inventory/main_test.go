package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadPairs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "os-release")
	if err := os.WriteFile(path, []byte("PRETTY_NAME=\"Example Linux\"\nID=example\n#comment\n"), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := readPairs(path, "=")
	if err != nil || values["PRETTY_NAME"] != "Example Linux" || values["ID"] != "example" {
		t.Fatalf("%v %v", values, err)
	}
}
