package formatter

import (
	"testing"
)

func TestEnvFormatter_Load(t *testing.T) {
	t.Setenv("APP_FOO", "bar")
	t.Setenv("APP_BAR_BAZ", "qux")
	t.Setenv("APP_NUMBER", "123")
	t.Setenv("APP_FLOAT", "12.34")
	t.Setenv("APP_BOOL_TRUE", "true")
	t.Setenv("APP_BOOL_FALSE", "false")
	t.Setenv("OTHER", "skip")

	f := NewEnvFormatter("APP")
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// Kiểm tra string
	if m["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", m["foo"])
	}
	if m["bar.baz"] != "qux" {
		t.Errorf("Expected bar.baz=qux, got %v", m["bar.baz"])
	}

	// Kiểm tra bỏ qua keys không có prefix
	if _, ok := m["other"]; ok {
		t.Error("Should not include keys without prefix")
	}

	// Kiểm tra chuyển đổi số
	number, ok := m["number"].(int)
	if !ok || number != 123 {
		t.Errorf("Expected number=123 (int), got %v (type %T)", m["number"], m["number"])
	}

	// Kiểm tra chuyển đổi boolean
	boolTrue, ok := m["bool.true"].(bool)
	if !ok || !boolTrue {
		t.Errorf("Expected bool.true=true (bool), got %v (type %T)", m["bool.true"], m["bool.true"])
	}

	boolFalse, ok := m["bool.false"].(bool)
	if !ok || boolFalse {
		t.Errorf("Expected bool.false=false (bool), got %v (type %T)", m["bool.false"], m["bool.false"])
	}
}

func TestEnvFormatter_Load_NoMatch(t *testing.T) {
	t.Setenv("APP_FOO", "")
	t.Setenv("APP_BAR_BAZ", "")
	t.Setenv("OTHER", "value")

	f := NewEnvFormatter("APP")
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if len(m) > 0 {
		t.Errorf("Expected empty map, got %v keys: %v", len(m), m)
	}
}

func TestEnvFormatter_Load_EmptyPrefix(t *testing.T) {
	t.Setenv("FOO", "bar")

	// Formatter không nên nạp biến môi trường khi không có prefix để tránh xung đột
	f := NewEnvFormatter("")
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if len(m) > 0 {
		t.Errorf("Expected empty map for empty prefix, got %v", m)
	}
}

func TestEnvFormatter_Name(t *testing.T) {
	f := NewEnvFormatter("APP")
	if f.Name() != "env" {
		t.Errorf("Name() = %v, want env", f.Name())
	}
}
