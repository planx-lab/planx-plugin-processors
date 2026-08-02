package processor

import (
	"context"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newTemplate(t *testing.T, tmpl string) *TextTemplate {
	t.Helper()
	tt := &TextTemplate{}
	cfg := `{"template":"` + escapeJSON(tmpl) + `"}`
	if err := tt.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return tt
}

func TestTextTemplate_RendersPerRow(t *testing.T) {
	tt := newTemplate(t, `{{.name}}`)
	in := []map[string]any{{"name": "x"}}
	out, err := tt.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows, ok := out.(sdk.Rows)
	if !ok {
		t.Fatalf("expected sdk.Rows output, got %T", out)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["_text"] != "x" {
		t.Fatalf(`expected _text="x", got %v`, rows[0]["_text"])
	}
}

func TestTextTemplate_MultiField(t *testing.T) {
	tt := newTemplate(t, `{{.first}}-{{.last}}`)
	in := []map[string]any{
		{"first": "jane", "last": "doe"},
		{"first": "bob", "last": "smith"},
	}
	out, _ := tt.Process(in)
	rows := out.(sdk.Rows)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0]["_text"] != "jane-doe" || rows[1]["_text"] != "bob-smith" {
		t.Fatalf("expected rendered names in _text, got %v / %v", rows[0]["_text"], rows[1]["_text"])
	}
	// original fields must survive
	if rows[0]["first"] != "jane" || rows[0]["last"] != "doe" {
		t.Errorf("input fields should be preserved, got %v", rows[0])
	}
}

// TestTextTemplate_CustomField verifies the configurable output Field.
func TestTextTemplate_CustomField(t *testing.T) {
	tt := &TextTemplate{}
	cfg := `{"template":"{{.name}}","field":"greeting"}`
	if err := tt.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	out, err := tt.Process(sdk.Rows{{"name": "alice"}})
	if err != nil {
		t.Fatal(err)
	}
	rows := out.(sdk.Rows)
	if rows[0]["greeting"] != "alice" {
		t.Fatalf(`expected greeting="alice", got %v`, rows[0]["greeting"])
	}
	if _, has := rows[0]["_text"]; has {
		t.Errorf("default _text field should not be set when field=greeting")
	}
}

func TestTextTemplate_InvalidTemplateErrorsOnInit(t *testing.T) {
	tt := &TextTemplate{}
	err := tt.Init(context.Background(), []byte(`{"template":"{{.name"}`))
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestTextTemplate_RequiresTemplateConfig(t *testing.T) {
	tt := &TextTemplate{}
	err := tt.Init(context.Background(), []byte(`{}`))
	if err == nil {
		t.Fatal("expected error when template is empty")
	}
}

func TestTextTemplate_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*TextTemplate)(nil)
}
