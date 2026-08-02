package processor

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

// TestToMaps_StringRows (CSV batch). CSV sources emit [][]string where the
// first row is the header (column names) and subsequent rows are values.
// ToMaps must turn this into []map[string]any keyed by the header names, so
// processors (filter/redact/etc.) can address fields by name.
func TestToMaps_StringRows(t *testing.T) {
	in := [][]string{
		{"name", "ssn", "city"}, // header
		{"alice", "123", "NYC"},
		{"bob", "987", "LA"},
	}
	out, err := ToMaps(in)
	if err != nil {
		t.Fatalf("[][]string: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 data rows (header consumed), got %d", len(out))
	}
	if out[0]["name"] != "alice" || out[0]["ssn"] != "123" || out[0]["city"] != "NYC" {
		t.Errorf("row 0 = %v, want alice/123/NYC keyed by header", out[0])
	}
	if out[1]["name"] != "bob" {
		t.Errorf("row 1 name = %v, want bob", out[1]["name"])
	}
}

// A single-row [][]string with just a header and no data → empty result.
func TestToMaps_StringRows_HeaderOnly(t *testing.T) {
	in := [][]string{{"a", "b"}}
	out, err := ToMaps(in)
	if err != nil {
		t.Fatalf("header-only: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 data rows, got %d", len(out))
	}
}

// An empty [][]string → empty result, no error.
func TestToMaps_StringRows_Empty(t *testing.T) {
	out, err := ToMaps([][]string{})
	if err != nil || len(out) != 0 {
		t.Fatalf("empty: out=%v err=%v, want empty/no-error", out, err)
	}
}
