package main

// cmd_build.go — build: compile graphify binary + auto version check on startup.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// binaryVersion is set at compile time via -ldflags "-X main.binaryVersion=X.Y.Z".
// Default "1.0.0" is used for local dev builds. Bump "graphify_version" in agents-resource/config.json
// to trigger auto-rebuild on next run.
var binaryVersion = "1.0.0"

// graphifySrcDir walks up from the binary to find the Go source directory (contains go.mod).
// Shared by cliBuild and checkAndAutoRebuild.
func graphifySrcDir() string {
	exe, _ := os.Executable()
	exeAbs, _ := filepath.EvalSymlinks(exe)
	dir := filepath.Dir(exeAbs)
	for d := dir; d != filepath.Dir(d); d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d
		}
	}
	return filepath.Join(root, "agents-resource/tools/graphify")
}

// checkAndAutoRebuild reads agents-resource/config.json and rebuilds + re-execs if
// the config's graphify_version differs from the binaryVersion baked into this binary.
// Only runs in CLI mode (not MCP), and skipped inside Docker.
func checkAndAutoRebuild() {
	if os.Getenv("ALPHA_IN_DOCKER") == "1" {
		return
	}

	configPath := filepath.Join(root, "agents-resource/config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	var cfg struct {
		GraphifyVersion string `json:"graphify_version"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.GraphifyVersion == "" {
		return
	}
	if cfg.GraphifyVersion == binaryVersion {
		return
	}

	fmt.Fprintf(os.Stderr, "[graphify] outdated (binary=%s, latest=%s) — rebuilding...\n",
		binaryVersion, cfg.GraphifyVersion)

	srcDir := graphifySrcDir()

	exe, _ := os.Executable()
	exeAbs, _ := filepath.EvalSymlinks(exe)
	tmp := exeAbs + ".new"

	ldflags := fmt.Sprintf("-X main.binaryVersion=%s", cfg.GraphifyVersion)
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", tmp, ".")
	cmd.Dir = srcDir
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp)
		fmt.Fprintf(os.Stderr, "[graphify] rebuild failed: %v\n", err)
		return
	}

	if err := os.Rename(tmp, exeAbs); err != nil {
		os.Remove(tmp)
		fmt.Fprintf(os.Stderr, "[graphify] replace binary: %v\n", err)
		return
	}
	os.Chmod(exeAbs, 0755)

	fmt.Fprintf(os.Stderr, "[graphify] updated to %s — restarting\n", cfg.GraphifyVersion)
	syscall.Exec(exeAbs, os.Args, os.Environ())
}

func cliBuild() {
	platform := runtime.GOOS
	if p := os.Getenv("GOOS"); p != "" {
		platform = p
	}

	srcDir := graphifySrcDir()
	dest := filepath.Join(root, "agents-resource/tools/bin/"+platform+"/graphify")

	tmp := filepath.Join(os.TempDir(), "graphify-build-"+platform)
	fmt.Printf("Building from %s ...\n", srcDir)
	cmd := exec.Command("go", "build", "-o", tmp, ".")
	cmd.Dir = srcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmp)

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", dest, err)
		os.Exit(1)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read tmp binary: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(dest, data, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", dest, err)
		os.Exit(1)
	}
	fmt.Printf("✅ %s (v%s)\n", dest, binaryVersion)
}
