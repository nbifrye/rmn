package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTimeEntries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/time_entries.json" {
			t.Errorf("expected path /time_entries.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entries": []map[string]interface{}{{"id": 1, "hours": 2.5, "comments": "Work"}},
			"total_count":  1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	entries, total, err := client.ListTimeEntries(context.Background(), TimeEntryListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if total != 1 {
		t.Errorf("expected total_count 1, got %d", total)
	}
}

func TestListTimeEntries_WithParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("project_id") != "test" {
			t.Errorf("expected project_id=test, got %s", r.URL.Query().Get("project_id"))
		}
		if r.URL.Query().Get("from") != "2024-01-01" {
			t.Errorf("expected from=2024-01-01, got %s", r.URL.Query().Get("from"))
		}
		if r.URL.Query().Get("to") != "2024-01-31" {
			t.Errorf("expected to=2024-01-31, got %s", r.URL.Query().Get("to"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entries": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListTimeEntries(context.Background(), TimeEntryListParams{
		ProjectID: "test", From: "2024-01-01", To: "2024-01-31",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTimeEntry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/time_entries/1.json" {
			t.Errorf("expected path /time_entries/1.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 1, "hours": 2.5},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	entry, err := client.GetTimeEntry(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Hours != 2.5 {
		t.Errorf("expected 2.5 hours, got %f", entry.Hours)
	}
}

func TestCreateTimeEntry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entry": map[string]interface{}{"id": 2, "hours": 1.0},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	entry, err := client.CreateTimeEntry(context.Background(), TimeEntryCreateParams{Hours: 1.0, IssueID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.ID != 2 {
		t.Errorf("expected ID 2, got %d", entry.ID)
	}
}

func TestUpdateTimeEntry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateTimeEntry(context.Background(), 1, TimeEntryUpdateParams{Hours: Float64Ptr(3.0)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTimeEntry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteTimeEntry(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTimeEntries_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListTimeEntries(context.Background(), TimeEntryListParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
