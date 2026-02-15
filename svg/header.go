package svg

import (
	"fmt"
	"math"
	"strings"

	"github.com/JeanGrijp/JeanGrijp/config"
)

const (
	HeaderWidth       = 850
	HeaderHeight      = 280
	HeaderCenterX     = 425.0
	HeaderCenterY     = 155.0
	HeaderMaxRadius   = 220.0
	HeaderSpiralTurns = 0.85
	HeaderNumPoints   = 30
	HeaderXScale      = 1.5
	HeaderYScale      = 0.38
)

var HeaderStartAngles = []float64{25, 150, 265}

func RenderGalaxyHeader(cfg *config.Config, theme map[string]string) string {
	width, height := HeaderWidth, HeaderHeight
	cx, cy := HeaderCenterX, HeaderCenterY

	username := cfg.Username
	if username == "" {
		username = "user"
	}
	profile := cfg.Profile
	name := profile.Name
	if name == "" {
		name = username
	}
	tagline := profile.Tagline
	philosophy := profile.Philosophy
	initial := "?"
	if len(name) > 0 {
		initial = string(name[0])
	}
	initial = strings.ToUpper(initial)

	galaxyArms := cfg.GalaxyArms
	projects := cfg.Projects
	armColors := config.ResolveArmColors(galaxyArms, theme)

	// Spiral geometry
	var allArmPoints [][]Point
	for i := range galaxyArms {
		angle := HeaderStartAngles[i%len(HeaderStartAngles)]
		pts := SpiralPoints(cx, cy, angle, HeaderNumPoints, HeaderMaxRadius, HeaderSpiralTurns, HeaderXScale, HeaderYScale)
		allArmPoints = append(allArmPoints, pts)
	}

	glowFiltersStr := buildHeaderGlowFilters(galaxyArms, armColors)

	labelGlowFilter := `    <filter id="label-glow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur stdDeviation="2" result="blur"/>
    </filter>`

	coreGlowFilter := `    <filter id="core-bright-glow" x="-100%" y="-100%" width="300%" height="300%">
      <feGaussianBlur stdDeviation="4"/>
    </filter>`

	starsStr := buildHeaderStarfield(username, width, height, theme)
	outerNebula, innerNebula := buildNebulae(cx, cy, theme)
	shootStarsStr := buildShootingStars()
	armPathsStr, armParticlesStr := buildSpiralArms(galaxyArms, armColors, allArmPoints)
	armDotsStr := buildTechLabels(galaxyArms, armColors, allArmPoints, cx, cy)
	projectStarsStr := buildHeaderProjectStars(projects, galaxyArms, armColors, allArmPoints)
	orbitalRings := buildOrbitalRings(cx, cy, theme)
	core := buildGalaxyCore(cx, cy, theme, initial)

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <defs>
    <style>
      .star-bg { animation: twinkle-slow 7s ease-in-out infinite; }
      .star-mid { animation: twinkle-mid 5s ease-in-out infinite; }
      .star-fg { animation: twinkle-fast 3s ease-in-out infinite; }
      @keyframes twinkle-slow { 0%%, 100%% { opacity: 0.08; } 50%% { opacity: 0.3; } }
      @keyframes twinkle-mid { 0%%, 100%% { opacity: 0.15; } 50%% { opacity: 0.5; } }
      @keyframes twinkle-fast { 0%%, 100%% { opacity: 0.4; } 50%% { opacity: 0.8; } }
      .core-ring { animation: pulse-core 3s ease-in-out infinite; }
      .core-ring-inner { animation: pulse-core 3s ease-in-out infinite 1.5s; }
      @keyframes pulse-core {
        0%%, 100%% { stroke-opacity: 0.3; transform: scale(1); transform-origin: %.1fpx %.1fpx; }
        50%% { stroke-opacity: 0.8; transform: scale(1.06); transform-origin: %.1fpx %.1fpx; }
      }
      .shooting-star { opacity: 0; animation: shoot linear infinite; }
      @keyframes shoot {
        0%% { opacity: 0; transform: translate(0, 0); }
        5%% { opacity: 0.9; }
        15%% { opacity: 0.6; transform: translate(var(--shoot-tx), var(--shoot-ty)); }
        20%% { opacity: 0; transform: translate(var(--shoot-tx), var(--shoot-ty)); }
        100%% { opacity: 0; }
      }
    </style>

    <filter id="nebula-outer">
      <feGaussianBlur stdDeviation="60"/>
    </filter>
    <filter id="nebula-inner">
      <feGaussianBlur stdDeviation="30"/>
    </filter>

%s
%s

    <radialGradient id="core-haze-gradient" cx="50%%" cy="50%%" r="50%%">
      <stop offset="0%%" stop-color="%s" stop-opacity="0.5"/>
      <stop offset="50%%" stop-color="%s" stop-opacity="0.2"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0"/>
    </radialGradient>

    <radialGradient id="core-inner-gradient" cx="50%%" cy="50%%" r="50%%">
      <stop offset="0%%" stop-color="#ffffff" stop-opacity="0.6"/>
      <stop offset="40%%" stop-color="%s" stop-opacity="0.3"/>
      <stop offset="100%%" stop-color="%s" stop-opacity="0"/>
    </radialGradient>

    <linearGradient id="shoot-grad" x1="0%%" y1="0%%" x2="100%%" y2="0%%">
      <stop offset="0%%" stop-color="#ffffff" stop-opacity="0.8"/>
      <stop offset="100%%" stop-color="#ffffff" stop-opacity="0"/>
    </linearGradient>

%s
  </defs>

  <!-- 1. Background -->
  <rect x="0" y="0" width="%d" height="%d" rx="12" ry="12" fill="%s"/>

  <!-- 2. Outer nebula -->
%s

  <!-- 3. Star field -->
%s

  <!-- 4. Inner nebula -->
%s

  <!-- 5. Shooting stars -->
%s

  <!-- 6. Spiral arm paths -->
%s

  <!-- 7. Arm particles -->
%s

  <!-- 8. Tech dots + leader lines + labels -->
%s

  <!-- 9. Project stars -->
%s

  <!-- 10. Orbital rings -->
%s

  <!-- 11. Galaxy core -->
%s

  <!-- 12. Profile text -->
  <text x="%.1f" y="26" text-anchor="middle" fill="%s" font-size="20" font-weight="bold" font-family="sans-serif">%s</text>
  <text x="%.1f" y="44" text-anchor="middle" fill="%s" font-size="12" font-family="sans-serif">%s</text>
  <text x="%.1f" y="%d" text-anchor="middle" fill="%s" font-size="11" font-family="monospace" font-style="italic">%s</text>
</svg>`, width, height, width, height, cx, cy, cx, cy, labelGlowFilter, coreGlowFilter,
		theme["synapse_cyan"], theme["dendrite_violet"], theme["synapse_cyan"],
		theme["synapse_cyan"], theme["synapse_cyan"],
		glowFiltersStr,
		width, height, theme["void"],
		outerNebula,
		starsStr,
		innerNebula,
		shootStarsStr,
		armPathsStr,
		armParticlesStr,
		armDotsStr,
		projectStarsStr,
		orbitalRings,
		core,
		cx, theme["text_bright"], Esc(name),
		cx, theme["text_dim"], Esc(tagline),
		cx, height-12, theme["text_faint"], Esc(philosophy))
}

func buildHeaderGlowFilters(galaxyArms []config.GalaxyArm, armColors []string) string {
	var parts []string
	for i := range galaxyArms {
		color := armColors[i]
		parts = append(parts, fmt.Sprintf(`    <filter id="star-glow-%d" x="-100%%" y="-100%%" width="300%%" height="300%%">
      <feGaussianBlur stdDeviation="3" result="blur"/>
      <feFlood flood-color="%s" flood-opacity="0.5" result="color"/>
      <feComposite in="color" in2="blur" operator="in" result="glow"/>
      <feMerge>
        <feMergeNode in="glow"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>`, i, color))
	}
	return strings.Join(parts, "\n")
}

func buildHeaderStarfield(username string, width, height int, theme map[string]string) string {
	type layer struct {
		Count  int
		RMin   float64
		RMax   float64
		OMin   float64
		OMax   float64
		DurMin float64
		DurMax float64
		Label  string
	}
	layers := []layer{
		{40, 0.3, 0.8, 0.08, 0.3, 5.0, 9.0, "bg"},
		{20, 0.6, 1.2, 0.15, 0.5, 3.5, 7.0, "mid"},
		{10, 1.0, 1.8, 0.4, 0.7, 2.0, 4.5, "fg"},
	}

	var stars []string
	for _, l := range layers {
		sx := DeterministicRandom(fmt.Sprintf("%s_sx_%s", username, l.Label), l.Count, 10, float64(width)-10)
		sy := DeterministicRandom(fmt.Sprintf("%s_sy_%s", username, l.Label), l.Count, 10, float64(height)-10)
		sr := DeterministicRandom(fmt.Sprintf("%s_sr_%s", username, l.Label), l.Count, l.RMin, l.RMax)
		so := DeterministicRandom(fmt.Sprintf("%s_so_%s", username, l.Label), l.Count, l.OMin, l.OMax)
		sd := DeterministicRandom(fmt.Sprintf("%s_sd_%s", username, l.Label), l.Count, l.DurMin, l.DurMax)

		for i := 0; i < l.Count; i++ {
			fill := "#ffffff"
			if i%12 == 0 {
				fill = theme["synapse_cyan"]
			} else if i%12 == 4 {
				fill = theme["dendrite_violet"]
			} else if i%12 == 8 {
				fill = theme["axon_amber"]
			}
			// Default logic fallback to white
			if fill == "" {
				fill = "#ffffff"
			}

			delay := fmt.Sprintf("%.1fs", sd[i]*0.3)
			stars = append(stars, fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="%.2f" fill="%s" opacity="%.2f" class="star-%s" style="animation-delay: %s"/>`,
				sx[i], sy[i], sr[i], fill, so[i], l.Label, delay))
		}
	}
	return strings.Join(stars, "\n")
}

