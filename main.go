package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JeanGrijp/JeanGrijp/config"
	"github.com/JeanGrijp/JeanGrijp/github"
	"github.com/JeanGrijp/JeanGrijp/svg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 1. Load config
	configPath := "config.yml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}

	cfg, err := config.ValidateAndApplyDefaults(data)
	if err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	log.Printf("Generating profile SVGs for @%s...", cfg.Username)

	// 2. Fetch GitHub data
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	client := github.NewClient(cfg.Username, token)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Println("Fetching stats...")
	stats, err := client.FetchStats(ctx)
	if err != nil {
		log.Printf("Warning: Could not fetch stats (%v). Using defaults.", err)
		stats = &github.Stats{}
	}

	log.Println("Fetching languages...")
	languages, err := client.FetchLanguages(ctx)
	if err != nil {
		log.Printf("Warning: Could not fetch languages (%v). Using defaults.", err)
		languages = make(map[string]int)
	}

	log.Printf("Stats: %+v", stats)
	log.Printf("Languages: %d found", len(languages))

	// 3. Build SVGs
	outputDir := filepath.Join("assets", "generated")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	svgs := map[string]string{
		"galaxy-header.svg":          svg.RenderGalaxyHeader(cfg, cfg.Theme),
		"stats-card.svg":             svg.RenderStatsCard(stats, cfg.Stats.Metrics, cfg.Theme),
		"tech-stack.svg":             svg.RenderTechStack(languages, cfg.GalaxyArms, cfg.Theme, cfg.Languages.Exclude, cfg.Languages.MaxDisplay),
		"projects-constellation.svg": svg.RenderProjectsConstellation(cfg.Projects, cfg.GalaxyArms, cfg.Theme),
	}

	for filename, content := range svgs {
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Printf("Failed to write %s: %v", path, err)
		} else {
			log.Printf("Wrote %s", path)
		}
	}

	log.Println("Done! 4 SVGs generated.")
}
