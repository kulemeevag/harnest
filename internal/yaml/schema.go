// Package yaml defines the schema for harnest.yaml — the declarative
// configuration file that serves as the source of truth for the Harnest CLI.
package yaml

// HarnestConfig is the top-level structure of a harnest.yaml file.
type HarnestConfig struct {
	Version      int           `yaml:"version"`
	Project      ProjectInfo   `yaml:"project,omitempty"`
	Stacks       []StackEntry  `yaml:"stacks,omitempty"`
	Agents       AgentsBlock   `yaml:"agents"`
	Harnesses    []string      `yaml:"harnesses"`
	DesignSystem string        `yaml:"design_system,omitempty"`
	Profiles     ProfilesBlock `yaml:"profiles,omitempty"`
	Settings     SettingsBlock `yaml:"settings,omitempty"`
}

// ProjectInfo holds optional human-readable metadata about the project.
type ProjectInfo struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// StackEntry represents a single detected or manually specified technology stack.
type StackEntry struct {
	Name     string `yaml:"name"`
	Lang     string `yaml:"lang"`
	Category string `yaml:"category"`
	Path     string `yaml:"path"`
}

// AgentsBlock configures both consilium (advisory) and executing agents.
type AgentsBlock struct {
	Consilium map[string]string `yaml:"consilium"`        // role -> agent name
	Executing []ExecEntry       `yaml:"executing"`
	Models    map[string]string `yaml:"models,omitempty"` // role -> capability tier (high/medium/low)
}

// ExecEntry maps a specific agent to a file glob scope for executing tasks.
type ExecEntry struct {
	Agent string `yaml:"agent"`
	Scope string `yaml:"scope"`
}

// ProfilesBlock controls which profiles are active and allows custom profile definitions.
type ProfilesBlock struct {
	Enabled []string        `yaml:"enabled,omitempty"`
	Custom  []CustomProfile `yaml:"custom,omitempty"`
}

// CustomProfile references a user-defined profile by name and file path.
type CustomProfile struct {
	Name string `yaml:"name"`
	File string `yaml:"file"`
}

// SettingsBlock holds global behavioral settings for the Harnest CLI.
type SettingsBlock struct {
	// AutoDetect enables automatic stack detection when stacks are empty or unset.
	AutoDetect bool `yaml:"auto_detect,omitempty"`
	// StackStrategy controls how detected stacks are merged with stacks declared in config.
	// Valid values: "replace" (default) or "merge".
	StackStrategy string `yaml:"stack_strategy,omitempty"`
	// LockFile enables writing a harnest.lock file after generation.
	LockFile bool `yaml:"lock_file,omitempty"`
}
