package batch

import "fmt"

// ToMaps converts any supported source batch type to []map[string]any.
// Supported: []map[string]any (native), map[string]any (single row),
// map[string]string / []map[string]string (value widen).
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
