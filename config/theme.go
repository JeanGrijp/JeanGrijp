package config

// DefaultTheme defines the deep-space palette
var DefaultTheme = map[string]string{
	"void":            "#080c14",
	"nebula":          "#0f1623",
	"star_dust":       "#1a2332",
	"synapse_cyan":    "#00d4ff",
	"dendrite_violet": "#a78bfa",
	"axon_amber":      "#ffb020",
	"text_bright":     "#f1f5f9",
	"text_dim":        "#94a3b8",
	"text_faint":      "#64748b",
}

// ResolveTheme merges user overrides with default theme
func ResolveTheme(userTheme map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range DefaultTheme {
		merged[k] = v
	}
	for k, v := range userTheme {
		merged[k] = v
	}
	return merged
}

// ResolveArmColors returns a list of hex color strings, one per arm
func ResolveArmColors(arms []GalaxyArm, theme map[string]string) []string {
	colors := make([]string, len(arms))
	for i, arm := range arms {
		// Try to find the color in the theme (e.g. "synapse_cyan")
		// If not found in theme, assume it's a direct hex code or fallback
		if val, ok := theme[arm.Color]; ok {
			colors[i] = val
		} else {
			// Fallback to synapse_cyan if the key config is weird,
			// or maybe the user provided a raw hex?
			// The python code did: theme.get(arm.color, theme.get("synapse_cyan"))
			// which implies arm.color is a key in theme.
			if defaultColor, ok := theme["synapse_cyan"]; ok {
				colors[i] = defaultColor
			} else {
				colors[i] = "#00d4ff"
			}
		}
	}
	return colors
}
