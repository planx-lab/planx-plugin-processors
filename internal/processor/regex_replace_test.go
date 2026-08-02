package processor

import (
	"context"
	"testing"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

func newReplacer(t *testing.T, cfg string) *RegexReplace {
	t.Helper()
	r := &RegexReplace{}
	if err := r.Init(context.Background(), []byte(cfg)); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return r
}

func TestBasicReplace(t *testing.T) {
	r := newReplacer(t, `{"field":"email","pattern":"@.*","replacement":"@redacted.com"}`)
	in := []map[string]any{{"email": "alice@example.com", "name": "alice"}}
	out, err := r.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if rows[0]["email"] != "alice@redacted.com" {
		t.Errorf("expected alice@redacted.com, got %v", rows[0]["email"])
	}
	if rows[0]["name"] != "alice" {
		t.Error("name should be untouched")
	}
}

func TestCaptureGroup(t *testing.T) {
	r := newReplacer(t, `{"field":"phone","pattern":"(\\d{3})-(\\d{4})","replacement":"${1}XXXX"}`)
	in := []map[string]any{{"phone": "123-4567"}}
	out, _ := r.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["phone"] != "123XXXX" {
		t.Errorf("expected 123XXXX, got %v", rows[0]["phone"])
	}
}

func TestNonStringField(t *testing.T) {
	r := newReplacer(t, `{"field":"age","pattern":"\\d","replacement":"X"}`)
	in := []map[string]any{{"age": int64(30), "name": "bob"}}
	out, err := r.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if rows[0]["age"] != int64(30) {
		t.Error("non-string field should be unchanged")
	}
}

func TestRegexReplace_MissingField(t *testing.T) {
	r := newReplacer(t, `{"field":"nonexistent","pattern":"x","replacement":"y"}`)
	in := []map[string]any{{"other": "value"}}
	out, err := r.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.([]map[string]any)
	if rows[0]["other"] != "value" {
		t.Error("other field should be unchanged")
	}
}

func TestInvalidRegex(t *testing.T) {
	r := &RegexReplace{}
	err := r.Init(context.Background(), []byte(`{"field":"x","pattern":"[invalid","replacement":"y"}`))
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestMultipleRows(t *testing.T) {
	r := newReplacer(t, `{"field":"domain","pattern":"\\.com$","replacement":".org"}`)
	in := []map[string]any{
		{"domain": "a.com"},
		{"domain": "b.org"},
		{"domain": "c.com"},
	}
	out, _ := r.Process(in)
	rows := out.([]map[string]any)
	if rows[0]["domain"] != "a.org" || rows[2]["domain"] != "c.org" {
		t.Errorf("expected a.org and c.org, got %v and %v", rows[0]["domain"], rows[2]["domain"])
	}
	if rows[1]["domain"] != "b.org" {
		t.Error("b.org should be unchanged (no match)")
	}
}

func TestRegexReplace_SPIConformance(t *testing.T) {
	var _ sdk.ProcessorSPI = (*RegexReplace)(nil)
}
