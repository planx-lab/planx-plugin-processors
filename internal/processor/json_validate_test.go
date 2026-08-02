package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newValidate(t *testing.T, required string) *JSONValidate {
	t.Helper()
	v := &JSONValidate{}
	cfg, err := json.Marshal(struct {
		RequiredFields string `json:"required_fields"`
	}{RequiredFields: required})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := v.Init(context.Background(), cfg); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return v
}

func TestJSONValidate_ValidBatchPassesThrough(t *testing.T) {
	v := newValidate(t, `["id","name"]`)
	in := []map[string]any{
		{"id": 1, "name": "alice"},
		{"id": 2, "name": "bob"},
	}
	out, err := v.Process(in)
	if err != nil {
		t.Fatalf("valid batch should pass, got %v", err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 2 {
		t.Fatalf("expected batch unchanged (2 rows), got %d", len(rows))
	}
}

func TestJSONValidate_MissingRequiredFieldErrors(t *testing.T) {
	v := newValidate(t, `["id","name"]`)
	in := []map[string]any{
		{"id": 1, "name": "alice"},
		{"id": 2}, // missing "name"
	}
	_, err := v.Process(in)
	if err == nil {
		t.Fatal("expected error when a required field is missing")
	}
}

func TestJSONValidate_NilValueCountsAsMissing(t *testing.T) {
	v := newValidate(t, `["name"]`)
	in := []map[string]any{{"name": nil}}
	_, err := v.Process(in)
	if err == nil {
		t.Fatal("expected error when required field value is nil")
	}
}

func TestJSONValidate_EmptyBatchPasses(t *testing.T) {
	v := newValidate(t, `["id"]`)
	out, err := v.Process([]map[string]any{})
	if err != nil {
		t.Fatalf("empty batch should pass, got %v", err)
	}
	rows := out.([]map[string]any)
	if len(rows) != 0 {
		t.Fatalf("empty in → empty out, got %d", len(rows))
	}
}

func TestJSONValidate_RequiresFieldsConfig(t *testing.T) {
	v := &JSONValidate{}
	err := v.Init(context.Background(), []byte(`{}`))
	if err == nil {
		t.Fatal("expected error when required_fields is empty")
	}
}

func TestJSONValidate_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*JSONValidate)(nil)
}
