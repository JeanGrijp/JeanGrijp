package svg

import (
	"fmt"
	"math"
	"strings"

	"github.com/JeanGrijp/JeanGrijp/config"
)

const (
	TechWidth       = 850
	TechRadarRadius = 65.0
)

func RenderTechStack(languages map[string]int, galaxyArms []config.GalaxyArm, theme map[string]string, exclude []string, maxDisplay int) string {
	width := TechWidth
	langData := CalculateLanguagePercentages(languages, exclude, maxDisplay)

	// Left side: Language bars
	leftX := 30.0
	startY := 65.0

	barsStr := buildLanguageBars(langData, theme, leftX, startY)

	// Right side: Focus Sectors radar
	armColors := config.ResolveArmColors(galaxyArms, theme)

	var sectors []sectorData
	for i, arm := range galaxyArms {
		color := armColors[i]
		detected := 0
		for _, item := range arm.Items {
			if _, ok := languages[item]; ok {
				detected++
			}
		}
		sectors = append(sectors, sectorData{
			Name:     arm.Name,
			Color:    color,
			Items:    len(arm.Items),
			Detected: detected,
			StartDeg: float64(i)*120 + 1,
			EndDeg:   float64(i+1)*120 - 1,
		})
	}

	// Radar geometry
	rcx := 637.0 // center of right half
	badgeStartY := 65.0
	rcy := badgeStartY + TechRadarRadius + 10.0
	gridRings := []float64{22, 44, 65}

	// Dynamic height
	langHeight := startY + float64(len(langData))*22 + 20
	radarHeight := rcy + TechRadarRadius + 35
	height := int(math.Max(200, math.Max(langHeight, radarHeight)))

	radarParts := []string{}
	radarParts = append(radarParts, buildRadarGrid(rcx, rcy, gridRings, theme)...)
	radarParts = append(radarParts, buildRadarSectors(sectors, rcx, rcy, TechRadarRadius, theme)...)
	radarParts = append(radarParts, buildRadarNeedle(rcx, rcy, TechRadarRadius, theme)...)
	radarParts = append(radarParts, buildRadarLabelsAndDots(sectors, galaxyArms, rcx, rcy, TechRadarRadius, theme)...)

	radarStr := strings.Join(radarParts, "\n")

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <defs/>

  <!-- Card background -->
  <rect x="0.5" y="0.5" width="%d" height="%d" rx="12" ry="12"
        fill="%s" stroke="%s" stroke-width="1"/>

  <!-- Left: Language Telemetry -->
  <text x="30" y="38" fill="%s" font-size="11" font-family="monospace" letter-spacing="3">LANGUAGE TELEMETRY</text>

  <!-- Vertical divider -->
  <line x1="425" y1="25" x2="425" y2="%d" stroke="%s" stroke-width="1" opacity="0.4"/>

  <!-- Right: Focus Sectors -->
  <text x="460" y="38" fill="%s" font-size="11" font-family="monospace" letter-spacing="3">FOCUS SECTORS</text>

%s

%s
</svg>`, width, height, width, height, width-1, height-1, theme["nebula"], theme["star_dust"], theme["text_faint"], height-25, theme["star_dust"], theme["text_faint"], barsStr, radarStr)
}

func buildLanguageBars(langData []LangStat, theme map[string]string, leftX, startY float64) string {
	var lines []string
	barMaxWidth := 200.0

	for i, lang := range langData {
		y := startY + float64(i)*22
		barW := math.Max(4, (lang.Percentage/100)*barMaxWidth)
		delay := fmt.Sprintf("%.1fs", float64(i)*0.1)

		lines = append(lines, fmt.Sprintf(`    <g transform="translate(%.1f, %.1f)">
      <text x="0" y="0" fill="%s" font-size="11" font-family="sans-serif" dominant-baseline="middle">%s</text>
      <rect x="110" y="-6" width="%.1f" height="12" rx="3" fill="%s" opacity="0.85">
        <animate attributeName="width" from="0" to="%.1f" dur="0.8s" begin="%s" fill="freeze"/>
      </rect>
      <text x="320" y="0" fill="%s" font-size="10" font-family="monospace" dominant-baseline="middle">%.1f%%</text>
    </g>`, leftX, y, theme["text_dim"], Esc(lang.Name), barW, lang.Color, barW, delay, theme["text_faint"], lang.Percentage))
	}
	return strings.Join(lines, "\n")
}

func buildRadarGrid(rcx, rcy float64, rings []float64, theme map[string]string) []string {
	var parts []string
	for _, r := range rings {
		parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="%s" stroke-width="0.5" stroke-dasharray="3,3" opacity="0.25"/>`, rcx, rcy, r, theme["text_faint"]))
	}
	return parts
}

type sectorData struct {
	Name     string
	Color    string
	Items    int
	Detected int
	StartDeg float64
	EndDeg   float64
}

