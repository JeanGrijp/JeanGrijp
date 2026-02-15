package svg

import (
	"fmt"
	"strings"

	"github.com/JeanGrijp/JeanGrijp/github"
)

const (
	StatsWidth  = 850
	StatsHeight = 180
)

func RenderStatsCard(stats *github.Stats, metrics []string, theme map[string]string) string {
	width, height := StatsWidth, StatsHeight
	cellWidth := float64(width) / float64(len(metrics))

	var cells, dividers []string
	for i, key := range metrics {
		cx := float64(i)*cellWidth + cellWidth/2.0

		// Icon color
		iconColor := theme["synapse_cyan"]
		if val, ok := MetricColors[key]; ok {
			if layoutColor, ok := theme[val]; ok {
				iconColor = layoutColor
			}
		}

		// Value
		val := 0
		switch key {
		case MetricCommits:
			val = stats.Commits
		case MetricStars:
			val = stats.Stars
		case MetricPRs:
			val = stats.PRs
		case MetricIssues:
			val = stats.Issues
		case MetricRepos:
			val = stats.Repos
		}
		valueStr := FormatNumber(val)

		// Label
		label := MetricLabels[key]
		if label == "" {
			label = strings.Title(key)
		}

		// Icon
		iconPath := MetricIcons[key]

		delay := fmt.Sprintf("%.1fs", float64(i)*0.3)

		cells = append(cells, fmt.Sprintf(`    <g class="metric-cell" transform="translate(%.1f, 95)">
      <g transform="translate(-8, -30) scale(1)">
        <svg viewBox="0 0 16 16" width="16" height="16" fill="%s" class="metric-icon" style="animation-delay: %s">
          %s
        </svg>
      </g>
      <text x="0" y="2" text-anchor="middle" fill="%s" font-size="28" font-weight="bold" font-family="sans-serif" opacity="0.35" filter="url(#num-glow)">%s</text>
      <text x="0" y="2" text-anchor="middle" fill="%s" font-size="28" font-weight="bold" font-family="sans-serif">%s</text>
      <text x="0" y="20" text-anchor="middle" fill="%s" font-size="11" font-family="monospace" letter-spacing="1">%s</text>
    </g>`, cx, iconColor, delay, iconPath, iconColor, valueStr, theme["text_bright"], valueStr, theme["text_faint"], label))

		if i < len(metrics)-1 {
			dx := cellWidth * float64(i+1)
			dividers = append(dividers, fmt.Sprintf(`    <line x1="%.1f" y1="55" x2="%.1f" y2="155" stroke="%s" stroke-width="1" opacity="0.5"/>`, dx, dx, theme["star_dust"]))
		}
	}

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <defs>
    <style>
      .metric-icon {
        animation: count-glow 4s ease-in-out infinite;
      }
      @keyframes count-glow {
        0%%, 100%% { fill-opacity: 0.7; }
        50%% { fill-opacity: 1; }
      }
    </style>
    <filter id="num-glow" x="-30%%" y="-30%%" width="160%%" height="160%%">
      <feGaussianBlur stdDeviation="3"/>
    </filter>
  </defs>

  <!-- Card background -->
  <rect x="0.5" y="0.5" width="%d" height="%d" rx="12" ry="12"
        fill="%s" stroke="%s" stroke-width="1"/>

  <!-- Section title -->
  <text x="30" y="38" fill="%s" font-size="11" font-family="monospace" letter-spacing="3">MISSION TELEMETRY</text>

  <!-- Dividers -->
%s

  <!-- Metric cells -->
%s
</svg>`, width, height, width, height, width-1, height-1, theme["nebula"], theme["star_dust"], theme["text_faint"], strings.Join(dividers, "\n"), strings.Join(cells, "\n"))
}
