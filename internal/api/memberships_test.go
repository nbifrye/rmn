package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMemberships(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/memberships.json" {
			t.Errorf("expected path /projects/test/memberships.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"memberships": []map[string]interface{}{
				{"id": 1, "user": map[string]interface{}{"id": 1, "name": "Admin"}, "roles": []map[string]interface{}{{"id": 1, "name": "Manager"}}},
			},
			"total_count": 1,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	memberships, total, err := client.ListMemberships(context.Background(), "test", MembershipListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(memberships) != 1 {
		t.Fatalf("expected 1 membership, got %d", len(memberships))
	}
	if total != 1 {
		t.Errorf("expected total_count 1, got %d", total)
	}
}

func TestGetMembership(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/memberships/1.json" {
			t.Errorf("expected path /memberships/1.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{"id": 1, "user": map[string]interface{}{"id": 1, "name": "Admin"}},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	membership, err := client.GetMembership(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if membership.ID != 1 {
		t.Errorf("expected ID 1, got %d", membership.ID)
	}
}

func TestCreateMembership(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"membership": map[string]interface{}{"id": 2},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	membership, err := client.CreateMembership(context.Background(), "test", MembershipCreateParams{UserID: 1, RoleIDs: []int{1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if membership.ID != 2 {
		t.Errorf("expected ID 2, got %d", membership.ID)
	}
}

func TestUpdateMembership(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateMembership(context.Background(), 1, MembershipUpdateParams{RoleIDs: []int{1, 2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteMembership(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteMembership(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListMemberships_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListMemberships(context.Background(), "nonexistent", MembershipListParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListMemberships_WithLimitOffset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("offset") != "5" {
			t.Errorf("expected offset=5, got %s", r.URL.Query().Get("offset"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"memberships": []map[string]interface{}{}, "total_count": 0,
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, _, err := client.ListMemberships(context.Background(), "test", MembershipListParams{Limit: 10, Offset: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetMembership_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetMembership(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateMembership_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":["Invalid"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.CreateMembership(context.Background(), "test", MembershipCreateParams{UserID: 1, RoleIDs: []int{1}})
	if err == nil {
		t.Fatal("expected error")
	}
}
