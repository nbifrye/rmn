package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWikiPages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/wiki/index.json" {
			t.Errorf("expected path /projects/test/wiki/index.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_pages": []map[string]interface{}{
				{"title": "Wiki", "version": 1},
			},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	pages, err := client.ListWikiPages(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if pages[0].Title != "Wiki" {
		t.Errorf("expected 'Wiki', got %q", pages[0].Title)
	}
}

func TestGetWikiPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/wiki/Start.json" {
			t.Errorf("expected path /projects/test/wiki/Start.json, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "Start", "text": "Hello", "version": 1},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	page, err := client.GetWikiPage(context.Background(), "test", "Start", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "Start" {
		t.Errorf("expected 'Start', got %q", page.Title)
	}
	if page.Text != "Hello" {
		t.Errorf("expected 'Hello', got %q", page.Text)
	}
}

func TestGetWikiPage_WithVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("version") != "2" {
			t.Errorf("expected version=2, got %s", r.URL.Query().Get("version"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "Start", "version": 2},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetWikiPage(context.Background(), "test", "Start", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWikiPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "NewPage", "text": "Content", "version": 1},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	page, err := client.CreateWikiPage(context.Background(), "test", "NewPage", WikiPageCreateParams{Text: "Content"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "NewPage" {
		t.Errorf("expected 'NewPage', got %q", page.Title)
	}
}

func TestUpdateWikiPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateWikiPage(context.Background(), "test", "Start", WikiPageUpdateParams{Text: StringPtr("Updated")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteWikiPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteWikiPage(context.Background(), "test", "Start")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListWikiPages_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.ListWikiPages(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetWikiPage_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Not found"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.GetWikiPage(context.Background(), "test", "Missing", 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetWikiPage_TitleWithSpaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/wiki/My Page.json" {
			t.Errorf("expected path with decoded title, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "My Page", "text": "Content", "version": 1},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	page, err := client.GetWikiPage(context.Background(), "test", "My Page", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "My Page" {
		t.Errorf("expected 'My Page', got %q", page.Title)
	}
}

func TestCreateWikiPage_TitleWithSpaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test/wiki/New Page.json" {
			t.Errorf("expected path with decoded title, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wiki_page": map[string]interface{}{"title": "New Page", "text": "Content", "version": 1},
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	page, err := client.CreateWikiPage(context.Background(), "test", "New Page", WikiPageCreateParams{Text: "Content"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Title != "New Page" {
		t.Errorf("expected 'New Page', got %q", page.Title)
	}
}

func TestUpdateWikiPage_TitleWithSpaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/wiki/My Page.json" {
			t.Errorf("expected path with decoded title, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.UpdateWikiPage(context.Background(), "test", "My Page", WikiPageUpdateParams{Text: StringPtr("Updated")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteWikiPage_TitleWithSpaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test/wiki/My Page.json" {
			t.Errorf("expected path with decoded title, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	err := client.DeleteWikiPage(context.Background(), "test", "My Page")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWikiPage_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":["Invalid"]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "test-key")
	_, err := client.CreateWikiPage(context.Background(), "test", "Page", WikiPageCreateParams{Text: "Hello"})
	if err == nil {
		t.Fatal("expected error")
	}
}
