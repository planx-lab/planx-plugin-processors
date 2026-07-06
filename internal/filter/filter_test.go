package filter

import (
	"context"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newFilter(t *testing.T, cfg string) *Filter {
	t.Helper()
	f := &Filter{}
	if err := f.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return f
}

func TestFilter_BoolEq(t *testing.T) {
	f := newFilter(t, `{"field":"active","operator":"eq","value":"true"}`)
	in := []map[string]any{
		{"active": true, "name": "alice"},
		{"active": false, "name": "bob"},
	}
	out, err := f.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 1 || rows[0]["name"] != "alice" {
		t.Fatalf("expected 1 row (alice), got %v", rows)
	}
}

func TestFilter_IntGt(t *testing.T) {
	f := newFilter(t, `{"field":"age","operator":"gt","value":"30"}`)
	in := []map[string]any{
		{"age": int64(25)},
		{"age": int64(35)},
		{"age": int64(40)},
	}
	out, _ := f.Process(in)
	rows := out.([]map[string]any)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows (35,40), got %d", len(rows))
	}
}

func TestFilter_StringContains(t *testing.T) {
	f := newFilter(t, `{"field":"email","operator":"contains","value":"@example.com"}`)
	in := []map[string]any{
		{"email": "alice@example.com"},
		{"email": "bob@other.com"},
	}
	out, _ := f.Process(in)
	rows := out.([]map[string]any)
	if len(rows) != 1 || rows[0]["email"] != "alice@example.com" {
		t.Fatalf("expected alice only, got %v", rows)
	}
}

func TestFilter_MissingField(t *testing.T) {
	f := newFilter(t, `{"field":"missing","operator":"eq","value":"x"}`)
	in := []map[string]any{{"other": "y"}}
	out, _ := f.Process(in)
	rows := out.([]map[string]any)
	if len(rows) != 0 {
		t.Fatalf("missing field → 0 rows, got %d", len(rows))
	}
}

func TestFilter_EmptyBatch(t *testing.T) {
	f := newFilter(t, `{"field":"x","operator":"eq","value":"1"}`)
	out, err := f.Process([]map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 0 {
		t.Fatalf("empty in → empty out, got %d", len(rows))
	}
}

func TestFilter_FloatGe(t *testing.T) {
	f := newFilter(t, `{"field":"score","operator":"ge","value":"90.5"}`)
	in := []map[string]any{
		{"score": 85.5},
		{"score": 90.5},
		{"score": 95.0},
	}
	out, _ := f.Process(in)
	rows := out.([]map[string]any)
	if len(rows) != 2 {
		t.Fatalf("ge 90.5 → 2 rows, got %d", len(rows))
	}
}

func TestFilter_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*Filter)(nil)
}
