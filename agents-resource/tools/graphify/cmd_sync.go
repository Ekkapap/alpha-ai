package main

// cmd_sync.go — sync command: update graph + save session file (CLI + MCP).
// session-summary.md is agent-curated; use update_session_summary MCP tool to write it.

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cliSync() {
	summary := "Manual sync"
	graphifyTarget := os.Getenv("ALPHA_ROOT")
	if graphifyTarget == "" {
		graphifyTarget = filepath.Dir(root)
	}
	for i, v := range os.Args {
		if (v == "--summary" || v == "-s") && i+1 < len(os.Args) {
			summary = os.Args[i+1]
		}
		if v == "--path" && i+1 < len(os.Args) {
			graphifyTarget = os.Args[i+1]
		}
	}

	fmt.Println("Syncing graphify...")
	if result, err := buildGraph(graphifyTarget, root, filepath.Dir(root), false); err != nil {
		fmt.Fprintf(os.Stderr, "warning: graph update: %v\n", err)
	} else {
		fmt.Printf("[graphify watch] Rebuilt: %d nodes, %d edges, %d communities\n",
			result.Nodes, result.Edges, result.Communities)
		if result.NeedLabels {
			fmt.Print("\n" + result.LabelPrompt)
		}
	}

	ts := time.Now().Format("2006-01-02 15:04")
	sessionPath := writeSessionFile(root, summary)
	appendSessionToSummary(root, ts, summary)
	fmt.Printf("Saved: %s\n", sessionPath)
	fmt.Println("session-summary.md updated.")

	exec.Command("open", filepath.Join(graphifyDataDir(root), "graph.html")).Run()
}

func registerMCPSync(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "sync",
		Description: "Incremental graph sync. Saves session-[ts].md and returns the entry. Merge the returned entry into session-summary.md via update_session_summary.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"summary": map[string]any{
					"type":        "string",
					"description": "Session summary to save",
				},
			},
			Required: []string{"summary"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		summary, _ := args["summary"].(string)

		slog.Info("sync: updating graph", "root", root)
		projectRoot := filepath.Dir(root)
		var graphMsg string
		if result, err := buildGraph(projectRoot, root, projectRoot, false); err != nil {
			slog.Error("graph update failed", "error", err)
			graphMsg = fmt.Sprintf("graph update error: %v", err)
		} else {
			graphMsg = fmt.Sprintf("graph: %d nodes, %d edges, %d communities", result.Nodes, result.Edges, result.Communities)
			if result.NeedLabels {
				graphMsg += "\n" + result.LabelPrompt
			}
		}

		ts := time.Now().Format("2006-01-02 15:04")
		sessionPath := writeSessionFile(root, summary)
		appendSessionToSummary(root, ts, summary)

		return mcp.NewToolResultText(fmt.Sprintf(
			"%s\nSession saved: %s\nsession-summary.md updated.",
			graphMsg, filepath.Base(sessionPath),
		)), nil
	})
}

// writeSessionFile writes memories/session-[ts].md and returns its path.
func writeSessionFile(alphaRoot, summary string) string {
	dir := memoriesDir(alphaRoot)
	os.MkdirAll(dir, 0755)
	ts := time.Now().Format("2006-01-02 15:04")
	tsFile := time.Now().Format("20060102-1504")
	path := filepath.Join(dir, fmt.Sprintf("session-%s.md", tsFile))
	os.WriteFile(path, []byte(fmt.Sprintf("# Session %s\n\n%s\n", ts, summary)), 0644)
	return path
}

// appendSessionToSummary appends a new session block to session-summary.md and updates the timestamp line.
func appendSessionToSummary(alphaRoot, ts, summary string) {
	summaryPath := filepath.Join(memoriesDir(alphaRoot), "session-summary.md")
	existing, _ := os.ReadFile(summaryPath)

	// Update อัปเดตล่าสุด timestamp.
	content := string(existing)
	date := ts[:10]
	if idx := strings.Index(content, "> อัปเดตล่าสุด:"); idx >= 0 {
		end := strings.Index(content[idx:], "\n")
		if end < 0 {
			end = len(content) - idx
		}
		content = content[:idx] + fmt.Sprintf("> อัปเดตล่าสุด: %s", date) + content[idx+end:]
	}

	// Append session block.
	entry := fmt.Sprintf("\n---\n\n## Session %s\n\n%s\n", ts, summary)
	content += entry

	os.WriteFile(summaryPath, []byte(content), 0644)
}

func registerMCPUpdateSessionSummary(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name: "update_session_summary",
		Description: "Write agent-rewritten session-summary.md to disk. " +
			"Call only when semantic cleanup is needed (AI rewrote the full content). " +
			"For normal sync, Go handles dedup automatically — no need to call this.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"content": map[string]any{
					"type":        "string",
					"description": "Full rewritten content for session-summary.md",
				},
			},
			Required: []string{"content"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		content, _ := args["content"].(string)
		summaryPath := filepath.Join(memoriesDir(root), "session-summary.md")
		if err := os.WriteFile(summaryPath, []byte(content), 0644); err != nil {
			return mcp.NewToolResultError("write session-summary.md: " + err.Error()), nil
		}
		return mcp.NewToolResultText("session-summary.md updated."), nil
	})
}
