package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// getGraphStats reads graph.json and returns a summary string
func getGraphStats(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, ".graphify-out/graph.json"))
	if err != nil {
		return "", err
	}
	var g struct {
		Nodes []any `json:"nodes"`
		Links []any `json:"links"`
	}
	if err := json.Unmarshal(data, &g); err != nil {
		return "", err
	}

	analysisData, _ := os.ReadFile(filepath.Join(root, ".graphify-out/.graphify_analysis.json"))
	var a struct {
		Communities map[string]any `json:"communities"`
	}
	json.Unmarshal(analysisData, &a)

	return fmt.Sprintf("Nodes: %d, Edges: %d, Communities: %d", len(g.Nodes), len(g.Links), len(a.Communities)), nil
}

// getGodNodes reads .graphify_analysis.json and returns top nodes
func getGodNodes(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, ".graphify-out/.graphify_analysis.json"))
	if err != nil {
		return "", err
	}
	var a struct {
		Gods []struct {
			Label  string `json:"label"`
			Degree int    `json:"degree"`
		} `json:"gods"`
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return "", err
	}

	var res strings.Builder
	for i, g := range a.Gods {
		if i >= 5 {
			break
		}
		res.WriteString(fmt.Sprintf("%d. %s (%d edges)\n", i+1, g.Label, g.Degree))
	}
	return res.String(), nil
}

// ── Two-phase query types ────────────────────────────────────────────────────

type gNode struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	SourceFile     string `json:"source_file"`
	SourceLocation string `json:"source_location"`
	Community      int    `json:"community"`
	FileType       string `json:"file_type"`
}

type gLink struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Relation string `json:"relation"`
}

type fullGraph struct {
	Nodes []gNode `json:"nodes"`
	Links []gLink `json:"links"`
}

func loadFullGraph(root string) (*fullGraph, error) {
	data, err := os.ReadFile(filepath.Join(root, ".graphify-out/graph.json"))
	if err != nil {
		return nil, err
	}
	var g fullGraph
	return &g, json.Unmarshal(data, &g)
}

func scoreNode(n gNode, terms []string) float64 {
	label := strings.ToLower(n.Label)
	file := strings.ToLower(n.SourceFile)
	var s float64
	for _, t := range terms {
		t = strings.ToLower(t)
		switch {
		case strings.EqualFold(n.Label, t):
			s += 3
		case strings.Contains(label, t):
			s += 2
		case strings.Contains(file, t):
			s += 1
		}
	}
	return s
}

// sketchGraph runs BFS from top-scoring seed nodes and returns compact JSON.
func sketchGraph(g *fullGraph, query string, depth int) string {
	nodeIdx := make(map[string]gNode, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeIdx[n.ID] = n
	}
	adj := make(map[string][]gLink, len(g.Nodes))
	for _, l := range g.Links {
		adj[l.Source] = append(adj[l.Source], l)
		rev := gLink{Source: l.Target, Target: l.Source, Relation: l.Relation}
		adj[l.Target] = append(adj[l.Target], rev)
	}

	terms := strings.Fields(query)
	type scored struct {
		n gNode
		s float64
	}
	var seeds []scored
	for _, n := range g.Nodes {
		if sc := scoreNode(n, terms); sc >= 1 {
			seeds = append(seeds, scored{n, sc})
		}
	}
	sort.Slice(seeds, func(i, j int) bool { return seeds[i].s > seeds[j].s })
	if len(seeds) > 5 {
		seeds = seeds[:5]
	}

	type qItem struct {
		id  string
		dep int
		via string
	}
	visited := make(map[string]bool)
	queue := make([]qItem, 0, len(seeds))
	seedLabels := make([]string, 0, len(seeds))
	for _, s := range seeds {
		visited[s.n.ID] = true
		queue = append(queue, qItem{s.n.ID, 0, ""})
		seedLabels = append(seedLabels, s.n.Label)
	}

	type outNode struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		File      string `json:"file"`
		Community int    `json:"community"`
		Depth     int    `json:"depth"`
		Via       string `json:"via,omitempty"`
	}
	var nodes []outNode
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		n, ok := nodeIdx[cur.id]
		if !ok {
			continue
		}
		loc := n.SourceFile
		if n.SourceLocation != "" {
			loc += ":" + n.SourceLocation
		}
		nodes = append(nodes, outNode{n.ID, n.Label, loc, n.Community, cur.dep, cur.via})
		if cur.dep < depth {
			for _, l := range adj[cur.id] {
				if !visited[l.Target] {
					visited[l.Target] = true
					queue = append(queue, qItem{l.Target, cur.dep + 1, l.Relation})
				}
			}
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Depth != nodes[j].Depth {
			return nodes[i].Depth < nodes[j].Depth
		}
		return nodes[i].Label < nodes[j].Label
	})

	out := map[string]any{"query": query, "seeds": seedLabels, "depth": depth, "total": len(nodes), "nodes": nodes}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b)
}

