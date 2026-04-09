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

func TestFloat64Ptr(t *testing.T) {
	v := 3.14
	p := Float64Ptr(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %f, got %f", v, *p)
	}
}

func TestBoolPtr(t *testing.T) {
	v := true
	p := BoolPtr(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}
