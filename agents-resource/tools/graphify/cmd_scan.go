package main

// cmd_scan.go — "graphify update <target> [--force]" CLI + registerMCPUpdate.
// Go-only AST extraction: no Python, no LLM required for code graph.
// Semantic labels (community names) use .env API key or agent fallback.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ── Graph JSON types (matches graph.go gNode/gLink + analysis format) ────────

type graphNode struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	NormLabel      string `json:"norm_label"`
	FileType       string `json:"file_type"`
	SourceFile     string `json:"source_file"`
	SourceLocation string `json:"source_location"`
	Community      int    `json:"community"`
}

type graphLink struct {
	Source          string  `json:"source"`
	Target          string  `json:"target"`
	Relation        string  `json:"relation"`
	Confidence      string  `json:"confidence"`
	SourceFile      string  `json:"source_file"`
	SourceLocation  string  `json:"source_location"`
	Weight          float64 `json:"weight"`
	ConfidenceScore float64 `json:"confidence_score"`
}

type graphJSON struct {
	Directed    bool        `json:"directed"`
	Multigraph  bool        `json:"multigraph"`
	Graph       struct{}    `json:"graph"`
	Nodes       []graphNode `json:"nodes"`
	Links       []graphLink `json:"links"`
	Hyperedges  []any       `json:"hyperedges"`
}

type analysisJSON struct {
	Communities map[string][]string `json:"communities"`
	Cohesion    map[string]float64  `json:"cohesion"`
	Gods        []GodNode           `json:"gods"`
	Surprises   []any               `json:"surprises"`
	Tokens      int                 `json:"tokens"`
}

// ── Build ─────────────────────────────────────────────────────────────────────

// BuildResult holds the outcome of a scan.
type BuildResult struct {
	Nodes       int
	Edges       int
	Communities int
	NeedLabels  bool   // true when no API key was found
	LabelPrompt string // agent fallback prompt
}

func buildGraph(scanRoot, alphaDir, projectRoot string, force bool) (*BuildResult, error) {
	outDir := graphifyDataDir(alphaDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}

	// Check existing graph
	graphPath := filepath.Join(outDir, "graph.json")
	if !force {
		if _, err := os.Stat(graphPath); err == nil {
			return runIncrementalUpdate(scanRoot, alphaDir, projectRoot, outDir)
		}
	}
	return runFullBuild(scanRoot, alphaDir, projectRoot, outDir)
}

func runFullBuild(scanRoot, alphaDir, projectRoot, outDir string) (*BuildResult, error) {
	files, err := scanProject(scanRoot, alphaDir)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return assembleGraph(files, scanRoot, alphaDir, projectRoot, outDir)
}

func runIncrementalUpdate(scanRoot, alphaDir, projectRoot, outDir string) (*BuildResult, error) {
	// For now incremental = full rebuild (fast enough for typical projects)
	return runFullBuild(scanRoot, alphaDir, projectRoot, outDir)
}

