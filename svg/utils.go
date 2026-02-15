package svg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// Constants for metrics
const (
	MetricCommits = "commits"
	MetricStars   = "stars"
	MetricPRs     = "prs"
	MetricIssues  = "issues"
	MetricRepos   = "repos"
)

var MetricLabels = map[string]string{
	MetricCommits: "Commits",
	MetricStars:   "Stars",
	MetricPRs:     "PRs",
	MetricIssues:  "Issues",
	MetricRepos:   "Repos",
}

var MetricColors = map[string]string{
	MetricCommits: "synapse_cyan",
	MetricStars:   "axon_amber",
	MetricPRs:     "dendrite_violet",
	MetricIssues:  "synapse_cyan",
	MetricRepos:   "dendrite_violet",
}

var MetricIcons = map[string]string{
	MetricCommits: `<path d="M8 1a7 7 0 1 0 0 14A7 7 0 0 0 8 1zm0 2.5a4.5 4.5 0 0 1 4.473 4H14.5a.5.5 0 0 1 0 1h-2.027A4.5 4.5 0 0 1 8 12.5a4.5 4.5 0 0 1-4.473-4H1.5a.5.5 0 0 1 0-1h2.027A4.5 4.5 0 0 1 8 3.5zm0 1.5a3 3 0 1 0 0 6 3 3 0 0 0 0-6z"/>`,
	MetricStars:   `<path d="M8 .25a.75.75 0 0 1 .673.418l1.882 3.815 4.21.612a.75.75 0 0 1 .416 1.279l-3.046 2.97.719 4.192a.75.75 0 0 1-1.088.791L8 12.347l-3.766 1.98a.75.75 0 0 1-1.088-.79l.72-4.194L.818 6.374a.75.75 0 0 1 .416-1.28l4.21-.611L7.327.668A.75.75 0 0 1 8 .25z"/>`,
	MetricPRs:     `<path d="M5 3.254V3.25v.005a.75.75 0 1 1 0-.005zm6.5 8a.75.75 0 1 1 0 1.5.75.75 0 0 1 0-1.5zM5 12.75a.75.75 0 1 1 0 1.5.75.75 0 0 1 0-1.5zm-1.5.75a1.5 1.5 0 1 0 1.5 1.5v-8.5a1.5 1.5 0 1 0-1.5-1.5v8.5a1.5 1.5 0 0 0 0 0zm8.5-2.5a1.5 1.5 0 0 0-1.5 1.5 1.5 1.5 0 1 0 3 0v-3.133l.025-.05A3.252 3.252 0 0 0 11 5.25V3.5h1.25a.75.75 0 0 0 .53-1.28l-2-2a.75.75 0 0 0-1.06 0l-2 2A.75.75 0 0 0 8.25 3.5H9.5v1.75a1.75 1.75 0 0 0 1.75 1.75h.244a1.75 1.75 0 0 1 1.006.319V11a1.5 1.5 0 0 0-1.5-1.5z"/>`,
	MetricIssues:  `<path d="M8 9.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z"/><path d="M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0zm0 1.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13z"/>`,
	MetricRepos:   `<path d="M2 2.5A2.5 2.5 0 0 1 4.5 0h8.75a.75.75 0 0 1 .75.75v12.5a.75.75 0 0 1-.75.75h-2.5a.75.75 0 0 1 0-1.5h1.75v-2h-8a1 1 0 0 0-.714 1.7.75.75 0 0 1-1.072 1.05A2.495 2.495 0 0 1 2 11.5zm10.5-1h-8a1 1 0 0 0-1 1v6.708A2.486 2.486 0 0 1 4.5 9h8zM5 12.25a.25.25 0 0 1 .25-.25h3.5a.25.25 0 0 1 .25.25v3.25a.25.25 0 0 1-.4.2l-1.45-1.087a.25.25 0 0 0-.3 0L5.4 15.7a.25.25 0 0 1-.4-.2z"/>`,
}

// GitHub Linguist colors for popular languages (incomplete list, but matches python basics)
var LanguageColors = map[string]string{
	"Python": "#3572A5", "JavaScript": "#f1e05a", "TypeScript": "#3178c6", "Java": "#b07219",
	"C#": "#178600", "C++": "#f34b7d", "C": "#555555", "Go": "#00ADD8", "Rust": "#dea584",
	"Ruby": "#701516", "PHP": "#4F5D95", "Swift": "#F05138", "Kotlin": "#A97BFF",
	"html": "#e34c26", "HTML": "#e34c26", "css": "#563d7c", "CSS": "#563d7c",
}

