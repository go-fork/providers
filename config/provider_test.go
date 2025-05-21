package config

import (
	"os"
	"testing"

	"github.com/go-fork/di"
)

type mockApp struct {
	container *di.Container
	basePath  string
}

func (a *mockApp) Container() *di.Container { return a.container }
func (a *mockApp) BasePath(path ...string) string {
	if len(path) > 0 && path[0] != "" {
		if a.basePath == "" {
			return ""
		}
		return a.basePath + "/" + path[0]
	}
	return a.basePath
}

func TestServiceProvider_Register_Basic(t *testing.T) {
	c := di.New()
	app := &mockApp{container: c, basePath: "notfound"}
	sp := NewServiceProvider()
	sp.Register(app)
	cfg, err := c.Make("config")
	if err != nil || cfg == nil {
		t.Errorf("config manager not registered in container, err=%v", err)
	}
}

func TestServiceProvider_Register_WithConfigFiles(t *testing.T) {
	dir := t.TempDir()
	// Tạo file cấu hình YAML hợp lệ
	file := dir + "/configs/a.yaml"

	// Tạo thư mục configs
	if err := os.MkdirAll(dir+"/configs", 0755); err != nil {
		t.Fatal(err)
	}

	content := "foo: 1\nbar: 2\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	c := di.New()
	app := &mockApp{container: c, basePath: dir}
	sp := NewServiceProvider()
	sp.Register(app)

	cfg, err := c.Make("config")
	if err != nil {
		t.Fatalf("config manager not found: %v", err)
	}

	mgr, ok := cfg.(Manager)
	if !ok {
		t.Fatal("config manager wrong type")
	}

	// Kiểm tra các giá trị đã nạp từ file YAML
	if mgr.GetInt("foo") != 1 || mgr.GetInt("bar") != 2 {
		t.Errorf("Config values not loaded from YAML: got foo=%v bar=%v", mgr.Get("foo"), mgr.Get("bar"))
	}
}

func TestServiceProvider_Boot(t *testing.T) {
	sp := NewServiceProvider()
	// Gọi Boot với nil và với mock app để đạt coverage 100%
	sp.Boot(nil)
	app := &mockApp{}
	sp.Boot(app)
}

func TestServiceProvider_Boot_Coverage(t *testing.T) {
	sp := NewServiceProvider()
	// Boot với nil
	sp.Boot(nil)
	// Boot với app bất kỳ
	sp.Boot(struct{}{})
}
