package formatter

import (
	"reflect"
	"testing"
)

func TestFlattenMapRecursive(t *testing.T) {
	input := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 2,
			"d": map[string]interface{}{"e": 3},
		},
	}
	want := map[string]interface{}{
		"a":     1,
		"b.c":   2,
		"b.d.e": 3,
	}
	out := make(map[string]interface{})
	flattenMapRecursive(out, input, "")
	if !reflect.DeepEqual(out, want) {
		t.Errorf("flattenMapRecursive failed, got %v, want %v", out, want)
	}
}

func TestFlattenMapRecursive_MapInterfaceInterface(t *testing.T) {
	out := make(map[string]interface{})
	input := map[string]interface{}{
		"foo": map[interface{}]interface{}{"bar": 1, "baz": 2},
	}
	FlattenMapRecursiveForTest(out, input, "")
	if out["foo.bar"] != 1 || out["foo.baz"] != 2 {
		t.Errorf("flattenMapRecursive map[interface{}]interface{} failed: %v", out)
	}
}

func TestFlattenMapRecursive_NilAndEmpty(t *testing.T) {
	FlattenMapRecursiveForTest(nil, nil, "prefix")
	FlattenMapRecursiveForTest(map[string]interface{}{}, nil, "prefix")
	FlattenMapRecursiveForTest(nil, map[string]interface{}{}, "prefix")
}

func TestFlattenMapRecursive_SkipEmptyKey(t *testing.T) {
	input := map[string]interface{}{"": 1, "x": 2}
	out := make(map[string]interface{})
	FlattenMapRecursiveForTest(out, input, "")
	if _, ok := out[""]; ok || out["x"] != 2 {
		t.Errorf("flattenMapRecursive skip empty key failed: %v", out)
	}
}
