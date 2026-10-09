package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaginationAndReadiness(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			fmt.Fprint(w, `{"metadata":{"continue":"page2"},"items":[{"metadata":{"name":"one"},"status":{"phase":"Running","conditions":[{"type":"Ready","status":"True"}]}}]}`)
		} else {
			if r.URL.Query().Get("continue") != "page2" {
				t.Error("missing continuation token")
			}
			fmt.Fprint(w, `{"items":[{"metadata":{"name":"two"},"status":{"phase":"Pending"}}]}`)
		}
	}))
	defer server.Close()
	pods, err := listPods(server.Client(), server.URL+"?limit=1")
	if err != nil || len(pods) != 2 {
		t.Fatalf("%+v %v", pods, err)
	}
	if !ready(pods[0]) || ready(pods[1]) {
		t.Fatal("incorrect readiness")
	}
	completed := pod{}
	completed.Status.Phase = "Succeeded"
	if !ready(completed) {
		t.Fatal("completed jobs should pass")
	}
}
func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(403) }))
	defer server.Close()
	if _, err := listPods(server.Client(), server.URL); err == nil {
		t.Fatal("403 must fail")
	}
}
