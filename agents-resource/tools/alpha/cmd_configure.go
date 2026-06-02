package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// configurePaths writes .mcp.json and, in local mode, creates project-root symlinks.
// Global mode (ALPHA_GLOBAL=1): writes global .mcp.json with ALPHA_HOME/ALPHA_PROJECT_ID.
// Local mode: writes local .mcp.json pointing to α/ + creates graphify-out/.understand-anything symlinks.
// Safe to re-run.
func configurePaths(alphaDir, projectRoot string) string {
	if os.Getenv("ALPHA_GLOBAL") == "1" {
		return configureGlobal(alphaDir, projectRoot)
	}
	return configureLocal(alphaDir, projectRoot)
}

// configureLocal — local mode (α/ inside project)
func configureLocal(alphaDir, projectRoot string) string {
	var out strings.Builder

	// HOST-side paths — inside Docker, projectRoot=/workspace but the host path
	// is in HOST_PROJECT_ROOT env or α/.env. Natively, projectRoot is the host path.
	hostProjectRoot := os.Getenv("HOST_PROJECT_ROOT")
	if hostProjectRoot == "" {
		hostProjectRoot = readEnvFile(alphaDir, "HOST_PROJECT_ROOT")
	}
	if hostProjectRoot == "" {
		hostProjectRoot = projectRoot
	}
	hostAlphaDir := filepath.Join(hostProjectRoot, "α")

	// ── .mcp.json ──────────────────────────────────────────────────────────
	templatePath := filepath.Join(alphaDir, "agents-resource/.mcp.json")
	dstPath := filepath.Join(projectRoot, ".mcp.json")

	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		out.WriteString("⚠️  .mcp.json template not found: " + err.Error() + "\n")
	} else {
		raw := strings.ReplaceAll(string(tmpl), "[ALPHA_DIR]", hostAlphaDir)
		raw = strings.ReplaceAll(raw, "${env:PWD}", hostProjectRoot)

		var parsed map[string]any
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			out.WriteString("⚠️  parse .mcp.json template: " + err.Error() + "\n")
		} else {
			stripNoteFields(parsed)
			clean, _ := json.MarshalIndent(parsed, "", "  ")
			if err := os.WriteFile(dstPath, append(clean, '\n'), 0644); err != nil {
				out.WriteString("⚠️  write .mcp.json: " + err.Error() + "\n")
			} else {
				out.WriteString("✅ .mcp.json written → " + dstPath + "\n")
			}
		}
	}

	// ── Symlinks (relative — work inside Docker too) ───────────────────────
	symlinks := [][2]string{
		{"graphify-out", "α/knowledge-graph/graphify-out"},
		{".understand-anything", "α/knowledge-graph/understand-anything"},
	}
	for _, sl := range symlinks {
		dst := filepath.Join(projectRoot, sl[0])
		target := sl[1]

		if fi, err := os.Lstat(dst); err == nil {
			if fi.Mode()&os.ModeSymlink != 0 {
				if current, _ := os.Readlink(dst); current == target {
					out.WriteString("✅ " + sl[0] + " → already correct\n")
					continue
				}
				os.Remove(dst)
			} else {
				out.WriteString("⚠️  " + sl[0] + " exists but is not a symlink — skipped\n")
				continue
			}
		}

		targetAbs := filepath.Join(projectRoot, target)
		if err := os.MkdirAll(targetAbs, 0755); err != nil {
			out.WriteString("⚠️  mkdir " + target + ": " + err.Error() + "\n")
			continue
		}

		if err := os.Symlink(target, dst); err != nil {
			out.WriteString("⚠️  symlink " + sl[0] + ": " + err.Error() + "\n")
		} else {
			out.WriteString("✅ " + sl[0] + " → " + target + "\n")
		}
	}

	out.WriteString("\n⚡ Restart your AI agent to load the new MCP config.")
	return strings.TrimSpace(out.String())
}

// configureGlobal — global mode (~/.alpha-ai/ shared installation, called when ALPHA_GLOBAL=1)
func configureGlobal(alphaDir, projectRoot string) string {
	hostProjectRoot := os.Getenv("HOST_PROJECT_ROOT")
	if hostProjectRoot == "" {
		hostProjectRoot = readEnvFile(alphaDir, "HOST_PROJECT_ROOT")
	}
	if hostProjectRoot == "" {
		hostProjectRoot = projectRoot
	}

	alphaHome := os.Getenv("ALPHA_HOME")
	if alphaHome == "" {
		alphaHome = readEnvFile(alphaDir, "ALPHA_HOME")
	}
	if alphaHome == "" {
		home, _ := os.UserHomeDir()
		alphaHome = filepath.Join(home, ".alpha-ai")
	}

	projectID := os.Getenv("ALPHA_PROJECT_ID")
	if projectID == "" {
		projectID = alphaProjectID(hostProjectRoot)
	}

	return projectInit(alphaHome, hostProjectRoot, projectID)
}

