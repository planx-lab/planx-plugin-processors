package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type TextTemplateConfig struct {
	Template string `json:"template"`
	Field    string `json:"field"` // optional output field; defaults to "_text"
}

// defaultTextField is the field text-template writes its rendered result into
// when the config omits "field".
const defaultTextField = "_text"

// TextTemplate renders a Go text/template against each row (a map[string]any)
// and writes the rendered string BACK into the same row under Field (default
// "_text"). The output batch stays sdk.Rows so the pipeline type contract
// (Rows end-to-end) is never broken. The previous design returned []string,
// which no downstream processor/sink could consume.
type TextTemplate struct {
	tmpl  *template.Template
	field string
}

// NewTextTemplate builds the text-template processor.
func NewTextTemplate() sdk.ProcessorSPI { return &TextTemplate{} }

func (t *TextTemplate) Init(_ context.Context, cfg []byte) error {
	var c TextTemplateConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return fmt.Errorf("text-template: config: %w", err)
	}
	if c.Template == "" {
		return fmt.Errorf("text-template: template is required")
	}
	tmpl, err := template.New("text-template").Parse(c.Template)
	if err != nil {
		return fmt.Errorf("text-template: parse template: %w", err)
	}
	t.tmpl = tmpl
	t.field = c.Field
	if t.field == "" {
		t.field = defaultTextField
	}
	return nil
}

func (t *TextTemplate) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := sdk.ToRows(b)
	if err != nil {
		return nil, fmt.Errorf("text-template: %w", err)
	}
	for i, row := range rows {
		var buf bytes.Buffer
		if err := t.tmpl.Execute(&buf, row); err != nil {
			return nil, fmt.Errorf("text-template: render row %d: %w", i, err)
		}
		row[t.field] = buf.String()
	}
	return rows, nil
}

func (t *TextTemplate) Close() error { return nil }
