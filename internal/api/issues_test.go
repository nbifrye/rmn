package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issues.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("project_id") != "my-project" {
			t.Errorf("expected project_id=my-project, got %s", r.URL.Query().Get("project_id"))
		}
		if r.URL.Query().Get("status_id") != "open" {
			t.Errorf("expected status_id=open, got %s", r.URL.Query().Get("status_id"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		resp := issuesResponse{
			Issues: []Issue{
				{ID: 1, Subject: "Test issue"},
				{ID: 2, Subject: "Another issue"},
			},
			TotalCount: 2,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	issues, total, err := c.ListIssues(context.Background(), IssueListParams{
		ProjectID: "my-project",
		StatusID:  "open",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Subject != "Test issue" {
		t.Errorf("unexpected subject: %s", issues[0].Subject)
	}
}

func TestGetIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := issueResponse{
			Issue: Issue{ID: 42, Subject: "Found issue"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	issue, err := c.GetIssue(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.ID != 42 {
		t.Errorf("expected ID 42, got %d", issue.ID)
	}
	if issue.Subject != "Found issue" {
		t.Errorf("unexpected subject: %s", issue.Subject)
	}
}

func TestCreateIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body issueCreateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.Subject != "New issue" {
			t.Errorf("unexpected subject: %s", body.Issue.Subject)
		}
		w.WriteHeader(http.StatusCreated)
		resp := issueResponse{
			Issue: Issue{ID: 100, Subject: "New issue"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	issue, err := c.CreateIssue(context.Background(), IssueCreateParams{
		ProjectID: 1,
		Subject:   "New issue",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.ID != 100 {
		t.Errorf("expected ID 100, got %d", issue.ID)
	}
}

func TestUpdateIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body issueUpdateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Issue.Notes != "Updated via test" {
			t.Errorf("unexpected notes: %s", body.Issue.Notes)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.UpdateIssue(context.Background(), 42, IssueUpdateParams{
		Notes: "Updated via test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/issues/42.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	err := c.DeleteIssue(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListIssues_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := issuesResponse{
			Issues:     []Issue{},
			TotalCount: 0,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	issues, total, err := c.ListIssues(context.Background(), IssueListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestCreateIssue_WithStringProjectID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		issue := body["issue"].(map[string]interface{})
		if issue["project_id"] != "my-project" {
			t.Errorf("expected string project_id, got %v", issue["project_id"])
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(issueResponse{Issue: Issue{ID: 1, Subject: "test"}})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "key")
	_, err := c.CreateIssue(context.Background(), IssueCreateParams{
		ProjectID: "my-project",
		Subject:   "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