func buildNebulae(cx, cy float64, theme map[string]string) (string, string) {
	outer := fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="120" fill="%s" opacity="0.015" filter="url(#nebula-outer)"/>
    <circle cx="%.1f" cy="%.1f" r="100" fill="%s" opacity="0.012" filter="url(#nebula-outer)"/>
    <circle cx="%.1f" cy="%.1f" r="140" fill="%s" opacity="0.01" filter="url(#nebula-outer)"/>`,
		cx-180, cy-30, theme["dendrite_violet"],
		cx+200, cy+20, theme["axon_amber"],
		cx, cy+40, theme["synapse_cyan"])

	inner := fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="70" fill="%s" opacity="0.04" filter="url(#nebula-inner)"/>
    <circle cx="%.1f" cy="%.1f" r="50" fill="%s" opacity="0.035" filter="url(#nebula-inner)"/>
    <circle cx="%.1f" cy="%.1f" r="45" fill="%s" opacity="0.03" filter="url(#nebula-inner)"/>`,
		cx, cy, theme["synapse_cyan"],
		cx-60, cy-20, theme["dendrite_violet"],
		cx+70, cy+15, theme["axon_amber"])

	return outer, inner
}

func buildShootingStars() string {
	data := []struct {
		sx, sy, tx, ty, dur float64
	}{
		{120, 30, 200, 80, 6},
		{650, 20, 180, 70, 8},
		{400, 250, 160, 60, 7},
	}
	var stars []string
	for i, d := range data {
		stars = append(stars, fmt.Sprintf(`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="url(#shoot-grad)" stroke-width="1.2" stroke-linecap="round" class="shooting-star" style="animation-delay: %.1fs; --shoot-tx: %.1fpx; --shoot-ty: %.1fpx; animation-duration: %.1fs"/>`,
			d.sx, d.sy, d.sx+20, d.sy+5, float64(i)*2.5, d.tx, d.ty, d.dur))
	}
	return strings.Join(stars, "\n")
}

