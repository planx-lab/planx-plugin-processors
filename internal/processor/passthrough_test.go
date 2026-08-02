package processor

import (
	"context"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newPassthrough(t *testing.T, cfg string) *Passthrough {
	t.Helper()
	p := &Passthrough{}
	if err := p.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return p
}

func TestPassthrough_ReturnsBatchUnchanged(t *testing.T) {
	p := newPassthrough(t, `{}`)
	in := []map[string]any{{"a": 1}}
	out, err := p.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows, ok := out.([]map[string]any)
	if !ok {
		t.Fatalf("expected []map[string]any, got %T", out)
	}
	if len(rows) != 1 || len(rows[0]) != 1 || rows[0]["a"] != 1 {
		t.Fatalf("expected 1:1 passthrough of {a:1}, got %v", rows)
	}
}

func TestPassthrough_EmptyConfig(t *testing.T) {
	p := newPassthrough(t, ``)
	if err := p.Init(context.Background(), nil); err != nil {
		t.Fatalf("Init with nil config should succeed, got %v", err)
	}
	in := []map[string]any{{"x": "y"}}
	out, err := p.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 1 || rows[0]["x"] != "y" {
		t.Fatalf("expected unchanged, got %v", rows)
	}
}

func TestPassthrough_MultiRow(t *testing.T) {
	p := newPassthrough(t, `{}`)
	in := []map[string]any{{"i": 1}, {"i": 2}, {"i": 3}}
	out, _ := p.Process(in)
	rows := out.([]map[string]any)
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows out, got %d", len(rows))
	}
}

func TestPassthrough_EmptyBatch(t *testing.T) {
	p := newPassthrough(t, `{}`)
	out, err := p.Process([]map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 0 {
		t.Fatalf("empty in → empty out, got %d", len(rows))
	}
}

func TestPassthrough_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*Passthrough)(nil)
}
