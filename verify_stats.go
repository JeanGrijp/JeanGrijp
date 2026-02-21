package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/JeanGrijp/JeanGrijp/github"
	"github.com/JeanGrijp/JeanGrijp/svg"
)

func init() {
	// 1. Load config (to get username/token logic if needed, though we mainly need env var)
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	username := "JeanGrijp" // Default or grab from config if we wanted to be fancy

	fmt.Printf("🔍 Verifying language stats for @%s...\n", username)
	fmt.Println("Connecting to GitHub API...")

	client := github.NewClient(username, token)
	languages, err := client.FetchLanguages(context.Background())
	if err != nil {
		fmt.Printf("Error fetching languages: %v\n", err)
		return
	}

	fmt.Println("\n📊 Raw Data (Bytes per Language):")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Language\tBytes\tPercentage")
	fmt.Fprintln(w, "--------\t-----\t----------")

	// Calculate total
	totalBytes := 0
	for _, bytes := range languages {
		totalBytes += bytes
	}

	// Sort for display
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range languages {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		pct := float64(kv.Value) / float64(totalBytes) * 100
		fmt.Fprintf(w, "%s\t%d\t%.2f%%\n", kv.Key, kv.Value, pct)
	}
	w.Flush()

	fmt.Printf("\nTotal Bytes: %d\n", totalBytes)

	// Compare with what the SVG generator would see (using the utils logic)
	fmt.Println("\n🎨 Filtered Stats (as seen in SVG):")
	// Using default exclude list and max display from common config
	stats := svg.CalculateLanguagePercentages(languages, []string{}, 8)

	w2 := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w2, "Language\tPercentage")
	fmt.Fprintln(w2, "--------\t----------")
	for _, s := range stats {
		fmt.Fprintf(w2, "%s\t%.1f%%\n", s.Name, s.Percentage)
	}
	w2.Flush()
}
