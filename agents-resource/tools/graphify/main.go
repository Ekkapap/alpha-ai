// Package graphify — MCP server and CLI tool for knowledge graph operations.
//
// File structure (split by command — add new commands as new files, reuse shared code):
//
//	graph.go        — graph types + operations (reused by all commands)
//	display.go      — display helpers: alphaDisplay, printHelp, findProjectRoot (reused)
//	cmd_awake.go    — awake: restore session context (CLI + MCP)
//	cmd_sync.go     — sync: save session summary + update graph (CLI + MCP)
//	cmd_query.go    — overview/sketch/detail/update: Phase 0/1/2 graph queries (CLI + MCP)
//	cmd_forget.go   — forget: delete session memories (CLI)
//	cmd_focus.go    — focus: search file sections + debug_info (CLI + MCP)
//	cmd_alpha.go    — alpha: system identity display (CLI + MCP)
//	cmd_build.go    — build: compile binary to α/bin/<platform>/ (CLI)
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
)

// root is the α/ directory (PROJECT_ROOT). Set in main() and shared across all command files.
var root string

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// Root detection: env override → walk up from binary → CWD fallback
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

	// CLI mode — dispatch to command handler
	if len(os.Args) > 1 {
		checkAndAutoRebuild()
		switch os.Args[1] {
		case "alpha":
			cliAlpha()
		case "awake":
			cliAwake()
		case "sync":
			cliSync()
		case "forget":
			cliForget()
		case "update":
			cliUpdate()
		case "overview":
			cliOverview()
		case "sketch":
			cliSketch()
		case "detail":
			cliDetail()
		case "focus":
			cliFocus()
		case "build":
			cliBuild()
		default:
			if isInputTTY() {
				fmt.Fprintf(os.Stderr, "Error: Unrecognized command or flag '%s'\n\n", os.Args[1])
				printHelp()
				os.Exit(1)
			}
			// Non-TTY unrecognized arg: fall through to MCP server mode
			goto mcpMode
		}
		os.Exit(0)
	}

mcpMode:
	slog.Info("Starting server", "root", root)
	projectName := filepath.Base(root)
	s := server.NewMCPServer(
		fmt.Sprintf("Graphify Helper (%s)", projectName),
		"1.1.1",
	)

	registerMCPDebugInfo(s)
	registerMCPFocus(s)
	registerMCPAlpha(s)
	registerMCPAwake(s)
	registerMCPSync(s)
	registerMCPUpdateSessionSummary(s)
	registerMCPSketch(s)
	registerMCPDetail(s)
	registerMCPOverview(s)
	registerMCPUpdate(s)
	registerMCPSetLabels(s)

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