func buildSpiralArms(galaxyArms []config.GalaxyArm, armColors []string, allArmPoints [][]Point) (string, string) {
	var paths, particles []string
	segmentCount := 4
	opacitySteps := []float64{0.50, 0.40, 0.30, 0.20}
	widthSteps := []float64{2.0, 1.7, 1.4, 1.1}

	for armIdx, armPoints := range allArmPoints {
		if len(armPoints) < 2 {
			continue
		}
		color := armColors[armIdx]

		fullPathD := fmt.Sprintf("M %.1f %.1f", armPoints[0].X, armPoints[0].Y)
		for j := 1; j < len(armPoints); j++ {
			px, py := armPoints[j-1].X, armPoints[j-1].Y
			x, y := armPoints[j].X, armPoints[j].Y
			cpx := (px + x) / 2
			cpy := (py + y) / 2
			fullPathD += fmt.Sprintf(" Q %.1f %.1f %.1f %.1f", px, py, cpx, cpy)
		}
		fullPathD += fmt.Sprintf(" L %.1f %.1f", armPoints[len(armPoints)-1].X, armPoints[len(armPoints)-1].Y)

		ptsPerSeg := len(armPoints) / segmentCount
		for seg := 0; seg < segmentCount; seg++ {
			startI := seg * ptsPerSeg
			endI := startI + ptsPerSeg + 1
			if endI > len(armPoints) {
				endI = len(armPoints)
			}
			segPts := armPoints[startI:endI]
			if len(segPts) < 2 {
				continue
			}

			segD := fmt.Sprintf("M %.1f %.1f", segPts[0].X, segPts[0].Y)
			for j := 1; j < len(segPts); j++ {
				ppx, ppy := segPts[j-1].X, segPts[j-1].Y
				sxp, syp := segPts[j].X, segPts[j].Y
				cpx := (ppx + sxp) / 2
				cpy := (ppy + syp) / 2
				segD += fmt.Sprintf(" Q %.1f %.1f %.1f %.1f", ppx, ppy, cpx, cpy)
			}
			segD += fmt.Sprintf(" L %.1f %.1f", segPts[len(segPts)-1].X, segPts[len(segPts)-1].Y)

			op := opacitySteps[seg]
			sw := widthSteps[seg]
			paths = append(paths, fmt.Sprintf(`    <path d="%s" fill="none" stroke="%s" stroke-width="%.1f" opacity="%.2f" stroke-linecap="round">
      <animate attributeName="opacity" values="%.2f;%.2f;%.2f" dur="8s" begin="%ds" repeatCount="indefinite"/>
    </path>`, segD, color, sw, op, op-0.1, op+0.1, op-0.1, armIdx))
		}

		for pIdx := 0; pIdx < 2; pIdx++ {
			delay := float64(armIdx*4 + pIdx*6)
			particles = append(particles, fmt.Sprintf(`    <circle r="1.5" fill="%s" opacity="0.6">
      <animateMotion dur="12s" begin="%.1fs" repeatCount="indefinite" path="%s"/>
      <animate attributeName="opacity" values="0;0.7;0.3;0" dur="12s" begin="%.1fs" repeatCount="indefinite"/>
    </circle>`, color, delay, fullPathD, delay))
		}
	}
	return strings.Join(paths, "\n"), strings.Join(particles, "\n")
}