func assembleGraph(files []ExtractedFile, scanRoot, alphaDir, projectRoot, outDir string) (*BuildResult, error) {
	gcfg := loadGraphConfig(alphaDir)
	resolvedIncludes := resolveIncludes(gcfg.Include, scanRoot, alphaDir)
	// ── Step 1: collect all nodes + edges, assign IDs ─────────────────────────
	type nodeKey struct{ id string }

	nodeByID := map[string]*graphNode{}
	var nodeOrder []string
	var rawLinks []graphLink

	ensureNode := func(id, label, normLabel, fileType, srcFile, srcLoc string) {
		if _, ok := nodeByID[id]; !ok {
			nodeByID[id] = &graphNode{
				ID: id, Label: label, NormLabel: normLabel,
				FileType: fileType, SourceFile: srcFile, SourceLocation: srcLoc,
			}
			nodeOrder = append(nodeOrder, id)
		}
	}

	for _, ef := range files {
		if len(ef.Nodes) == 0 {
			// doc/config file: emit as single file node
			fileLabel := lastPathSegment(ef.RelPath)
			fileID := makeID(ef.RelPath, "")
			srcFile := relToAlpha(ef.RelPath, scanRoot, projectRoot)
			ensureNode(fileID, fileLabel, normalize(fileLabel), "doc", srcFile, "L1")
			continue
		}

		srcFile := relToAlpha(ef.RelPath, scanRoot, projectRoot)
		dirID := dirPrefix(ef.RelPath)
		isKnowledge := isIncluded(ef.RelPath, resolvedIncludes)

		for _, n := range ef.Nodes {
			id := makeID(ef.RelPath, n.Label)
			if n.Kind == "file" {
				id = dirID + "_" + slugify(n.Label)
			}
			fileType := "code"
			if isKnowledge {
				fileType = "knowledge"
			} else if n.Kind == "file" && (ef.Language == "md" || ef.Language == "text") {
				fileType = "doc"
			}
			ensureNode(id, n.Label, normalize(n.Label), fileType, srcFile, n.Location)
		}

		for _, e := range ef.Edges {
			srcID := makeID(ef.RelPath, e.FromLabel)
			if e.FromLabel == lastPathSegment(ef.RelPath) {
				srcID = dirID + "_" + slugify(e.FromLabel)
			}
			tgtID := makeID(ef.RelPath, e.ToLabel)

			// For cross-file references (imports), target may not exist — skip
			if e.Relation == "references" || e.Relation == "imports" {
				rawLinks = append(rawLinks, graphLink{
					Source: srcID, Target: tgtID,
					Relation: e.Relation, Confidence: "EXTRACTED",
					SourceFile: srcFile, SourceLocation: e.Location,
					Weight: 1.0, ConfidenceScore: 1.0,
				})
				continue
			}
			rawLinks = append(rawLinks, graphLink{
				Source: srcID, Target: tgtID,
				Relation: e.Relation, Confidence: "EXTRACTED",
				SourceFile: srcFile, SourceLocation: e.Location,
				Weight: 1.0, ConfidenceScore: 1.0,
			})
		}
	}

	// ── Step 1b: build label→IDs index for cross-file call resolution ────────
	// Maps slugified label (e.g. "buildawakeoverview") → list of node IDs that have that label.
	labelToIDs := map[string][]string{}
	for id, n := range nodeByID {
		key := slugify(n.Label)
		labelToIDs[key] = append(labelToIDs[key], id)
	}

	// Re-resolve unresolved call targets: tgtID was built with source file prefix,
	// but the callee may live in a different file. Extract the label slug from tgtID's
	// last segment and look it up globally. Only resolve when unambiguous (1 match).
	for i := range rawLinks {
		l := &rawLinks[i]
		if l.Relation != "calls" {
			continue
		}
		if _, ok := nodeByID[l.Target]; ok {
			continue // already points to a real node
		}
		parts := strings.Split(l.Target, "_")
		funcSlug := parts[len(parts)-1]
		if candidates, ok := labelToIDs[funcSlug]; ok && len(candidates) == 1 {
			l.Target = candidates[0]
		}
	}

	// ── Step 2: filter links to only reference existing nodes ────────────────
	var links []graphLink
	for _, l := range rawLinks {
		_, srcOK := nodeByID[l.Source]
		_, tgtOK := nodeByID[l.Target]
		if srcOK && tgtOK {
			links = append(links, l)
		}
	}

	// ── Step 3: community detection ──────────────────────────────────────────
	idxByID := map[string]int{}
	nodeIDs := make([]string, len(nodeOrder))
	nodeLabels := make([]string, len(nodeOrder))
	for i, id := range nodeOrder {
		idxByID[id] = i
		nodeIDs[i] = id
		nodeLabels[i] = nodeByID[id].Label
	}

	edges := make([][2]int, 0, len(links))
	for _, l := range links {
		a, aOK := idxByID[l.Source]
		b, bOK := idxByID[l.Target]
		if aOK && bOK {
			edges = append(edges, [2]int{a, b})
		}
	}

	communities := AssignCommunities(nodeIDs, edges, 30)
	if communities == nil {
		communities = make([]int, len(nodeOrder))
	}
	for i, id := range nodeOrder {
		nodeByID[id].Community = communities[i]
	}

	// ── Step 4: build output nodes list ──────────────────────────────────────
	nodes := make([]graphNode, len(nodeOrder))
	for i, id := range nodeOrder {
		nodes[i] = *nodeByID[id]
	}

	// ── Step 5: god nodes ─────────────────────────────────────────────────────
	gods := FindGodNodes(nodeIDs, nodeLabels, edges, 10)

	// ── Step 6: communities map for analysis.json ─────────────────────────────
	commMap := map[string][]string{}
	for i, id := range nodeOrder {
		cid := fmt.Sprintf("%d", communities[i])
		commMap[cid] = append(commMap[cid], id)
	}

	// ── Step 7: write graph.json ──────────────────────────────────────────────
	g := graphJSON{
		Directed: true, Multigraph: false,
		Nodes: nodes, Links: links,
		Hyperedges: []any{},
	}
	if err := writeJSON(filepath.Join(outDir, "graph.json"), g); err != nil {
		return nil, err
	}

	// ── Step 8: write .graphify_analysis.json ────────────────────────────────
	cohesion := map[string]float64{}
	for cid := range commMap {
		cohesion[cid] = 1.0
	}
	a := analysisJSON{
		Communities: commMap,
		Cohesion:    cohesion,
		Gods:        gods,
		Surprises:   []any{},
		Tokens:      0,
	}
	if err := writeJSON(filepath.Join(outDir, ".graphify_analysis.json"), a); err != nil {
		return nil, err
	}

	// ── Step 9: semantic labels ───────────────────────────────────────────────
	nodeLabelMap := map[string]string{}
	for id, n := range nodeByID {
		nodeLabelMap[id] = n.Label
	}

	result := &BuildResult{
		Nodes: len(nodes), Edges: len(links), Communities: len(commMap),
	}

	apiKey := LoadAPIKey(projectRoot)
	if apiKey != nil {
		labels, err := GenerateCommunityLabels(apiKey, commMap, nodeLabelMap)
		if err == nil {
			_ = writeJSON(filepath.Join(outDir, ".graphify_labels.json"), labels)
		}
	} else {
		result.NeedLabels = true
		result.LabelPrompt = AgentFallbackPrompt(commMap, nodeLabelMap)
	}

	return result, nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

// ── ID generation ─────────────────────────────────────────────────────────────

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, "()")
	s = nonAlnum.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	return s
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSuffix(s, "()"))
}