// projectInit initialises a project directory to use a global alpha-ai installation.
// alphaHome: path to ~/.alpha-ai, projectRoot: target project directory, projectID: stable ID.
// Safe to re-run.
func projectInit(alphaHome, projectRoot, projectID string) string {
	var out strings.Builder

	// ── α/config.json ─────────────────────────────────────────────────────
	alphaConfigDir := filepath.Join(projectRoot, "α")
	if err := os.MkdirAll(alphaConfigDir, 0755); err == nil {
		cfg := map[string]string{"project_id": projectID, "alpha_home": alphaHome}
		data, _ := json.MarshalIndent(cfg, "", "  ")
		if err := os.WriteFile(filepath.Join(alphaConfigDir, "config.json"), append(data, '\n'), 0644); err != nil {
			out.WriteString("⚠️  write α/config.json: " + err.Error() + "\n")
		} else {
			out.WriteString("✅ α/config.json → project_id=" + projectID + "\n")
		}
	}

	// ── Per-project data dirs ──────────────────────────────────────────────
	projectDataDir := filepath.Join(alphaHome, "knowledge-graph", "projects", projectID)
	os.MkdirAll(filepath.Join(projectDataDir, "graphify-out"), 0755)
	os.MkdirAll(filepath.Join(projectDataDir, "understand-anything"), 0755)
	out.WriteString("✅ Global data dirs → " + projectDataDir + "\n")

	// ── .mcp.json (global template) ───────────────────────────────────────
	templatePath := filepath.Join(alphaHome, "agents-resource/.mcp.global.json")
	dstPath := filepath.Join(projectRoot, ".mcp.json")

	tmpl, err := os.ReadFile(templatePath)
	if err != nil {
		out.WriteString("⚠️  .mcp.global.json template not found: " + err.Error() + "\n")
	} else {
		raw := strings.ReplaceAll(string(tmpl), "[ALPHA_HOME]", alphaHome)
		raw = strings.ReplaceAll(raw, "[HOST_PROJECT_ROOT]", projectRoot)
		raw = strings.ReplaceAll(raw, "[ALPHA_PROJECT_ID]", projectID)

		var parsed map[string]any
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			out.WriteString("⚠️  parse .mcp.global.json: " + err.Error() + "\n")
		} else {
			stripNoteFields(parsed)
			clean, _ := json.MarshalIndent(parsed, "", "  ")
			if err := os.WriteFile(dstPath, append(clean, '\n'), 0644); err != nil {
				out.WriteString("⚠️  write .mcp.json: " + err.Error() + "\n")
			} else {
				out.WriteString("✅ .mcp.json (global) written → " + dstPath + "\n")
			}
		}
	}

	out.WriteString("\n⚡ Restart your AI agent to load the new MCP config.")
	return strings.TrimSpace(out.String())
}

// runProjectInit is the entry point for `alpha --project-init` and the project_init MCP tool.
// It auto-detects ~/.alpha-ai/ and uses $PWD as the project root.
func runProjectInit(cwdOverride string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "⚠️  cannot determine home dir: " + err.Error()
	}
	alphaHome := filepath.Join(home, ".alpha-ai")
	if _, err := os.Stat(alphaHome); os.IsNotExist(err) {
		return "⚠️  ~/.alpha-ai not found — run install.sh --global first"
	}

	projectRoot := cwdOverride
	if projectRoot == "" {
		projectRoot, _ = os.Getwd()
	}
	projectID := alphaProjectID(projectRoot)
	return projectInit(alphaHome, projectRoot, projectID)
}

// stripNoteFields recursively removes "_note" keys from parsed JSON maps.
func stripNoteFields(v any) {
	switch m := v.(type) {
	case map[string]any:
		delete(m, "_note")
		for _, val := range m {
			stripNoteFields(val)
		}
	case []any:
		for _, item := range m {
			stripNoteFields(item)
		}
	}
}

// readEnvFile reads a key=value from alphaDir/.env (handles comments and quotes).
func readEnvFile(alphaDir, key string) string {
	f, err := os.Open(filepath.Join(alphaDir, ".env"))
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	prefix := key + "="
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, prefix) {
			val := strings.TrimPrefix(line, prefix)
			val = strings.Trim(val, `"'`)
			return val
		}
	}
	return ""
}
