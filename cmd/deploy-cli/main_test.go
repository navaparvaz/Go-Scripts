package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendPayloadAndToken(t *testing.T) {
	t.Setenv("DEPLOY_TOKEN", "test-token")
	payload := []byte(`{"application":"web","version":"v2","environment":"staging"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != string(payload) || r.Method != "POST" || r.Header.Get("Authorization") != "Bearer test-token" || r.Header.Get("Content-Type") != "application/json" {
			t.Error("unexpected deployment request")
		}
		w.WriteHeader(202)
	}))
	defer server.Close()
	if err := send(server.Client(), server.URL, payload); err != nil {
		t.Fatal(err)
	}
}
func TestSendRejectsFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer server.Close()
	if err := send(server.Client(), server.URL, []byte(`{}`)); err == nil {
		t.Fatal("500 must fail")
	}
}
