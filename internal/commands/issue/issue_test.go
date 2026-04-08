package issue

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCmdIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	f := newTestFactory(srv)
	cmd := NewCmdIssue(f)

	if cmd.Use != "issue" {
		t.Errorf("expected Use 'issue', got %q", cmd.Use)
	}

	// Verify all 6 subcommands exist
	expected := []string{"list", "view", "create", "update", "close", "delete"}
	for _, name := range expected {
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}
