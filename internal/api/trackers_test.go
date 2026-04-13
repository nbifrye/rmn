package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTrackers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trackers.json" {
			t.Errorf("expected path /trackers.json, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"trackers": []map[string]interface{}{
				{"id": 1, "name": "Bug"},
				{"id": 2, "name": "Feature"},
			},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	trackers, err := client.ListTrackers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trackers) != 2 {
		t.Fatalf("expected 2 trackers, got %d", len(trackers))
	}
	if trackers[0].Name != "Bug" {
		t.Errorf("expected 'Bug', got %q", trackers[0].Name)
	}
}

func TestListTrackers_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.ListTrackers(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
