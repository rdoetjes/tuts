# Template Engine

A small command-line template rendering tool written in Go. It reads a JSON configuration that defines a set of default values and optional environment-specific overrides, then renders a template by replacing placeholders with values from the merged configuration.

This README documents:
- usage
- configuration file format
- placeholder syntax and rendering behavior
- special environment rules
- error handling
- testing and development tooling
- suggestions for editor/CI setup to avoid common problems

---

## Quick usage

Build:

    go build

Run:

    ./template_engine <config.json> <environment> <template.file> <output.file>

Example:

    ./template_engine storage_account_configs.json d template.tpl out.txt

This will:
1. Parse `storage_account_configs.json`.
2. Find the override for environment `d` (if present).
3. Merge defaults with overrides (overrides are optional, but may only touch keys that exist in defaults).
4. Read `template.tpl`, replace placeholders, and write `out.txt`.
5. Print the generated dictionary and the rendered template to stdout.

---

## Config file format

The config file must be strict JSON (no trailing commas, not YAML/JSONC). It must have the top-level structure matching the Go `Config` type:

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
- JSON must be strict. Trailing commas are not allowed (they will cause errors like `invalid character ']' looking for beginning of value`).
- If you have editor tooling that inserts trailing commas on save, see the "Tooling" section below.

---

## Placeholder syntax

Two placeholder styles are supported in templates:

- `$key$` — replaced with the value for `key`
- `$<key>$` — alternative syntax that is treated equivalently

All replacements are simple string substitutions. If a placeholder is present in the template but no value exists for it in the final map, the placeholder is left unchanged in the output (so you can detect missing values in the rendered output).

Example:

Template:

    sku=$sku$, location=$<location>$

Given final map:

    { "sku": "Standard_LRS", "location": "westeurope" }

Rendered:

    sku=Standard_LRS, location=westeurope

---

## Merge semantics and validation

- The final rendering map is produced by starting from `defaults` and then applying any keys present in the environment `values`.
- Overrides are optional and may be partial — you only need to provide the keys you want to change.
- You may **not** override a key that does not exist in `defaults`. Attempting to provide an override for a non-existent key returns an error, e.g.:
  `override for environment 'd' contains unknown key 'i_don_exist' (not present in defaults)`

This protects against typos and accidental introduction of unexpected keys.

---

## Special environment handling

The tool has two convenience behaviors:

1. Prod-like environments — when the requested environment is `p`, `prod`, or `prd` the tool will implicitly use the `defaults` map and will not require an explicit override entry.
2. Missing override fallback — if you request an environment (other than prod-like) that does not have an override entry in the config, the tool will fall back to using `defaults` (same as prod-like behavior). This makes it convenient to render templates for staging/other ephemeral envs without adding explicit override entries when they are not needed.

---

## Error messages you may see

- `invalid character 'c' looking for beginning of value` — your config is not valid JSON (maybe it starts with `config:` or is YAML). Make sure it's strict JSON.
- `invalid character ']' looking for beginning of value` — often due to trailing commas. Remove trailing commas.
- `no override found for environment 'x'` — older behavior; current tool falls back to defaults. If you still see an error like this, make sure you are using the updated binary.
- `override for environment 'd' contains unknown key 'x' (not present in defaults)` — you attempted to override a key not present in defaults (typo or incorrect key).

When an error occurs, the CLI prints the message to stderr and exits with a non-zero status.

---

## Testing

The project includes a small test suite that exercises:
- config parsing
- finding overrides
- merging defaults with partial overrides
- rejecting unknown override keys
- placeholder replacement
- prod-like and missing-override fallbacks

Run the tests:

    go test ./...

The tests are fast and should pass in a normal Go toolchain.

---

## Tooling & editor advice

Common problem: editors/formatters inserting trailing commas in JSON, making it invalid for strict JSON parsers.

Recommendations:

1. Prettier config (prevents trailing commas)
   Add a project `.prettierrc.json` with:

   ```json
   {
     "trailingComma": "none"
   }
   ```

   If your editor uses Prettier on save, this will stop Prettier from inserting trailing commas in JSON files.

2. Validate JSON on save (or before commit)
   - Use `jq . file.json` or `python -m json.tool file.json` to lint files quickly.
   - Consider adding a pre-commit hook that runs a JSON linter/validator to catch invalid JSON before it gets committed.

3. If you want to allow YAML or JSON5/JSONC input, consider switching to a more permissive parser or adding a small pre-processing step (not currently implemented).

---

## CI suggestions

- Run `gofmt` / `go vet` and `go test`.
- Validate JSON files in the repository (or at least any sample config files) with `jq` or `python -m json.tool`.
- Optionally run `prettier --check` if you adopt Prettier.

Example GitHub Actions steps:
- setup Go
- `gofmt -l .` and fail if anything is unformatted
- `go test ./...`
- validate JSON samples: `python -m json.tool path/to/config.json`

---

## Extending the project

If you'd like to make the tool more permissive or integrate with other formats:

- Add YAML support:
  - Use a YAML library to parse YAML, or detect file type by extension and parse accordingly.
- Add JSON5/JSONC support:
  - Use a JSON5 parser to accept trailing commas/comments.
- Make the renderer pluggable:
  - Replace simple placeholder replacement with a templating engine (text/template) if you need conditional logic, loops, etc.

---

## Contact / Contributing

- Feel free to open issues or PRs if you want additional features, CI setup, or if you want the program to accept different input formats.
- Tests are small and straightforward — adding coverage for any new behavior is welcome.

---

License
- No license file is included by default. If you intend to share this project, add a LICENSE to the repository describing the intended license.