func GetLanguageColor(lang string) string {
	if c, ok := LanguageColors[lang]; ok {
		return c
	}
	// fallback matching logic or default
	return "#8b949e"
}

func Esc(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "&", "&amp;"), "<", "&lt;")
}

func FormatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func WrapText(text string, maxChars int) []string {
	words := strings.Fields(text)
	var lines []string
	var current string
	for _, word := range words {
		if len(current)+1+len(word) > maxChars {
			lines = append(lines, current)
			current = word
		} else {
			if current == "" {
				current = word
			} else {
				current += " " + word
			}
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func DeterministicRandom(seedStr string, count int, minVal, maxVal float64) []float64 {
	values := make([]float64, count)
	for i := 0; i < count; i++ {
		h := md5.Sum([]byte(fmt.Sprintf("%s_%d", seedStr, i)))
		hexStr := hex.EncodeToString(h[:])
		// Use first 8 chars (32 bits)
		valUint := int64(0)
		fmt.Sscanf(hexStr[:8], "%x", &valUint)
		normalized := float64(valUint) / 0xFFFFFFFF
		values[i] = minVal + normalized*(maxVal-minVal)
	}
	return values
}

type Point struct {
	X, Y float64
}

func SpiralPoints(cx, cy, startAngle float64, numPoints int, maxRadius, turns, xScal, yScal float64) []Point {
	points := make([]Point, numPoints)
	for i := 0; i < numPoints; i++ {
		t := float64(i) / math.Max(float64(numPoints-1), 1)
		angle := (startAngle * math.Pi / 180) + t*turns*2*math.Pi
		r := t * maxRadius
		x := cx + r*math.Cos(angle)*xScal
		y := cy + r*math.Sin(angle)*yScal
		points[i] = Point{X: x, Y: y}
	}
	return points
}

func SvgArcPath(cx, cy, r, startDeg, endDeg float64) string {
	startRad := (startDeg - 90) * math.Pi / 180
	endRad := (endDeg - 90) * math.Pi / 180
	x1 := cx + r*math.Cos(startRad)
	y1 := cy + r*math.Sin(startRad)
	x2 := cx + r*math.Cos(endRad)
	y2 := cy + r*math.Sin(endRad)
	largeArc := 0
	if endDeg-startDeg > 180 {
		largeArc = 1
	}
	return fmt.Sprintf("M %.1f %.1f L %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f Z", cx, cy, x1, y1, r, r, largeArc, x2, y2)
}

type LangStat struct {
	Name       string
	Bytes      int
	Percentage float64
	Color      string
}

func CalculateLanguagePercentages(languages map[string]int, exclude []string, maxDisplay int) []LangStat {
	excludedSet := make(map[string]bool)
	for _, e := range exclude {
		excludedSet[e] = true
	}

	filtered := make(map[string]int)
	total := 0
	for k, v := range languages {
		if !excludedSet[k] {
			filtered[k] = v
			total += v
		}
	}

	if total == 0 {
		return []LangStat{}
	}

	// Simple sort mostly by value descending
	// In Go, map iteration is random, so we need to convert to slice to sort
	type kv struct {
		Key string
		Val int
	}
	var ss []kv
	for k, v := range filtered {
		ss = append(ss, kv{k, v})
	}
	// Bubble sort or whatever simple sort
	for i := 0; i < len(ss); i++ {
		for j := i + 1; j < len(ss); j++ {
			if ss[j].Val > ss[i].Val {
				ss[i], ss[j] = ss[j], ss[i]
			}
		}
	}

	if len(ss) > maxDisplay {
		ss = ss[:maxDisplay]
	}

	var stats []LangStat
	for _, item := range ss {
		pct := float64(item.Val) / float64(total) * 100
		stats = append(stats, LangStat{
			Name:       item.Key,
			Bytes:      item.Val,
			Percentage: math.Round(pct*10) / 10,
			Color:      GetLanguageColor(item.Key),
		})
	}
	return stats
}

func Title(s string) string {
	return strings.Title(s)
}
