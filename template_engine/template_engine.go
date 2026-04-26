package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	ConfigPath   string
	Env          string
	TemplatePath string
	OutPath      string
	Strict       bool
}

// usage prints a short help message and exits with code 2.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s -c <config.json> -e <environment> -t <template.file> -o <output.file>\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func parseArguments(args *Args) {
	// define short and long flags (both point to the same variables on the args struct)
	flag.StringVar(&args.ConfigPath, "c", "", "path to config JSON (required)")
	flag.StringVar(&args.ConfigPath, "config", "", "path to config JSON (required)")
	flag.StringVar(&args.Env, "e", "", "environment name (required)")
	flag.StringVar(&args.Env, "environment", "", "environment name (required)")
	flag.StringVar(&args.TemplatePath, "t", "", "template file path (required)")
	flag.StringVar(&args.TemplatePath, "template", "", "template file path (required)")
	flag.StringVar(&args.OutPath, "o", "", "output file path (required)")
	flag.StringVar(&args.OutPath, "out", "", "output file path (required)")
	flag.BoolVar(&args.Strict, "s", false, "fail if unreplaced placeholders remain")
	flag.BoolVar(&args.Strict, "strict-placeholders", false, "fail if unreplaced placeholders remain")

	// keep the same usage behavior
	flag.Usage = func() { usage() }

	flag.Parse()

	// require all flags
	if args.ConfigPath == "" || args.Env == "" || args.TemplatePath == "" || args.OutPath == "" {
		usage()
	}
}

// main is a thin CLI glue layer that delegates parsing, validation and rendering
// to helper functions in other files. It now uses named flags for clarity:
// -c/--config, -e/--environment, -t/--template, -o/--out
func main() {
	args := Args{}
	parseArguments(&args)

	cfg, err := ParseConfig(args.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var finalMap map[string]string

	// For prod-like environments (p, prod, prd) use defaults implicitly.
	// This means no override entry is required — we just render using defaults.
	if args.Env == "p" || args.Env == "prod" || args.Env == "prd" {
		finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, nil, args.Env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Try to find an explicit override. If none exists, fall back to defaults.
		override, findErr := FindOverride(cfg, args.Env)
		if findErr != nil {
			// No override found; use defaults
			finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, nil, args.Env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// override found; override.Values may be nil; MergeDefaultsWithOverride handles nil (uses defaults)
			finalMap, err = MergeDefaultsWithOverride(cfg.Defaults, override.Values, args.Env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}
	}

	templateStr, err := ReadTemplate(args.TemplatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rendered := ReplaceTemplatePlaceholders(templateStr, finalMap)

	if args.Strict {
		placeholders := DetectUnreplacedPlaceholders(rendered)
		if len(placeholders) > 0 {
			fmt.Fprintf(os.Stderr, "error: unreplaced placeholders: %v\n", placeholders)
			os.Exit(1)
		}
	}

	if err := WriteOutput(args.OutPath, rendered); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Print results to stdout for visibility (same behavior as original)
	PrintResults(finalMap, rendered)
}