func buildTechLabels(galaxyArms []config.GalaxyArm, armColors []string, allArmPoints [][]Point, cx, cy float64) string {
	var dots []string
	outerStart := 8

	for armIdx, arm := range galaxyArms {
		color := armColors[armIdx]
		points := allArmPoints[armIdx]
		items := arm.Items
		if len(items) == 0 {
			continue
		}

		available := len(points) - outerStart - 2
		spacing := 1
		if len(items) > 0 {
			spacing = available / len(items)
		}
		if spacing < 1 {
			spacing = 1
		}

		for i, item := range items {
			ptIdx := outerStart + i*spacing
			if ptIdx >= len(points) {
				ptIdx = len(points) - 1
			}
			px, py := points[ptIdx].X, points[ptIdx].Y

			dx := px - cx
			dy := py - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist == 0 {
				dist = 1
			}
			nx := dx / dist
			ny := dy / dist

			labelX := px + nx*18
			labelY := py + ny*18

			anchor := "middle"
			if dx > 20 {
				anchor = "start"
			} else if dx < -20 {
				anchor = "end"
			}

			dots = append(dots, fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="2.5" fill="%s" opacity="0.85">
      <animate attributeName="opacity" values="0.85;1;0.85" dur="5s" begin="%.1fs" repeatCount="indefinite"/>
    </circle>`, px, py, color, float64(i)*0.7))

			dots = append(dots, fmt.Sprintf(`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="0.5" opacity="0.25" stroke-dasharray="2 2"/>`, px, py, labelX, labelY, color))

			dots = append(dots, fmt.Sprintf(`    <text x="%.1f" y="%.1f" text-anchor="%s" fill="%s" font-size="9" font-family="monospace" opacity="0.2" filter="url(#label-glow)">%s</text>`, labelX, labelY+3, anchor, color, Esc(item)))

			dots = append(dots, fmt.Sprintf(`    <text x="%.1f" y="%.1f" text-anchor="%s" fill="%s" font-size="9" font-family="monospace" opacity="0.85">%s</text>`, labelX, labelY+3, anchor, color, Esc(item)))
		}
	}
	return strings.Join(dots, "\n")
}

func buildHeaderProjectStars(projects []config.Project, galaxyArms []config.GalaxyArm, armColors []string, allArmPoints [][]Point) string {
	var stars []string
	count := 0
	for _, proj := range projects {
		if count >= 3 {
			break
		}
		armIdx := proj.Arm % len(galaxyArms)
		color := armColors[armIdx]
		points := allArmPoints[armIdx]

		ptIdx := len(points) - 3
		if ptIdx > 24 {
			ptIdx = 24
		}
		if ptIdx < 0 {
			ptIdx = 0
		}

		px, py := points[ptIdx].X, points[ptIdx].Y
		delay := fmt.Sprintf("%.1fs", float64(armIdx)*0.8)

		stars = append(stars, fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="4" fill="%s" filter="url(#star-glow-%d)">
      <animate attributeName="opacity" values="0.6;1;0.6" dur="4s" begin="%s" repeatCount="indefinite"/>
    </circle>`, px, py, color, armIdx, delay))
		count++
	}
	return strings.Join(stars, "\n")
}

