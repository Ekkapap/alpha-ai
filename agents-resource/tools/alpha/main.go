// Package alpha — unified MCP entry point for all knowledge graph operations.
//
// File structure (split by command — add new commands as new files, reuse shared code):
//
//	main.go                   — findRoots, binPath, runTool, MCP tools, server
//	cmd_knowledge_graph.go    — alpha --knowledge-graph: docker + graph management (CLI)
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Alpha MCP Server — unified entry point for all knowledge graph operations.
//
// Two-root design:
//   alphaDir  (PROJECT_ROOT) = α/ directory — where knowledge-graph/, memories/ live
//   projectRoot (ALPHA_ROOT) = parent of α/ — where Python graphify CLI analyzes code
//
// In Docker: PROJECT_ROOT=/workspace/α, ALPHA_ROOT=/workspace (set in docker-compose.yml)
// Native:    alpha binary detects both automatically

func inDocker() bool { return os.Getenv("ALPHA_IN_DOCKER") == "1" }

// findRoots returns (alphaDir, projectRoot).
//
//   alphaDir    = the α/ (or ~/.alpha-ai/) directory — where knowledge-graph/ lives
//   projectRoot = the project being analyzed (parent of α/ in local mode;
//                 HOST_PROJECT_ROOT in global mode)
//
// Global mode is active when ALPHA_GLOBAL=1.
func findRoots() (string, string) {
	// Global mode: alpha dir is ~/.alpha-ai/, project root from HOST_PROJECT_ROOT
	if os.Getenv("ALPHA_GLOBAL") == "1" {
		alphaDir := os.Getenv("PROJECT_ROOT")
		if alphaDir == "" {
			home, _ := os.UserHomeDir()
			alphaDir = filepath.Join(home, ".alpha-ai")
		}
		projectRoot := os.Getenv("ALPHA_ROOT")
		if projectRoot == "" {
			projectRoot = os.Getenv("HOST_PROJECT_ROOT")
		}
		if projectRoot == "" {
			projectRoot = filepath.Dir(alphaDir)
		}
		return alphaDir, projectRoot
	}

	// Local mode: explicit env vars (Docker or manual override)
	alphaDir := os.Getenv("PROJECT_ROOT")
	projectRoot := os.Getenv("ALPHA_ROOT")

	if alphaDir != "" && projectRoot != "" {
		return alphaDir, projectRoot
	}

	// Walk up from executable
	dir, _ := os.Executable()
	dir = filepath.Dir(dir)
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		// We are inside a project that has α/
		if _, err := os.Stat(filepath.Join(dir, "α")); err == nil {
			a := filepath.Join(dir, "α")
			return a, dir
		}
		// We are inside α/ itself (knowledge-graph/ is here)
		if _, err := os.Stat(filepath.Join(dir, "knowledge-graph")); err == nil {
			return dir, filepath.Dir(dir)
		}
		dir = parent
	}

	// CWD fallback
	cwd, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(cwd, "α")); err == nil {
		return filepath.Join(cwd, "α"), cwd
	}
	return cwd, filepath.Dir(cwd)
}

// alphaProjectDataDir returns the per-project data directory for graphify/understand.
// Global mode: alphaDir/knowledge-graph/projects/<id>/
// Local mode:  alphaDir/knowledge-graph/
func alphaProjectDataDir(alphaDir, projectRoot string) string {
	if os.Getenv("ALPHA_GLOBAL") == "1" && projectRoot != "" {
		id := alphaProjectID(projectRoot)
		return filepath.Join(alphaDir, "knowledge-graph", "projects", id)
	}
	return filepath.Join(alphaDir, "knowledge-graph")
}

