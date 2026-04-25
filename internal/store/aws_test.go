package store

import (
	"testing"
)

func TestNewAWSStore_CopiesVars(t *testing.T) {
	src := map[string]string{"/app/prod/KEY": "val"}
	s, err := NewAWSStore("/app/prod", src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src["/app/prod/KEY"] = "mutated"
	v, _ := s.Get("KEY")
	if v != "val" {
		t.Errorf("expected original value, got %q", v)
	}
}

func TestAWSStore_EmptyNamespace(t *testing.T) {
	_, err := NewAWSStore("", nil)
	if err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestAWSStore_Namespace(t *testing.T) {
	s, _ := NewAWSStore("/app/prod", nil)
	if s.Namespace() != "/app/prod" {
		t.Errorf("unexpected namespace: %s", s.Namespace())
	}
}

func TestAWSStore_GetMissingKey(t *testing.T) {
	s, _ := NewAWSStore("/app/prod", nil)
	_, err := s.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAWSStore_SetAndGet(t *testing.T) {
	s, _ := NewAWSStore("/app/prod", nil)
	if err := s.Set("DB_HOST", "localhost"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	v, err := s.Get("DB_HOST")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v != "localhost" {
		t.Errorf("expected %q, got %q", "localhost", v)
	}
}

func TestAWSStore_SetEmptyKey(t *testing.T) {
	s, _ := NewAWSStore("/app/prod", nil)
	if err := s.Set("", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestAWSStore_Delete(t *testing.T) {
	s, _ := NewAWSStore("/app/prod", nil)
	_ = s.Set("TO_DELETE", "bye")
	_ = s.Delete("TO_DELETE")
	_, err := s.Get("TO_DELETE")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestAWSStore_ToMap(t *testing.T) {
	initial := map[string]string{
		"/ns/A": "1",
		"/ns/B": "2",
	}
	s, _ := NewAWSStore("/ns", initial)
	m := s.ToMap()
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
	if _, ok := m["/ns/A"]; ok {
		t.Error("full key should be stripped from ToMap output")
	}
}