// dirPrefix turns "α/agents-resource/tools/alpha/main.go" → "alpha"
func dirPrefix(rel string) string {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) >= 2 {
		return slugify(parts[len(parts)-2])
	}
	return slugify(parts[0])
}

// makeID builds a node ID from file path + label.
// "α/tools/alpha/main.go" + "runTool()" → "alpha_main_runtool"
func makeID(relPath, label string) string {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	var segments []string

	// Take last 2 meaningful path segments
	meaningful := filterEmpty(parts[:len(parts)-1])
	if len(meaningful) > 0 {
		segments = append(segments, slugify(meaningful[len(meaningful)-1]))
	}

	// File stem
	base := parts[len(parts)-1]
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	if stem != "" {
		segments = append(segments, slugify(stem))
	}

	// Label
	if label != "" {
		segments = append(segments, slugify(label))
	}

	id := strings.Join(segments, "_")
	id = nonAlnum.ReplaceAllString(id, "_")
	id = strings.Trim(id, "_")
	return id
}

func filterEmpty(ss []string) []string {
	var out []string
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// relToAlpha converts scanRoot-relative path to a display path.
func relToAlpha(rel, scanRoot, projectRoot string) string {
	full := filepath.Join(scanRoot, rel)
	if r, err := filepath.Rel(filepath.Dir(projectRoot), full); err == nil {
		return r
	}
	return rel
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// ── CLI command ───────────────────────────────────────────────────────────────

func cliUpdate() {
	target := ""
	force := false
	for _, a := range os.Args[2:] {
		switch a {
		case "--force":
			force = true
		default:
			if !strings.HasPrefix(a, "-") {
				target, _ = filepath.Abs(a)
			}
		}
	}

	// Resolve alphaDir (the α/ dir containing knowledge-graph/)
	// Priority: PROJECT_ROOT env → α/ inside target → root global → target itself
	alphaDir := root
	if target != "" {
		// Check if target contains α/ subdir
		candidate := filepath.Join(target, "α")
		if _, err := os.Stat(filepath.Join(candidate, "knowledge-graph")); err == nil {
			alphaDir = candidate
		} else if _, err := os.Stat(filepath.Join(target, "knowledge-graph")); err == nil {
			// target IS the α/ dir
			alphaDir = target
		}
	}
	if target == "" {
		// Global mode: scan the project specified by ALPHA_ROOT
		if isGlobalMode() {
			if r := os.Getenv("ALPHA_ROOT"); r != "" {
				target = r
			}
		}
		if target == "" {
			target = filepath.Dir(alphaDir) // default scan = project root
		}
	}

	result, err := buildGraph(target, alphaDir, filepath.Dir(alphaDir), force)
	if err != nil {
		fmt.Fprintln(os.Stderr, "update error:", err)
		os.Exit(1)
	}

	fmt.Printf("[graphify watch] Rebuilt: %d nodes, %d edges, %d communities\n",
		result.Nodes, result.Edges, result.Communities)
	fmt.Println("[graphify watch] graph.json and .graphify_analysis.json updated")

	if result.NeedLabels {
		fmt.Print("\n" + result.LabelPrompt)
	}
}

// cliExtract runs scan+build+semantic labeling with an optional explicit backend override.
// Usage: graphify extract [path] [--backend gemini|anthropic|openai] [--force]
func cliExtract() {
	target := ""
	force := false
	backendOverride := ""
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--force":
			force = true
		case a == "--backend" && i+1 < len(args):
			backendOverride = args[i+1]
			i++
		case strings.HasPrefix(a, "--backend="):
			backendOverride = strings.TrimPrefix(a, "--backend=")
		case !strings.HasPrefix(a, "-"):
			target, _ = filepath.Abs(a)
		}
	}

	alphaDir := root
	if target != "" {
		candidate := filepath.Join(target, "α")
		if _, err := os.Stat(filepath.Join(candidate, "knowledge-graph")); err == nil {
			alphaDir = candidate
		} else if _, err := os.Stat(filepath.Join(target, "knowledge-graph")); err == nil {
			alphaDir = target
		}
	}
	if target == "" {
		if isGlobalMode() {
			if r := os.Getenv("ALPHA_ROOT"); r != "" {
				target = r
			}
		}
		if target == "" {
			target = filepath.Dir(alphaDir)
		}
	}

	// Build graph (scan + AST extract + community detection)
	result, err := buildGraph(target, alphaDir, filepath.Dir(alphaDir), force)
	if err != nil {
		fmt.Fprintln(os.Stderr, "extract error:", err)
		os.Exit(1)
	}
	fmt.Printf("Extracted: %d nodes, %d edges, %d communities\n", result.Nodes, result.Edges, result.Communities)

	// Semantic labeling — use explicit backend or auto-detect from .env
	projectRoot := filepath.Dir(alphaDir)
	apiKey := LoadAPIKey(projectRoot)
	if apiKey == nil {
		fmt.Fprintln(os.Stderr, "⚠️  No API key found in .env — skipping semantic labels")
		if result.NeedLabels {
			fmt.Print("\n" + result.LabelPrompt)
		}
		return
	}
	// Override provider if --backend specified
	if backendOverride != "" {
		apiKey.Provider = backendOverride
		fmt.Printf("Backend: %s (forced)\n", backendOverride)
	} else {
		fmt.Printf("Backend: %s (from .env)\n", apiKey.Provider)
	}

	// Load existing graph to get community map for labeling
	outDir := graphifyDataDir(alphaDir)
	g, err := loadFullGraph(alphaDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "⚠️  Could not load graph for labeling:", err)
		return
	}
	commMap := map[string][]string{}
	nodeLabelMap := map[string]string{}
	for _, n := range g.Nodes {
		cid := fmt.Sprintf("%d", n.Community)
		commMap[cid] = append(commMap[cid], n.ID)
		nodeLabelMap[n.ID] = n.Label
	}
	labels, err := GenerateCommunityLabels(apiKey, commMap, nodeLabelMap)
	if err != nil {
		fmt.Fprintln(os.Stderr, "⚠️  Label generation failed:", err)
		return
	}
	// Write labels to analysis JSON
	analysisPath := filepath.Join(outDir, ".graphify_analysis.json")
	existing := analysisJSON{}
	if data, err := os.ReadFile(analysisPath); err == nil {
		json.Unmarshal(data, &existing)
	}
	for cid, name := range labels {
		existing.Communities[cid] = append(existing.Communities[cid], name)
	}
	if data, err := json.MarshalIndent(existing, "", "  "); err == nil {
		os.WriteFile(analysisPath, data, 0644)
	}
	fmt.Printf("✅ %d community labels written\n", len(labels))
}

