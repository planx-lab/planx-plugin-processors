package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type JSONRedactConfig struct {
	Fields string `json:"fields"`
}

const redactMask = "***"

// JSONRedact overwrites each listed field on every row with "***". Fields that
// are absent from a row are a no-op (not an error). Other fields are left
// untouched.
type JSONRedact struct {
	fields []string
}

// NewJSONRedact builds the json-redact processor.
func NewJSONRedact() sdk.ProcessorSPI { return &JSONRedact{} }

func (r *JSONRedact) Init(_ context.Context, cfg []byte) error {
	var c JSONRedactConfig
	if len(cfg) > 0 {
		if err := json.Unmarshal(cfg, &c); err != nil {
			return fmt.Errorf("json-redact: config: %w", err)
		}
	}
	if c.Fields == "" {
		return fmt.Errorf("json-redact: fields is required")
	}
	if err := json.Unmarshal([]byte(c.Fields), &r.fields); err != nil {
		return fmt.Errorf("json-redact: parse fields JSON: %w", err)
	}
	if len(r.fields) == 0 {
		return fmt.Errorf("json-redact: fields must list at least one field")
	}
	return nil
}

func (r *JSONRedact) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("json-redact: %w", err)
	}
	for _, row := range rows {
		for _, field := range r.fields {
			if _, ok := row[field]; ok {
				row[field] = redactMask
			}
		}
	}
	return rows, nil
}

func (r *JSONRedact) Close() error { return nil }
