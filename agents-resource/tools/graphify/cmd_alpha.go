package main

// cmd_alpha.go — alpha command: display system identity (CLI + MCP).

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func cliAlpha() {
	fmt.Print(alphaDisplay(root))
}

func registerMCPAlpha(s *server.MCPServer) {
	s.AddTool(mcp.Tool{
		Name:        "alpha",
		Description: "Display the Alpha System identity — logo, stats, commands.",
		InputSchema: mcp.ToolInputSchema{Type: "object", Properties: map[string]any{}},
	}, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		raw := alphaDisplay(root)
		ansiStrip := func(s string) string {
			var out strings.Builder
			skip := false
			for _, r := range s {
				if r == '\033' {
					skip = true
					continue
				}
				if skip {
					if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
						skip = false
					}
					continue
				}
				out.WriteRune(r)
			}
			return out.String()
		}
		return mcp.NewToolResultText(ansiStrip(raw)), nil
	})
}
