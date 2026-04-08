package api

import "testing"

func TestIntPtr(t *testing.T) {
	v := 42
	p := IntPtr(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %d, got %d", v, *p)
	}
}

func TestStringPtr(t *testing.T) {
	v := "hello"
	p := StringPtr(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %q, got %q", v, *p)
	}
}
