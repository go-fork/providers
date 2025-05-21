package formatter

import (
	"os"
	"testing"
)

func writeTempJSON(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestJsonFormatter_Load(t *testing.T) {
	file := writeTempJSON(t, `{"foo": {"bar": 1}, "baz": 2}`)
	defer os.Remove(file)

	f := NewJsonFormatter(file)
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if m["foo.bar"] != float64(1) {
		t.Errorf("Expected foo.bar=1, got %v", m["foo.bar"])
	}
	if m["baz"] != float64(2) {
		t.Errorf("Expected baz=2, got %v", m["baz"])
	}
}

func TestJsonFormatter_Load_FileNotFound(t *testing.T) {
	f := NewJsonFormatter("notfound.json")
	_, err := f.Load()
	if err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestJsonFormatter_Load_EmptyPath(t *testing.T) {
	f := NewJsonFormatter("")
	_, err := f.Load()
	if err == nil {
		t.Error("Expected error for empty path")
	}
}

func TestJsonFormatter_Load_Directory(t *testing.T) {
	// Tạo thư mục tạm
	tempDir, err := os.MkdirTemp("", "json-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	f := NewJsonFormatter(tempDir)
	_, err = f.Load()
	if err == nil {
		t.Error("Expected error when path is a directory")
	}
}

func TestJsonFormatter_Name(t *testing.T) {
	f := NewJsonFormatter("abc.json")
	if f.Name() != "json:abc.json" {
		t.Errorf("Name() = %v, want json:abc.json", f.Name())
	}
}

func TestJsonFormatter_Load_BadJSON(t *testing.T) {
	file := writeTempJSON(t, "{bad json}")
	defer os.Remove(file)
	f := NewJsonFormatter(file)
	_, err := f.Load()
	if err == nil {
		t.Error("Expected error for bad json")
	}
}

func TestJsonFormatter_Load_EmptyFile(t *testing.T) {
	file := writeTempJSON(t, "")
	defer os.Remove(file)

	f := NewJsonFormatter(file)
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error for empty file: %v", err)
	}

	if len(m) != 0 {
		t.Errorf("Expected empty map for empty file, got %v", m)
	}
}

func TestJsonFormatter_Load_EmptyObject(t *testing.T) {
	file := writeTempJSON(t, "{}")
	defer os.Remove(file)

	f := NewJsonFormatter(file)
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error for empty object: %v", err)
	}

	if len(m) != 0 {
		t.Errorf("Expected empty map for empty object, got %v", m)
	}
}

func TestJsonFormatter_Load_ArrayOrPrimitive(t *testing.T) {
	file := writeTempJSON(t, "[1,2,3]")
	defer os.Remove(file)
	f := NewJsonFormatter(file)
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error for array: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("Expected empty map for array root, got %v", m)
	}

	file2 := writeTempJSON(t, "123")
	defer os.Remove(file2)
	f2 := NewJsonFormatter(file2)
	m2, err := f2.Load()
	if err != nil {
		t.Fatalf("Load error for primitive: %v", err)
	}
	if len(m2) != 0 {
		t.Errorf("Expected empty map for primitive root, got %v", m2)
	}
}
