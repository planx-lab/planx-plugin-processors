package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type JSONValidateConfig struct {
	RequiredFields string `json:"required_fields"`
}

// JSONValidate checks that every row carries all required_fields. If any row is
// missing a required field (absent or nil), Process returns an error
// (fail-fast); otherwise the batch is returned unchanged.
type JSONValidate struct {
	required []string
}

// NewJSONValidate builds the json-validate processor.
func NewJSONValidate() sdk.ProcessorSPI { return &JSONValidate{} }

func (v *JSONValidate) Init(_ context.Context, cfg []byte) error {
	var c JSONValidateConfig
	// An empty config is allowed only if it still carries required_fields; a
	// nil/empty byte slice decodes to a zero-value struct and fails the
	// required check below.
	if len(cfg) > 0 {
		if err := json.Unmarshal(cfg, &c); err != nil {
			return fmt.Errorf("json-validate: config: %w", err)
		}
	}
	if c.RequiredFields == "" {
		return fmt.Errorf("json-validate: required_fields is required")
	}
	if err := json.Unmarshal([]byte(c.RequiredFields), &v.required); err != nil {
		return fmt.Errorf("json-validate: parse required_fields JSON: %w", err)
	}
	if len(v.required) == 0 {
		return fmt.Errorf("json-validate: required_fields must list at least one field")
	}
	return nil
}

func (v *JSONValidate) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("json-validate: %w", err)
	}
	for _, row := range rows {
		for _, field := range v.required {
			val, ok := row[field]
			if !ok || val == nil {
				return nil, fmt.Errorf("json-validate: row missing required field %q", field)
			}
		}
	}
	return rows, nil
}

func (v *JSONValidate) Close() error { return nil }
