package main

import (
	"github.com/planx-lab/planx-plugin-processors/internal/processor"
	"github.com/planx-lab/planx-sdk-go/sdk"
)

func main() {
	sdk.Serve(sdk.Plugin{
		ID:          "processors",
		Version:     "1.0.0",
		DisplayName: "Processors",
		Description: "Common ETL processors: passthrough, filter, field-mapper, json-transform, json-validate, json-redact, regex-replace, text-template.",
		Components: []sdk.ComponentSpec{
			{
				ID:          "passthrough",
				Kind:        sdk.KindProcessor,
				DisplayName: "Passthrough",
				Processor:   processor.NewPassthrough,
			},
			{
				ID:          "filter",
				Kind:        sdk.KindProcessor,
				DisplayName: "Filter",
				Processor:   processor.NewFilter,
				ConfigSchema: sdk.Schema(
					sdk.StringField("field", sdk.Required(), sdk.WithDescription("Field name to filter on")),
					sdk.EnumField("operator", []string{"eq", "ne", "gt", "lt", "ge", "le", "contains"},
						sdk.WithDefault(sdk.StringValue("eq")), sdk.WithDescription("Comparison operator")),
					sdk.StringField("value", sdk.Required(), sdk.WithDescription("Value to compare (auto-parsed: bool/number/string)")),
				),
			},
			{
				ID:          "field-mapper",
				Kind:        sdk.KindProcessor,
				DisplayName: "Field Mapper",
				Processor:   processor.NewFieldMapper,
				ConfigSchema: sdk.Schema(
					sdk.StringField("mappings", sdk.Required(),
						sdk.WithDescription(`JSON: [{"from":"old","to":"new","action":"rename"}]`),
						sdk.WithExample(`[{"from":"name","to":"username","action":"rename"}]`)),
				),
			},
			{
				ID:          "json-transform",
				Kind:        sdk.KindProcessor,
				DisplayName: "JSON Transform",
				Processor:   processor.NewJSONTransform,
				ConfigSchema: sdk.Schema(
					sdk.StringField("operations", sdk.Required(),
						sdk.WithDescription(`JSON: [{"op":"extract","path":"data.name","to":"name"}]`),
						sdk.WithExample(`[{"op":"extract","path":"data.name","to":"name"}]`)),
				),
			},
			{
				ID:          "json-validate",
				Kind:        sdk.KindProcessor,
				DisplayName: "JSON Validate",
				Processor:   processor.NewJSONValidate,
				ConfigSchema: sdk.Schema(
					sdk.StringField("required_fields", sdk.Required(),
						sdk.WithDescription(`JSON array of field names every row must carry; fail-fast on any missing field`),
						sdk.WithExample(`["id","name"]`)),
				),
			},
			{
				ID:          "json-redact",
				Kind:        sdk.KindProcessor,
				DisplayName: "JSON Redact",
				Processor:   processor.NewJSONRedact,
				ConfigSchema: sdk.Schema(
					sdk.StringField("fields", sdk.Required(),
						sdk.WithDescription(`JSON array of field names to overwrite with "***"`),
						sdk.WithExample(`["ssn","email"]`)),
				),
			},
			{
				ID:          "regex-replace",
				Kind:        sdk.KindProcessor,
				DisplayName: "Regex Replace",
				Processor:   processor.NewRegexReplace,
				ConfigSchema: sdk.Schema(
					sdk.StringField("field", sdk.Required()),
					sdk.StringField("pattern", sdk.Required(), sdk.WithDescription("Go regexp (RE2 syntax)")),
					sdk.StringField("replacement", sdk.Required(), sdk.WithDescription("Replacement (${1} for capture groups)")),
				),
			},
			{
				ID:          "text-template",
				Kind:        sdk.KindProcessor,
				DisplayName: "Text Template",
				Processor:   processor.NewTextTemplate,
				ConfigSchema: sdk.Schema(
					sdk.StringField("template", sdk.Required(),
						sdk.WithDescription("Go text/template rendered per row (dot = row map); output is []string"),
						sdk.WithExample("{{.name}}")),
				),
			},
		},
	})
}
