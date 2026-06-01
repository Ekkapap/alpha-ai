package main

// cmd_query.go — overview/sketch/detail/update commands: graph Phase 0/1/2 queries (CLI + MCP).

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cliOverview() {
	ov, err := graphOverview(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	b, _ := json.MarshalIndent(ov, "", "  ")
	fmt.Println(string(b))
}

func cliSketch() {
	query := ""
	depth := 3
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--query" && i+1 < len(os.Args) {
			query = os.Args[i+1]
		}
		if os.Args[i] == "--depth" && i+1 < len(os.Args) {
			fmt.Sscanf(os.Args[i+1], "%d", &depth)
		}
	}
	if query == "" {
		fmt.Fprintln(os.Stderr, "Usage: sketch --query <text> [--depth N]")
		os.Exit(2)
	}
	g, err := loadFullGraph(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(sketchGraph(g, query, depth))
}

func cliDetail() {
	ids := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--ids" && i+1 < len(os.Args) {
			ids = os.Args[i+1]
		}
	}
	if ids == "" {
		fmt.Fprintln(os.Stderr, "Usage: detail --ids <id1,id2,...>")
		os.Exit(2)
	}
	g, err := loadFullGraph(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(detailNodes(g, strings.Split(ids, ",")))
}

func registerMCPOverview(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "overview",
		Description: "Phase 0: Returns a compact graph summary (nodes, edges, community count, god nodes, top communities) under 200 tokens. Call this first so the agent knows what to query before running sketch.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ov, err := graphOverview(root)
		if err != nil {
			return mcp.NewToolResultError("overview: " + err.Error()), nil
		}
		b, _ := json.MarshalIndent(ov, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	})
}

func registerMCPSketch(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "sketch",
		Description: "Phase 1: BFS traversal of the knowledge graph. Returns a compact node list (id, label, file, depth, relation) so the agent can assess relevance before asking for full details.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Natural language query, e.g. 'login flow' or 'authentication'",
				},
				"depth": map[string]any{
					"type":        "integer",
					"description": "BFS depth from seed nodes (default: 3)",
				},
			},
			Required: []string{"query"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		query, _ := args["query"].(string)
		depth := 3
		if d, ok := args["depth"].(float64); ok {
			depth = int(d)
		}
		g, err := loadFullGraph(root)
		if err != nil {
			return mcp.NewToolResultError("load graph: " + err.Error()), nil
		}
		return mcp.NewToolResultText(sketchGraph(g, query, depth)), nil
	})
}

func registerMCPDetail(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "detail",
		Description: "Phase 2: Given node IDs from sketch, returns full callers, callees, file type, and community. Call this only for nodes confirmed relevant in Phase 1.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"ids": map[string]any{
					"type":        "string",
					"description": "Comma-separated node IDs from sketch output, e.g. 'authservice,authservice_attempt'",
				},
			},
			Required: []string{"ids"},
		},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		idsStr, _ := args["ids"].(string)
		ids := strings.Split(idsStr, ",")
		g, err := loadFullGraph(root)
		if err != nil {
			return mcp.NewToolResultError("load graph: " + err.Error()), nil
		}
		return mcp.NewToolResultText(detailNodes(g, ids)), nil
	})
}