// detailNodes returns callers, callees, and file info for given node IDs.
func detailNodes(g *fullGraph, ids []string) string {
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[strings.TrimSpace(id)] = true
	}
	nodeIdx := make(map[string]gNode, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeIdx[n.ID] = n
	}

	callers := make(map[string][]string)
	callees := make(map[string][]string)
	for _, l := range g.Links {
		if l.Relation != "calls" {
			continue
		}
		callers[l.Target] = append(callers[l.Target], l.Source)
		callees[l.Source] = append(callees[l.Source], l.Target)
	}

	resolve := func(nodeIDs []string) []string {
		seen := make(map[string]bool)
		var out []string
		for _, id := range nodeIDs {
			if seen[id] {
				continue
			}
			seen[id] = true
			if n, ok := nodeIdx[id]; ok {
				out = append(out, n.Label+" ("+n.SourceFile+")")
			}
		}
		return out
	}

	type outNode struct {
		ID        string   `json:"id"`
		Label     string   `json:"label"`
		File      string   `json:"file"`
		FileType  string   `json:"file_type"`
		Community int      `json:"community"`
		Callers   []string `json:"callers"`
		Callees   []string `json:"callees"`
	}
	var nodes []outNode
	for _, n := range g.Nodes {
		if !idSet[n.ID] {
			continue
		}
		loc := n.SourceFile
		if n.SourceLocation != "" {
			loc += ":" + n.SourceLocation
		}
		nodes = append(nodes, outNode{
			ID: n.ID, Label: n.Label, File: loc,
			FileType: n.FileType, Community: n.Community,
			Callers: resolve(callers[n.ID]),
			Callees: resolve(callees[n.ID]),
		})
	}
	b, _ := json.MarshalIndent(map[string]any{"nodes": nodes}, "", "  ")
	return string(b)
}

// ── Phase 0: overview ────────────────────────────────────────────────────────

type godNodeInfo struct {
	Label     string `json:"label"`
	Edges     int    `json:"edges"`
	Community string `json:"community"`
}

type overviewResult struct {
	Nodes          int           `json:"nodes"`
	Edges          int           `json:"edges"`
	Communities    int           `json:"communities"`
	GodNodes       []godNodeInfo `json:"god_nodes"`
	TopCommunities []string      `json:"top_communities"`
}

func graphOverview(root string) (*overviewResult, error) {
	// node + edge counts from graph.json
	graphData, err := os.ReadFile(filepath.Join(root, ".graphify-out/graph.json"))
	if err != nil {
		return nil, err
	}
	var g struct {
		Nodes []json.RawMessage `json:"nodes"`
		Links []json.RawMessage `json:"links"`
	}
	if err := json.Unmarshal(graphData, &g); err != nil {
		return nil, err
	}

	// analysis for gods + community membership
	analysisData, err := os.ReadFile(filepath.Join(root, ".graphify-out/.graphify_analysis.json"))
	if err != nil {
		return nil, err
	}
	var a struct {
		Communities map[string][]string `json:"communities"`
		Gods        []struct {
			ID     string `json:"id"`
			Label  string `json:"label"`
			Degree int    `json:"degree"`
		} `json:"gods"`
	}
	if err := json.Unmarshal(analysisData, &a); err != nil {
		return nil, err
	}

	// build reverse map: node_id → community_id
	nodeToCommunity := make(map[string]string, len(g.Nodes))
	for commID, members := range a.Communities {
		for _, nodeID := range members {
			nodeToCommunity[nodeID] = commID
		}
	}

	// community labels (may be generic "Community N")
	labelsData, _ := os.ReadFile(filepath.Join(root, ".graphify-out/.graphify_labels.json"))
	var labels map[string]string
	json.Unmarshal(labelsData, &labels)
	communityName := func(id string) string {
		if labels != nil {
			if name, ok := labels[id]; ok {
				return name
			}
		}
		return "Community " + id
	}

	// top 10 god nodes with community
	gods := make([]godNodeInfo, 0, len(a.Gods))
	for _, gn := range a.Gods {
		commID := nodeToCommunity[gn.ID]
		gods = append(gods, godNodeInfo{
			Label:     gn.Label,
			Edges:     gn.Degree,
			Community: communityName(commID),
		})
	}

	// top communities: pick communities that contain at least one god node, deduped
	seen := make(map[string]bool)
	var topComms []string
	for _, gn := range a.Gods {
		commID := nodeToCommunity[gn.ID]
		name := communityName(commID)
		if !seen[name] {
			seen[name] = true
			topComms = append(topComms, name)
		}
	}

	return &overviewResult{
		Nodes:          len(g.Nodes),
		Edges:          len(g.Links),
		Communities:    len(a.Communities),
		GodNodes:       gods,
		TopCommunities: topComms,
	}, nil
}

