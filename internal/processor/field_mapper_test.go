package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newMapper(t *testing.T, mappings string) *FieldMapper {
	t.Helper()
	m := &FieldMapper{}
	// Per spec §5.2, "mappings" is a JSON-encoded string field. Marshal the
	// config so the inner quotes in `mappings` are properly escaped.
	cfg, err := json.Marshal(struct {
		Mappings string `json:"mappings"`
	}{Mappings: mappings})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := m.Init(context.Background(), cfg); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return m
}

func TestRename(t *testing.T) {
	m := newMapper(t, `[{"from":"old_name","to":"new_name","action":"rename"}]`)
	in := []map[string]any{{"old_name": "alice", "age": int64(30)}}
	out, err := m.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if _, hasOld := rows[0]["old_name"]; hasOld {
		t.Error("old_name should be gone")
	}
	if rows[0]["new_name"] != "alice" {
		t.Errorf("new_name should be alice, got %v", rows[0]["new_name"])
	}
	if rows[0]["age"] != int64(30) {
		t.Error("age should be untouched")
	}
}

func TestDrop(t *testing.T) {
	m := newMapper(t, `[{"from":"temp","action":"drop"}]`)
	in := []map[string]any{{"temp": "x", "keep": "y"}}
	out, _ := m.Process(in)
	rows := out.([]map[string]any)
	if _, has := rows[0]["temp"]; has {
		t.Error("temp should be dropped")
	}
	if rows[0]["keep"] != "y" {
		t.Error("keep should survive")
	}
}

func TestAdd(t *testing.T) {
	m := newMapper(t, `[{"to":"source","action":"add","value":"planx"}]`)
	in := []map[string]any{{"name": "alice"}}
	out, _ := m.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["source"] != "planx" {
		t.Errorf("source should be planx, got %v", rows[0]["source"])
	}
	if rows[0]["name"] != "alice" {
		t.Error("name should survive")
	}
}

func TestMultipleMappings(t *testing.T) {
	m := newMapper(t, `[{"from":"a","to":"alpha","action":"rename"},{"from":"b","action":"drop"},{"to":"tag","action":"add","value":"v1"}]`)
	in := []map[string]any{{"a": 1, "b": 2, "c": 3}}
	out, _ := m.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["alpha"] != 1 {
		t.Error("a→alpha rename failed")
	}
	if _, has := rows[0]["a"]; has {
		t.Error("a should be gone")
	}
	if _, has := rows[0]["b"]; has {
		t.Error("b should be dropped")
	}
	if rows[0]["c"] != 3 {
		t.Error("c should survive")
	}
	if rows[0]["tag"] != "v1" {
		t.Error("tag add failed")
	}
}

func TestMissingSourceKey(t *testing.T) {
	m := newMapper(t, `[{"from":"nonexistent","to":"x","action":"rename"}]`)
	in := []map[string]any{{"other": "y"}}
	out, err := m.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if _, has := rows[0]["x"]; has {
		t.Error("rename of missing key should not create target")
	}
}

func TestFieldMapper_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*FieldMapper)(nil)
}
