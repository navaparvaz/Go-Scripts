package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMonitorEmitsOnlyTransitions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 4 {
			cancel()
		}
		if calls >= 3 {
			w.WriteHeader(503)
		} else {
			w.WriteHeader(200)
		}
	}))
	defer server.Close()
	var output bytes.Buffer
	if err := monitor(ctx, []string{server.URL}, time.Millisecond, json.NewEncoder(&output)); err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(&output)
	var first, second event
	if err := decoder.Decode(&first); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&second); err != nil {
		t.Fatal(err)
	}
	if !first.Healthy || second.Healthy {
		t.Fatal("incorrect transitions")
	}
	if decoder.Decode(&event{}) == nil {
		t.Fatal("unchanged results should be suppressed")
	}
}
