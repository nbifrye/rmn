package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://redmine.example.com/", "test-key")
	if c.BaseURL != "https://redmine.example.com" {
		t.Errorf("expected trailing slash to be trimmed, got %s", c.BaseURL)
	}
	if c.APIKey != "test-key" {
		t.Errorf("expected APIKey to be test-key, got %s", c.APIKey)
	}
}

func TestGet_SetsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Redmine-API-Key") != "my-key" {
			t.Errorf("expected API key header, got %s", r.Header.Get("X-Redmine-API-Key"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "my-key")
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["Issue not found"]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Get(context.Background(), "/issues/999.json", nil, &result)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if err.Error() != "API error (status 404): Issue not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestGet_APIErrorUnstructured(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", nil, &result)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	if err.Error() != "API error (status 500): Internal Server Error" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestPost_SendsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["key"] != "value" {
			t.Errorf("expected body key=value, got %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Post(context.Background(), "/test.json", map[string]string{"key": "value"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPut_NoResponseBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.Put(context.Background(), "/test.json", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.Delete(context.Background(), "/test.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRequest_MarshalError(t *testing.T) {
	c := NewClient("http://example.com", "key")
	// Channels cannot be marshaled to JSON
	_, err := c.newRequest(context.Background(), http.MethodPost, "/test", make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable body")
	}
}

func TestGet_WithParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("foo") != "bar" {
			t.Errorf("expected query param foo=bar, got %s", r.URL.Query().Get("foo"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	params := make(map[string][]string)
	params["foo"] = []string{"bar"}
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", params, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGet_NetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close() // Close immediately to simulate network failure

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", nil, &result)
	if err == nil {
		t.Fatal("expected error for closed server")
	}
}

func TestPost_NetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Post(context.Background(), "/test.json", map[string]string{"key": "val"}, &result)
	if err == nil {
		t.Fatal("expected error for closed server")
	}
}

func TestPut_NetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.Put(context.Background(), "/test.json", map[string]string{"key": "val"})
	if err == nil {
		t.Fatal("expected error for closed server")
	}
}

func TestDelete_NetworkFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.Delete(context.Background(), "/test.json")
	if err == nil {
		t.Fatal("expected error for closed server")
	}
}

func TestDo_EmptyBody204(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	// Pass a non-nil result - should be OK with empty body
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRequest_NilBody(t *testing.T) {
	c := NewClient("http://example.com", "key")
	req, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Body != nil && req.Body != http.NoBody {
		t.Error("expected nil body for GET request")
	}
}

func TestNewRequest_InvalidURL(t *testing.T) {
	c := NewClient("://invalid", "key")
	_, err := c.newRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestGet_NewRequestError(t *testing.T) {
	c := NewClient("://invalid", "key")
	err := c.Get(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestPost_NewRequestError(t *testing.T) {
	c := NewClient("://invalid", "key")
	err := c.Post(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestPut_NewRequestError(t *testing.T) {
	c := NewClient("://invalid", "key")
	err := c.Put(context.Background(), "/test", nil)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestDelete_NewRequestError(t *testing.T) {
	c := NewClient("://invalid", "key")
	err := c.Delete(context.Background(), "/test")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) { return 0, fmt.Errorf("read error") }
func (e errReader) Close() error                { return nil }

type errBodyTransport struct{}

func (t errBodyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       errReader{},
		Header:     make(http.Header),
	}, nil
}

func TestDo_ReadAllError(t *testing.T) {
	c := NewClient("http://example.com", "key")
	c.HTTPClient = &http.Client{Transport: errBodyTransport{}}

	err := c.Get(context.Background(), "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error for body read failure")
	}
}

func TestDo_InvalidResponseJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	var result map[string]interface{}
	err := c.Get(context.Background(), "/test.json", nil, &result)
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
}
