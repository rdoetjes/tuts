package main

// Override represents environment-specific values that override defaults.
type Override struct {
	Environment string                 `json:"environment"`
	Values      map[string]interface{} `json:"values"`
}

// Config is the top-level configuration structure containing defaults
// and a list of environment-specific overrides.
type Config struct {
	Defaults  map[string]interface{} `json:"defaults"`
	Overrides []Override             `json:"overrides"`
}
