package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	return &cfg, nil
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
