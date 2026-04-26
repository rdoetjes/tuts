package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// usage prints a short help message and exits with code 2.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s -c <config.json> -e <environment> -t <template.file> -o <output.file>\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

// main is a thin CLI glue layer that delegates parsing, validation and rendering
// to helper functions in other files. It now uses named flags for clarity:
// -c/--config, -e/--environment, -t/--template, -o/--out
func main() {
	var configPath string
	var env string
	var templatePath string
	var outPath string

	// define short and long flags (both point to the same variables)
	flag.StringVar(&configPath, "c", "", "path to config JSON (required)")
	flag.StringVar(&configPath, "config", "", "path to config JSON (required)")
	flag.StringVar(&env, "e", "", "environment name (required)")
	flag.StringVar(&env, "environment", "", "environment name (required)")
	flag.StringVar(&templatePath, "t", "", "template file path (required)")
	flag.StringVar(&templatePath, "template", "", "template file path (required)")
	flag.StringVar(&outPath, "o", "", "output file path (required)")
	flag.StringVar(&outPath, "out", "", "output file path (required)")

	// keep the same usage behavior
	flag.Usage = func() { usage() }

	flag.Parse()

	// require all flags
	if configPath == "" || env == "" || templatePath == "" || outPath == "" {
		usage()
	}

	cfg, err := ParseConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var finalMap map[string]string

	// For prod-like environments (p, prod, prd) use defaults implicitly.
	// This means no override entry is required — we just render using defaults.
	if env == "p" || env == "prod" || env == "prd" {
		finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, nil, env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Try to find an explicit override. If none exists, fall back to defaults.
		override, findErr := FindOverride(cfg, env)
		if findErr != nil {
			// No override found; use defaults
			finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, nil, env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// override found; override.Values may be nil; MergeDefaultsWithOverride handles nil (uses defaults)
			finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, override.Values, env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}
	}

	templateStr, err := ReadTemplate(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rendered := ReplaceTemplatePlaceholders(templateStr, finalMap)

	if err := WriteOutput(outPath, rendered); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Print results to stdout for visibility (same behavior as original)
	PrintResults(finalMap, rendered)
}
