package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-fork/providers/config/formatter"
)

// mockFormatter implements formatter.Formatter for testing
// Allows simulating different config sources

type mockFormatter struct {
	values map[string]interface{}
	fail   bool
}

func (m *mockFormatter) Load() (map[string]interface{}, error) {
	if m.fail {
		return nil, errors.New("load error")
	}
	return m.values, nil
}
func (m *mockFormatter) Name() string { return "mock" }

func TestNewManager(t *testing.T) {
	mgr := NewManager()
	if mgr == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestManagerSetGetHas(t *testing.T) {
	mgr := NewManager()
	mgr.Set("foo", 123)
	if v := mgr.Get("foo"); v != 123 {
		t.Errorf("Get failed, got %v", v)
	}
	if !mgr.Has("foo") {
		t.Error("Has should return true for existing key")
	}
	if mgr.Has("bar") {
		t.Error("Has should return false for missing key")
	}
}

func TestManagerGetString(t *testing.T) {
	mgr := NewManager()
	mgr.Set("str", "abc")
	if got := mgr.GetString("str"); got != "abc" {
		t.Errorf("GetString failed, got %v", got)
	}
	if got := mgr.GetString("notfound", "def"); got != "def" {
		t.Errorf("GetString default failed, got %v", got)
	}
}

func TestManagerGetString_TypeMismatch(t *testing.T) {
	mgr := NewManager()
	mgr.Set("num", 123)
	if got := mgr.GetString("num", "x"); got != "x" {
		t.Errorf("GetString type mismatch, want default, got %v", got)
	}
}

type myStringer struct{}

func (myStringer) String() string { return "stringer" }

func TestManagerGetString_TypeStringer(t *testing.T) {
	mgr := NewManager()
	mgr.Set("s", myStringer{})
	if got := mgr.GetString("s"); got != "stringer" {
		t.Errorf("GetString with Stringer failed, got %v", got)
	}
}

func TestManagerGetInt(t *testing.T) {
	mgr := NewManager()
	mgr.Set("i1", 42)
	mgr.Set("i2", int64(43))
	mgr.Set("i3", float64(44))
	if mgr.GetInt("i1") != 42 || mgr.GetInt("i2") != 43 || mgr.GetInt("i3") != 44 {
		t.Error("GetInt type conversion failed")
	}
	if mgr.GetInt("notfound", 99) != 99 {
		t.Error("GetInt default failed")
	}
}

func TestManagerGetInt_TypeMismatch(t *testing.T) {
	mgr := NewManager()
	mgr.Set("str", "abc")
	if got := mgr.GetInt("str", 7); got != 7 {
		t.Errorf("GetInt type mismatch, want default, got %v", got)
	}
}

func TestManagerGetInt_StringParseFail(t *testing.T) {
	mgr := NewManager()
	mgr.Set("badint", "abc")
	if got := mgr.GetInt("badint", 9); got != 9 {
		t.Errorf("GetInt with bad string should return default, got %v", got)
	}
}

func TestManagerGetInt_Branches(t *testing.T) {
	mgr := NewManager()
	// int32
	mgr.Set("i32", int32(10))
	if mgr.GetInt("i32") != 10 {
		t.Errorf("GetInt int32 failed")
	}
	// float32
	mgr.Set("f32", float32(11.9))
	if mgr.GetInt("f32") != 11 {
		t.Errorf("GetInt float32 failed")
	}
	// string parse fail, no default
	mgr.Set("badstr", "abc")
	if mgr.GetInt("badstr") != 0 {
		t.Errorf("GetInt bad string, no default, should return 0")
	}
}

func TestManagerGetBool(t *testing.T) {
	mgr := NewManager()
	mgr.Set("b1", true)
	if !mgr.GetBool("b1") {
		t.Error("GetBool failed")
	}
	if mgr.GetBool("notfound", true) != true {
		t.Error("GetBool default failed")
	}
}

func TestManagerGetBool_TypeMismatch(t *testing.T) {
	mgr := NewManager()
	mgr.Set("str", "abc")
	if got := mgr.GetBool("str", true); got != true {
		t.Errorf("GetBool type mismatch, want default, got %v", got)
	}
}

func TestManagerGetBool_StringCases(t *testing.T) {
	mgr := NewManager()
	mgr.Set("btrue", "yes")
	mgr.Set("bfalse", "no")
	if !mgr.GetBool("btrue") {
		t.Errorf("GetBool with 'yes' should be true")
	}
	if mgr.GetBool("bfalse") {
		t.Errorf("GetBool with 'no' should be false")
	}
}

func TestManagerGetStringMap(t *testing.T) {
	mgr := NewManager()
	m := map[string]interface{}{"a": 1}
	mgr.Set("map", m)
	if !reflect.DeepEqual(mgr.GetStringMap("map"), m) {
		t.Error("GetStringMap failed")
	}
	if got := mgr.GetStringMap("notfound"); len(got) != 0 {
		t.Error("GetStringMap default failed")
	}
}

func TestManagerGetStringMap_ComplexCases(t *testing.T) {
	mgr := NewManager()
	// map[interface{}]interface{}
	m := map[interface{}]interface{}{"a": 1, "b": 2}
	mgr.Set("miface", m)
	got := mgr.GetStringMap("miface")
	if got["a"] != 1 || got["b"] != 2 {
		t.Errorf("GetStringMap with map[interface{}]interface{} failed: %v", got)
	}
	// JSON marshalling fallback
	mgr.Set("jsonmap", struct{ X int }{X: 5})
	got2 := mgr.GetStringMap("jsonmap")
	if got2["X"] != float64(5) {
		t.Errorf("GetStringMap with struct fallback failed: %v", got2)
	}
}

func TestManagerGetStringMap_JSONFallbackFail(t *testing.T) {
	mgr := NewManager()
	mgr.Set("bad", func() {})
	m := mgr.GetStringMap("bad")
	if len(m) != 0 {
		t.Errorf("GetStringMap with marshal error should return empty map, got %v", m)
	}
}

func TestManagerGetStringSlice(t *testing.T) {
	mgr := NewManager()
	s := []string{"a", "b"}
	mgr.Set("slice", s)
	if !reflect.DeepEqual(mgr.GetStringSlice("slice"), s) {
		t.Error("GetStringSlice failed")
	}
	mgr.Set("iface", []interface{}{"x", "y"})
	if got := mgr.GetStringSlice("iface"); !reflect.DeepEqual(got, []string{"x", "y"}) {
		t.Errorf("GetStringSlice iface failed, got %v", got)
	}
	if got := mgr.GetStringSlice("notfound"); len(got) != 0 {
		t.Error("GetStringSlice default failed")
	}
}

func TestManagerGetStringSlice_ComplexCases(t *testing.T) {
	mgr := NewManager()
	// []interface{} with mixed types
	mgr.Set("mix", []interface{}{"a", 2, 3.5})
	got := mgr.GetStringSlice("mix")
	if len(got) != 3 || got[0] != "a" || got[1] != "2" || got[2] != "3.5" {
		t.Errorf("GetStringSlice with []interface{} failed: %v", got)
	}
	// string as JSON array
	mgr.Set("jsonarr", `["x","y"]`)
	got2 := mgr.GetStringSlice("jsonarr")
	if len(got2) != 2 || got2[0] != "x" || got2[1] != "y" {
		t.Errorf("GetStringSlice with JSON array string failed: %v", got2)
	}
	// string as plain value
	mgr.Set("plain", "z")
	got3 := mgr.GetStringSlice("plain")
	if len(got3) != 1 || got3[0] != "z" {
		t.Errorf("GetStringSlice with plain string failed: %v", got3)
	}
	// struct fallback
	mgr.Set("struct", struct{ Y int }{Y: 7})
	got4 := mgr.GetStringSlice("struct")
	if len(got4) == 0 || got4[0] == "" {
		t.Errorf("GetStringSlice with struct fallback failed: %v", got4)
	}
}

func TestManagerGetStringSlice_JSONFallbackFail(t *testing.T) {
	mgr := NewManager()
	mgr.Set("bad", func() {})
	s := mgr.GetStringSlice("bad")
	if len(s) != 0 {
		t.Errorf("GetStringSlice with marshal error should return empty slice, got %v", s)
	}
}

func TestManagerGetStringSlice_BytesAndStruct(t *testing.T) {
	mgr := NewManager()
	// []byte chứa JSON array
	mgr.Set("bytesarr", []byte(`["a","b"]`))
	arr := mgr.GetStringSlice("bytesarr")
	if len(arr) != 2 || arr[0] != "a" || arr[1] != "b" {
		t.Errorf("GetStringSlice with []byte JSON array failed: %v", arr)
	}
	// []byte chứa JSON interface
	mgr.Set("bytesiface", []byte(`[1,2]`))
	arr2 := mgr.GetStringSlice("bytesiface")
	if len(arr2) != 2 || arr2[0] != "1" || arr2[1] != "2" {
		t.Errorf("GetStringSlice with []byte JSON interface failed: %v", arr2)
	}
	// []byte không hợp lệ
	mgr.Set("bytesbad", []byte(`notjson`))
	arr3 := mgr.GetStringSlice("bytesbad")
	if len(arr3) != 1 || arr3[0] != "[110 111 116 106 115 111 110]" {
		t.Errorf("GetStringSlice with []byte not json failed: %v", arr3)
	}
	// struct lồng nhau
	type Inner struct{ Z int }
	type Outer struct{ Y Inner }
	mgr.Set("nest", Outer{Y: Inner{Z: 9}})
	arr4 := mgr.GetStringSlice("nest")
	if len(arr4) == 0 || arr4[0] == "" {
		t.Errorf("GetStringSlice with nested struct failed: %v", arr4)
	}
}

func TestManagerGetInt_Float64Decimal(t *testing.T) {
	mgr := NewManager()
	mgr.Set("f64", 12.99)
	if mgr.GetInt("f64") != 12 {
		t.Errorf("GetInt float64 with decimal failed")
	}
}

func TestManagerUnmarshal(t *testing.T) {
	type S struct{ X int }
	mgr := NewManager()
	mgr.Set("s", map[string]interface{}{"X": 5})
	var s S
	if err := mgr.Unmarshal("s", &s); err != nil || s.X != 5 {
		t.Errorf("Unmarshal failed: %v, s=%+v", err, s)
	}
	if err := mgr.Unmarshal("notfound", &s); err == nil {
		t.Error("Unmarshal should fail for missing key")
	}
}

func TestManagerUnmarshal_MarshalError(t *testing.T) {
	mgr := NewManager()
	mgr.Set("bad", func() {}) // func cannot be marshaled
	var v interface{}
	err := mgr.Unmarshal("bad", &v)
	if err == nil {
		t.Error("Unmarshal should fail for marshal error")
	}
}

func TestManagerLoad(t *testing.T) {
	mgr := NewManager()
	f := &mockFormatter{values: map[string]interface{}{"a": 1}}
	if err := mgr.Load(f); err != nil {
		t.Errorf("Load failed: %v", err)
	}
	if mgr.Get("a") != 1 {
		t.Error("Load did not set value")
	}
	f2 := &mockFormatter{fail: true}
	if err := mgr.Load(f2); err == nil {
		t.Error("Load should fail on formatter error")
	}
}

func TestManagerLoad_NilFormatter(t *testing.T) {
	mgr := NewManager()
	err := mgr.Load(nil)
	if err == nil {
		t.Error("Load should fail for nil formatter")
	}
}

func TestManager_Get_DefaultValueVariants(t *testing.T) {
	mgr := NewManager()
	// Không có key, không truyền default
	if v := mgr.Get("notfound"); v != nil {
		t.Errorf("Get without default should return nil, got %v", v)
	}
	// Không có key, truyền nhiều default
	if v := mgr.Get("notfound", 1, 2, 3); v != 1 {
		t.Errorf("Get with multiple defaults should return first, got %v", v)
	}
}

func TestManager_GetString_DefaultValueVariants(t *testing.T) {
	mgr := NewManager()
	// Không có key, không truyền default
	if v := mgr.GetString("notfound"); v != "" {
		t.Errorf("GetString without default should return empty string, got %v", v)
	}
	// Không có key, truyền nhiều default
	if v := mgr.GetString("notfound", "a", "b"); v != "a" {
		t.Errorf("GetString with multiple defaults should return first, got %v", v)
	}
}

func TestManager_GetInt_DefaultValueVariants(t *testing.T) {
	mgr := NewManager()
	// Không có key, không truyền default
	if v := mgr.GetInt("notfound"); v != 0 {
		t.Errorf("GetInt without default should return 0, got %v", v)
	}
	// Không có key, truyền nhiều default
	if v := mgr.GetInt("notfound", 7, 8); v != 7 {
		t.Errorf("GetInt with multiple defaults should return first, got %v", v)
	}
}

func TestManager_GetBool_DefaultValueVariants(t *testing.T) {
	mgr := NewManager()
	// Không có key, không truyền default
	if v := mgr.GetBool("notfound"); v != false {
		t.Errorf("GetBool without default should return false, got %v", v)
	}
	// Không có key, truyền nhiều default
	if v := mgr.GetBool("notfound", true, false); v != true {
		t.Errorf("GetBool with multiple defaults should return first, got %v", v)
	}
}

func TestManager_Unmarshal_NonMapValue(t *testing.T) {
	mgr := NewManager()
	mgr.Set("foo", 123)
	var v int
	err := mgr.Unmarshal("foo", &v)
	if err != nil {
		t.Errorf("Unmarshal should work for primitive types, got error: %v", err)
	}
}

func TestManager_Unmarshal_Errors(t *testing.T) {
	mgr := NewManager()
	// out == nil
	err := mgr.Unmarshal("foo", nil)
	if err == nil || err.Error() != "output pointer cannot be nil" {
		t.Errorf("Unmarshal nil out: want error, got %v", err)
	}
	// out không phải pointer
	var notPtr int
	err = mgr.Unmarshal("foo", notPtr)
	if err == nil || err.Error() != "output must be a non-nil pointer" {
		t.Errorf("Unmarshal non-pointer: want error, got %v", err)
	}
	// out là nil pointer
	var nilPtr *int
	err = mgr.Unmarshal("foo", nilPtr)
	if err == nil || err.Error() != "output must be a non-nil pointer" {
		t.Errorf("Unmarshal nil pointer: want error, got %v", err)
	}
	// key không tồn tại
	var v int
	err = mgr.Unmarshal("notfound", &v)
	if err == nil || err.Error() != "key 'notfound' not found in configuration" {
		t.Errorf("Unmarshal missing key: want error, got %v", err)
	}
	// marshal error
	mgr.Set("bad", func() {})
	err = mgr.Unmarshal("bad", &v)
	if err == nil || err.Error() != "failed to marshal configuration: json: unsupported type: func()" {
		t.Errorf("Unmarshal marshal error: want error, got %v", err)
	}
}

func TestFlattenMapRecursive_AllBranches(t *testing.T) {
	// nil result or nested
	call := formatter.FlattenMapRecursiveForTest
	call(nil, nil, "prefix")
	call(map[string]interface{}{}, nil, "prefix")
	call(nil, map[string]interface{}{}, "prefix")
	// map[interface{}]interface{} branch
	input := map[string]interface{}{
		"a": map[interface{}]interface{}{"b": 2},
	}
	out := make(map[string]interface{})
	call(out, input, "")
	if out["a.b"] != 2 {
		t.Errorf("flattenMapRecursive with map[interface{}]interface{} failed: %v", out)
	}
	// skip empty key
	input2 := map[string]interface{}{"": 1, "x": 2}
	out2 := make(map[string]interface{})
	call(out2, input2, "")
	if _, ok := out2[""]; ok || out2["x"] != 2 {
		t.Errorf("flattenMapRecursive skip empty key failed: %v", out2)
	}
}

func TestManager_Get_EdgeCases(t *testing.T) {
	mgr := NewManager()
	// key rỗng, defaultValue
	if v := mgr.Get(""); v != nil {
		t.Errorf("Get empty key without default should return nil, got %v", v)
	}
	if v := mgr.Get("", 123); v != 123 {
		t.Errorf("Get empty key with default should return default, got %v", v)
	}
}

func TestManager_Has_EmptyKey(t *testing.T) {
	mgr := NewManager()
	if mgr.Has("") {
		t.Errorf("Has empty key should return false")
	}
}
