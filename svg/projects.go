package svg

import (
	"fmt"
	"math"
	"strings"

	"github.com/JeanGrijp/JeanGrijp/config"
)

const (
	ProjectsWidth  = 850
	ProjectsHeight = 220
)

func RenderProjectsConstellation(projects []config.Project, galaxyArms []config.GalaxyArm, theme map[string]string) string {
	width, height := ProjectsWidth, ProjectsHeight
	armColors := config.ResolveArmColors(galaxyArms, theme)

	n := len(projects)
	if n > 3 {
		n = 3
	}
	// Adaptive card sizing
	cardWidth := 240.0
	if n == 2 {
		cardWidth = 340.0
	}
	totalCardsWidth := cardWidth * float64(n)
	gap := (float64(width) - totalCardsWidth) / float64(n+1)

	var cardColors []string
	for i := 0; i < n; i++ {
		proj := projects[i]
		armIdx := proj.Arm
		if armIdx < 0 || armIdx >= len(galaxyArms) {
			armIdx = 0
		}
		cardColors = append(cardColors, armColors[armIdx])
	}

	defsStr := buildDefs(n, cardWidth, gap, cardColors, theme)
	bg := fmt.Sprintf(`<rect x="0.5" y="0.5" width="%d" height="%d" rx="12" ry="12" fill="%s" stroke="%s" stroke-width="1"/>`, width-1, height-1, theme["nebula"], theme["star_dust"])
	starsStr := buildStarfield(n, width, height, cardColors, theme)
	gridStr := buildGridOverlay(width, height, theme)
	connStr := buildConnections(n, cardWidth, gap, cardColors)
	titleStr := buildTitleArea(n, width, height, theme)

	var cards []string
	for i := 0; i < n; i++ {
		proj := projects[i]
		armIdx := proj.Arm
		var arm config.GalaxyArm
		if armIdx >= 0 && armIdx < len(galaxyArms) {
			arm = galaxyArms[armIdx]
		} else {
			arm = galaxyArms[0]
		}
		color := cardColors[i]

		cardX := gap + float64(i)*(cardWidth+gap)
		cardCx := cardX + cardWidth/2.0
		cards = append(cards, buildProjectCard(i, proj, arm, color, cardWidth, cardCx, cardX, theme))
	}
	cardsStr := strings.Join(cards, "\n")
	scanLine := buildScanLine(width, theme)

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <defs>
%s
  </defs>

  <!-- Background -->
%s

  <!-- Star field -->
%s

  <!-- Grid overlay -->
%s

  <!-- Connection lines -->
%s

  <!-- Title area -->
%s

  <!-- Project cards -->
%s

  <!-- Global scan line -->
%s
</svg>`, width, height, width, height, defsStr, bg, starsStr, gridStr, connStr, titleStr, cardsStr, scanLine)
}

func buildDefs(n int, cardWidth, gap float64, cardColors []string, theme map[string]string) string {
	var parts []string
	for i := 0; i < n; i++ {
		color := cardColors[i]
		parts = append(parts, fmt.Sprintf(`    <filter id="proj-glow-%d" x="-80%%" y="-80%%" width="260%%" height="260%%">
      <feGaussianBlur stdDeviation="4" in="SourceGraphic" result="blur"/>
      <feFlood flood-color="%s" flood-opacity="0.6" result="color"/>
      <feComposite in="color" in2="blur" operator="in" result="glow"/>
      <feMerge>
        <feMergeNode in="glow"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>`, i, color))
	}
	parts = append(parts, `    <filter id="card-nebula" x="-50%" y="-50%" width="200%" height="200%">
      <feGaussianBlur stdDeviation="15"/>
    </filter>`)

	for i := 0; i < n; i++ {
		parts = append(parts, fmt.Sprintf(`    <linearGradient id="card-bg-%d" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.6"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.9"/>
    </linearGradient>`, i, theme["star_dust"], theme["nebula"]))
	}

	if n >= 2 {
		parts = append(parts, fmt.Sprintf(`    <linearGradient id="conn-grad" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.4"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0.4"/>
    </linearGradient>`, cardColors[0], cardColors[n-1]))
	}

	for i := 0; i < n; i++ {
		cardX := gap + float64(i)*(cardWidth+gap)
		parts = append(parts, fmt.Sprintf(`    <clipPath id="card-clip-%d">
      <rect x="%.1f" y="55" width="%.1f" height="140" rx="8" ry="8"/>
    </clipPath>`, i, cardX, cardWidth))
	}

	parts = append(parts, `    <style>
      @keyframes twinkle {
        0%, 100% { opacity: 0.1; }
        50% { opacity: 0.6; }
      }
      @keyframes orbit {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
      }
      @keyframes card-appear {
        from { opacity: 0; transform: translateY(8px); }
        to { opacity: 1; transform: translateY(0); }
      }
      @keyframes scan-sweep {
        0% { transform: translateY(0); }
        100% { transform: translateY(160px); }
      }
    </style>`)

	return strings.Join(parts, "\n")
}

func buildStarfield(n, width, height int, cardColors []string, theme map[string]string) string {
	var stars []string
	// 15 faint bg stars
	sx := DeterministicRandom("proj-star-x", 15, 10, float64(width)-10)
	sy := DeterministicRandom("proj-star-y", 15, 10, float64(height)-10)
	sr := DeterministicRandom("proj-star-r", 15, 0.3, 0.9)
	so := DeterministicRandom("proj-star-o", 15, 0.05, 0.25)
	sd := DeterministicRandom("proj-star-d", 15, 5.0, 8.0)

	for i := 0; i < 15; i++ {
		fill := theme["text_dim"]
		if i%4 == 0 {
			fill = cardColors[i%n]
		}
		stars = append(stars, fmt.Sprintf(`  <circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" opacity="%.2f">
<animate attributeName="opacity" values="%.2f;%.2f;%.2f" dur="%.1fs" repeatCount="indefinite"/>
</circle>`, sx[i], sy[i], sr[i], fill, so[i], so[i], math.Min(so[i]*3, 0.6), so[i], sd[i]))
	}

	// 10 mid-ground stars
	mx := DeterministicRandom("proj-mstar-x", 10, 15, float64(width)-15)
	my := DeterministicRandom("proj-mstar-y", 10, 15, float64(height)-15)
	mr := DeterministicRandom("proj-mstar-r", 10, 0.5, 1.2)
	mo := DeterministicRandom("proj-mstar-o", 10, 0.10, 0.40)
	md := DeterministicRandom("proj-mstar-d", 10, 3.0, 6.0)

	for i := 0; i < 10; i++ {
		fill := theme["text_dim"]
		if i%4 == 0 {
			fill = cardColors[i%n]
		}
		stars = append(stars, fmt.Sprintf(`  <circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" opacity="%.2f">
<animate attributeName="opacity" values="%.2f;%.2f;%.2f" dur="%.1fs" repeatCount="indefinite"/>
</circle>`, mx[i], my[i], mr[i], fill, mo[i], mo[i], math.Min(mo[i]*2.5, 0.8), mo[i], md[i]))
	}
	return strings.Join(stars, "\n")
}

func buildGridOverlay(width, height int, theme map[string]string) string {
	var lines []string
	for y := 40; y < height; y += 40 {
		lines = append(lines, fmt.Sprintf(`  <line x1="12" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="0.5" stroke-dasharray="4,8" opacity="0.12"/>`, y, width-12, y, theme["text_faint"]))
	}
	for x := 80; x < width; x += 80 {
		lines = append(lines, fmt.Sprintf(`  <line x1="%d" y1="12" x2="%d" y2="%d" stroke="%s" stroke-width="0.5" stroke-dasharray="4,8" opacity="0.08"/>`, x, x, height-12, theme["text_faint"]))
	}
	return strings.Join(lines, "\n")
}

func buildConnections(n int, cardWidth, gap float64, cardColors []string) string {
	var lines []string
	if n >= 2 {
		for i := 0; i < n-1; i++ {
			x1 := gap + float64(i)*(cardWidth+gap) + cardWidth/2.0
			x2 := gap + float64(i+1)*(cardWidth+gap) + cardWidth/2.0
			lines = append(lines, fmt.Sprintf(`  <line x1="%.1f" y1="85" x2="%.1f" y2="85" stroke="url(#conn-grad)" stroke-width="1" stroke-dasharray="6,4" opacity="0.5"/>`, x1, x2))
		}
	}
	return strings.Join(lines, "\n")
}

func buildTitleArea(n, width, height int, theme map[string]string) string {
	var parts []string
	bk := theme["text_faint"]
	bl := 16.0
	// Corner brackets
	parts = append(parts, fmt.Sprintf(`  <g opacity="0.4">
    <polyline points="5,%.1f 5,5 %.1f,5" fill="none" stroke="%s" stroke-width="1.5"/>
    <polyline points="%d,5 %d,5 %d,%.1f" fill="none" stroke="%s" stroke-width="1.5"/>
    <polyline points="5,%.1f 5,%d %.1f,%d" fill="none" stroke="%s" stroke-width="1.5"/>
    <polyline points="%d,%d %d,%d %d,%.1f" fill="none" stroke="%s" stroke-width="1.5"/>
  </g>`, bl+5, bl+5, bk, width-int(bl)-5, width-5, width-5, bl+5, bk, float64(height)-bl-5, height-5, bl+5, height-5, bk, width-int(bl)-5, height-5, width-5, height-5, width-5, float64(height)-bl-5, bk))

	parts = append(parts, fmt.Sprintf(`  <text x="30" y="38" fill="%s" font-size="11" font-family="monospace" letter-spacing="3">FEATURED SYSTEMS</text>`, theme["text_faint"]))
	cyan := theme["synapse_cyan"]
	parts = append(parts, fmt.Sprintf(`  <circle cx="218" cy="34" r="3" fill="%s" opacity="0.8">
<animate attributeName="opacity" values="0.4;1;0.4" dur="2s" repeatCount="indefinite"/>
</circle>`, cyan))
	parts = append(parts, fmt.Sprintf(`  <text x="%d" y="38" fill="%s" font-size="10" font-family="monospace" text-anchor="end" opacity="0.5">SYS %d/%d ONLINE</text>`, width-30, theme["text_faint"], n, n))

	return strings.Join(parts, "\n")
}

func buildProjectCard(i int, proj config.Project, arm config.GalaxyArm, color string, cardWidth, cardCx, cardX float64, theme map[string]string) string {
	repoName := proj.Repo
	if strings.Contains(repoName, "/") {
		parts := strings.Split(repoName, "/")
		repoName = parts[len(parts)-1]
	}
	desc := proj.Description
	maxChars := int(cardWidth / 7.5)
	descLines := WrapText(desc, maxChars)

	delay := fmt.Sprintf("%.1fs", float64(i)*0.3)
	parts := []string{}
	parts = append(parts, fmt.Sprintf(`  <g opacity="0" style="animation: card-appear 0.6s ease %s forwards">`, delay))
	parts = append(parts, fmt.Sprintf(`    <rect x="%.1f" y="55" width="%.1f" height="140" rx="8" ry="8" fill="url(#card-bg-%d)" stroke="%s" stroke-width="1"/>`, cardX, cardWidth, i, theme["star_dust"]))

	parts = append(parts, fmt.Sprintf(`    <g clip-path="url(#card-clip-%d)">`, i))
	parts = append(parts, fmt.Sprintf(`      <circle cx="%.1f" cy="90" r="50" fill="%s" opacity="0.025" filter="url(#card-nebula)"/>`, cardX+cardWidth*0.3, color))
	parts = append(parts, fmt.Sprintf(`      <circle cx="%.1f" cy="150" r="40" fill="%s" opacity="0.03" filter="url(#card-nebula)"/>`, cardX+cardWidth*0.7, color))
	parts = append(parts, fmt.Sprintf(`      <rect x="%.1f" y="55" width="%.1f" height="2" fill="%s" opacity="0.1">
<animateTransform attributeName="transform" type="translate" from="0 0" to="0 140" dur="6s" repeatCount="indefinite"/>
</rect>`, cardX, cardWidth, color))
	parts = append(parts, `    </g>`)

	parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="85" r="14" fill="none" stroke="%s" stroke-width="0.8" stroke-dasharray="4,3" opacity="0.5">
<animateTransform attributeName="transform" type="rotate" from="0 %.1f 85" to="360 %.1f 85" dur="12s" repeatCount="indefinite"/>
</circle>`, cardCx, color, cardCx, cardCx))
	parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="85" r="8" fill="%s" opacity="0.15" filter="url(#proj-glow-%d)"/>`, cardCx, color, i))

	parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="85" r="5" fill="%s" opacity="0.7">
<animate attributeName="opacity" values="0.5;0.9;0.5" dur="3s" begin="%s" repeatCount="indefinite"/>
<animate attributeName="r" values="4.5;5.5;4.5" dur="3s" begin="%s" repeatCount="indefinite"/>
</circle>`, cardCx, color, delay, delay))
	parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="85" r="2" fill="#ffffff" opacity="0.9"/>`, cardCx))

	parts = append(parts, fmt.Sprintf(`    <text x="%.1f" y="111" fill="%s" font-size="14" font-weight="bold" font-family="sans-serif" text-anchor="middle">%s</text>`, cardCx, theme["text_bright"], Esc(repoName)))

	for j := 0; j < len(descLines) && j < 2; j++ {
		line := descLines[j]
		yPos := 129 + j*15
		parts = append(parts, fmt.Sprintf(`    <text x="%.1f" y="%d" fill="%s" font-size="11" font-family="sans-serif" text-anchor="middle">%s</text>`, cardCx, yPos, theme["text_dim"], Esc(line)))
	}

	tagText := arm.Name
	tagWidth := float64(len(tagText)*7 + 16)
	tagX := cardCx - tagWidth/2.0
	parts = append(parts, fmt.Sprintf(`    <rect x="%.1f" y="163" width="%.1f" height="18" rx="9" ry="9" fill="%s" opacity="0.12"/>`, tagX, tagWidth, color))
	parts = append(parts, fmt.Sprintf(`    <text x="%.1f" y="175" fill="%s" font-size="9" font-family="monospace" text-anchor="middle" opacity="0.85">%s</text>`, cardCx, color, Esc(tagText)))

	parts = append(parts, `  </g>`)
	return strings.Join(parts, "\n")
}

func buildScanLine(width int, theme map[string]string) string {
	cyan := theme["synapse_cyan"]
	return fmt.Sprintf(`  <rect x="12" y="50" width="%d" height="1.5" fill="%s" opacity="0.08">
<animateTransform attributeName="transform" type="translate" from="0 0" to="0 160" dur="6s" repeatCount="indefinite"/>
</rect>`, width-24, cyan)
}
