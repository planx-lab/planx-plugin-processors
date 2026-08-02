package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type FieldMapperConfig struct {
	Mappings string `json:"mappings"`
}

type Mapping struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Action string `json:"action"` // rename | drop | add
	Value  string `json:"value"`  // for add
}

type FieldMapper struct {
	mappings []Mapping
}

// NewFieldMapper builds the field-mapper processor.
func NewFieldMapper() sdk.ProcessorSPI { return &FieldMapper{} }

func (m *FieldMapper) Init(_ context.Context, cfg []byte) error {
	var c FieldMapperConfig
	if err := json.Unmarshal(cfg, &c); err != nil {
		return fmt.Errorf("field-mapper: config: %w", err)
	}
	if c.Mappings == "" {
		return fmt.Errorf("field-mapper: mappings is required")
	}
	if err := json.Unmarshal([]byte(c.Mappings), &m.mappings); err != nil {
		return fmt.Errorf("field-mapper: parse mappings JSON: %w", err)
	}
	return nil
}

func (m *FieldMapper) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := sdk.ToRows(b)
	if err != nil {
		return nil, fmt.Errorf("field-mapper: %w", err)
	}
	for _, row := range rows {
		for _, mp := range m.mappings {
			switch mp.Action {
			case "rename":
				if val, ok := row[mp.From]; ok {
					delete(row, mp.From)
					row[mp.To] = val
				}
			case "drop":
				delete(row, mp.From)
			case "add":
				row[mp.To] = mp.Value
			}
		}
	}
	return rows, nil
}

func (m *FieldMapper) Close() error { return nil }
