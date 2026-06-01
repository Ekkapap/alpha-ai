package main

// cmd_forget.go — forget command: delete session memories (CLI).
// MCP forget is handled by alpha/main.go (calls graphify CLI via exec).

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func cliForget() {
	pattern := ""
	autoConfirm := false

	for i := 2; i < len(os.Args); i++ {
		v := os.Args[i]
		if v == "-y" || v == "--yes" {
			autoConfirm = true
		} else if pattern == "" {
			pattern = v
		}
	}

	memDir := filepath.Join(root, "knowledge-graph/memories")
	var targets []string

	if pattern == "" {
		matches, _ := filepath.Glob(filepath.Join(memDir, "session_summary_*.md"))
		if len(matches) > 0 {
			sort.Strings(matches)
			targets = append(targets, matches[len(matches)-1])
		}
	} else {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			targets = matches
		} else {
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

		parts := strings.Split(filepath.Base(f), "_")
		if len(parts) >= 4 {
			ts := strings.TrimSuffix(parts[2]+"_"+parts[3], ".md")
			gMatches, _ := filepath.Glob(filepath.Join(root, "knowledge-graph/graphify-out/memory/*"+ts+"*"))
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
}
