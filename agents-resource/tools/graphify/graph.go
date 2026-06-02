package main

// graph.go — graph types and operations, reused by all commands.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ── Two-phase query types ────────────────────────────────────────────────────

type gNode struct {
	ID             string `json:"id"`
	Label          string `json:"label"`
	SourceFile     string `json:"source_file"`
	SourceLocation string `json:"source_location"`
	Community      int    `json:"community"`
	FileType       string `json:"file_type"`
}

type gLink struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Relation string `json:"relation"`
}

type fullGraph struct {
	Nodes []gNode `json:"nodes"`
	Links []gLink `json:"links"`
}

func loadFullGraph(r string) (*fullGraph, error) {
	data, err := os.ReadFile(filepath.Join(graphifyDataDir(r), "graph.json"))
	if err != nil {
		return nil, err
	}
	var g fullGraph
	return &g, json.Unmarshal(data, &g)
}

func scoreNode(n gNode, terms []string) float64 {
	label := strings.ToLower(n.Label)
	file := strings.ToLower(n.SourceFile)
	var s float64
	for _, t := range terms {
		t = strings.ToLower(t)
		switch {
		case strings.EqualFold(n.Label, t):
			s += 3
		case strings.Contains(label, t):
			s += 2
		case strings.Contains(file, t):
			s += 1
		}
	}
	return s
}

// sketchGraph runs BFS from top-scoring seed nodes and returns compact JSON.
func sketchGraph(g *fullGraph, query string, depth int) string {
	nodeIdx := make(map[string]gNode, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeIdx[n.ID] = n
	}
	adj := make(map[string][]gLink, len(g.Nodes))
	for _, l := range g.Links {
		adj[l.Source] = append(adj[l.Source], l)
		rev := gLink{Source: l.Target, Target: l.Source, Relation: l.Relation}
		adj[l.Target] = append(adj[l.Target], rev)
	}

	terms := strings.Fields(query)
	type scored struct {
		n gNode
		s float64
	}
	var seeds []scored
	for _, n := range g.Nodes {
		if sc := scoreNode(n, terms); sc >= 1 {
			seeds = append(seeds, scored{n, sc})
		}
	}
	sort.Slice(seeds, func(i, j int) bool { return seeds[i].s > seeds[j].s })
	if len(seeds) > 5 {
		seeds = seeds[:5]
	}

	type qItem struct {
		id  string
		dep int
		via string
	}
	visited := make(map[string]bool)
	queue := make([]qItem, 0, len(seeds))
	seedLabels := make([]string, 0, len(seeds))
	for _, s := range seeds {
		visited[s.n.ID] = true
		queue = append(queue, qItem{s.n.ID, 0, ""})
		seedLabels = append(seedLabels, s.n.Label)
	}

	type outNode struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		File      string `json:"file"`
		Community int    `json:"community"`
		Depth     int    `json:"depth"`
		Via       string `json:"via,omitempty"`
	}
	var nodes []outNode
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		n, ok := nodeIdx[cur.id]
		if !ok {
			continue
		}
		loc := n.SourceFile
		if n.SourceLocation != "" {
			loc += ":" + n.SourceLocation
		}
		nodes = append(nodes, outNode{n.ID, n.Label, loc, n.Community, cur.dep, cur.via})
		if cur.dep < depth {
			for _, l := range adj[cur.id] {
				if !visited[l.Target] {
					visited[l.Target] = true
					queue = append(queue, qItem{l.Target, cur.dep + 1, l.Relation})
				}
			}
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Depth != nodes[j].Depth {
			return nodes[i].Depth < nodes[j].Depth
		}
		return nodes[i].Label < nodes[j].Label
	})

	out := map[string]any{"query": query, "seeds": seedLabels, "depth": depth, "total": len(nodes), "nodes": nodes}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b)
}

