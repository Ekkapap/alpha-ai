package main

// cmd_focus.go — focus command: search file sections (CLI + MCP). Also: debug_info MCP tool.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cliFocus() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: focus <path> <term> [max_lines]")
		os.Exit(1)
	}
	path := os.Args[2]
	term := os.Args[3]
	maxLines := 20
	if len(os.Args) > 4 {
		fmt.Sscanf(os.Args[4], "%d", &maxLines)
	}

	absPath := path
	if !filepath.IsAbs(path) {
		absPath = filepath.Join(root, path)
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("❌ Could not read file %s: %v\n", absPath, err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	foundLine := -1
	for i, line := range lines {
		if strings.Contains(line, term) {
			foundLine = i + 1
			break
		}
	}

	if foundLine == -1 {
		fmt.Printf("❌ Term '%s' not found in %s\n", term, path)
		os.Exit(1)
	}

	start := foundLine - 1
	end := foundLine + maxLines
	if end > len(lines) {
		end = len(lines)
	}

	fmt.Printf("🎯 Found '%s' at line %d in %s\n\n", term, foundLine, path)
	for i := start; i < end; i++ {
		fmt.Printf("%d: %s\n", i+1, lines[i])
	}
}

func registerMCPDebugInfo(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "debug_info",
		Description: "Check server environment for debugging.",
		InputSchema: mcp.ToolInputSchema{Type: "object"},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wd, _ := os.Getwd()
		absWd, _ := filepath.Abs(wd)
		return mcp.NewToolResultText(fmt.Sprintf("Root: %s\nWD: %s\nExe: %s", root, absWd, os.Args[0])), nil
	})
}

func registerMCPFocus(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "focus",
		Description: "Search for a term in a file and read a range of lines around it.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "Path to the file to read (relative to project root)",
				},
				"term": map[string]any{
					"type":        "string",
					"description": "Search term to find the starting line",
				},
				"max_lines": map[string]any{
					"type":        "integer",
					"description": "Number of lines to read after the match (default: 20)",
				},
			},
			Required: []string{"path", "term"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("invalid arguments format"), nil
		}

		path, _ := args["path"].(string)
		term, _ := args["term"].(string)
		maxLines := 20
		if val, ok := args["max_lines"].(float64); ok {
			maxLines = int(val)
		}

		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(root, path)
		}
		content, err := os.ReadFile(absPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("could not read file %s: %v", absPath, err)), nil
		}

		lines := strings.Split(string(content), "\n")
		foundLine := -1
		for i, line := range lines {
			if strings.Contains(line, term) {
				foundLine = i + 1
				break
			}
		}

		if foundLine == -1 {
			return mcp.NewToolResultError(fmt.Sprintf("term '%s' not found", term)), nil
		}

		start := foundLine - 1
		end := foundLine + maxLines
		if end > len(lines) {
			end = len(lines)
		}

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Found '%s' at line %d in %s\n\n", term, foundLine, path))
		for i := start; i < end; i++ {
			result.WriteString(fmt.Sprintf("%d: %s\n", i+1, lines[i]))
		}

		return mcp.NewToolResultText(result.String()), nil
	})
}
