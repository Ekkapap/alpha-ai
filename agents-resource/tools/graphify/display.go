package main

// display.go — display helpers, reused by alpha and build commands.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// findProjectRoot walks up from the binary looking for project root markers.
func findProjectRoot(binaryPath string) string {
	dir := filepath.Dir(binaryPath)
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		if _, err := os.Stat(filepath.Join(dir, "knowledge-graph", "graphify-out")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "α")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		dir = parent
	}
	return dir
}

type alphaConfig struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MCP         []string `json:"mcp"`
	Tools       map[string]struct {
		Source string `json:"source"`
		Bin    string `json:"bin"`
	} `json:"tools"`
}

func loadAlphaConfig(r string) alphaConfig {
	cfg := alphaConfig{Version: "0.0.0", Name: "Alpha System", Description: "AI-Native Intelligence Toolchain"}
	data, err := os.ReadFile(filepath.Join(r, "α/alpha.json"))
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	return cfg
}

type readyCheck struct {
	label  string
	ok     bool
	detail string
}

func alphaReadyChecks(r string, cfg alphaConfig) []readyCheck {
	platform := runtime.GOOS
	var checks []readyCheck

	binPath := filepath.Join(r, "α/bin/"+platform+"/graphify")
	if _, err := os.Stat(binPath); err == nil {
		checks = append(checks, readyCheck{"binary", true, "α/bin/" + platform + "/graphify"})
	} else {
		checks = append(checks, readyCheck{"binary", false, "not found — run setup-hooks.sh"})
	}

	if _, err := os.Stat(filepath.Join(graphifyDataDir(r), "graph.json")); err == nil {
		checks = append(checks, readyCheck{"graph", true, "graphify-out/graph.json"})
	} else {
		checks = append(checks, readyCheck{"graph", false, "not built — run /graphify"})
	}

	memPath := filepath.Join(r, "knowledge-graph/memories/session-summary.md")
	if stat, err := os.Stat(memPath); err == nil {
		checks = append(checks, readyCheck{"memory", true, "last sync " + stat.ModTime().Format("2006-01-02 15:04")})
	} else {
		checks = append(checks, readyCheck{"memory", false, "no session yet — run /sync"})
	}

	if _, err := os.Stat(filepath.Join(r, ".mcp.json")); err == nil {
		checks = append(checks, readyCheck{"mcp", true, ".mcp.json"})
	} else {
		checks = append(checks, readyCheck{"mcp", false, ".mcp.json not found"})
	}

	if _, err := os.Stat(filepath.Join(r, "α/hooks/bin/awake")); err == nil {
		checks = append(checks, readyCheck{"hooks", true, "α/hooks/bin/"})
	} else {
		checks = append(checks, readyCheck{"hooks", false, "not installed — run setup-hooks.sh"})
	}

	return checks
}

func alphaDisplay(r string) string {
	const (
		yellow = "\033[33m"
		green  = "\033[32m"
		red    = "\033[31m"
		cyan   = "\033[96m"
		white  = "\033[97m"
		dim    = "\033[2m"
		bold   = "\033[1m"
		reset  = "\033[0m"
	)

	cfg := loadAlphaConfig(r)

	nodes, edges, comms := 0, 0, 0
	if data, err := os.ReadFile(filepath.Join(graphifyDataDir(r), "graph.json")); err == nil {
		var g struct {
			Nodes []json.RawMessage `json:"nodes"`
			Links []json.RawMessage `json:"links"`
		}
		if json.Unmarshal(data, &g) == nil {
			nodes, edges = len(g.Nodes), len(g.Links)
		}
	}
	if data, err := os.ReadFile(filepath.Join(graphifyDataDir(r), ".graphify_analysis.json")); err == nil {
		var a struct {
			Communities map[string]json.RawMessage `json:"communities"`
		}
		if json.Unmarshal(data, &a) == nil {
			comms = len(a.Communities)
		}
	}
	graphLine := dim + "not initialized" + reset
	if nodes > 0 {
		graphLine = fmt.Sprintf("%d nodes · %d edges · %d communities", nodes, edges, comms)
	}

	logo := []string{
		yellow + "  ████████████  " + reset,
		yellow + " ██" + reset + "          " + yellow + "██ " + reset,
		yellow + "██" + reset + "            " + yellow + "██" + reset,
		yellow + "██" + reset + "            " + yellow + "██" + reset,
		yellow + " ██" + reset + "          " + yellow + "██ " + reset,
		yellow + "  ██" + reset + "        " + yellow + "██  " + reset,
		yellow + "  ██" + reset + "        " + yellow + "██  " + reset,
		yellow + " ████" + reset + "      " + yellow + "████ " + reset,
	}

	info := []string{
		bold + white + "α  " + cfg.Name + reset + "  " + dim + "v" + cfg.Version + reset,
		dim + "────────────────────────────────────" + reset,
		cfg.Description,
		cyan + "Human in the Loop" + reset + " · Always",
		"",
		dim + "The model thinks. You decide." + reset,
		"",
		dim + "────────────────────────────────────" + reset,
	}

	var sb strings.Builder
	sb.WriteString("\n")
	maxL := len(logo)
	if len(info) > maxL {
		maxL = len(info)
	}
	for i := 0; i < maxL; i++ {
		l := strings.Repeat(" ", 16)
		if i < len(logo) {
			l = logo[i]
		}
		r := ""
		if i < len(info) {
			r = info[i]
		}
		sb.WriteString(fmt.Sprintf("  %s   %s\n", l, r))
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s\n", "Root", r))
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s\n", "Graph", graphLine))
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s · graphify\n", "Runtime", runtime.GOOS+"/"+runtime.GOARCH))

	sb.WriteString("\n")
	checks := alphaReadyChecks(r, cfg)
	allOk := true
	for _, c := range checks {
		icon := green + "✓" + reset
		if !c.ok {
			icon = red + "✗" + reset
			allOk = false
		}
		sb.WriteString(fmt.Sprintf("  %s  "+dim+"%-8s"+reset+"  %s\n", icon, c.label, c.detail))
	}

	sb.WriteString("\n")
	if allOk {
		sb.WriteString("  " + green + bold + "● System Ready" + reset + "\n")
	} else {
		sb.WriteString("  " + red + bold + "● Needs Setup" + reset + "  — resolve ✗ items above\n")
	}

	sb.WriteString("\n")
	sb.WriteString("  " + dim + "❝ Intelligentia humana semper in centro ❞" + reset + "\n\n")

	return sb.String()
}

func isInputTTY() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func printHelp() {
	fmt.Println(`Graphify Helper — CLI Tool

Usage:
  graphify <command> [args]

Commands:
  awake                    Show project memory and latest graph status
  sync [--summary <text>]  Sync graph database, save memory, and open visualization
  forget [pattern] [-y]   Delete past session summaries and matching graph memories
  overview                 Show compact Phase 0 architecture stats
  sketch --query <text>    Run Phase 1 BFS query to assess relevance
  detail --ids <id1,id2>   Run Phase 2 detailed caller/callee analysis
  focus <path> <term>      Search and view file lines around match
  alpha                    Display Alpha System identity — logo, stats, commands
  build                    Compile and copy binary to α/bin/<platform>/`)
}
