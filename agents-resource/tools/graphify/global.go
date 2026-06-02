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
		projectRoot := os.Getenv("ALPHA_ROOT")
		if projectRoot != "" {
			return filepath.Join(alphaDir, "knowledge-graph", "projects", projectID(projectRoot), "graphify-out")
		}
	}
	return filepath.Join(alphaDir, "knowledge-graph", "graphify-out")
}

// memoriesDir returns the memories directory path.
//   - Global mode: ALPHA_ROOT/α/knowledge-graph/memories (local to project)
//                  fallback: ALPHA_ROOT/.alpha/memories
//   - Local mode:  alphaDir/knowledge-graph/memories (unchanged)
func memoriesDir(alphaDir string) string {
	if isGlobalMode() {
		projectRoot := os.Getenv("ALPHA_ROOT")
		if projectRoot != "" {
			// Prefer α/ subdir if it exists
			alphaSub := filepath.Join(projectRoot, "α", "knowledge-graph", "memories")
			if _, err := os.Stat(filepath.Join(projectRoot, "α")); err == nil {
				return alphaSub
			}
			return filepath.Join(projectRoot, ".alpha", "memories")
		}
	}
	return filepath.Join(alphaDir, "knowledge-graph", "memories")
}