// ── MCP tool ──────────────────────────────────────────────────────────────────

func registerMCPUpdate(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "update",
		Description: "Rebuild AST knowledge graph (Go-only, no Python). Extracts functions, types, call edges. Returns label prompt if no API key found in .env.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"force": map[string]any{"type": "boolean", "description": "Force full rebuild even if graph.json exists"},
			},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		force, _ := args["force"].(bool)

		projectRoot := filepath.Dir(root)
		result, err := buildGraph(projectRoot, root, projectRoot, force)
		if err != nil {
			return mcp.NewToolResultError("update: " + err.Error()), nil
		}

		msg := fmt.Sprintf("Graph updated: %d nodes, %d edges, %d communities.", result.Nodes, result.Edges, result.Communities)
		if result.NeedLabels {
			msg += "\n\n" + result.LabelPrompt
		}
		return mcp.NewToolResultText(msg), nil
	})
}

// registerMCPSetLabels allows the agent to provide community names when no API key.
func registerMCPSetLabels(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "set_labels",
		Description: "Write community labels to .graphify_labels.json. Call after 'update' returns a label prompt.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"labels": map[string]any{
					"type":        "string",
					"description": `JSON string mapping community IDs to names, e.g. {"0":"Auth Layer","1":"Graph Core"}`,
				},
			},
			Required: []string{"labels"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := req.Params.Arguments.(map[string]any)
		labelsStr, _ := args["labels"].(string)

		var labels map[string]string
		if err := json.Unmarshal([]byte(labelsStr), &labels); err != nil {
			return mcp.NewToolResultError("invalid labels JSON: " + err.Error()), nil
		}

		outPath := filepath.Join(graphifyDataDir(root), ".graphify_labels.json")
		if err := writeJSON(outPath, labels); err != nil {
			return mcp.NewToolResultError("write labels: " + err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Labels saved: %d communities named.", len(labels))), nil
	})
}
