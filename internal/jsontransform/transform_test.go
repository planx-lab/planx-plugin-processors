package jsontransform

import (
	"context"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newTransform(t *testing.T, ops string) *JSONTransform {
	t.Helper()
	tr := &JSONTransform{}
	cfg := `{"operations":"` + escapeJSON(ops) + `"}`
	if err := tr.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return tr
}

// escapeJSON properly escapes inner quotes for embedding in a JSON string value
func escapeJSON(s string) string {
	out := make([]byte, 0, len(s))
	for _, b := range []byte(s) {
		if b == '"' {
			out = append(out, '\\', '"')
		} else if b == '\\' {
			out = append(out, '\\', '\\')
		} else {
			out = append(out, b)
		}
	}
	return string(out)
}

func TestExtract(t *testing.T) {
	tr := newTransform(t, `[{"op":"extract","path":"data.name","to":"name"}]`)
	in := []map[string]any{
		{"data": map[string]any{"name": "alice", "age": float64(30)}},
	}
	out, err := tr.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if rows[0]["name"] != "alice" {
		t.Errorf("expected name=alice, got %v", rows[0]["name"])
	}
	// data should still be there (extract is copy, not move)
	if rows[0]["data"] == nil {
		t.Error("data should still exist after extract")
	}
}

func TestFlatten(t *testing.T) {
	tr := newTransform(t, `[{"op":"flatten","prefix":"addr."}]`)
	in := []map[string]any{
		{"addr.street": "123 Main St", "addr.city": "NYC", "name": "alice"},
	}
	out, _ := tr.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["street"] != "123 Main St" {
		t.Errorf("expected street=123 Main St, got %v", rows[0]["street"])
	}
	if rows[0]["city"] != "NYC" {
		t.Errorf("expected city=NYC, got %v", rows[0]["city"])
	}
	if _, has := rows[0]["addr.street"]; has {
		t.Error("addr.street should be gone after flatten")
	}
}

func TestRemove(t *testing.T) {
	tr := newTransform(t, `[{"op":"remove","path":"temp"}]`)
	in := []map[string]any{{"temp": "x", "keep": "y"}}
	out, _ := tr.Process(in)
	rows := out.([]map[string]any)
	if _, has := rows[0]["temp"]; has {
		t.Error("temp should be removed")
	}
	if rows[0]["keep"] != "y" {
		t.Error("keep should survive")
	}
}

func TestExtractNestedMissing(t *testing.T) {
	tr := newTransform(t, `[{"op":"extract","path":"data.nonexistent","to":"result"}]`)
	in := []map[string]any{{"data": map[string]any{"name": "alice"}}}
	out, err := tr.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if _, has := rows[0]["result"]; has {
		t.Error("result should not exist for missing path")
	}
}

func TestMultipleOps(t *testing.T) {
	tr := newTransform(t, `[{"op":"extract","path":"user.name","to":"name"},{"op":"remove","path":"user"},{"op":"flatten","prefix":"meta."}]`)
	in := []map[string]any{
		{"user": map[string]any{"name": "bob", "pass": "x"}, "meta.tag": "v1"},
	}
	out, _ := tr.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["name"] != "bob" {
		t.Error("extract user.name→name failed")
	}
	if _, has := rows[0]["user"]; has {
		t.Error("remove user failed")
	}
	if rows[0]["tag"] != "v1" {
		t.Error("flatten meta.tag→tag failed")
	}
}

func TestSPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*JSONTransform)(nil)
}
