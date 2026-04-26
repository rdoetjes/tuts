package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ParseConfig reads and parses the JSON config at the given path.
// It returns a pointer to Config on success or an error.
func ParseConfig(configPath string) (*Config, error) {
	b, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config json: %w", err)
	}

	if cfg.Defaults == nil {
		return nil, fmt.Errorf("config json missing 'defaults' object")
	}

	// Expand "include_*" keys in defaults: load included files and merge their defaults.
	if err := expandIncludes(&cfg, filepath.Dir(configPath), make(map[string]struct{})); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// expandIncludes looks for keys in cfg.Defaults that start with "include_"
// and treats their string value as a path to another JSON file containing
// defaults to merge in. Included files can be either a full config JSON
// (with a "defaults" object) or a plain object containing key/value pairs.
// Includes are resolved relative to baseDir if not absolute. We track visited
// files to avoid cycles. When merging, values from the included file are only
// added when the key does not already exist in the current cfg.Defaults
// (i.e. explicit defaults in the main config override included values).
func expandIncludes(cfg *Config, baseDir string, visited map[string]struct{}) error {
	// Collect include keys up-front to avoid mutating the map while iterating.
	includeKeys := make([]string, 0, 4)
	for k := range cfg.Defaults {
		if strings.HasPrefix(k, "include_") {
			includeKeys = append(includeKeys, k)
		}
	}

	for _, incKey := range includeKeys {
		rawVal, ok := cfg.Defaults[incKey]
		if !ok {
			continue
		}
		includePath, ok := rawVal.(string)
		if !ok || includePath == "" {
			// ignore non-string include values
			delete(cfg.Defaults, incKey)
			continue
		}

		// Resolve relative paths against baseDir
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(baseDir, includePath)
		}

		// Security: only allow includes with a .json extension
		ext := strings.ToLower(filepath.Ext(includePath))
		if ext != ".json" {
			// remove the include directive and fail loudly to avoid silently accepting unexpected formats
			delete(cfg.Defaults, incKey)
			return fmt.Errorf("included defaults file %s has disallowed extension %q; only .json is allowed", includePath, ext)
		}

		// Prevent cycles
		if _, seen := visited[includePath]; seen {
			// remove the include key and skip
			delete(cfg.Defaults, incKey)
			continue
		}
		visited[includePath] = struct{}{}

		// Read included file
		b, err := os.ReadFile(includePath)
		if err != nil {
			return fmt.Errorf("error reading included defaults file %s: %w", includePath, err)
		}

		// Try to parse as a full Config with "defaults", otherwise treat as a raw map of defaults.
		// Note: json.Unmarshal may succeed when unmarshalling a plain object into a struct
		// (it will simply not populate the `Defaults` field). Therefore we must check both
		// cases: (a) unmarshal into Config produced Defaults, or (b) treat the JSON as a raw map.
		var incCfg Config
		if err := json.Unmarshal(b, &incCfg); err != nil {
			// If unmarshalling into Config failed, fallback to parsing as a raw map.
			var raw map[string]interface{}
			if err2 := json.Unmarshal(b, &raw); err2 != nil {
				return fmt.Errorf("error parsing included defaults json %s: %w", includePath, err)
			}
			// If the file contains a top-level "defaults" object, use it; otherwise treat the object as defaults directly.
			if d, ok := raw["defaults"].(map[string]interface{}); ok {
				incCfg.Defaults = d
			} else {
				incCfg.Defaults = raw
			}
		} else {
			// Unmarshal into Config succeeded. If Defaults is still nil it means the file
			// was probably a plain object (not wrapped in "defaults"). Try to parse into
			// a raw map and use it as defaults in that case.
			if incCfg.Defaults == nil {
				var raw map[string]interface{}
				if err2 := json.Unmarshal(b, &raw); err2 == nil {
					if d, ok := raw["defaults"].(map[string]interface{}); ok {
						incCfg.Defaults = d
					} else {
						incCfg.Defaults = raw
					}
				}
			}
		}

		// If included file has its own includes, expand them recursively.
		if incCfg.Defaults != nil {
			if err := expandIncludes(&incCfg, filepath.Dir(includePath), visited); err != nil {
				return err
			}
		}

		// Merge: do not override existing keys in the main cfg.Defaults
		for k, v := range incCfg.Defaults {
			if _, exists := cfg.Defaults[k]; !exists {
				cfg.Defaults[k] = v
			}
		}

		// Remove the include_* directive after processing so it does not remain a default key.
		delete(cfg.Defaults, incKey)
	}

	return nil
}

