package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newRedact(t *testing.T, fields string) *JSONRedact {
	t.Helper()
	r := &JSONRedact{}
	cfg, err := json.Marshal(struct {
		Fields string `json:"fields"`
	}{Fields: fields})
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := r.Init(context.Background(), cfg); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return r
}

func TestJSONRedact_ReplacesListedFields(t *testing.T) {
	r := newRedact(t, `["ssn","email"]`)
	in := []map[string]any{
		{"ssn": "111-22-3333", "email": "a@x.com", "name": "alice"},
	}
	out, err := r.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if rows[0]["ssn"] != "***" {
		t.Errorf("ssn should be ***, got %v", rows[0]["ssn"])
	}
	if rows[0]["email"] != "***" {
		t.Errorf("email should be ***, got %v", rows[0]["email"])
	}
	if rows[0]["name"] != "alice" {
		t.Errorf("name should be untouched, got %v", rows[0]["name"])
	}
}

func TestJSONRedact_MissingFieldIsNoOp(t *testing.T) {
	r := newRedact(t, `["missing"]`)
	in := []map[string]any{{"name": "alice"}}
	out, err := r.Process(in)
	if err != nil {
		t.Fatalf("missing field should not error, got %v", err)
	}
	rows := out.([]map[string]any)
	if _, has := rows[0]["missing"]; has {
		t.Errorf("missing field should not be created, got %v", rows[0]["missing"])
	}
	if rows[0]["name"] != "alice" {
		t.Errorf("existing field should be untouched, got %v", rows[0]["name"])
	}
}

func TestJSONRedact_MultiRow(t *testing.T) {
	r := newRedact(t, `["secret"]`)
	in := []map[string]any{
		{"secret": "a", "keep": "1"},
		{"secret": "b", "keep": "2"},
	}
	out, _ := r.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["secret"] != "***" || rows[1]["secret"] != "***" {
		t.Errorf("both rows redacted, got %v %v", rows[0]["secret"], rows[1]["secret"])
	}
	if rows[0]["keep"] != "1" || rows[1]["keep"] != "2" {
		t.Error("non-listed fields should be unchanged")
	}
}

func TestJSONRedact_RequiresFieldsConfig(t *testing.T) {
	r := &JSONRedact{}
	err := r.Init(context.Background(), []byte(`{}`))
	if err == nil {
		t.Fatal("expected error when fields is empty")
	}
}

func TestJSONRedact_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*JSONRedact)(nil)
}
