package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to write a file in the test temp dir
func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
	return path
}

func TestParseConfig_ValidAndInvalid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		td := t.TempDir()
		cfgJSON := `{
  "defaults": { "a": "1", "b": "2" },
  "overrides": [
    { "environment": "x", "values": { "a": "9" } }
  ]
}`
		path := writeFile(t, td, "config.json", cfgJSON)

		cfg, err := ParseConfig(path)
		if err != nil {
			t.Fatalf("ParseConfig returned error for valid json: %v", err)
		}
		if cfg == nil {
			t.Fatal("ParseConfig returned nil cfg for valid json")
		}
		if len(cfg.Defaults) != 2 {
			t.Fatalf("expected 2 defaults, got %d", len(cfg.Defaults))
		}
		if len(cfg.Overrides) != 1 {
			t.Fatalf("expected 1 override, got %d", len(cfg.Overrides))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		td := t.TempDir()
		path := writeFile(t, td, "bad.json", `not a json`)
		_, err := ParseConfig(path)
		if err == nil {
			t.Fatal("ParseConfig did not return error for invalid json")
		}
	})
}

func TestFindOverride(t *testing.T) {
	cfg := &Config{
		Defaults: map[string]interface{}{"x": "1"},
		Overrides: []Override{
			{Environment: "one", Values: map[string]interface{}{"x": "9"}},
			{Environment: "two", Values: map[string]interface{}{"x": "8"}},
		},
	}

	o, err := FindOverride(cfg, "two")
	if err != nil {
		t.Fatalf("FindOverride returned unexpected error: %v", err)
	}
	if o.Environment != "two" {
		t.Fatalf("FindOverride returned wrong override: got %q want %q", o.Environment, "two")
	}

	_, err = FindOverride(cfg, "missing")
	if err == nil {
		t.Fatal("FindOverride did not return error for missing environment")
	}
}

