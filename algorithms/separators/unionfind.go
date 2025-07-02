package separators

// Union find
// combine roots of two nodes to one root
func union(parent map[int]int, x, y int) {
	rootX := find(parent, x)
	rootY := find(parent, y)
	if rootX != rootY {
		parent[rootY] = rootX
	}
}

// find root node
func find(parent map[int]int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}
