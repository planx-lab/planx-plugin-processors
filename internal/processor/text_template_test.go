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
	rows, ok := out.([]string)
	if !ok {
		t.Fatalf("expected []string output, got %T", out)
	}
	if len(rows) != 1 || rows[0] != "x" {
		t.Fatalf(`expected ["x"], got %v`, rows)
	}
}

func TestTextTemplate_MultiField(t *testing.T) {
	tt := newTemplate(t, `{{.first}}-{{.last}}`)
	in := []map[string]any{
		{"first": "jane", "last": "doe"},
		{"first": "bob", "last": "smith"},
	}
	out, _ := tt.Process(in)
	rows := out.([]string)
	if len(rows) != 2 || rows[0] != "jane-doe" || rows[1] != "bob-smith" {
		t.Fatalf("expected rendered names, got %v", rows)
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
