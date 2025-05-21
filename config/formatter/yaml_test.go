package formatter

import (
	"os"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.yaml")
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

func TestYamlFormatter_Load(t *testing.T) {
	file := writeTempYAML(t, "foo:\n  bar: 1\nbaz: 2\n")
	defer os.Remove(file)
	f := NewYamlFormatter(file)
	m, err := f.Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if m["foo.bar"] != 1 {
		t.Errorf("Expected foo.bar=1, got %v", m["foo.bar"])
	}
	if m["baz"] != 2 {
		t.Errorf("Expected baz=2, got %v", m["baz"])
	}
}

func TestYamlFormatter_Load_FileNotFound(t *testing.T) {
	f := NewYamlFormatter("notfound.yaml")
	_, err := f.Load()
	if err == nil {
		t.Error("Expected error for missing file")
	}
}

func TestYamlFormatter_Name(t *testing.T) {
	f := NewYamlFormatter("abc.yaml")
	if f.Name() != "yaml:abc.yaml" {
		t.Errorf("Name() = %v, want yaml:abc.yaml", f.Name())
	}
}

func TestLoadFromDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.yaml", []byte("foo: 1\nbar: 2\n"), 0644)
	os.WriteFile(dir+"/b.yml", []byte("baz: 3\n"), 0644)
	m, err := LoadFromDirectory(dir)
	if err != nil {
		t.Fatalf("LoadFromDirectory error: %v", err)
	}
	if m["foo"] != 1 || m["bar"] != 2 || m["baz"] != 3 {
		t.Errorf("Unexpected result: %v", m)
	}
}

func TestLoadFromDirectory_NotExist(t *testing.T) {
	_, err := LoadFromDirectory("/not/exist/dir")
	if err == nil {
		t.Error("Expected error for not exist directory")
	}
}

func TestLoadFromDirectory_BadYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/bad.yaml", []byte(":bad yaml"), 0644)
	_, err := LoadFromDirectory(dir)
	if err == nil {
		t.Error("Expected error for bad yaml file")
	}
}

func TestLoadFromDirectory_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	m, err := LoadFromDirectory(dir)
	if err != nil {
		t.Fatalf("Empty dir should not error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("Expected empty map for empty dir, got %v", m)
	}
}

func TestLoadFromDirectory_OnlyNonYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/a.txt", []byte("foo: 1"), 0644)
	os.WriteFile(dir+"/b.json", []byte("{}"), 0644)
	m, err := LoadFromDirectory(dir)
	if err != nil {
		t.Fatalf("Non-YAML files should not error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("Expected empty map for only non-YAML files, got %v", m)
	}
}

func TestLoadFromDirectory_MixValidAndInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(dir+"/good.yaml", []byte("foo: 1\n"), 0644)
	os.WriteFile(dir+"/bad.yaml", []byte(":bad yaml"), 0644)
	_, err := LoadFromDirectory(dir)
	if err == nil {
		t.Error("Expected error for directory with invalid YAML file")
	}
}
