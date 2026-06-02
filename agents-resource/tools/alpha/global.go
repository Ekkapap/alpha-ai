package main

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlphaNumAlpha = regexp.MustCompile(`[^a-z0-9]+`)

// alphaProjectID returns a stable human-readable ID for a project path.
// Format: <sanitized-basename>-<8char-sha256>
// Example: "/Users/neo/work/my-project" → "my-project-a1b2c3d4"
func alphaProjectID(absPath string) string {
	name := strings.ToLower(filepath.Base(absPath))
	name = nonAlphaNumAlpha.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if len(name) > 30 {
		name = name[:30]
	}
	h := sha256.Sum256([]byte(absPath))
	return fmt.Sprintf("%s-%x", name, h[:4])
}
