package main

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlphaNumU = regexp.MustCompile(`[^a-z0-9]+`)

// understandProjectID returns a stable human-readable ID for a project path.
// Format: <sanitized-basename>-<8char-sha256>
func understandProjectID(absPath string) string {
	name := strings.ToLower(filepath.Base(absPath))
	name = nonAlphaNumU.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if len(name) > 30 {
		name = name[:30]
	}
	h := sha256.Sum256([]byte(absPath))
	return fmt.Sprintf("%s-%x", name, h[:4])
}
