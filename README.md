# planx-plugin-processors

A Planx Connector plugin bundling common ETL **processor** components. One
self-describing binary declares multiple `sdk.ComponentSpec`s (ADR-008/009); each
processor implements `sdk.ProcessorSPI` and lives under
[`internal/processor/`](./internal/processor/).

## Components

| ID              | Factory                      | What it does                                                                                  |
| --------------- | ---------------------------- | --------------------------------------------------------------------------------------------- |
| `passthrough`   | `processor.NewPassthrough`   | Returns the batch unchanged (1:1). No config.                                                 |
| `filter`        | `processor.NewFilter`        | Keeps rows where `field <op> value` holds (`eq`, `ne`, `gt`, `lt`, `ge`, `le`, `contains`).   |
| `field-mapper`  | `processor.NewFieldMapper`   | Renames / drops / adds fields per a JSON mappings spec.                                       |
| `json-transform`| `processor.NewJSONTransform` | `extract` (dot-path), `flatten` (prefix strip), `remove` operations per a JSON operations spec.|
| `json-validate` | `processor.NewJSONValidate`  | Fail-fast if any row is missing a required field; otherwise passes the batch through.         |
| `json-redact`   | `processor.NewJSONRedact`    | Overwrites listed fields with `"***"`; absent fields are a no-op.                             |
| `regex-replace` | `processor.NewRegexReplace`  | Applies a Go RE2 regexp replace to a string field.                                            |
| `text-template` | `processor.NewTextTemplate`  | Renders a Go `text/template` per row (dot = row map); output is `[]string`.                   |

## Layout

```
cmd/plugin/main.go        # sdk.Serve(...) registering all 8 processors
internal/processor/       # package processor — one file per component + convert.go
  convert.go              # ToMaps helper (shared decode)
  passthrough.go
  filter.go
  field_mapper.go
  json_transform.go
  json_validate.go
  json_redact.go
  regex_replace.go
  text_template.go
```

All processors are stateless, deterministic, and fail-fast on bad input. Batch
is treated as opaque; row-based batches are decoded via `processor.ToMaps`.
