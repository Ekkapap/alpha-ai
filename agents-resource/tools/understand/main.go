package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	projectRoot  = "" // α/ dir — scripts + knowledge-graph/ base
	dataRoot     = "" // per-project data dir (global: ~/.alpha-ai/knowledge-graph/projects/<id>/, local: projectRoot)
	gitRoot      = "" // project root — git operations (.git lives here)
	pluginRoot   = ""
	scriptFolder = ""
)

func logInfo(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func logInfof(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// Metadata schema for Meta file
type Meta struct {
	LastAnalyzedAt string `json:"lastAnalyzedAt"`
	GitCommitHash  string `json:"gitCommitHash"`
	Version        string `json:"version"`
	AnalyzedFiles  int    `json:"analyzedFiles"`
}

// Graph structures
type GraphNode struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	FilePath   string   `json:"filePath,omitempty"`
	Summary    string   `json:"summary"`
	Tags       []string `json:"tags"`
	Complexity string   `json:"complexity"`
}

type GraphEdge struct {
	Source    string  `json:"source"`
	Target    string  `json:"target"`
	Type      string  `json:"type"`
	Direction string  `json:"direction"`
	Weight    float64 `json:"weight"`
}

type Layer struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	NodeIds     []string `json:"nodeIds"`
}

