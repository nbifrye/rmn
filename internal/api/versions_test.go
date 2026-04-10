package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListVersions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/versions.json" {
			t.Errorf("expected path /projects/test/versions.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"versions":    []map[string]interface{}{{"id": 1, "name": "v1.0", "status": "open"}},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	versions, total, err := client.ListVersions(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}
	if total != 1 {
		t.Errorf("expected total_count 1, got %d", total)
	}
}

func TestGetVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/versions/1.json" {
			t.Errorf("expected path /versions/1.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": map[string]interface{}{"id": 1, "name": "v1.0"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	version, err := client.GetVersion(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version.Name != "v1.0" {
		t.Errorf("expected 'v1.0', got %q", version.Name)
	}
}

func TestCreateVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version": map[string]interface{}{"id": 2, "name": "v2.0"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	version, err := client.CreateVersion(context.Background(), "test", VersionCreateParams{Name: "v2.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version.ID != 2 {
		t.Errorf("expected ID 2, got %d", version.ID)
	}
}

func TestUpdateVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateVersion(context.Background(), 1, VersionUpdateParams{Name: StringPtr("v1.1")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteVersion(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListVersions_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListVersions(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}