func buildRadarSectors(sectors []sectorData, rcx, rcy, radius float64, theme map[string]string) []string {
	var parts []string
	for _, sec := range sectors {
		d := SvgArcPath(rcx, rcy, radius, sec.StartDeg, sec.EndDeg)
		parts = append(parts, fmt.Sprintf(`    <path d="%s" fill="%s" fill-opacity="0.10" stroke="%s" stroke-opacity="0.3" stroke-width="0.5"/>`, d, sec.Color, sec.Color))
	}

	for i := 0; i < len(sectors); i++ {
		angleRad := (float64(i)*120 - 90) * math.Pi / 180
		lx := rcx + radius*math.Cos(angleRad)
		ly := rcy + radius*math.Sin(angleRad)
		parts = append(parts, fmt.Sprintf(`    <line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="0.5" opacity="0.3"/>`, rcx, rcy, lx, ly, theme["text_faint"]))
	}
	return parts
}

func buildRadarNeedle(rcx, rcy, radius float64, theme map[string]string) []string {
	scanColor := theme["synapse_cyan"] // default/fallback
	if scanColor == "" {
		scanColor = "#00d4ff"
	}
	tipX := rcx
	tipY := rcy - radius
	sweepD := SvgArcPath(rcx, rcy, radius, 330, 360)
	outerHaw := 2.5
	innerHw := 0.8

	needle := fmt.Sprintf(`    <g>
      <!-- Sweep trail -->
      <path d="%s" fill="%s" fill-opacity="0.07"/>
      <!-- Outer wedge -->
      <polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" opacity="0.25"/>
      <!-- Inner bright core -->
      <polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" opacity="0.5"/>
      <!-- Tip glow -->
      <circle cx="%.1f" cy="%.1f" r="2" fill="%s" opacity="0.6">
        <animate attributeName="opacity" values="0.4;0.8;0.4" dur="2s" repeatCount="indefinite"/>
      </circle>
      <animateTransform attributeName="transform" type="rotate" from="0 %.1f %.1f" to="360 %.1f %.1f" dur="8s" repeatCount="indefinite"/>
    </g>`, sweepD, scanColor,
		rcx-outerHaw, rcy, tipX, tipY, rcx+outerHaw, rcy, scanColor,
		rcx-innerHw, rcy, tipX, tipY, rcx+innerHw, rcy, scanColor,
		tipX, tipY, scanColor,
		rcx, rcy, rcx, rcy)

	return []string{needle}
}

func buildRadarLabelsAndDots(sectors []sectorData, galaxyArms []config.GalaxyArm, rcx, rcy, radius float64, theme map[string]string) []string {
	var parts []string
	for _, sec := range sectors {
		midDeg := (sec.StartDeg + sec.EndDeg) / 2
		midRad := (midDeg - 90) * math.Pi / 180
		labelR := radius + 18
		lx := rcx + labelR*math.Cos(midRad)
		ly := rcy + labelR*math.Sin(midRad)

		anchor := "middle"
		if math.Abs(lx-rcx) >= 5 {
			if lx > rcx {
				anchor = "start"
			} else {
				anchor = "end"
			}
		}

		parts = append(parts, fmt.Sprintf(`    <text x="%.1f" y="%.1f" fill="%s" font-size="9" font-family="monospace" text-anchor="%s" dominant-baseline="middle">%s</text>`, lx, ly, sec.Color, anchor, Esc(sec.Name)))

		countY := ly + 12
		parts = append(parts, fmt.Sprintf(`    <text x="%.1f" y="%.1f" fill="%s" font-size="8" font-family="monospace" text-anchor="%s" dominant-baseline="middle">(%d)</text>`, lx, countY, theme["text_faint"], anchor, sec.Items))
	}

	radiiCycle := []float64{24, 40, 56}
	for secI, sec := range sectors {
		arm := galaxyArms[secI]
		items := arm.Items
		itemCount := len(items)
		edgePad := 10.0

		for j := 0; j < itemCount; j++ {
			usableStart := sec.StartDeg + edgePad
			usableEnd := sec.EndDeg - edgePad
			itemAngle := 0.0
			if itemCount == 1 {
				itemAngle = (usableStart + usableEnd) / 2
			} else {
				itemAngle = usableStart + (usableEnd-usableStart)*float64(j)/float64(itemCount-1)
			}
			itemRad := (itemAngle - 90) * math.Pi / 180
			dotR := radiiCycle[j%3]
			dx := rcx + dotR*math.Cos(itemRad)
			dy := rcy + dotR*math.Sin(itemRad)

			pulseBegin := (itemAngle/360)*8 - 0.3
			if pulseBegin < 0 {
				pulseBegin += 8
			}

			parts = append(parts, fmt.Sprintf(`    <circle cx="%.1f" cy="%.1f" r="3" fill="%s" opacity="0.35">
      <animate attributeName="opacity" values="0.35;0.35;1.0;0.35;0.35" keyTimes="0;0.04;0.06;0.10;1" dur="8s" begin="%.2fs" repeatCount="indefinite"/>
    </circle>`, dx, dy, sec.Color, pulseBegin))
		}
	}
	return parts
}