func buildOrbitalRings(cx, cy float64, theme map[string]string) string {
	return fmt.Sprintf(`    <ellipse cx="%.1f" cy="%.1f" rx="55" ry="18" fill="none" stroke="%s" stroke-width="0.6" opacity="0.15" stroke-dasharray="4 6">
      <animateTransform attributeName="transform" type="rotate" from="0 %.1f %.1f" to="360 %.1f %.1f" dur="20s" repeatCount="indefinite"/>
    </ellipse>
    <ellipse cx="%.1f" cy="%.1f" rx="75" ry="24" fill="none" stroke="%s" stroke-width="0.5" opacity="0.1" stroke-dasharray="3 8">
      <animateTransform attributeName="transform" type="rotate" from="360 %.1f %.1f" to="0 %.1f %.1f" dur="30s" repeatCount="indefinite"/>
    </ellipse>`, cx, cy, theme["synapse_cyan"], cx, cy, cx, cy, cx, cy, theme["dendrite_violet"], cx, cy, cx, cy)
}

func buildGalaxyCore(cx, cy float64, theme map[string]string, initial string) string {
	return fmt.Sprintf(`    <!-- Outer haze -->
    <circle cx="%.1f" cy="%.1f" r="40" fill="url(#core-haze-gradient)" opacity="0.4"/>
    <!-- Inner glow -->
    <circle cx="%.1f" cy="%.1f" r="24" fill="url(#core-inner-gradient)" opacity="0.6"/>
    <!-- Outer ring -->
    <ellipse cx="%.1f" cy="%.1f" rx="20" ry="18" fill="none" stroke="%s" stroke-width="1.2" opacity="0.55" stroke-dasharray="5 3" class="core-ring"/>
    <!-- Inner ring -->
    <circle cx="%.1f" cy="%.1f" r="14" fill="none" stroke="%s" stroke-width="0.8" opacity="0.4" class="core-ring-inner"/>
    <!-- Solid core -->
    <circle cx="%.1f" cy="%.1f" r="11" fill="%s" stroke="%s" stroke-width="0.5"/>
    <!-- Bright center dot -->
    <circle cx="%.1f" cy="%.1f" r="3" fill="%s" filter="url(#core-bright-glow)" opacity="0.9"/>
    <!-- Initial -->
    <text x="%.1f" y="%.1f" text-anchor="middle" fill="%s" font-size="14" font-weight="bold" font-family="monospace">%s</text>`,
		cx, cy, cx, cy, cx, cy, theme["synapse_cyan"], cx, cy, theme["dendrite_violet"], cx, cy, theme["nebula"], theme["star_dust"], cx, cy, theme["synapse_cyan"], cx, cy+5, theme["synapse_cyan"], initial)
}
