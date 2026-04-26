# Template Engine

A small command-line template rendering tool written in Go. It reads a strict JSON configuration that defines a set of default values and optional environment-specific overrides, then renders a template by replacing placeholders with values from the merged configuration.

This file clarifies a few gotchas you ran into recently:
- templates can be any file (HTML, text, etc.) — they are not parsed as JSON
- argument order matters (config must be the first argument)
- unreplaced placeholders remain in the output so you can spot typos
- tips to detect common mistakes (typos, wrong argument order, trailing commas)

---

## Quick usage

Build:

```/dev/null/build.sh#L1-1
go build
```

Run (argument order matters):

```/dev/null/run.sh#L1-1
./template_engine <config.json> <environment> <template.file> <output.file>
```

Example:

```/dev/null/run-example.sh#L1-1
./template_engine storage_account_configs.json d templates/html_example/index.template.html out.html
```

Important: the first argument must be the JSON config file. If you pass the template file as the first argument by mistake, the program will attempt to parse the template as JSON and fail with errors like `invalid character '<' looking for beginning of value`. See "Argument order gotcha" below.

---

## Argument order gotcha

You must supply arguments in this order:

1. Config (strict JSON)
2. Environment (string)
3. Template file (any file; HTML/text/etc.)
4. Output file (path to write)

If you swap the config and template paths (or otherwise pass the wrong file as the config), the program will try to parse whatever you gave as the config as JSON and fail. Double-check the ordering when invoking the CLI.

---

## Templates: any file, no JSON validation

The template (3rd argument) is treated as an opaque text file. It can be HTML, XML, plain text, or any other textual format. The program does not parse or validate this file as JSON — it simply reads the file contents and performs string substitutions.

This means:
- Your template may start with `<` (HTML) and that's fine.
- If the program fails with a JSON parsing error that mentions `<`, it usually means you passed the template file path in the config position by mistake (see previous section).

---

## Config file format

The config file must be strict JSON (no trailing commas, no YAML or JSONC). It must have the top-level structure matching the Go `Config` type:

- `defaults`: object — required. A map of keys to values that represent default values.
- `overrides`: array — optional. Each element is an object with:
  - `environment`: string
  - `values`: object (map of keys to override values)

Example valid config:

```/dev/null/example.json#L1-20
{
  "defaults": {
    "sku": "Standard_GRS",
    "kind": "StorageV2",
    "location": "westeurope"
  },
  "overrides": [
    { "environment": "d", "values": { "sku": "Standard_LRS" } },
    { "environment": "t", "values": { "sku": "Standard_LRS" } }
  ]
}
```

Notes:
- Trailing commas are not permitted in strict JSON. If your editor inserts trailing commas on save, configure it not to or use the Prettier config suggested below.
- The program expects strict JSON for the config only. Templates are not validated as JSON.

---

## Placeholder syntax and behavior

Supported placeholder syntaxes (simple string replacement):

- `$key$` — replaced with the value for `key`
- `$<key>$` — alternative syntax treated equivalently

Replacement rules:
- The final rendering map is produced by starting from `defaults` and applying any keys present in the environment `values`.
- Overrides are optional and may be partial — you only need to provide the keys you want to change.
- You may **not** override a key that does not exist in `defaults`. Attempting to provide an override for a non-existent key returns an error.
- If a placeholder exists in the template but the final map does not contain a value for that key, the placeholder is left unchanged in the output. This behavior intentionally surfaces typos and missing keys (you will see the placeholder verbatim in the rendered output).

Example template:

```/dev/null/tpl.example#L1-2
sku=$sku$, location=$<location>$
```

Given final map:

```
{ "sku": "Standard_LRS", "location": "westeurope" }
```

Rendered:

```
sku=Standard_LRS, location=westeurope
```

If the template contains `$wiki_ulr$` while your key is `wiki_url`, the output will still contain `$wiki_ulr$` — this is a helpful signal that the template and config do not match (typo).

---

## Example with the HTML template

```
./template_engine templates/html_example/index_configs.json p templates/html_example/index.template.html index.html
```

## Validation and detection tips

Common mistakes and how to detect/fix them:

- Wrong argument order (passing the template as the config):
  - Symptom: JSON parsing error referencing `<` or other non-JSON characters.
  - Fix: Ensure the first argument is the JSON config file.

- Trailing commas in JSON:
  - Symptom: errors like `invalid character ']' looking for beginning of value`.
  - Fix: Remove trailing commas; use a Prettier config to avoid them (example below).

- Typos between template placeholders and config keys:
  - Symptom: placeholder remains unreplaced in output (e.g., `$wiki_ulr$`).
  - Fix: search the template for unreplaced placeholders, or add a post-render diagnostic (see below) to list placeholders still present.

Quick local validators:
- Python built-in JSON tool:
```/dev/null/json-validate.sh#L1-1
python -m json.tool storage_account_configs.json
```
- jq (if installed):
```/dev/null/jq-validate.sh#L1-1
jq . storage_account_configs.json
```

---

## Suggestions (optional improvements)

If you'd like better visibility and to avoid common mistakes, consider:

- Add a light diagnostic after rendering that reports unreplaced placeholders (the program intentionally leaves them so you can see typos — a diagnostic can help locate them).
- Add a small CLI helper that validates `<config>` looks like JSON (best-effort) and suggests checking argument order when it fails. This must be non-invasive; templates are allowed to start with `<`.
- Add a `-config <path>` flag to make argument ordering explicit and less error-prone.

I can implement any of these if you want — tell me which and I'll add it.

---

## Tooling (editor/CI advice)

- Prettier: to prevent trailing commas in JSON, add a project `.prettierrc.json`:
```/dev/null/.prettierrc.json#L1-3
{
  "trailingComma": "none"
}
```

- CI: Run `gofmt`, `go vet`, `go test ./...` and validate sample JSON files with `python -m json.tool` or `jq`.

- Editor: disable format-on-save for JSON if your editor's formatter introduces trailing commas or configure the formatter to follow the project Prettier settings.

---

## Testing

The repository includes unit tests that cover:
- config parsing
- finding overrides
- merging defaults with partial overrides
- rejecting unknown override keys
- placeholder replacement
- prod-like and missing-override fallbacks

Run:

```/dev/null/test.sh#L1-1
go test ./...
```