func TestMergeDefaultsWithOverride_Behavior(t *testing.T) {
	defaults := map[string]interface{}{
		"sku":      "Standard_GRS",
		"location": "westeurope",
		"kind":     "StorageV2",
	}

	t.Run("nil override uses defaults", func(t *testing.T) {
		m, err := MergeDefaultsWithOverride(defaults, nil, "prod")
		if err != nil {
			t.Fatalf("unexpected error merging nil override: %v", err)
		}
		if m["sku"] != "Standard_GRS" || m["location"] != "westeurope" || m["kind"] != "StorageV2" {
			t.Fatalf("merged map does not equal defaults: %#v", m)
		}
	})

	t.Run("partial override replaces single key", func(t *testing.T) {
		override := map[string]interface{}{"sku": "Standard_LRS"}
		m, err := MergeDefaultsWithOverride(defaults, override, "d")
		if err != nil {
			t.Fatalf("unexpected error merging partial override: %v", err)
		}
		if got := m["sku"]; got != "Standard_LRS" {
			t.Fatalf("expected sku to be overridden to Standard_LRS, got %q", got)
		}
		// ensure other keys are preserved
		if got := m["location"]; got != "westeurope" {
			t.Fatalf("expected location to be preserved as westeurope, got %q", got)
		}
	})

	t.Run("override with unknown key errors", func(t *testing.T) {
		override := map[string]interface{}{"i_don_exist": "x"}
		_, err := MergeDefaultsWithOverride(defaults, override, "weird")
		if err == nil {
			t.Fatal("expected error when overriding a non-existent default key, got nil")
		}
		if !strings.Contains(err.Error(), "unknown key") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func TestReplaceTemplatePlaceholders(t *testing.T) {
	tpl := "sku:$sku$, loc:$<location>$, missing:$missing$"
	final := map[string]string{
		"sku":      "Standard_LRS",
		"location": "westeurope",
	}
	out := ReplaceTemplatePlaceholders(tpl, final)

	if !strings.Contains(out, "sku:Standard_LRS") {
		t.Fatalf("expected sku replacement, got %q", out)
	}
	if !strings.Contains(out, "loc:westeurope") {
		t.Fatalf("expected location replacement, got %q", out)
	}
	// missing key should leave placeholder as-is
	if !strings.Contains(out, "$missing$") {
		t.Fatalf("expected missing placeholder to remain, got %q", out)
	}
}

func TestIntegration_ProdLikeUsesDefaults(t *testing.T) {
	td := t.TempDir()
	cfgJSON := `{
  "defaults": { "sku": "S_GRS", "kind": "StorageV2" },
  "overrides": [
    { "environment": "d", "values": { "sku": "S_LRS" } }
  ]
}`
	cfgPath := writeFile(t, td, "config.json", cfgJSON)
	tplPath := writeFile(t, td, "tpl.txt", "sku=$sku$,kind=$kind$")

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}
	// prod-like env should use defaults even if no override entry exists
	final, err := MergeDefaultsWithOverride(cfg.Defaults, nil, "prod")
	if err != nil {
		t.Fatalf("MergeDefaultsWithOverride failed for prod: %v", err)
	}
	if final["sku"] != "S_GRS" {
		t.Fatalf("expected sku S_GRS for prod, got %q", final["sku"])
	}

	// also exercise template read/replace/write path functions (without checking disk output)
	tplStr, err := ReadTemplate(tplPath)
	if err != nil {
		t.Fatalf("ReadTemplate failed: %v", err)
	}
	out := ReplaceTemplatePlaceholders(tplStr, final)
	if !strings.Contains(out, "sku=S_GRS") || !strings.Contains(out, "kind=StorageV2") {
		t.Fatalf("rendering with defaults failed, got %q", out)
	}
}

func TestIntegration_MissingOverrideFallsBack(t *testing.T) {
	td := t.TempDir()
	cfgJSON := `{
  "defaults": { "sku": "S_GRS", "kind": "StorageV2", "location": "westeurope" },
  "overrides": [
    { "environment": "d", "values": { "sku": "S_LRS" } }
  ]
}`
	cfgPath := writeFile(t, td, "config.json", cfgJSON)

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	// request an environment that has no override entry; should fall back to defaults
	final, err := MergeDefaultsWithOverride(cfg.Defaults, nil, "staging")
	if err != nil {
		t.Fatalf("MergeDefaultsWithOverride failed for missing override: %v", err)
	}
	if final["sku"] != "S_GRS" {
		t.Fatalf("expected sku S_GRS for missing override, got %q", final["sku"])
	}
	if final["location"] != "westeurope" {
		t.Fatalf("expected location westeurope for missing override, got %q", final["location"])
	}
}

func TestDetectUnreplacedPlaceholders(t *testing.T) {
	// Template contains three placeholders; only two are provided in final map.
	tpl := "Welcome $name$, check $<missing>$ and $present$."
	final := map[string]string{
		"name":    "Alice",
		"present": "yes",
	}
	rendered := ReplaceTemplatePlaceholders(tpl, final)

	got := DetectUnreplacedPlaceholders(rendered)
	// Expect exactly one unreplaced placeholder: "missing"
	if len(got) != 1 || got[0] != "missing" {
		t.Fatalf("expected unreplaced placeholders [\"missing\"], got %#v; rendered: %q", got, rendered)
	}
}

func TestStrictModeSimulation(t *testing.T) {
	// Simulate strict placeholders behavior: when placeholders remain after rendering,
	// strict mode should be considered a failure. We test the detection logic here.
	tpl := "A:$a$ B:$b$ C:$c$"
	final := map[string]string{
		"a": "1",
		"c": "3",
		// "b" is intentionally missing
	}
	rendered := ReplaceTemplatePlaceholders(tpl, final)

	placeholders := DetectUnreplacedPlaceholders(rendered)
	if len(placeholders) == 0 {
		t.Fatalf("expected unreplaced placeholders, got none; rendered: %q", rendered)
	}
	// ensure the missing placeholder is reported
	found := false
	for _, p := range placeholders {
		if p == "b" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected placeholder 'b' to be reported, got %#v", placeholders)
	}
}

func TestExpandIncludes_SimpleInclude(t *testing.T) {
	td := t.TempDir()

	// Create a base include file containing defaults
	baseJSON := `{
  "sku": "Included_Sku",
  "location": "included-loc"
}`
	basePath := writeFile(t, td, "base.json", baseJSON)

	// Main config references the include; also overrides sku to ensure override precedence.
	cfgJSON := fmt.Sprintf(`{
  "defaults": {
    "include_base": "%s",
    "sku": "Main_Sku"
  }
}`, basePath)
	cfgPath := writeFile(t, td, "config.json", cfgJSON)

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		t.Fatalf("ParseConfig failed for include test: %v", err)
	}

	// include_* directive should be removed after expansion
	if _, ok := cfg.Defaults["include_base"]; ok {
		t.Fatalf("include key should have been removed from defaults after expansion")
	}

	// sku should reflect the main config override
	if got := fmt.Sprintf("%v", cfg.Defaults["sku"]); got != "Main_Sku" {
		t.Fatalf("expected sku to be Main_Sku, got %q", got)
	}

	// location should be pulled from the included file
	if got := fmt.Sprintf("%v", cfg.Defaults["location"]); got != "included-loc" {
		t.Fatalf("expected location from included defaults, got %q", got)
	}
}

func TestExpandIncludes_RecursiveInclude(t *testing.T) {
	td := t.TempDir()

	// more.json provides b
	moreJSON := `{
  "b": "value-b"
}`
	morePath := writeFile(t, td, "more.json", moreJSON)

	// base.json includes more.json and provides a
	baseJSON := fmt.Sprintf(`{
  "include_more": "%s",
  "a": "value-a"
}`, morePath)
	basePath := writeFile(t, td, "base1.json", baseJSON)

	// main config includes base1.json
	cfgJSON := fmt.Sprintf(`{
  "defaults": {
    "include_base": "%s"
  }
}`, basePath)
	cfgPath := writeFile(t, td, "config.json", cfgJSON)

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		t.Fatalf("ParseConfig failed for recursive include test: %v", err)
	}

	// ensure both a and b are present after recursive expansion
	if got := fmt.Sprintf("%v", cfg.Defaults["a"]); got != "value-a" {
		t.Fatalf("expected a from included base, got %q", got)
	}
	if got := fmt.Sprintf("%v", cfg.Defaults["b"]); got != "value-b" {
		t.Fatalf("expected b from recursive include, got %q", got)
	}

	// include_* directives should have been removed
	for k := range cfg.Defaults {
		if strings.HasPrefix(k, "include_") {
			t.Fatalf("expected include_* keys to be removed after expansion, found %q", k)
		}
	}
}