// FindOverride returns the override entry for the requested environment.
// If no matching override exists it returns a descriptive error.
func FindOverride(cfg *Config, env string) (*Override, error) {
	for i := range cfg.Overrides {
		if cfg.Overrides[i].Environment == env {
			return &cfg.Overrides[i], nil
		}
	}
	return nil, fmt.Errorf("no override found for environment '%s'", env)
}

// ToStringMap converts a map[string]interface{} into map[string]string
// by formatting non-nil values with fmt.Sprintf("%v", v). Nil values map to empty string.
func ToStringMap(src map[string]interface{}) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		if v == nil {
			out[k] = ""
			continue
		}
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

// MergeDefaultsWithOverride produces a final map by starting with the defaults
// and applying any keys provided in overrideValues. Overrides are optional:
// you do not need to override every default. However, any key present in the
// override must exist in the defaults; attempting to override a non-existent
// default key is an error.
func MergeDefaultsWithOverride(defaults map[string]interface{}, overrideValues map[string]interface{}, env string) (map[string]string, error) {
	// Convert defaults to string map
	defaultMap := ToStringMap(defaults)

	// If there are no overrides, simply return a stringified copy of defaults
	if overrideValues == nil {
		final := make(map[string]string, len(defaultMap))
		for k, v := range defaultMap {
			final[k] = v
		}
		return final, nil
	}

	overrideMap := ToStringMap(overrideValues)

	// Validate that all override keys exist in defaults
	for k := range overrideMap {
		if _, ok := defaultMap[k]; !ok {
			return nil, fmt.Errorf("override for environment '%s' contains unknown key '%s' (not present in defaults)", env, k)
		}
	}

	// Start with defaults, then apply overrides
	finalMap := make(map[string]string, len(defaultMap))
	for k, v := range defaultMap {
		finalMap[k] = v
	}
	for k, v := range overrideMap {
		finalMap[k] = v
	}
	return finalMap, nil
}

// ReadTemplate reads the template file from disk and returns its contents as a string.
func ReadTemplate(templatePath string) (string, error) {
	tplB, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %w", err)
	}
	return string(tplB), nil
}

// ReplaceTemplatePlaceholders replaces placeholders in the template with values from finalMap.
// It supports both $key$ and $<key>$ variants (i.e. $<key>$).
func ReplaceTemplatePlaceholders(tpl string, finalMap map[string]string) string {
	for k, v := range finalMap {
		tpl = strings.ReplaceAll(tpl, "$"+k+"$", v)
		tpl = strings.ReplaceAll(tpl, "$<"+k+">$", v)
	}
	return tpl
}

// DetectUnreplacedPlaceholders scans the rendered template and returns a sorted,
// unique list of placeholder names that remain in the text. It recognizes both
// $key$ and $<key>$ forms. This function is useful for a "strict placeholders"
// mode where the program should fail if anything remains unreplaced.
func DetectUnreplacedPlaceholders(tpl string) []string {
	// regex captures either the <name> form in group 1 or the plain name in group 2
	re := regexp.MustCompile(`\$(?:<([^>]+)>|([A-Za-z0-9_]+))\$`)
	matches := re.FindAllStringSubmatch(tpl, -1)
	set := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		var name string
		if m[1] != "" {
			name = m[1]
		} else {
			name = m[2]
		}
		if name != "" {
			set[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// WriteOutput writes the rendered template to the specified output path using
// file mode 0644.
func WriteOutput(outPath, rendered string) error {
	if err := os.WriteFile(outPath, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}
	return nil
}

// PrintResults prints the generated dictionary (final map) and the rendered template
// to stdout in a human-friendly format.
func PrintResults(finalMap map[string]string, rendered string) {
	jsonFinal, _ := json.MarshalIndent(finalMap, "", "  ")
	fmt.Println("Generated dictionary:")
	fmt.Println(string(jsonFinal))
	fmt.Println("\nRendered template:")
	fmt.Println(rendered)
}
