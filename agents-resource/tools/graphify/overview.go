package main

// overview.go — buildAwakeOverview: reused by cliAwake and registerMCPAwake.

import (
	"fmt"
	"sort"
	"strings"
)

// buildAwakeOverview returns a formatted overview string with graph stats and god nodes.
// Tries graphOverview (requires .graphify_analysis.json) first; falls back to computing
// god nodes directly from graph.json link degrees so it never silently skips.
func buildAwakeOverview(r string) string {
	var sb strings.Builder

	ov, err := graphOverview(r)
	if err == nil {
		// Happy path: analysis file exists
		sb.WriteString(fmt.Sprintf("### GRAPH STATS\nNodes: %d | Edges: %d | Communities: %d\n\n",
			ov.Nodes, ov.Edges, ov.Communities))

		if len(ov.GodNodes) > 0 {
			sb.WriteString("### GOD NODES (Architecture Pillars)\n")
			for i, gn := range ov.GodNodes {
				sb.WriteString(fmt.Sprintf("%d. %s (%d edges) — %s\n", i+1, gn.Label, gn.Edges, gn.Community))
			}
			sb.WriteString("\n")
		}

		if len(ov.TopCommunities) > 0 {
			sb.WriteString("### TOP COMMUNITIES\n")
			for _, c := range ov.TopCommunities {
				sb.WriteString("- ")
			sb.WriteString(c)
			sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	} else {
		// Fallback: compute stats directly from graph.json
		g, err := loadFullGraph(r)
		if err != nil {
			sb.WriteString("### GRAPH\n(not initialized)\n\n")
			return sb.String()
		}

		sb.WriteString(fmt.Sprintf("### GRAPH STATS\nNodes: %d | Edges: %d\n\n", len(g.Nodes), len(g.Links)))

		// Compute degree from links
		degree := make(map[string]int, len(g.Nodes))
		for _, l := range g.Links {
			degree[l.Source]++
			degree[l.Target]++
		}

		// Index nodes by ID
		nodeIdx := make(map[string]gNode, len(g.Nodes))
		for _, n := range g.Nodes {
			nodeIdx[n.ID] = n
		}

		type pair struct {
			id  string
			deg int
		}
		pairs := make([]pair, 0, len(degree))
		for id, d := range degree {
			pairs = append(pairs, pair{id, d})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].deg > pairs[j].deg })
		if len(pairs) > 10 {
			pairs = pairs[:10]
		}

		if len(pairs) > 0 {
			sb.WriteString("### GOD NODES (Architecture Pillars)\n")
			for i, p := range pairs {
				label := p.id
				if n, ok := nodeIdx[p.id]; ok {
					label = n.Label
				}
				sb.WriteString(fmt.Sprintf("%d. %s (%d edges)\n", i+1, label, p.deg))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("STRATEGY: Use overview→sketch→detail flow. Call sketch(query) for Phase 1 BFS, then detail(ids) for Phase 2 callers/callees.\n\n")
	return sb.String()
}
