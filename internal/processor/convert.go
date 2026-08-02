package processor

import "fmt"

// ToMaps converts any supported source batch type to []map[string]any.
// Supported: []map[string]any (native), map[string]any (single row),
// map[string]string / []map[string]string (value widen),
// [][]string (CSV rows: first row is the header, rest are data rows).
func ToMaps(b any) ([]map[string]any, error) {
	switch v := b.(type) {
	case nil:
		return nil, fmt.Errorf("batch is nil")
	case []map[string]any:
		return v, nil
	case map[string]any:
		return []map[string]any{v}, nil
	case map[string]string:
		return []map[string]any{widenStringMap(v)}, nil
	case []map[string]string:
		out := make([]map[string]any, len(v))
		for i, m := range v {
			out[i] = widenStringMap(m)
		}
		return out, nil
	case [][]string:
		// CSV-shaped batch: first row is the header (column names), subsequent
		// rows are values. Map each value row to a map keyed by the header.
		if len(v) == 0 {
			return []map[string]any{}, nil
		}
		header := v[0]
		out := make([]map[string]any, 0, len(v)-1)
		for _, row := range v[1:] {
			m := make(map[string]any, len(header))
			for i, col := range header {
				if i < len(row) {
					m[col] = row[i]
				}
			}
			out = append(out, m)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported batch type %T", b)
	}
}

func widenStringMap(m map[string]string) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
