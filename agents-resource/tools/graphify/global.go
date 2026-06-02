package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

// projectID returns a stable, human-readable project identifier from an absolute path.
// Format: <sanitized-basename>-<8char-sha256-of-path>
// Example: "/Users/neo/work/my-project" → "my-project-a1b2c3d4"
func projectID(absPath string) string {
	name := strings.ToLower(filepath.Base(absPath))
	name = nonAlphaNum.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if len(name) > 30 {
		name = name[:30]
	}
	h := sha256.Sum256([]byte(absPath))
	return fmt.Sprintf("%s-%x", name, h[:4])
}

// isGlobalMode reports whether we are running in global mode (ALPHA_GLOBAL=1).
func isGlobalMode() bool {
	return os.Getenv("ALPHA_GLOBAL") == "1"
}

// graphifyDataDir returns the graphify-out directory for the current run.
//   - Global mode: alphaDir/knowledge-graph/projects/<project-id>/graphify-out
//   - Local mode:  alphaDir/knowledge-graph/graphify-out  (unchanged)
func graphifyDataDir(alphaDir string) string {
	if isGlobalMode() {
		id := os.Getenv("ALPHA_PROJECT_ID")
		if id == "" {
			id = projectID(os.Getenv("ALPHA_ROOT"))
		}
		if id != "" && id != "-" {
			return filepath.Join(alphaDir, "knowledge-graph", "projects", id, "graphify-out")
		}
	}
	return filepath.Join(alphaDir, "knowledge-graph", "graphify-out")
}

// memoriesDir returns the memories directory path.
//   - Global mode: ALPHA_ROOT/memories  (local per-project directory)
//   - Local mode:  alphaDir/knowledge-graph/memories
func memoriesDir(alphaDir string) string {
	if isGlobalMode() {
		// In Docker: ALPHA_ROOT=/workspace (project mount) — memories are local to project
		projectRoot := os.Getenv("ALPHA_ROOT")
		if projectRoot != "" {
			return filepath.Join(projectRoot, "memories")
		}
	}
	return filepath.Join(alphaDir, "knowledge-graph", "memories")
}
