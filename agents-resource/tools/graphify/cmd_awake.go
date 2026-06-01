package main

// cmd_awake.go — awake command: restore session context (CLI + MCP).

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cliAwake() {
	graphPath := filepath.Join(root, "knowledge-graph/graphify-out/graph.json")
	if _, err := os.Stat(graphPath); os.IsNotExist(err) {
		projectRoot := filepath.Dir(root)
		fmt.Printf("No knowledge graph found.\n\nWould you like to initialize it now?\n  1. Yes — scan from %s\n  2. Yes, specify path\n  3. No — skip\n\nChoice [1/2/3]: ", projectRoot)
		var choice string
		fmt.Scanln(&choice)
		graphifyBin := "graphify"
		if os.Getenv("ALPHA_IN_DOCKER") == "1" {
			graphifyBin = "graphify-core"
		}
		switch choice {
		case "1":
			fmt.Println("Running scan...")
			cmd := exec.Command(graphifyBin, "update", projectRoot)
			cmd.Dir = projectRoot
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		case "2":
			fmt.Print("Enter path to scan: ")
			var customPath string
			fmt.Scanln(&customPath)
			if customPath != "" {
				cmd := exec.Command(graphifyBin, "update", customPath)
				cmd.Dir = customPath
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		default:
			fmt.Println("Skipped. Run 'alpha --update' to initialize later.")
			os.Exit(0)
		}
	}

	// Optional focused path: "alpha --awake path/to/area"
	if len(os.Args) > 2 {
		focusPath := os.Args[2]
		g, err := loadFullGraph(root)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Printf("[AGENT_CONTEXT_START]\n\n### FOCUSED CONTEXT: %s\n%s\n\n[AGENT_CONTEXT_END]\n",
			focusPath, sketchGraph(g, focusPath, 3))
		return
	}

	fmt.Print(buildAwakeOverview(root))

	latestStatePath := filepath.Join(root, "knowledge-graph/memories/session-summary.md")
	if state, err := os.ReadFile(latestStatePath); err == nil {
		fmt.Printf("### PREVIOUS SESSION SUMMARY\n%s\n\n", string(state))
	}
}

func registerMCPAwake(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "awake",
		Description: "Restore session context. Optional path narrows context to a specific area (reduces tokens).",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Optional path or keyword to focus context on a specific area",
				},
			},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		graphPath := filepath.Join(root, "knowledge-graph/graphify-out/graph.json")
		if _, err := os.Stat(graphPath); os.IsNotExist(err) {
			msg := `No knowledge graph found yet.

Please ask the user:
> The knowledge graph hasn't been initialized. Would you like to scan the project now?
>
> 1. **Yes** — scan from project root (` + filepath.Dir(root) + `)
> 2. **Yes, custom path** — user will specify the path to scan
> 3. **No** — skip for now (you can run ` + "`mcp__ALPHA__update`" + ` later)

If user chooses 1: call ` + "`mcp__ALPHA__update`" + `
If user chooses 2: call ` + "`mcp__ALPHA__update`" + ` after confirming the path with user
If user chooses 3: stop here`
			return mcp.NewToolResultText(msg), nil
		}

		args, _ := request.Params.Arguments.(map[string]any)
		focusPath, _ := args["path"].(string)

		var out strings.Builder
		out.WriteString("[AGENT_CONTEXT_START]\n\n")

		if focusPath != "" {
			g, err := loadFullGraph(root)
			if err == nil {
				out.WriteString(fmt.Sprintf("### FOCUSED CONTEXT: %s\n", focusPath))
				out.WriteString(sketchGraph(g, focusPath, 3))
				out.WriteString("\n\nSTRATEGY: Focused context loaded. Call detail(ids) for deeper Phase 2 analysis.\n\n")
			}
		} else {
			out.WriteString(buildAwakeOverview(root))
		}

		if content, err := os.ReadFile(filepath.Join(root, "knowledge-graph/memories/session-summary.md")); err == nil {
			out.WriteString("### PREVIOUS SESSION SUMMARY\n")
			out.WriteString(string(content))
			out.WriteString("\n")
		}

		out.WriteString("[AGENT_CONTEXT_END]\n\n")
		return mcp.NewToolResultText(out.String()), nil
	})
}
