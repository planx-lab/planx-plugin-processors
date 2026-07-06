package jsontransform

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/planx-lab/planx-plugin-processors/internal/batch"
	"github.com/planx-lab/planx-sdk-go/sdk"
)

type Config struct {
	Operations string `json:"operations"`
}

type Operation struct {
	Op     string `json:"op"`     // extract | flatten | remove
	Path   string `json:"path"`   // dot-path for extract/remove
	To     string `json:"to"`     // target field name for extract
	Prefix string `json:"prefix"` // key prefix for flatten (e.g. "addr.")
}

type JSONTransform struct {
	ops []Operation
}

func New() sdk.ProcessorSPI { return &JSONTransform{} }

func (t *JSONTransform) Init(_ context.Context, cfg []byte) error {
	var c Config
	if err := json.Unmarshal(cfg, &c); err != nil {
		return fmt.Errorf("json-transform: config: %w", err)
	}
	if c.Operations == "" {
		return fmt.Errorf("json-transform: operations is required")
	}
	if err := json.Unmarshal([]byte(c.Operations), &t.ops); err != nil {
		return fmt.Errorf("json-transform: parse operations JSON: %w", err)
	}
	return nil
}

func (t *JSONTransform) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := batch.ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("json-transform: %w", err)
	}
	for _, row := range rows {
		for _, op := range t.ops {
			switch op.Op {
			case "extract":
				val, ok := getByPath(row, op.Path)
				if ok {
					row[op.To] = val
				}
			case "flatten":
				for key, val := range row {
					if strings.HasPrefix(key, op.Prefix) {
						flatKey := strings.TrimPrefix(key, op.Prefix)
						row[flatKey] = val
						delete(row, key)
					}
				}
			case "remove":
				delete(row, op.Path)
			}
		}
	}
	return rows, nil
}

// getByPath resolves a dot-separated path (e.g. "data.name") in a nested map.
// Returns (value, true) if found, (nil, false) otherwise.
func getByPath(row map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = row
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func (t *JSONTransform) Close() error { return nil }
