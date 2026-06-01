package main

// community.go — label propagation community detection + god node detection.

import (
	"math/rand"
	"sort"
)

// AssignCommunities runs label propagation and returns community ID per node index.
func AssignCommunities(nodeIDs []string, edges [][2]int, iterations int) []int {
	n := len(nodeIDs)
	if n == 0 {
		return nil
	}

	// Build adjacency list
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		if a >= 0 && a < n && b >= 0 && b < n {
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
	}

	// Init: each node is its own community
	labels := make([]int, n)
	for i := range labels {
		labels[i] = i
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}

	for iter := 0; iter < iterations; iter++ {
		rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
		changed := false
		for _, i := range order {
			if len(adj[i]) == 0 {
				continue
			}
			// Count neighbor labels
			freq := make(map[int]int, len(adj[i]))
			for _, nb := range adj[i] {
				freq[labels[nb]]++
			}
			best, bestCount := labels[i], 0
			for lbl, cnt := range freq {
				if cnt > bestCount || (cnt == bestCount && lbl < best) {
					best, bestCount = lbl, cnt
				}
			}
			if best != labels[i] {
				labels[i] = best
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// Normalize community IDs to 0..k-1
	idMap := map[int]int{}
	next := 0
	result := make([]int, n)
	for i, lbl := range labels {
		if _, ok := idMap[lbl]; !ok {
			idMap[lbl] = next
			next++
		}
		result[i] = idMap[lbl]
	}
	return result
}

// GodNode is a high-degree node.
type GodNode struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Degree int    `json:"degree"`
}

// FindGodNodes returns top N nodes by degree.
func FindGodNodes(nodeIDs, nodeLabels []string, edges [][2]int, topN int) []GodNode {
	degree := make(map[int]int)
	for _, e := range edges {
		degree[e[0]]++
		degree[e[1]]++
	}
	type pair struct {
		idx int
		deg int
	}
	pairs := make([]pair, 0, len(degree))
	for idx, deg := range degree {
		pairs = append(pairs, pair{idx, deg})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].deg > pairs[j].deg })
	if topN > len(pairs) {
		topN = len(pairs)
	}
	result := make([]GodNode, topN)
	for i, p := range pairs[:topN] {
		result[i] = GodNode{
			ID:     nodeIDs[p.idx],
			Label:  nodeLabels[p.idx],
			Degree: p.deg,
		}
	}
	return result
}
