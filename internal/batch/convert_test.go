package batch

import (
	"testing"
)

func TestToMaps_NativeSlice(t *testing.T) {
	in := []map[string]any{{"a": 1}, {"b": 2}}
	out, err := ToMaps(in)
	if err != nil || len(out) != 2 || out[0]["a"] != 1 {
		t.Fatalf("native passthrough: %v %v", out, err)
	}
}

func TestToMaps_SingleRow(t *testing.T) {
	in := map[string]any{"x": "y"}
	out, err := ToMaps(in)
	if err != nil || len(out) != 1 || out[0]["x"] != "y" {
		t.Fatalf("single row wrap: %v %v", out, err)
	}
}

func TestToMaps_StringMap(t *testing.T) {
	in := map[string]string{"name": "alice"}
	out, err := ToMaps(in)
	if err != nil || len(out) != 1 || out[0]["name"] != "alice" {
		t.Fatalf("string map widen: %v %v", out, err)
	}
}

func TestToMaps_StringMapSlice(t *testing.T) {
	in := []map[string]string{{"a": "1"}, {"a": "2"}}
	out, err := ToMaps(in)
	if err != nil || len(out) != 2 || out[1]["a"] != "2" {
		t.Fatalf("string map slice widen: %v %v", out, err)
	}
}

func TestToMaps_Unsupported(t *testing.T) {
	_, err := ToMaps(42)
	if err == nil {
		t.Fatal("expected error for int")
	}
}

func TestToMaps_NilBatch(t *testing.T) {
	_, err := ToMaps(nil)
	if err == nil {
		t.Fatal("expected error for nil")
	}
}
