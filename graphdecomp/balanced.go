package graphdecomp

import (
	"bachelor-project/config"
)

// returns map of nodeids which stores root node (connected components)
// parent map
func unionFind(adjlist map[int][]int) map[int]int {
	parent := make(map[int]int)

	// initialize map with every node pointing to itself
	for key := range adjlist {
		parent[key] = key
	}

	// unite every key with its neighbor nodes
	for node := range adjlist {
		for _, neighbor := range adjlist[node] {
			if node < neighbor {
				union(parent, node, neighbor)
			}
		}
	}

	//finalize every node to its root
	for node := range parent {
		parent[node] = find(parent, node)
	}

	return parent
}

// recursively find root of nodeid
func find(parent map[int]int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

// combine roots of two nodes to one root
func union(parent map[int]int, x, y int) {
	rootX := find(parent, x)
	rootY := find(parent, y)
	if rootX != rootY {
		parent[rootY] = rootX // Einfacher Merge
	}
}

// return number of connected components
func countComponents(parent map[int]int, components map[int]bool) int {
	for id := range parent {
		root := find(parent, id)
		components[root] = true
	}
	return len(components)
}

// checks if new subgraphs are alpha balanced
// expects parent map
func checkBalanced(parent map[int]int, nodeCount int) bool {
	components := make(map[int]bool)

	// check if graph is made of 2 or more subgraphes
	if countComponents(parent, components) < 2 {
		return false
	}

	// upper boundary for each subgraph
	limit := int(float64(nodeCount) * config.Alpha)

	// compute number of node per connected component
	componentSizes := make(map[int]int)
	for _, root := range parent {
		componentSizes[root]++
	}
	// check every connected component for valid size
	for _, size := range componentSizes {
		if size > limit {
			return false
		}
	}

	return true
}