type TourStep struct {
	Order       int      `json:"order"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	NodeIds     []string `json:"nodeIds"`
}

type ProjectInfo struct {
	Name          string   `json:"name"`
	Languages     []string `json:"languages"`
	Frameworks    []string `json:"frameworks"`
	Description   string   `json:"description"`
	AnalyzedAt    string   `json:"analyzedAt"`
	GitCommitHash string   `json:"gitCommitHash"`
}

type KnowledgeGraph struct {
	Version string      `json:"version"`
	Project ProjectInfo `json:"project"`
	Nodes   []GraphNode `json:"nodes"`
	Edges   []GraphEdge `json:"edges"`
	Layers  []Layer     `json:"layers"`
	Tour    []TourStep  `json:"tour"`
}

// apiEnv reads .env from gitRoot and returns API key env vars to inject into subprocesses.
// The Node.js analysis scripts (generate-batches.mjs etc.) read these from process.env.
func apiEnv() []string {
	envFile := filepath.Join(gitRoot, ".env")
	data, err := os.ReadFile(envFile)
	if err != nil {
		return nil
	}
	var extras []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		for _, key := range []string{"ANTHROPIC_API_KEY", "GEMINI_API_KEY", "GOOGLE_API_KEY", "OPENAI_API_KEY"} {
			if strings.HasPrefix(line, key+"=") {
				extras = append(extras, line)
			}
		}
	}
	return extras
}

func runCmd(command string, args []string, dir string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	// Inject API keys from .env so Node.js analysis scripts can call LLM APIs
	cmd.Env = append(os.Environ(), apiEnv()...)
	return cmd.Run()
}

func startPipeline() error {
	logInfo("🎬 Starting full analysis pipeline...")

	// 1. Scan Project
	logInfo("\n[Phase 1/5] Scanning project files...")
	scanScript := filepath.Join(pluginRoot, "skills/understand/scan-project.mjs")
	scanResult := filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate/scan-result.json")

	// Clean stale batches from previous runs
	os.RemoveAll(filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate"))
	os.MkdirAll(filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate"), 0755)
	if err := runCmd("node", []string{scanScript, projectRoot, scanResult}, projectRoot); err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// 2. Louvain Batching
	logInfo("\n[Phase 2/5] Computing Louvain batches...")
	batchScript := filepath.Join(pluginRoot, "skills/understand/compute-batches.mjs")
	if err := runCmd("node", []string{batchScript, projectRoot}, projectRoot); err != nil {
		return fmt.Errorf("batching failed: %w", err)
	}

	// 3. Get total batches
	batchesPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate/batches.json")
	data, err := os.ReadFile(batchesPath)
	if err != nil {
		return err
	}
	var bData struct {
		Batches []any `json:"batches"`
	}
	if err := json.Unmarshal(data, &bData); err != nil {
		return err
	}
	totalBatches := len(bData.Batches)

	// 4. Batch Analysis
	logInfof("\n[Phase 3/5] Running analysis for %d batches...", totalBatches)
	analysisScript := filepath.Join(projectRoot, "scripts/generate-batches.mjs")
	if err := runCmd("node", []string{analysisScript, "1", fmt.Sprintf("%d", totalBatches)}, projectRoot); err != nil {
		return fmt.Errorf("batch analysis failed: %w", err)
	}

	// 5. Merge
	logInfo("\n[Phase 4/5] Merging batch files...")
	mergeScript := filepath.Join(pluginRoot, "skills/understand/merge-batch-graphs.py")
	if err := runCmd("python3", []string{mergeScript, projectRoot}, projectRoot); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	// 6. Synthesis Graph and Meta
	logInfo("\n[Phase 5/5] Synthesizing final knowledge-graph.json and meta.json...")
	return finalizeGraph()
}

func updatePipeline() error {
	logInfo("🔄 Starting incremental update pipeline...")

	// Read last meta to get commit hash
	metaPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/meta.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		logInfo("⚠️  Could not read meta.json, falling back to full --start.")
		return startPipeline()
	}
	var lastMeta Meta
	if err := json.Unmarshal(metaData, &lastMeta); err != nil {
		return err
	}

	// Git commit hash
	gitHashCmd := exec.Command("git", "rev-parse", "HEAD")
	gitHashCmd.Dir = gitRoot
	currentHashBytes, err := gitHashCmd.Output()
	if err != nil {
		logInfo("⚠️  No git repo found, skipping hash check — running full update.")
		return startPipeline()
	}
	currentHash := strings.TrimSpace(string(currentHashBytes))

	if lastMeta.GitCommitHash == currentHash {
		logInfo("✅ Codebase is already up-to-date with git commit:", currentHash)
		return nil
	}

	// Git diff changed files
	logInfo("Detecting changed files since:", lastMeta.GitCommitHash)
	diffCmd := exec.Command("git", "diff", lastMeta.GitCommitHash+"..HEAD", "--name-only")
	diffCmd.Dir = gitRoot
	changedFilesBytes, err := diffCmd.Output()
	if err != nil {
		return fmt.Errorf("failed git diff: %w", err)
	}
	changedFilesStr := strings.TrimSpace(string(changedFilesBytes))
	if changedFilesStr == "" {
		logInfo("✅ No files changed.")
		return nil
	}

	tmpDir := filepath.Join(dataRoot, "knowledge-graph/understand-anything/tmp")
	os.MkdirAll(tmpDir, 0755)
	changedFilesTxt := filepath.Join(tmpDir, "changed-files.txt")
	if err := os.WriteFile(changedFilesTxt, []byte(changedFilesStr), 0644); err != nil {
		return err
	}

	// Louvain Batching with --changed-files
	logInfo("Computing batches for changed files...")
	batchScript := filepath.Join(pluginRoot, "skills/understand/compute-batches.mjs")
	if err := runCmd("node", []string{batchScript, projectRoot, "--changed-files=" + changedFilesTxt}, projectRoot); err != nil {
		return fmt.Errorf("batching failed: %w", err)
	}

	// Read batches.json
	batchesPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate/batches.json")
	data, err := os.ReadFile(batchesPath)
	if err != nil {
		return err
	}
	var bData struct {
		Batches []any `json:"batches"`
	}
	if err := json.Unmarshal(data, &bData); err != nil {
		return err
	}
	totalBatches := len(bData.Batches)

	// Batch Analysis
	logInfof("Running analysis for %d incremental batches...", totalBatches)
	analysisScript := filepath.Join(projectRoot, "scripts/generate-batches.mjs")
	if err := runCmd("node", []string{analysisScript, "1", fmt.Sprintf("%d", totalBatches)}, projectRoot); err != nil {
		return fmt.Errorf("batch analysis failed: %w", err)
	}

	// Pruning logic
	logInfo("Pruning unchanged elements and generating batch-existing.json...")
	if err := pruneAndPreserveExisting(changedFilesStr); err != nil {
		return err
	}

	// Merge
	logInfo("Merging all batch files...")
	mergeScript := filepath.Join(pluginRoot, "skills/understand/merge-batch-graphs.py")
	if err := runCmd("python3", []string{mergeScript, projectRoot}, projectRoot); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	// Finalize
	return finalizeGraph()
}

func pruneAndPreserveExisting(changedFilesStr string) error {
	graphPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/knowledge-graph.json")
	graphData, err := os.ReadFile(graphPath)
	if err != nil {
		return err
	}
	var graph KnowledgeGraph
	if err := json.Unmarshal(graphData, &graph); err != nil {
		return err
	}

	changedFilesList := strings.Split(changedFilesStr, "\n")
	changedMap := make(map[string]bool)
	for _, f := range changedFilesList {
		changedMap[strings.TrimSpace(f)] = true
	}

	// Prune nodes matching changed files
	var keptNodes []GraphNode
	keptNodeIds := make(map[string]bool)
	for _, node := range graph.Nodes {
		if node.FilePath != "" && changedMap[node.FilePath] {
			continue // Skip/delete
		}
		keptNodes = append(keptNodes, node)
		keptNodeIds[node.ID] = true
	}

	// Prune edges referencing deleted nodes
	var keptEdges []GraphEdge
	for _, edge := range graph.Edges {
		if keptNodeIds[edge.Source] && keptNodeIds[edge.Target] {
			keptEdges = append(keptEdges, edge)
		}
	}

	// Write batch-existing.json
	batchExisting := map[string]any{
		"nodes": keptNodes,
		"edges": keptEdges,
	}
	prunedBytes, err := json.MarshalIndent(batchExisting, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate/batch-existing.json"), prunedBytes, 0644)
}

func finalizeGraph() error {
	assembledPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/intermediate/assembled-graph.json")
	outGraphPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/knowledge-graph.json")
	outMetaPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/meta.json")

	assembledData, err := os.ReadFile(assembledPath)
	if err != nil {
		return err
	}

	var assembled struct {
		Nodes []GraphNode `json:"nodes"`
		Edges []GraphEdge `json:"edges"`
	}
	if err := json.Unmarshal(assembledData, &assembled); err != nil {
		return err
	}

	// Get git commit hash
	gitHashCmd := exec.Command("git", "rev-parse", "HEAD")
	gitHashCmd.Dir = gitRoot
	currentHashBytes, _ := gitHashCmd.Output()
	currentHash := strings.TrimSpace(string(currentHashBytes))
	if currentHash == "" {
		currentHash = "local-draft"
	}

	// Derive file-to-file imports edges from call graph (cross-file method calls)
	derivedImports := make(map[string]bool)
	var finalEdges []GraphEdge
	finalEdges = append(finalEdges, assembled.Edges...)

	for _, edge := range assembled.Edges {
		if edge.Type == "calls" {
			srcFile := getFilePathFromNodeID(edge.Source)
			tgtFile := getFilePathFromNodeID(edge.Target)

			if srcFile != "" && tgtFile != "" && srcFile != tgtFile {
				importKey := fmt.Sprintf("file:%s -> file:%s", srcFile, tgtFile)
				if !derivedImports[importKey] {
					derivedImports[importKey] = true
					finalEdges = append(finalEdges, GraphEdge{
						Source:    "file:" + srcFile,
						Target:    "file:" + tgtFile,
						Type:      "imports",
						Direction: "forward",
						Weight:    0.7,
					})
				}
			}
		}
	}

	projectName := filepath.Base(projectRoot)
	detectedLangs := detectLanguages(assembled.Nodes)
	detectedFrameworks := detectFrameworks(projectRoot)

	// Create a friendly description dynamically
	desc := fmt.Sprintf("%s-based project named %s, dynamically analyzed and indexed.", formatList(detectedLangs), projectName)

	graph := KnowledgeGraph{
		Version: "1.0.0",
		Project: ProjectInfo{
			Name:          projectName,
			Languages:     detectedLangs,
			Frameworks:    detectedFrameworks,
			Description:   desc,
			AnalyzedAt:    time.Now().UTC().Format(time.RFC3339) + "Z",
			GitCommitHash: currentHash,
		},
		Nodes: assembled.Nodes,
		Edges: finalEdges,
		Layers: []Layer{
			{
				ID:          "layer:default",
				Name:        "Default Layer",
				Description: "Standard application layer containing analyzed components",
				NodeIds:     getFileNodeIds(assembled.Nodes),
			},
		},
		Tour: []TourStep{},
	}

	// Enrich manually for views/configs in app/ if they exist
	enrichNodes(graph.Nodes)

	graphBytes, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(outGraphPath, graphBytes, 0644); err != nil {
		return err
	}

	meta := Meta{
		LastAnalyzedAt: time.Now().UTC().Format(time.RFC3339) + "Z",
		GitCommitHash:  currentHash,
		Version:        "1.0.0",
		AnalyzedFiles:  countFiles(assembled.Nodes),
	}
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(outMetaPath, metaBytes, 0644); err != nil {
		return err
	}

	logInfo("🚀 Synthesized knowledge-graph.json and meta.json successfully!")
	return nil
}

func detectLanguages(nodes []GraphNode) []string {
	langMap := make(map[string]bool)
	for _, n := range nodes {
		path := n.FilePath
		if path == "" {
			path = n.Name
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".php":
			langMap["PHP"] = true
		case ".go":
			langMap["Go"] = true
		case ".js", ".mjs", ".jsx":
			langMap["JavaScript"] = true
		case ".ts", ".tsx":
			langMap["TypeScript"] = true
		case ".py":
			langMap["Python"] = true
		case ".java":
			langMap["Java"] = true
		case ".rb":
			langMap["Ruby"] = true
		case ".rs":
			langMap["Rust"] = true
		case ".c", ".h", ".cpp", ".hpp":
			langMap["C/C++"] = true
		case ".cs":
			langMap["C#"] = true
		case ".sh":
			langMap["Shell"] = true
		case ".md":
			langMap["Markdown"] = true
		case ".html":
			langMap["HTML"] = true
		case ".css":
			langMap["CSS"] = true
		}
	}
	var langs []string
	for l := range langMap {
		langs = append(langs, l)
	}
	sort.Strings(langs)
	if len(langs) == 0 {
		langs = []string{"Unknown"}
	}
	return langs
}

func detectFrameworks(root string) []string {
	frameworkMap := make(map[string]bool)

	contains := func(filePath string, terms ...string) bool {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return false
		}
		content := strings.ToLower(string(data))
		for _, term := range terms {
			if strings.Contains(content, strings.ToLower(term)) {
				return true
			}
		}
		return false
	}

	packageJSON := filepath.Join(root, "package.json")
	if _, err := os.Stat(packageJSON); err == nil {
		if contains(packageJSON, "react") {
			frameworkMap["React"] = true
		}
		if contains(packageJSON, "vue") {
			frameworkMap["Vue"] = true
		}
		if contains(packageJSON, "angular") {
			frameworkMap["Angular"] = true
		}
		if contains(packageJSON, "next") {
			frameworkMap["Next.js"] = true
		}
		if contains(packageJSON, "nuxt") {
			frameworkMap["Nuxt.js"] = true
		}
		if contains(packageJSON, "svelte") {
			frameworkMap["Svelte"] = true
		}
		if contains(packageJSON, "tailwind") {
			frameworkMap["Tailwind"] = true
		}
		if contains(packageJSON, "apexcharts") {
			frameworkMap["ApexCharts"] = true
		}
		if contains(packageJSON, "express") {
			frameworkMap["Express"] = true
		}
	}

	composerJSON := filepath.Join(root, "composer.json")
	if _, err := os.Stat(composerJSON); err == nil {
		if contains(composerJSON, "laravel/framework", "laravel/lumen") {
			frameworkMap["Laravel"] = true
		} else if contains(composerJSON, "symfony/") {
			frameworkMap["Symfony"] = true
		} else if contains(composerJSON, "yiisoft/yii2") {
			frameworkMap["Yii2"] = true
		} else if contains(composerJSON, "codeigniter/") {
			frameworkMap["CodeIgniter"] = true
		}
		if len(frameworkMap) == 0 {
			frameworkMap["MVC"] = true
		}
	}

	goMod := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goMod); err == nil {
		if contains(goMod, "github.com/gin-gonic/gin") {
			frameworkMap["Gin"] = true
		}
		if contains(goMod, "github.com/labstack/echo") {
			frameworkMap["Echo"] = true
		}
		if contains(goMod, "github.com/fiber/fiber") {
			frameworkMap["Fiber"] = true
		}
	}

	if len(frameworkMap) == 0 {
		if _, err := os.Stat(filepath.Join(root, "app")); err == nil {
			frameworkMap["MVC"] = true
		} else {
			frameworkMap["Custom"] = true
		}
	}

	var frameworks []string
	for f := range frameworkMap {
		frameworks = append(frameworks, f)
	}
	sort.Strings(frameworks)
	return frameworks
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "software"
	}
	if len(items) == 1 {
		return items[0]
	}
	if len(items) == 2 {
		return items[0] + " and " + items[1]
	}
	return strings.Join(items[:len(items)-1], ", ") + ", and " + items[len(items)-1]
}

func getFileNodeIds(nodes []GraphNode) []string {
	fileTypes := map[string]bool{
		"file": true, "config": true, "document": true, "service": true,
		"pipeline": true, "table": true, "schema": true, "resource": true, "endpoint": true,
	}
	var ids []string
	for _, n := range nodes {
		if fileTypes[n.Type] {
			ids = append(ids, n.ID)
		}
	}
	return ids
}

func getFilePathFromNodeID(id string) string {
	if !strings.HasPrefix(id, "function:") && !strings.HasPrefix(id, "class:") {
		return ""
	}
	parts := strings.Split(id, ":")
	if len(parts) >= 3 {
		return parts[1]
	}
	return ""
}

func countFiles(nodes []GraphNode) int {
	c := 0
	for _, n := range nodes {
		if n.Type == "file" {
			c++
		}
	}
	return c
}

func enrichNodes(nodes []GraphNode) {
	for i := range nodes {
		if nodes[i].Type != "file" || !strings.HasPrefix(nodes[i].FilePath, "app/") {
			continue
		}

		absPath := filepath.Join(projectRoot, nodes[i].FilePath)
		contentBytes, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}

		content := string(contentBytes)
		summary := ""

		// Heuristic 1: Look for PHP Docstrings / Block Comments (/** ... */ or /* ... */)
		if startIdx := strings.Index(content, "/*"); startIdx != -1 {
			if endIdx := strings.Index(content[startIdx:], "*/"); endIdx != -1 {
				comment := content[startIdx+2 : startIdx+endIdx]
				summary = cleanComment(comment)
			}
		}

		// Heuristic 2: Look for HTML Comments (<!-- ... -->) commonly found at the top of views
		if summary == "" {
			if startIdx := strings.Index(content, "<!--"); startIdx != -1 {
				if endIdx := strings.Index(content[startIdx:], "-->"); endIdx != -1 {
					comment := content[startIdx+4 : startIdx+endIdx]
					summary = cleanComment(comment)
				}
			}
		}

		// Heuristic 3: Default fallback descriptive summaries
		if summary == "" {
			nameLower := strings.ToLower(nodes[i].Name)
			if strings.HasPrefix(nodes[i].FilePath, "app/config/") {
				summary = fmt.Sprintf("Configuration settings file for the %s module.", strings.TrimSuffix(nodes[i].Name, ".php"))
			} else if strings.HasPrefix(nodes[i].FilePath, "app/views/") {
				summary = fmt.Sprintf("View template rendering the %s component.", strings.TrimSuffix(nameLower, ".php"))
			} else {
				summary = fmt.Sprintf("PHP source file representing the %s component.", nodes[i].Name)
			}
		}

		// Set dynamic properties
		nodes[i].Summary = summary

		// Dynamic tags
		ext := filepath.Ext(nodes[i].FilePath)
		tags := []string{strings.TrimPrefix(ext, ".")}
		if strings.Contains(nodes[i].FilePath, "/views/") {
			tags = append(tags, "view", "markup")
		} else if strings.Contains(nodes[i].FilePath, "/config/") {
			tags = append(tags, "config")
		} else if strings.Contains(nodes[i].FilePath, "/Controllers/") {
			tags = append(tags, "controller", "api")
		} else if strings.Contains(nodes[i].FilePath, "/Repositories/") {
			tags = append(tags, "repository", "database")
		} else if strings.Contains(nodes[i].FilePath, "/Services/") {
			tags = append(tags, "service")
		}
		nodes[i].Tags = tags
	}
}

func cleanComment(comment string) string {
	lines := strings.Split(comment, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "@") {
			cleaned = append(cleaned, line)
		}
	}
	res := strings.Join(cleaned, " ")
	if len(res) > 200 {
		res = res[:197] + "..."
	}
	return res
}

func startDashboard() error {
	logInfo("🚀 Launching Understand Dashboard on http://localhost:5173 ...")
	cmd := exec.Command("pnpm", "--filter", "@understand-anything/dashboard", "dev")
	cmd.Dir = pluginRoot
	cmd.Env = append(os.Environ(), "GRAPH_DIR="+projectRoot, "NODE_OPTIONS=--no-deprecation")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateOnboarding() error {
	logInfo("📑 Generating Onboarding Guide...")
	graphPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/knowledge-graph.json")
	graphData, err := os.ReadFile(graphPath)
	if err != nil {
		return fmt.Errorf("knowledge graph not found: %w", err)
	}

	var graph KnowledgeGraph
	if err := json.Unmarshal(graphData, &graph); err != nil {
		return err
	}

	var md strings.Builder
	md.WriteString("# Onboarding Guide — ")
	md.WriteString(graph.Project.Name)
	md.WriteString("\n\n")
	md.WriteString(graph.Project.Description)
	md.WriteString("\n\n")
	md.WriteString("## 🚀 Tech Stack\n")
	md.WriteString(fmt.Sprintf("- **Languages**: %s\n", strings.Join(graph.Project.Languages, ", ")))
	md.WriteString(fmt.Sprintf("- **Frameworks**: %s\n\n", strings.Join(graph.Project.Frameworks, ", ")))

	md.WriteString("## 🏗️ Architecture Layers\n")
	for _, layer := range graph.Layers {
		md.WriteString(fmt.Sprintf("### %s\n", layer.Name))
		md.WriteString(layer.Description)
		md.WriteString("\n\n")
		md.WriteString("**Key Files**:\n")
		for _, nodeID := range layer.NodeIds {
			// Find node
			for _, node := range graph.Nodes {
				if node.ID == nodeID {
					md.WriteString(fmt.Sprintf("- [%s](file:///%s) — %s\n", node.Name, filepath.Join(projectRoot, node.FilePath), node.Summary))
				}
			}
		}
		md.WriteString("\n")
	}

	md.WriteString("## 🎯 Guided Walkthrough Tours\n")
	if len(graph.Tour) == 0 {
		md.WriteString("*No guided tours generated yet.*\n\n")
	} else {
		for _, tour := range graph.Tour {
			md.WriteString(fmt.Sprintf("### Step %d: %s\n", tour.Order, tour.Title))
			md.WriteString(tour.Description)
			md.WriteString("\n\n")
		}
	}

	md.WriteString("## ⚠️ Complexity Hotspots\n")
	md.WriteString("The following files are identified as highly complex and should be approached with care:\n\n")
	count := 0
	for _, node := range graph.Nodes {
		if node.Type == "file" && node.Complexity == "complex" {
			md.WriteString(fmt.Sprintf("- **[%s](file:///%s)** — %s\n", node.Name, filepath.Join(projectRoot, node.FilePath), node.Summary))
			count++
		}
	}
	if count == 0 {
		md.WriteString("*No high-complexity files found.*\n\n")
	}

	docsDir := filepath.Join(projectRoot, "docs")
	os.MkdirAll(docsDir, 0755)
	outPath := filepath.Join(docsDir, "ONBOARDING.md")
	if err := os.WriteFile(outPath, []byte(md.String()), 0644); err != nil {
		return err
	}

	logInfo("✅ Onboarding guide generated successfully at: " + outPath)
	return nil
}

func showDiff() error {
	fmt.Println("🔍 Analyzing uncommitted changes blast radius...")
	diffCmd := exec.Command("git", "status", "--porcelain")
	diffCmd.Dir = gitRoot
	out, err := diffCmd.Output()
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		fmt.Println("✅ No uncommitted changes.")
		return nil
	}

	var changedFiles []string
	for _, l := range lines {
		if len(l) > 3 {
			changedFiles = append(changedFiles, strings.TrimSpace(l[3:]))
		}
	}

	fmt.Println("Modified Files:")
	for _, f := range changedFiles {
		fmt.Println("  -", f)
	}

	graphPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/knowledge-graph.json")
	graphData, err := os.ReadFile(graphPath)
	if err != nil {
		return nil // Non-fatal if graph doesn't exist
	}
	var graph KnowledgeGraph
	json.Unmarshal(graphData, &graph)

	nodeIdsMap := make(map[string]string)
	for _, f := range changedFiles {
		nodeIdsMap[`file:`+f] = f
	}

	fmt.Println("\n💥 Blast Radius / Ripple Effects:")
	affected := false
	for _, edge := range graph.Edges {
		if srcVal, ok := nodeIdsMap[edge.Source]; ok {
			fmt.Printf("  - %s ──(%s)──> %s\n", srcVal, edge.Type, edge.Target)
			affected = true
		}
		if tgtVal, ok := nodeIdsMap[edge.Target]; ok {
			fmt.Printf("  - %s ──(%s)──> %s\n", edge.Source, edge.Type, tgtVal)
			affected = true
		}
	}
	if !affected {
		fmt.Println("  * No directly connected relationships in the graph. Minimal blast radius.")
	}
	return nil
}

func showExplain(filePath string) {
	graphPath := filepath.Join(dataRoot, "knowledge-graph/understand-anything/knowledge-graph.json")
	graphData, err := os.ReadFile(graphPath)
	if err != nil {
		fmt.Println("❌ Knowledge graph not found. Run --start first.")
		return
	}
	var graph KnowledgeGraph
	json.Unmarshal(graphData, &graph)

	found := false
	for _, node := range graph.Nodes {
		if node.FilePath == filePath {
			fmt.Printf("\n🎯 File Details: %s\n", node.FilePath)
			fmt.Printf("  - Summary: %s\n", node.Summary)
			fmt.Printf("  - Tags: %s\n", strings.Join(node.Tags, ", "))
			fmt.Printf("  - Complexity: %s\n", node.Complexity)
			found = true
		}
	}
	if !found {
		fmt.Printf("❌ File '%s' not found in knowledge graph.\n", filePath)
	}
}

func printHelp() {
	fmt.Println(`Understand Anything — CLI Helper

Usage:
  understand <flag> [args]

Flags:
  --start        Run the full scanning and AST synthesis pipeline
  --update       Run incremental updates on changed files since last scan
  --dashboard    Launch the interactive visual React dashboard
  --onboard      Generate docs/ONBOARDING.md onboarding guide
  --diff         Analyze uncommitted changes blast radius
  --explain <f>  Deep-dive explain a specific file's purpose
  --help         Show this help information`)
}

// findProjectRoot walks up from the binary location looking for project root markers.
// Works regardless of whether the binary lives under .agents/ or .claude/.
func findProjectRoot(binaryPath string) string {
	dir := filepath.Dir(binaryPath)
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = parent
	}
	return dir
}

func main() {
	// Initialize paths dynamically: env override → walk up from binary → CWD fallback
	// PROJECT_ROOT = α/ dir (knowledge-graph/ storage)
	// ALPHA_ROOT   = project root (git, code scanning)
	if envRoot := os.Getenv("PROJECT_ROOT"); envRoot != "" {
		projectRoot, _ = filepath.Abs(envRoot)
	} else if exe, err := os.Executable(); err == nil {
		exePath, _ := filepath.EvalSymlinks(exe)
		absExe, _ := filepath.Abs(exePath)
		projectRoot = findProjectRoot(absExe)
	}
	if projectRoot == "" || projectRoot == "/" {
		wd, _ := os.Getwd()
		projectRoot, _ = filepath.Abs(wd)
	}

	// gitRoot: ALPHA_ROOT if set, else parent of projectRoot (α/ → project)
	if alphaRoot := os.Getenv("ALPHA_ROOT"); alphaRoot != "" {
		gitRoot, _ = filepath.Abs(alphaRoot)
	} else {
		gitRoot = filepath.Dir(projectRoot)
		if gitRoot == "" || gitRoot == "." {
			gitRoot = projectRoot
		}
	}

	scriptFolder = filepath.Join(projectRoot, "scripts")

	// dataRoot: per-project data directory.
	// Global mode (ALPHA_GLOBAL=1): ~/.alpha-ai/knowledge-graph/projects/<id>/
	// Local mode: projectRoot/knowledge-graph/understand-anything/../  (i.e. projectRoot)
	if os.Getenv("ALPHA_GLOBAL") == "1" && gitRoot != "" {
		id := understandProjectID(gitRoot)
		dataRoot = filepath.Join(projectRoot, "knowledge-graph", "projects", id)
	} else {
		dataRoot = projectRoot
	}
	os.MkdirAll(filepath.Join(dataRoot, "knowledge-graph", "understand-anything"), 0755)

	home, err := os.UserHomeDir()
	if err == nil {
		pluginRoot = filepath.Join(home, ".understand-anything/repo/understand-anything-plugin")
	} else {
		// Fallback
		pluginRoot = "/Users/neo/.understand-anything/repo/understand-anything-plugin"
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "--start":
			if err := startPipeline(); err != nil {
				logInfo("❌ Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--update":
			if err := updatePipeline(); err != nil {
				logInfo("❌ Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--dashboard":
			if err := startDashboard(); err != nil {
				logInfo("❌ Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--onboard":
			if err := generateOnboarding(); err != nil {
				logInfo("❌ Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--diff":
			if err := showDiff(); err != nil {
				logInfo("❌ Error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "--explain":
			if len(os.Args) < 3 {
				logInfo("Usage: understand --explain <filepath>")
				os.Exit(2)
			}
			showExplain(os.Args[2])
			os.Exit(0)
		case "--help":
			printHelp()
			os.Exit(0)
		default:
			if isInputTTY() {
				fmt.Fprintf(os.Stderr, "Error: Unrecognized command or flag '%s'\n\n", arg)
				printHelp()
				os.Exit(1)
			}
		}
	}

	// Normal MCP Server Mode
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	s := server.NewMCPServer("Understand Anything Helper", "1.0.0")
	slog.Info("Starting MCP Server mode")

	// Tool: awake
	s.AddTool(mcp.Tool{
		Name:        "awake",
		Description: "Run the understand session start check.",
		InputSchema: mcp.ToolInputSchema{Type: "object"},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("Understand MCP initialized successfully."), nil
	})

	// Tool: start
	s.AddTool(mcp.Tool{
		Name:        "start",
		Description: "Start full project analysis.",
		InputSchema: mcp.ToolInputSchema{Type: "object"},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := startPipeline(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Project analysis completed successfully."), nil
	})

	// Tool: onboard
	s.AddTool(mcp.Tool{
		Name:        "onboard",
		Description: "Generate onboarding documentation.",
		InputSchema: mcp.ToolInputSchema{Type: "object"},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if err := generateOnboarding(); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText("Onboarding guide generated successfully."), nil
	})

	// Run MCP
	if err := server.ServeStdio(s); err != nil {
		slog.Error("MCP server failed", "error", err)
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
