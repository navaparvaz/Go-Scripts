package main

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestWholeWordMatches(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.log")
	if err := os.WriteFile(path, []byte("ERROR failed\nINFO errorless\npanic and fatal\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := analyze(path, regexp.MustCompile(`(?i)\b(error|fatal|panic)\b`))
	if err != nil || got.Lines != 3 || got.Matches != 2 {
		t.Fatalf("%+v %v", got, err)
	}
	if _, err := analyze(path+"missing", regexp.MustCompile("error")); err == nil {
		t.Fatal("missing file must fail")
	}
}
