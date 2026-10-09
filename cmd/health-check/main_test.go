package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeStatuses(t *testing.T) {
	for _, code := range []int{200, 204, 503} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(code) }))
		got := probe(server.Client(), server.URL)
		server.Close()
		if got.Healthy != (code < 300) {
			t.Fatalf("code %d: %+v", code, got)
		}
	}
	got := probe(http.DefaultClient, "not a URL")
	if got.Healthy || got.Detail == "" {
		t.Fatal("invalid URL must fail")
	}
}
