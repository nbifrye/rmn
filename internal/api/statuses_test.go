package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListStatuses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issue_statuses.json" {
			t.Errorf("expected path /issue_statuses.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"issue_statuses": []map[string]interface{}{
				{"id": 1, "name": "New", "is_closed": false},
				{"id": 5, "name": "Closed", "is_closed": true},
			},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	statuses, err := client.ListStatuses(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}
	if statuses[0].Name != "New" {
		t.Errorf("expected 'New', got %q", statuses[0].Name)
	}
	if statuses[1].IsClosed != true {
		t.Errorf("expected Closed to have is_closed=true")
	}
}

func TestListStatuses_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.ListStatuses(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
