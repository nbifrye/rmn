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

func TestListTimeEntries_AllParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("issue_id") != "5" {
			t.Errorf("expected issue_id=5, got %s", q.Get("issue_id"))
		}
		if q.Get("user_id") != "3" {
			t.Errorf("expected user_id=3, got %s", q.Get("user_id"))
		}
		if q.Get("spent_on") != "2024-01-15" {
			t.Errorf("expected spent_on=2024-01-15, got %s", q.Get("spent_on"))
		}
		if q.Get("activity_id") != "9" {
			t.Errorf("expected activity_id=9, got %s", q.Get("activity_id"))
		}
		if q.Get("limit") != "20" {
			t.Errorf("expected limit=20, got %s", q.Get("limit"))
		}
		if q.Get("offset") != "10" {
			t.Errorf("expected offset=10, got %s", q.Get("offset"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"time_entries": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListTimeEntries(context.Background(), TimeEntryListParams{
		IssueID: 5, UserID: 3, SpentOn: "2024-01-15", ActivityID: 9, Limit: 20, Offset: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTimeEntry_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetTimeEntry(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateTimeEntry_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":["Invalid"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.CreateTimeEntry(context.Background(), TimeEntryCreateParams{Hours: 1.0, IssueID: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}
