package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type FilterConfig struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Filter struct {
	cfg FilterConfig
}

// NewFilter builds the filter processor.
func NewFilter() sdk.ProcessorSPI { return &Filter{} }

func (f *Filter) Init(_ context.Context, cfg []byte) error {
	if err := json.Unmarshal(cfg, &f.cfg); err != nil {
		return fmt.Errorf("filter: config: %w", err)
	}
	if f.cfg.Field == "" {
		return fmt.Errorf("filter: field is required")
	}
	if f.cfg.Operator == "" {
		f.cfg.Operator = "eq"
	}
	return nil
}

func (f *Filter) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}
	var result []map[string]any
	for _, row := range rows {
		if f.matches(row) {
			result = append(result, row)
		}
	}
	if result == nil {
		result = []map[string]any{}
	}
	return result, nil
}

func (f *Filter) matches(row map[string]any) bool {
	val, ok := row[f.cfg.Field]
	if !ok || val == nil {
		return false // missing/null → excluded
	}
	switch v := val.(type) {
	case bool:
		target, err := strconv.ParseBool(f.cfg.Value)
		if err != nil {
			return false
		}
		switch f.cfg.Operator {
		case "eq":
			return v == target
		case "ne":
			return v != target
		}
	case float64: // JSON numbers come as float64
		target, err := strconv.ParseFloat(f.cfg.Value, 64)
		if err != nil {
			return false
		}
		return compareNumeric(v, target, f.cfg.Operator)
	case int64:
		target, err := strconv.ParseFloat(f.cfg.Value, 64)
		if err != nil {
			return false
		}
		return compareNumeric(float64(v), target, f.cfg.Operator)
	case int:
		target, err := strconv.ParseFloat(f.cfg.Value, 64)
		if err != nil {
			return false
		}
		return compareNumeric(float64(v), target, f.cfg.Operator)
	case string:
		switch f.cfg.Operator {
		case "eq":
			return v == f.cfg.Value
		case "ne":
			return v != f.cfg.Value
		case "contains":
			return strings.Contains(v, f.cfg.Value)
		}
	}
	return false
}

func compareNumeric(a, b float64, op string) bool {
	switch op {
	case "eq":
		return a == b
	case "ne":
		return a != b
	case "gt":
		return a > b
	case "lt":
		return a < b
	case "ge":
		return a >= b
	case "le":
		return a <= b
	}
	return false
}

func (f *Filter) Close() error { return nil }
