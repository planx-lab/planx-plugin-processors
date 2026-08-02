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
}

// TextTemplate renders a Go text/template against each row (a map[string]any),
// producing one rendered string per row. The output batch is []string.
type TextTemplate struct {
	tmpl *template.Template
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
	return nil
}

func (t *TextTemplate) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("text-template: %w", err)
	}
	out := make([]string, len(rows))
	for i, row := range rows {
		var buf bytes.Buffer
		if err := t.tmpl.Execute(&buf, row); err != nil {
			return nil, fmt.Errorf("text-template: render row %d: %w", i, err)
		}
		out[i] = buf.String()
	}
	return out, nil
}

func (t *TextTemplate) Close() error { return nil }
