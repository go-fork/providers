package utils

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
	FlattenMapRecursive(out, input, "")
	if out["foo.bar"] != 1 || out["foo.baz"] != 2 {
		t.Errorf("flattenMapRecursive map[interface{}]interface{} failed: %v", out)
	}
}

func TestFlattenMapRecursive_NilAndEmpty(t *testing.T) {
	FlattenMapRecursive(nil, nil, "prefix")
	FlattenMapRecursive(map[string]interface{}{}, nil, "prefix")
	FlattenMapRecursive(nil, map[string]interface{}{}, "prefix")
}

func TestFlattenMapRecursive_SkipEmptyKey(t *testing.T) {
	input := map[string]interface{}{"": 1, "x": 2}
	out := make(map[string]interface{})
	FlattenMapRecursive(out, input, "")
	if _, ok := out[""]; ok || out["x"] != 2 {
		t.Errorf("flattenMapRecursive skip empty key failed: %v", out)
	}
}

func TestUnflattenDotMapRecursive(t *testing.T) {
	// Trường hợp phẳng đơn giản
	flat := map[string]interface{}{
		"a.b": 1,
		"a.c": 2,
		"d":   3,
	}
	nested := UnflattenDotMapRecursive(flat)
	a, ok := nested["a"].(map[string]interface{})
	if !ok || a["b"] != 1 || a["c"] != 2 {
		t.Errorf("UnflattenDotMapRecursive failed for simple case: %v", nested)
	}
	if nested["d"] != 3 {
		t.Errorf("UnflattenDotMapRecursive failed for root key: %v", nested)
	}

	// Trường hợp lồng nhiều cấp
	flat2 := map[string]interface{}{
		"x.y.z": 5,
		"x.y.t": 6,
	}
	nested2 := UnflattenDotMapRecursive(flat2)
	x, ok := nested2["x"].(map[string]interface{})
	if !ok {
		t.Fatalf("UnflattenDotMapRecursive failed for nested: %v", nested2)
	}
	y, ok := x["y"].(map[string]interface{})
	if !ok || y["z"] != 5 || y["t"] != 6 {
		t.Errorf("UnflattenDotMapRecursive failed for deep nested: %v", nested2)
	}

	// Trường hợp empty input
	flat3 := map[string]interface{}{}
	nested3 := UnflattenDotMapRecursive(flat3)
	if len(nested3) != 0 {
		t.Errorf("UnflattenDotMapRecursive failed for empty input: %v", nested3)
	}

	// Trường hợp không có dot
	flat4 := map[string]interface{}{"foo": 42}
	nested4 := UnflattenDotMapRecursive(flat4)
	if !reflect.DeepEqual(nested4, flat4) {
		t.Errorf("UnflattenDotMapRecursive failed for no dot: %v", nested4)
	}

	// Trường hợp lồng 3 cấp
	flat5 := map[string]interface{}{
		"a.b.c": 7,
		"a.b.d": 8,
	}
	nested5 := UnflattenDotMapRecursive(flat5)
	a5, ok := nested5["a"].(map[string]interface{})
	if !ok {
		t.Fatalf("UnflattenDotMapRecursive failed for 3-level: %v", nested5)
	}
	b5, ok := a5["b"].(map[string]interface{})
	if !ok || b5["c"] != 7 || b5["d"] != 8 {
		t.Errorf("UnflattenDotMapRecursive failed for 3-level: %v", nested5)
	}
}