// binPath returns the correct binary path.
// In Docker the binaries are on PATH (/usr/local/bin/); natively look in tools/bin/.
// alphaDir = the α/ (or ~/.alpha-ai/) directory.
func binPath(alphaDir, name string) string {
	if inDocker() {
		return name // on PATH inside container
	}
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = filepath.Join(alphaDir, "agents-resource", "tools", "bin", "windows")
		return filepath.Join(dir, name+".exe")
	case "darwin":
		dir = filepath.Join(alphaDir, "agents-resource", "tools", "bin", "darwin")
	default:
		dir = filepath.Join(alphaDir, "agents-resource", "tools", "bin", "linux")
	}
	return filepath.Join(dir, name)
}

func runTool(alphaDir, projectRoot, bin string, args ...string) (*mcp.CallToolResult, error) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = alphaDir
	env := append(os.Environ(),
		"PROJECT_ROOT="+alphaDir,
		"ALPHA_ROOT="+projectRoot,
	)
	// Propagate global mode env vars to subprocesses
	if os.Getenv("ALPHA_GLOBAL") == "1" {
		env = append(env, "ALPHA_GLOBAL=1")
	}
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("%s error: %s\n%s", filepath.Base(bin), err, out)), nil
	}
	return mcp.NewToolResultText(strings.TrimSpace(string(out))), nil
}

func strParam(req mcp.CallToolRequest, key string) string {
	args, _ := req.Params.Arguments.(map[string]any)
	if args == nil {
		return ""
	}
	v, _ := args[key].(string)
	return v
}

