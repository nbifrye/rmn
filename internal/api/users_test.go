package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUsers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users.json" {
			t.Errorf("expected path /users.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users":       []map[string]interface{}{{"id": 1, "login": "admin", "firstname": "Admin", "lastname": "User"}},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	users, total, err := client.ListUsers(context.Background(), UserListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if total != 1 {
		t.Errorf("expected total_count 1, got %d", total)
	}
}

func TestListUsers_WithParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "1" {
			t.Errorf("expected status=1, got %s", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("name") != "admin" {
			t.Errorf("expected name=admin, got %s", r.URL.Query().Get("name"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListUsers(context.Background(), UserListParams{Status: 1, Name: "admin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/1.json" {
			t.Errorf("expected path /users/1.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": map[string]interface{}{"id": 1, "login": "admin"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	user, err := client.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Login != "admin" {
		t.Errorf("expected 'admin', got %q", user.Login)
	}
}

func TestGetCurrentUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/current.json" {
			t.Errorf("expected path /users/current.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user": map[string]interface{}{"id": 1, "login": "me"},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	user, err := client.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Login != "me" {
		t.Errorf("expected 'me', got %q", user.Login)
	}
}

func TestListUsers_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":["Forbidden"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListUsers(context.Background(), UserListParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}
