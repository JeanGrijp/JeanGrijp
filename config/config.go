package config

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

var hexColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

type Config struct {
	Username   string                 `yaml:"username"`
	Profile    Profile                `yaml:"profile"`
	GalaxyArms []GalaxyArm            `yaml:"galaxy_arms"`
	Projects   []Project              `yaml:"projects"`
	Theme      map[string]string      `yaml:"theme"`
	Stats      StatsConfig            `yaml:"stats"`
	Languages  LanguagesConfig        `yaml:"languages"`
	Social     map[string]string      `yaml:"social"`
}

type Profile struct {
	Name       string `yaml:"name"`
	Tagline    string `yaml:"tagline"`
	Company    string `yaml:"company"`
	Location   string `yaml:"location"`
	Bio        string `yaml:"bio"`
	Philosophy string `yaml:"philosophy"`
}

type GalaxyArm struct {
	Name  string   `yaml:"name"`
	Color string   `yaml:"color"`
	Items []string `yaml:"items"`
}

type Project struct {
	Repo        string `yaml:"repo"`
	Arm         int    `yaml:"arm"`
	Description string `yaml:"description"`
}

type StatsConfig struct {
	Metrics []string `yaml:"metrics"`
}

type LanguagesConfig struct {
	Exclude    []string `yaml:"exclude"`
	MaxDisplay int      `yaml:"max_display"`
}

func ValidateAndApplyDefaults(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Username
	if cfg.Username == "" {
		return nil, fmt.Errorf("'username' is required")
	}

	// Profile
	if cfg.Profile.Name == "" {
		return nil, fmt.Errorf("'profile.name' is required")
	}

	// Galaxy Arms
	if len(cfg.GalaxyArms) == 0 {
		return nil, fmt.Errorf("'galaxy_arms' must be a non-empty list")
	}
	for i, arm := range cfg.GalaxyArms {
		if arm.Name == "" {
			return nil, fmt.Errorf("galaxy_arms[%d].name is required", i)
		}
		if arm.Color == "" {
			return nil, fmt.Errorf("galaxy_arms[%d].color is required", i)
		}
	}

	// Projects
	for i, proj := range cfg.Projects {
		if proj.Repo == "" {
			return nil, fmt.Errorf("projects[%d].repo is required", i)
		}
		if proj.Arm < 0 || proj.Arm >= len(cfg.GalaxyArms) {
			return nil, fmt.Errorf("projects[%d].arm out of range", i)
		}
	}

	// Theme Validation
	if cfg.Theme == nil {
		cfg.Theme = make(map[string]string)
	}
	for k, v := range cfg.Theme {
		if !hexColorRe.MatchString(v) {
			return nil, fmt.Errorf("theme.%s must be a valid hex color, got '%s'", k, v)
		}
	}
	// Apply defaults
	cfg.Theme = ResolveTheme(cfg.Theme)

	// Other defaults
	if cfg.Stats.Metrics == nil {
		cfg.Stats.Metrics = []string{"commits", "stars", "prs", "issues", "repos"}
	}
	if cfg.Languages.MaxDisplay == 0 {
		cfg.Languages.MaxDisplay = 8
	}

	return &cfg, nil
}