func main() {
	alphaDir, projectRoot := findRoots()

	gfyBin := binPath(alphaDir, "graphify-core") // graphify-core inside Docker, Go binary natively
	uaBin := binPath(alphaDir, "understand")

	gfy := func(args ...string) (*mcp.CallToolResult, error) {
		return runTool(alphaDir, projectRoot, gfyBin, args...)
	}
	ua := func(args ...string) (*mcp.CallToolResult, error) {
		return runTool(alphaDir, projectRoot, uaBin, args...)
	}

	s := server.NewMCPServer("ALPHA", "1.0.0", server.WithToolCapabilities(false))

	// ── graphify tools ──────────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "awake",
		Description: "Restore session context from memories and knowledge graph state.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return gfy("awake")
	})

	s.AddTool(mcp.Tool{
		Name:        "overview",
		Description: "Return architecture overview: pillar nodes, God Nodes, stats, and top communities.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return gfy("overview")
	})

	s.AddTool(mcp.Tool{
		Name:        "sketch",
		Description: "BFS traversal from a query term to find candidate node IDs.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{"type": "string", "description": "Search term"},
				"depth": map[string]any{"type": "number", "description": "BFS depth (default 2)"},
			},
			Required: []string{"query"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"sketch", strParam(req, "query")}
		if d := strParam(req, "depth"); d != "" {
			args = append(args, "--depth", d)
		}
		return gfy(args...)
	})

	s.AddTool(mcp.Tool{
		Name:        "detail",
		Description: "Return full context for specific node IDs: callers, callees, cluster, dependencies.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"ids": map[string]any{"type": "string", "description": "Comma-separated node IDs"},
			},
			Required: []string{"ids"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return gfy("detail", strParam(req, "ids"))
	})

	s.AddTool(mcp.Tool{
		Name:        "sync",
		Description: "Incremental graph sync with an optional session summary.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"summary": map[string]any{"type": "string", "description": "Session summary to store"},
			},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := []string{"sync"}
		if s := strParam(req, "summary"); s != "" {
			args = append(args, "-s", s)
		}
		return gfy(args...)
	})

	s.AddTool(mcp.Tool{
		Name:        "focus",
		Description: "Extract a focused section of a file around a search term.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{"type": "string", "description": "File path"},
				"term": map[string]any{"type": "string", "description": "Search term"},
			},
			Required: []string{"path", "term"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return gfy("focus", strParam(req, "path"), strParam(req, "term"))
	})

	s.AddTool(mcp.Tool{
		Name:        "update_session_summary",
		Description: "Write agent-merged session-summary.md to disk. Call after merging sync output with existing session-summary.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"content": map[string]any{"type": "string", "description": "Full merged content to write to session-summary.md"},
			},
			Required: []string{"content"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		content := strParam(req, "content")
		summaryPath := filepath.Join(alphaDir, "knowledge-graph/memories/session-summary.md")
		if err := os.WriteFile(summaryPath, []byte(content), 0644); err != nil {
			return mcp.NewToolResultError("write session-summary.md: " + err.Error()), nil
		}
		return mcp.NewToolResultText("session-summary.md updated."), nil
	})

	s.AddTool(mcp.Tool{
		Name:        "forget",
		Description: "Remove a memory entry by timestamp or keyword.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"target": map[string]any{"type": "string", "description": "Timestamp or keyword to forget"},
			},
			Required: []string{"target"},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return gfy("forget", strParam(req, "target"))
	})

	s.AddTool(mcp.Tool{
		Name:        "update",
		Description: "Incremental graph update: runs Python graphify + understand update without opening the browser.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(strings.TrimSpace(runGraphifyUpdate(alphaDir, projectRoot, uaBin))), nil
	})

	// ── understand tools ────────────────────────────────────────────────────

	s.AddTool(mcp.Tool{
		Name:        "understand",
		Description: "Deep AST analysis. mode: start (full), update (incremental), diff (blast radius).",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"mode": map[string]any{
					"type": "string",
					"enum": []string{"start", "update", "diff"},
				},
			},
		},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		mode := strParam(req, "mode")
		if mode == "" {
			mode = "update"
		}
		return ua("--" + mode)
	})

	s.AddTool(mcp.Tool{
		Name:        "diff",
		Description: "Estimate blast radius of uncommitted changes using the understand graph.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return ua("--diff")
	})

	s.AddTool(mcp.Tool{
		Name:        "configure",
		Description: "Write .mcp.json and create project-root symlinks (graphify-out, .understand-anything) without re-running install.sh. Safe to re-run.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(configurePaths(alphaDir, projectRoot)), nil
	})

	s.AddTool(mcp.Tool{
		Name:        "project_init",
		Description: "Initialise the current directory as a project using the global ~/.alpha-ai installation. Creates α/config.json, .mcp.json, and per-project data dirs. Run from the project directory. Safe to re-run.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(runProjectInit("")), nil
	})

	// ── CLI passthrough ─────────────────────────────────────────────────────
	if len(os.Args) > 1 {
		arg := os.Args[1]

		if arg == "--knowledge-graph" {
			handleKnowledgeGraph(alphaDir, projectRoot, uaBin)
			return
		}
		if arg == "--update" {
			fmt.Print(runGraphifyUpdate(alphaDir, projectRoot, uaBin))
			return
		}
		if arg == "--configure" {
			fmt.Println(configurePaths(alphaDir, projectRoot))
			return
		}
		if arg == "--project-init" {
			fmt.Println(runProjectInit(""))
			return
		}

		var bin string
		var cliArgs []string
		if arg == "--understand" || arg == "--diff" {
			bin = uaBin
			cliArgs = os.Args[2:]
			if arg == "--diff" {
				cliArgs = append([]string{"--diff"}, cliArgs...)
			}
		} else {
			bin = gfyBin
			// Strip "--" prefix: "alpha --awake path" → "graphify awake path"
			cliArgs = append([]string{strings.TrimPrefix(arg, "--")}, os.Args[2:]...)
		}
		cmd := exec.Command(bin, cliArgs...)
		cmd.Dir = alphaDir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(),
			"PROJECT_ROOT="+alphaDir,
			"ALPHA_ROOT="+projectRoot,
		)
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
		return
	}

	// ── MCP server mode ─────────────────────────────────────────────────────
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		os.Exit(0)
	}()

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "alpha: %v\n", err)
		os.Exit(1)
	}
}
