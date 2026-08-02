package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

type RegexReplaceConfig struct {
	Field       string `json:"field"`
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

type RegexReplace struct {
	cfg RegexReplaceConfig
	re  *regexp.Regexp
}

// NewRegexReplace builds the regex-replace processor.
func NewRegexReplace() sdk.ProcessorSPI { return &RegexReplace{} }

func (r *RegexReplace) Init(_ context.Context, cfg []byte) error {
	if err := json.Unmarshal(cfg, &r.cfg); err != nil {
		return fmt.Errorf("regex-replace: config: %w", err)
	}
	if r.cfg.Field == "" {
		return fmt.Errorf("regex-replace: field is required")
	}
	if r.cfg.Pattern == "" {
		return fmt.Errorf("regex-replace: pattern is required")
	}
	re, err := regexp.Compile(r.cfg.Pattern)
	if err != nil {
		return fmt.Errorf("regex-replace: invalid pattern %q: %w", r.cfg.Pattern, err)
	}
	r.re = re
	return nil
}

func (r *RegexReplace) Process(b sdk.Batch) (sdk.Batch, error) {
	rows, err := ToMaps(b)
	if err != nil {
		return nil, fmt.Errorf("regex-replace: %w", err)
	}
	for _, row := range rows {
		val, ok := row[r.cfg.Field]
		if !ok {
			continue // missing field — skip
		}
		s, ok := val.(string)
		if !ok {
			continue // non-string — skip
		}
		row[r.cfg.Field] = r.re.ReplaceAllString(s, r.cfg.Replacement)
	}
	return rows, nil
}

func (r *RegexReplace) Close() error { return nil }