// generateProjectSummary writes .graphify-out/project_summary.json
func generateProjectSummary(root string) error {
	ov, err := graphOverview(root)
	if err != nil {
		return err
	}
	out := map[string]any{
		"generated_at":    time.Now().UTC().Format(time.RFC3339),
		"stats":           map[string]int{"nodes": ov.Nodes, "edges": ov.Edges},
		"god_nodes":       ov.GodNodes,
		"communities":     ov.Communities,
		"top_communities": ov.TopCommunities,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, ".graphify-out/project_summary.json"), b, 0644)
}

// extractNodesFromSummary tries to find paths or node labels in the summary
func extractNodesFromSummary(summary string) []string {
	var nodes []string
	lines := strings.Split(summary, "\n")
	for _, line := range lines {
		// Look for bullet points with paths or names
		if strings.HasPrefix(strings.TrimSpace(line), "- ") || strings.HasPrefix(strings.TrimSpace(line), "* ") {
			content := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line), "-*"))
			// Heuristic: if it looks like a path or a single word identifier
			if strings.Contains(content, "/") || (len(strings.Fields(content)) == 1 && len(content) > 3) {
				nodes = append(nodes, content)
			}
		}
	}
	return nodes
}

// findProjectRoot walks up from the binary location looking for project root markers.
// Works regardless of where the binary lives — searches upward for α/, .graphify-out/, or .git.
func findProjectRoot(binaryPath string) string {
	dir := filepath.Dir(binaryPath)
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		if _, err := os.Stat(filepath.Join(dir, ".graphify-out")); err == nil {
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
	Version     string   `json:"version"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MCP         []string `json:"mcp"`
	Tools       map[string]struct {
		Source string `json:"source"`
		Bin    string `json:"bin"`
	} `json:"tools"`
}

func loadAlphaConfig(root string) alphaConfig {
	cfg := alphaConfig{Version: "0.0.0", Name: "Alpha System", Description: "AI-Native Intelligence Toolchain"}
	data, err := os.ReadFile(filepath.Join(root, "α/alpha.json"))
	if err == nil {
		json.Unmarshal(data, &cfg)
	}
	return cfg
}

// check returns (label, ok, detail)
type readyCheck struct {
	label  string
	ok     bool
	detail string
}

func alphaReadyChecks(root string, cfg alphaConfig) []readyCheck {
	platform := runtime.GOOS
	checks := []readyCheck{}

	// Binary
	binPath := filepath.Join(root, "α/bin/"+platform+"/graphify")
	if _, err := os.Stat(binPath); err == nil {
		checks = append(checks, readyCheck{"binary", true, "α/bin/" + platform + "/graphify"})
	} else {
		checks = append(checks, readyCheck{"binary", false, "not found — run setup-hooks.sh"})
	}

	// Graph
	if _, err := os.Stat(filepath.Join(root, ".graphify-out/graph.json")); err == nil {
		checks = append(checks, readyCheck{"graph", true, ".graphify-out/graph.json"})
	} else {
		checks = append(checks, readyCheck{"graph", false, "not built — run /graphify"})
	}

	// Memories
	memPath := filepath.Join(root, "α/memories/latest_state.md")
	if stat, err := os.Stat(memPath); err == nil {
		checks = append(checks, readyCheck{"memory", true, "last sync " + stat.ModTime().Format("2006-01-02 15:04")})
	} else {
		checks = append(checks, readyCheck{"memory", false, "no session yet — run /sync"})
	}

	// MCP config
	if _, err := os.Stat(filepath.Join(root, ".mcp.json")); err == nil {
		checks = append(checks, readyCheck{"mcp", true, ".mcp.json"})
	} else {
		checks = append(checks, readyCheck{"mcp", false, ".mcp.json not found"})
	}

	// Hooks
	if _, err := os.Stat(filepath.Join(root, "α/hooks/bin/awake")); err == nil {
		checks = append(checks, readyCheck{"hooks", true, "α/hooks/bin/"})
	} else {
		checks = append(checks, readyCheck{"hooks", false, "not installed — run setup-hooks.sh"})
	}

	return checks
}

func alphaDisplay(root string) string {
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

	cfg := loadAlphaConfig(root)

	// ── Stats ────────────────────────────────────────────────────────────────
	nodes, edges, comms := 0, 0, 0
	if data, err := os.ReadFile(filepath.Join(root, ".graphify-out/graph.json")); err == nil {
		var g struct {
			Nodes []json.RawMessage `json:"nodes"`
			Links []json.RawMessage `json:"links"`
		}
		if json.Unmarshal(data, &g) == nil {
			nodes, edges = len(g.Nodes), len(g.Links)
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, ".graphify-out/.graphify_analysis.json")); err == nil {
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

	// ── α ASCII Logo (16 visual cols wide) ───────────────────────────────────
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

	// ── Info lines alongside logo ─────────────────────────────────────────────
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

	// ── System info ───────────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s\n", "Root", root))
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s\n", "Graph", graphLine))
	sb.WriteString(fmt.Sprintf("  "+dim+"%-9s"+reset+"  %s · graphify\n", "Runtime", runtime.GOOS+"/"+runtime.GOARCH))

	// ── Readiness checks ─────────────────────────────────────────────────────
	sb.WriteString("\n")
	checks := alphaReadyChecks(root, cfg)
	allOk := true
	for _, c := range checks {
		icon := green + "✓" + reset
		if !c.ok {
			icon = red + "✗" + reset
			allOk = false
		}
		sb.WriteString(fmt.Sprintf("  %s  "+dim+"%-8s"+reset+"  %s\n", icon, c.label, c.detail))
	}

	// ── Overall status ────────────────────────────────────────────────────────
	sb.WriteString("\n")
	if allOk {
		sb.WriteString("  " + green + bold + "● System Ready" + reset + "\n")
	} else {
		sb.WriteString("  " + red + bold + "● Needs Setup" + reset + "  — resolve ✗ items above\n")
	}

	// ── Motto ─────────────────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString("  " + dim + "❝ Intelligentia humana semper in centro ❞" + reset + "\n\n")

	return sb.String()
}

func main() {
	// Initialize structured logging to stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// Get current working directory (Full Path)
	wd, _ := os.Getwd()
	absWd, _ := filepath.Abs(wd)

	// Root detection: env override → walk up from binary → CWD fallback
	var root string
	if envRoot := os.Getenv("PROJECT_ROOT"); envRoot != "" {
		root, _ = filepath.Abs(envRoot)
	} else if exe, err := os.Executable(); err == nil {
		exePath, _ := filepath.EvalSymlinks(exe)
		absExe, _ := filepath.Abs(exePath)
		root = findProjectRoot(absExe)
	}

	if root == "" || root == "/" {
		wd, _ := os.Getwd()
		root, _ = filepath.Abs(wd)
	}

	projectName := filepath.Base(root)

	// Create a new MCP server dynamically named after the project
	s := server.NewMCPServer(
		fmt.Sprintf("Graphify Helper (%s)", projectName),
		"1.1.1",
	)

	// CLI MODE: Support manual execution
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "alpha" {
			fmt.Print(alphaDisplay(root))
			os.Exit(0)
		}
		if arg == "awake" {
			res, _ := getGraphStats(root)
			gods, _ := getGodNodes(root)
			fmt.Printf("### GRAPH STATS\n%s\n\n### GOD NODES (Architecture Pillars)\n%s\n", res, gods)

			// Read latest state
			latestStatePath := filepath.Join(root, "α/memories/latest_state.md")
			if state, err := os.ReadFile(latestStatePath); err == nil {
				fmt.Printf("### LATEST STATE\n%s\n\n", string(state))
			}

			// 3. Graph Overview (Phase 0 JSON - Token-Efficient)
			if ov, err := graphOverview(root); err == nil {
				b, _ := json.MarshalIndent(ov, "", "  ")
				fmt.Printf("### GRAPH OVERVIEW\n%s\n", string(b))
			}
			os.Exit(0)
		}
		if arg == "sync" {
			summary := "Manual Summary"
			for i, v := range os.Args {
				if (v == "--summary" || v == "-s") && i+1 < len(os.Args) {
					summary = os.Args[i+1]
				}
			}

			// Sync Graphify — run from project root (one level above α/)
			fmt.Println("Syncing graphify...")
			syncCmd := exec.Command("graphify", "update", root)
			syncCmd.Dir = root
			syncCmd.Run()

			// Generate project_summary.json
			if err := generateProjectSummary(root); err != nil {
				fmt.Fprintf(os.Stderr, "warning: project_summary: %v\n", err)
			} else {
				fmt.Println("Generated .graphify-out/project_summary.json")
			}

			// Merge summary into latest_state.md
			dir := filepath.Join(root, "α/memories")
			os.MkdirAll(dir, 0755)
			latestStatePath := filepath.Join(dir, "latest_state.md")
			ts := time.Now().Format("2006-01-02 15:04")
			entry := fmt.Sprintf("## %s\n\n%s\n\n---\n\n", ts, summary)
			existing, _ := os.ReadFile(latestStatePath)
			os.WriteFile(latestStatePath, append([]byte(entry), existing...), 0644)
			fmt.Printf("Saved to: %s\n", latestStatePath)

			// Always open browser as requested
			fmt.Println("Opening graph visualization...")
			exec.Command("open", filepath.Join(root, ".graphify-out/graph.html")).Run()
			os.Exit(0)
		}
		if arg == "forget" {
			pattern := ""
			autoConfirm := false

			// Simple flag parsing
			for i := 2; i < len(os.Args); i++ {
				v := os.Args[i]
				if v == "-y" || v == "--yes" {
					autoConfirm = true
				} else if pattern == "" {
					pattern = v
				}
			}

			memDir := filepath.Join(root, "α/memories")
			var targets []string

			if pattern == "" {
				matches, _ := filepath.Glob(filepath.Join(memDir, "session_summary_*.md"))
				if len(matches) > 0 {
					sort.Strings(matches)
					targets = append(targets, matches[len(matches)-1])
				}
			} else {
				// 1. Try as a direct path (relative to CWD or absolute)
				matches, _ := filepath.Glob(pattern)
				if len(matches) > 0 {
					targets = matches
				} else {
					// 2. Try relative to memories directory
					targets, _ = filepath.Glob(filepath.Join(memDir, pattern))
				}
			}

			if len(targets) == 0 {
				fmt.Println("❌ No files found matching your request.")
				os.Exit(1)
			}

			fmt.Println("⚠️  THE FOLLOWING WILL BE PERMANENTLY REMOVED:")
			var graphFiles []string
			for _, f := range targets {
				fmt.Printf("  - %s (Memory)\n", filepath.Base(f))

				// Find matching graph memory
				parts := strings.Split(filepath.Base(f), "_")
				if len(parts) >= 4 {
					// Format: session_summary_20260506_0715.md
					ts := strings.TrimSuffix(parts[2]+"_"+parts[3], ".md")
					gMatches, _ := filepath.Glob(filepath.Join(root, ".graphify-out/memory/*"+ts+"*"))
					for _, g := range gMatches {
						fmt.Printf("  - %s (Graph Memory)\n", filepath.Base(g))
						graphFiles = append(graphFiles, g)
					}
				}
			}

			if !autoConfirm {
				fmt.Print("\nConfirm deletion and Graph update? (y/n): ")
				var confirm string
				fmt.Scanln(&confirm)

				if strings.ToLower(confirm) != "y" {
					fmt.Println("🛑 Aborted.")
					os.Exit(0)
				}
			}

			for _, f := range targets {
				os.Remove(f)
			}
			for _, f := range graphFiles {
				os.Remove(f)
			}
			fmt.Println("✅ Files deleted. Syncing graph...")
			forgetSyncCmd := exec.Command("rtk", "run", "graphify", root, "--update", "--force")
			forgetSyncCmd.Dir = root
			forgetSyncCmd.Run()
			fmt.Println("✨ Knowledge Graph updated.")
			os.Exit(0)
		}
		if arg == "overview" {
			ov, err := graphOverview(root)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			b, _ := json.MarshalIndent(ov, "", "  ")
			fmt.Println(string(b))
			os.Exit(0)
		}
		if arg == "sketch" {
			query := ""
			depth := 3
			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] == "--query" && i+1 < len(os.Args) {
					query = os.Args[i+1]
				}
				if os.Args[i] == "--depth" && i+1 < len(os.Args) {
					fmt.Sscanf(os.Args[i+1], "%d", &depth)
				}
			}
			if query == "" {
				fmt.Fprintln(os.Stderr, "Usage: sketch --query <text> [--depth N]")
				os.Exit(2)
			}
			g, err := loadFullGraph(root)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			fmt.Println(sketchGraph(g, query, depth))
			os.Exit(0)
		}
		if arg == "detail" {
			ids := ""
			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] == "--ids" && i+1 < len(os.Args) {
					ids = os.Args[i+1]
				}
			}
			if ids == "" {
				fmt.Fprintln(os.Stderr, "Usage: detail --ids <id1,id2,...>")
				os.Exit(2)
			}
			g, err := loadFullGraph(root)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			fmt.Println(detailNodes(g, strings.Split(ids, ",")))
			os.Exit(0)
		}
		if arg == "focus" {
			if len(os.Args) < 4 {
				fmt.Println("Usage: focus <path> <term> [max_lines]")
				os.Exit(1)
			}
			path := os.Args[2]
			term := os.Args[3]
			maxLines := 20
			if len(os.Args) > 4 {
				fmt.Sscanf(os.Args[4], "%d", &maxLines)
			}

			absPath := path
			if !filepath.IsAbs(path) {
				absPath = filepath.Join(root, path)
			}
			content, err := os.ReadFile(absPath)
			if err != nil {
				fmt.Printf("❌ Could not read file %s: %v\n", absPath, err)
				os.Exit(1)
			}

			lines := strings.Split(string(content), "\n")
			foundLine := -1
			for i, line := range lines {
				if strings.Contains(line, term) {
					foundLine = i + 1
					break
				}
			}

			if foundLine == -1 {
				fmt.Printf("❌ Term '%s' not found in %s\n", term, path)
				os.Exit(1)
			}

			start := foundLine - 1
			end := foundLine + maxLines
			if end > len(lines) {
				end = len(lines)
			}

			fmt.Printf("🎯 Found '%s' at line %d in %s\n\n", term, foundLine, path)
			for i := start; i < end; i++ {
				fmt.Printf("%d: %s\n", i+1, lines[i])
			}
			os.Exit(0)
		}

		if arg == "build" {
			// Detect platform suffix
			platform := "darwin"
			if p := os.Getenv("GOOS"); p != "" {
				platform = p
			}

			// Source dir is next to this binary's source (walk up to find go.mod)
			exe, _ := os.Executable()
			exeAbs, _ := filepath.EvalSymlinks(exe)
			srcDir := ""
			dir := filepath.Dir(exeAbs)
			for d := dir; d != filepath.Dir(d); d = filepath.Dir(d) {
				if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
					srcDir = d
					break
				}
			}
			if srcDir == "" {
				srcDir = filepath.Join(root, "α/tools/graphify")
			}

			// Build to temp file
			tmp := filepath.Join(os.TempDir(), "graphify-build")
			fmt.Printf("Building from %s ...\n", srcDir)
			cmd := exec.Command("go", "build", "-o", tmp, ".")
			cmd.Dir = srcDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Build failed: %v\n", err)
				os.Exit(1)
			}

			// Output to α/bin/<platform>/
			dests := []string{
				filepath.Join(root, "α/bin/"+platform+"/graphify"),
			}
			data, err := os.ReadFile(tmp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Could not read built binary: %v\n", err)
				os.Exit(1)
			}
			for _, dest := range dests {
				if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
					fmt.Fprintf(os.Stderr, "❌ mkdir %s: %v\n", dest, err)
					continue
				}
				if err := os.WriteFile(dest, data, 0755); err != nil {
					fmt.Fprintf(os.Stderr, "❌ Write %s: %v\n", dest, err)
					continue
				}
				fmt.Printf("✅ %s\n", dest)
			}
			os.Remove(tmp)
			os.Exit(0)
		}

		// Fallback for interactive terminal on unrecognized CLI command/flag
		if isInputTTY() {
			fmt.Fprintf(os.Stderr, "Error: Unrecognized command or flag '%s'\n\n", arg)
			printHelp()
			os.Exit(1)
		}
	}

	// Normal MCP Server Mode
	slog.Info("Starting server", "root", root, "wd", absWd)

	// Tool: debug_info
	s.AddTool(mcp.Tool{
		Name:        "debug_info",
		Description: "Check server environment for debugging.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wd, _ := os.Getwd()
		absWd, _ := filepath.Abs(wd)
		return mcp.NewToolResultText(fmt.Sprintf("Root: %s\nWD: %s\nExe: %s", root, absWd, os.Args[0])), nil
	})
	s.AddTool(mcp.Tool{
		Name:        "focus",
		Description: "Search for a term in a file and read a range of lines around it.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the file to read (relative to project root)",
				},
				"term": map[string]any{
					"type":        "string",
					"description": "Search term to find the starting line",
				},
				"max_lines": map[string]any{
					"type":        "integer",
					"description": "Number of lines to read after the match (default: 20)",
				},
			},
			Required: []string{"path", "term"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments format"), nil
		}

		path, _ := args["path"].(string)
		term, _ := args["term"].(string)
		maxLines := 20
		if val, ok := args["max_lines"].(float64); ok {
			maxLines = int(val)
		}

		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(root, path)
		}
		content, err := os.ReadFile(absPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("could not read file %s: %v", absPath, err)), nil
		}

		lines := strings.Split(string(content), "\n")
		foundLine := -1
		for i, line := range lines {
			if strings.Contains(line, term) {
				foundLine = i + 1
				break
			}
		}

		if foundLine == -1 {
			return mcp.NewToolResultError(fmt.Sprintf("term '%s' not found", term)), nil
		}

		start := foundLine - 1
		end := foundLine + maxLines
		if end > len(lines) {
			end = len(lines)
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found '%s' at line %d in %s\n\n", term, foundLine, path))
		for i := start; i < end; i++ {
			result.WriteString(fmt.Sprintf("%d: %s\n", i+1, lines[i]))
		}

		return mcp.NewToolResultText(result.String()), nil
	})

	// Tool: alpha
	s.AddTool(mcp.Tool{
		Name:        "alpha",
		Description: "Display the Alpha System identity — logo, stats, commands.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Strip ANSI for MCP (plain text output to AI)
		raw := alphaDisplay(root)
		ansiStrip := func(s string) string {
			var out strings.Builder
			skip := false
			for _, r := range s {
				if r == '\033' {
					skip = true
					continue
				}
				if skip {
					if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
						skip = false
					}
					continue
				}
				out.WriteRune(r)
			}
			return out.String()
		}
		return mcp.NewToolResultText(ansiStrip(raw)), nil
	})

	// Tool: awake
	s.AddTool(mcp.Tool{
		Name:        "awake",
		Description: "Run the session start hook.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var out strings.Builder
		out.WriteString("[AGENT_CONTEXT_START]\n\n")

		// Phase 0: compact graph overview (replaces heavy GRAPH_REPORT.md read)
		if ov, err := graphOverview(root); err == nil {
			b, _ := json.MarshalIndent(ov, "", "  ")
			out.WriteString("### GRAPH OVERVIEW\n")
			out.WriteString(string(b))
			out.WriteString("\n\n")
			out.WriteString("KNOWLEDGE_GRAPH: Active (.graphify-out/graph.json)\n")
			out.WriteString("STRATEGY: Use overview→sketch→detail flow. overview() loaded. Call sketch(query) for Phase 1 BFS, then detail(ids) for Phase 2 callers/callees.\n\n")
		}

		if content, err := os.ReadFile(filepath.Join(root, "α/memories/latest_state.md")); err == nil {
			out.WriteString("### PREVIOUS SESSION SUMMARY\n")
			out.WriteString(string(content))
			out.WriteString("\n")
		}

		out.WriteString("[AGENT_CONTEXT_END]\n\n")
		return mcp.NewToolResultText(out.String()), nil
	})

	// Tool: sync
	s.AddTool(mcp.Tool{
		Name:        "sync",
		Description: "Run the session end hook.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"summary": map[string]any{
					"type":        "string",
					"description": "The session summary to save",
				},
			},
			Required: []string{"summary"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		summary, _ := args["summary"].(string)

		// 1. Synchronous Graphify Update — run from project root (one level above α/)
		slog.Info("Syncing graphify before saving memory", "path", root)
		updateCmd := exec.Command("rtk", "run", "graphify", root, "--update")
		updateCmd.Dir = root
		if err := updateCmd.Run(); err != nil {
			slog.Error("Graphify update failed", "error", err)
		}

		// 1b. Generate project_summary.json
		if err := generateProjectSummary(root); err != nil {
			slog.Warn("project_summary generation failed", "error", err)
		}

		// 2. Merge summary into latest_state.md
		dir := filepath.Join(root, "α/memories")
		os.MkdirAll(dir, 0755)
		latestStatePath := filepath.Join(dir, "latest_state.md")
		ts := time.Now().Format("2006-01-02 15:04")
		entry := fmt.Sprintf("## %s\n\n%s\n\n---\n\n", ts, summary)
		existing, _ := os.ReadFile(latestStatePath)
		os.WriteFile(latestStatePath, append([]byte(entry), existing...), 0644)

		// Always open browser as requested
		exec.Command("open", filepath.Join(root, ".graphify-out/graph.html")).Run()

		return mcp.NewToolResultText(fmt.Sprintf("Graph updated. Summary saved to latest_state.md.\nRoot: %s", root)), nil
	})

	// Tool: sketch (Phase 1) — compact BFS subgraph for a query
	s.AddTool(mcp.Tool{
		Name:        "sketch",
		Description: "Phase 1: BFS traversal of the knowledge graph. Returns a compact node list (id, label, file, depth, relation) so the agent can assess relevance before asking for full details.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Natural language query, e.g. 'login flow' or 'authentication'",
				},
				"depth": map[string]any{
					"type":        "integer",
					"description": "BFS depth from seed nodes (default: 3)",
				},
			},
			Required: []string{"query"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		query, _ := args["query"].(string)
		depth := 3
		if d, ok := args["depth"].(float64); ok {
			depth = int(d)
		}
		g, err := loadFullGraph(root)
		if err != nil {
			return mcp.NewToolResultError("load graph: " + err.Error()), nil
		}
		return mcp.NewToolResultText(sketchGraph(g, query, depth)), nil
	})

	// Tool: detail (Phase 2) — full callers/callees for chosen nodes
	s.AddTool(mcp.Tool{
		Name:        "detail",
		Description: "Phase 2: Given node IDs from sketch, returns full callers, callees, file type, and community. Call this only for nodes confirmed relevant in Phase 1.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"ids": map[string]any{
					"type":        "string",
					"description": "Comma-separated node IDs from sketch output, e.g. 'authservice,authservice_attempt'",
				},
			},
			Required: []string{"ids"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		idsStr, _ := args["ids"].(string)
		ids := strings.Split(idsStr, ",")
		g, err := loadFullGraph(root)
		if err != nil {
			return mcp.NewToolResultError("load graph: " + err.Error()), nil
		}
		return mcp.NewToolResultText(detailNodes(g, ids)), nil
	})

	// Tool: overview (Phase 0) — compact graph summary < 200 tokens
	s.AddTool(mcp.Tool{
		Name:        "overview",
		Description: "Phase 0: Returns a compact graph summary (nodes, edges, community count, god nodes, top communities) under 200 tokens. Call this first so the agent knows what to query before running sketch.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ov, err := graphOverview(root)
		if err != nil {
			return mcp.NewToolResultError("overview: " + err.Error()), nil
		}
		b, _ := json.MarshalIndent(ov, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	})

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		os.Exit(0)
	}()

	if err := server.ServeStdio(s); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
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
