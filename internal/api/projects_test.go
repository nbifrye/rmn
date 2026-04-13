package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects.json" {
			t.Errorf("expected path /projects.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects":    []map[string]interface{}{{"id": 1, "name": "Test", "identifier": "test"}},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	projects, total, err := client.ListProjects(context.Background(), ProjectListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if total != 1 {
		t.Errorf("expected total_count 1, got %d", total)
	}
}

func TestListProjects_WithParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "active" {
			t.Errorf("expected status=active, got %s", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListProjects(context.Background(), ProjectListParams{Status: "active", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test.json" {
			t.Errorf("expected path /projects/test.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 1, "name": "Test", "identifier": "test"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	project, err := client.GetProject(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.Name != "Test" {
		t.Errorf("expected 'Test', got %q", project.Name)
	}
}

func TestGetProject_WithInclude(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("include") != "trackers,issue_categories" {
			t.Errorf("expected include=trackers,issue_categories, got %s", r.URL.Query().Get("include"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 1, "name": "Test"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetProject(context.Background(), "test", []string{"trackers", "issue_categories"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"project": map[string]interface{}{"id": 2, "name": "New Project", "identifier": "new-project"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	project, err := client.CreateProject(context.Background(), ProjectCreateParams{
		Name: "New Project", Identifier: "new-project",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.ID != 2 {
		t.Errorf("expected ID 2, got %d", project.ID)
	}
}

func TestUpdateProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test.json" {
			t.Errorf("expected path /projects/test.json, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateProject(context.Background(), "test", ProjectUpdateParams{
		Name: StringPtr("Updated"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestArchiveProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/archive.json" {
			t.Errorf("expected path /projects/test/archive.json, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.ArchiveProject(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnarchiveProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/unarchive.json" {
			t.Errorf("expected path /projects/test/unarchive.json, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UnarchiveProject(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteProject(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListProjects_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":["Server error"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListProjects(context.Background(), ProjectListParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListProjects_WithOffset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("offset") != "20" {
			t.Errorf("expected offset=20, got %s", r.URL.Query().Get("offset"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListProjects(context.Background(), ProjectListParams{Offset: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetProject_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetProject(context.Background(), "nonexistent", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateProject_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":["Name cannot be blank"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.CreateProject(context.Background(), ProjectCreateParams{Name: "Test", Identifier: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
}