// detailNodes returns callers, callees, and file info for given node IDs.
func detailNodes(g *fullGraph, ids []string) string {
	idSet := make(map[string]bool, len(ids))
	for _, id := range ids {
		idSet[strings.TrimSpace(id)] = true
	}
	nodeIdx := make(map[string]gNode, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeIdx[n.ID] = n
	}

	callers := make(map[string][]string)
	callees := make(map[string][]string)
	for _, l := range g.Links {
		if l.Relation != "calls" {
			continue
		}
		callers[l.Target] = append(callers[l.Target], l.Source)
		callees[l.Source] = append(callees[l.Source], l.Target)
	}

	resolve := func(nodeIDs []string) []string {
		seen := make(map[string]bool)
		var out []string
		for _, id := range nodeIDs {
			if seen[id] {
				continue
			}
			seen[id] = true
			if n, ok := nodeIdx[id]; ok {
				out = append(out, n.Label+" ("+n.SourceFile+")")
			}
		}
		return out
	}

	type outNode struct {
		ID        string   `json:"id"`
		Label     string   `json:"label"`
		File      string   `json:"file"`
		FileType  string   `json:"file_type"`
		Community int      `json:"community"`
		Callers   []string `json:"callers"`
		Callees   []string `json:"callees"`
	}
	var nodes []outNode
	for _, n := range g.Nodes {
		if !idSet[n.ID] {
			continue
		}
		loc := n.SourceFile
		if n.SourceLocation != "" {
			loc += ":" + n.SourceLocation
		}
		nodes = append(nodes, outNode{
			ID: n.ID, Label: n.Label, File: loc,
			FileType: n.FileType, Community: n.Community,
			Callers: resolve(callers[n.ID]),
			Callees: resolve(callees[n.ID]),
		})
	}
	b, _ := json.MarshalIndent(map[string]any{"nodes": nodes}, "", "  ")
	return string(b)
}

// ── Phase 0: overview ────────────────────────────────────────────────────────

type godNodeInfo struct {
	Label     string `json:"label"`
	Edges     int    `json:"edges"`
	Community string `json:"community"`
}

type overviewResult struct {
	Nodes          int           `json:"nodes"`
	Edges          int           `json:"edges"`
	Communities    int           `json:"communities"`
	GodNodes       []godNodeInfo `json:"god_nodes"`
	TopCommunities []string      `json:"top_communities"`
}

func graphOverview(r string) (*overviewResult, error) {
	graphData, err := os.ReadFile(filepath.Join(graphifyDataDir(r), "graph.json"))
	if err != nil {
		return nil, err
	}
	var g struct {
		Nodes []json.RawMessage `json:"nodes"`
		Links []json.RawMessage `json:"links"`
	}
	if err := json.Unmarshal(graphData, &g); err != nil {
		return nil, err
	}

	analysisData, err := os.ReadFile(filepath.Join(graphifyDataDir(r), ".graphify_analysis.json"))
	if err != nil {
		return nil, err
	}
	var a struct {
		Communities map[string][]string `json:"communities"`
		Gods        []struct {
			ID     string `json:"id"`
			Label  string `json:"label"`
			Degree int    `json:"degree"`
		} `json:"gods"`
	}
	if err := json.Unmarshal(analysisData, &a); err != nil {
		return nil, err
	}

	nodeToCommunity := make(map[string]string, len(g.Nodes))
	for commID, members := range a.Communities {
		for _, nodeID := range members {
			nodeToCommunity[nodeID] = commID
		}
	}

	labelsData, _ := os.ReadFile(filepath.Join(graphifyDataDir(r), ".graphify_labels.json"))
	var labels map[string]string
	json.Unmarshal(labelsData, &labels)
	communityName := func(id string) string {
		if labels != nil {
			if name, ok := labels[id]; ok {
				return name
			}
		}
		return "Community " + id
	}

	gods := make([]godNodeInfo, 0, len(a.Gods))
	for _, gn := range a.Gods {
		commID := nodeToCommunity[gn.ID]
		gods = append(gods, godNodeInfo{
			Label:     gn.Label,
			Edges:     gn.Degree,
			Community: communityName(commID),
		})
	}

	seen := make(map[string]bool)
	var topComms []string
	for _, gn := range a.Gods {
		commID := nodeToCommunity[gn.ID]
		name := communityName(commID)
		if !seen[name] {
			seen[name] = true
			topComms = append(topComms, name)
		}
	}

	return &overviewResult{
		Nodes:          len(g.Nodes),
		Edges:          len(g.Links),
		Communities:    len(a.Communities),
		GodNodes:       gods,
		TopCommunities: topComms,
	}, nil
}
