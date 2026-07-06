package main

import (
	"github.com/planx-lab/planx-plugin-processors/internal/fieldmapper"
	"github.com/planx-lab/planx-plugin-processors/internal/filter"
	"github.com/planx-lab/planx-plugin-processors/internal/jsontransform"
	"github.com/planx-lab/planx-plugin-processors/internal/regexreplace"
	"github.com/planx-lab/planx-sdk-go/sdk"
)

func init() {
	sdk.RegisterType([]map[string]any{})
	sdk.RegisterType(map[string]any{})
}

func main() {
	sdk.Serve(sdk.Plugin{
		ID:          "processors",
		Version:     "1.0.0",
		DisplayName: "Processors",
		Description: "Common ETL processors: filter, field-mapper, json-transform, regex-replace.",
		Components: []sdk.ComponentSpec{
			{
				ID:   "filter",
				Kind: sdk.KindProcessor,
				DisplayName: "Filter",
				Processor:   filter.New,
				ConfigSchema: sdk.Schema(
					sdk.StringField("field", sdk.Required(), sdk.WithDescription("Field name to filter on")),
					sdk.EnumField("operator", []string{"eq", "ne", "gt", "lt", "ge", "le", "contains"},
						sdk.WithDefault(sdk.StringValue("eq")), sdk.WithDescription("Comparison operator")),
					sdk.StringField("value", sdk.Required(), sdk.WithDescription("Value to compare (auto-parsed: bool/number/string)")),
				),
			},
			{
				ID:   "field-mapper",
				Kind: sdk.KindProcessor,
				DisplayName: "Field Mapper",
				Processor:   fieldmapper.New,
				ConfigSchema: sdk.Schema(
					sdk.StringField("mappings", sdk.Required(),
						sdk.WithDescription(`JSON: [{"from":"old","to":"new","action":"rename"}]`),
						sdk.WithExample(`[{"from":"name","to":"username","action":"rename"}]`)),
				),
			},
			{
				ID:   "json-transform",
				Kind: sdk.KindProcessor,
				DisplayName: "JSON Transform",
				Processor:   jsontransform.New,
				ConfigSchema: sdk.Schema(
					sdk.StringField("operations", sdk.Required(),
						sdk.WithDescription(`JSON: [{"op":"extract","path":"data.name","to":"name"}]`),
						sdk.WithExample(`[{"op":"extract","path":"data.name","to":"name"}]`)),
				),
			},
			{
				ID:   "regex-replace",
				Kind: sdk.KindProcessor,
				DisplayName: "Regex Replace",
				Processor:   regexreplace.New,
				ConfigSchema: sdk.Schema(
					sdk.StringField("field", sdk.Required()),
					sdk.StringField("pattern", sdk.Required(), sdk.WithDescription("Go regexp (RE2 syntax)")),
					sdk.StringField("replacement", sdk.Required(), sdk.WithDescription("Replacement (${1} for capture groups)")),
				),
			},
		},
	})
}
