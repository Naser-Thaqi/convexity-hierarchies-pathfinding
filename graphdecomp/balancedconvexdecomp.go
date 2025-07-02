package graphdecomp

import (
	"bachelor-project/graph"
	"context"
)

// Decompose graph into convex alpha balanced subgraphs
func BalancedConvexDecomposition(g *graph.Graph, separator []int, ctx context.Context) ([]*graph.Graph, bool) {
	if len((separator)) > 0 {

		// Create adjacency list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		// To differentiate connected components (subgraphs)
		parent := unionFind(copyAdjlist)

		if checkBalanced(parent, len(g.AdjList)) {
			if checkConvexity(g, copyAdjlist, parent, ctx) {
				return decomposeGraph(g, parent), true
			}
		}
	}
	return nil, false
}

func ConvexDecomposition(g *graph.Graph, separator []int, ctx context.Context) ([]*graph.Graph, bool) {
	if len((separator)) > 0 {

		// Create adjancy list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		// To differentiate connected components (subgraphs)
		parent := unionFind(copyAdjlist)

		if checkConvexity(g, copyAdjlist, parent, ctx) {
			return decomposeGraph(g, parent), true
		}
	}
	return nil, false
}

// Decompose Graph into balanced subgraphs
func BalancedDecomposition(g *graph.Graph, separator []int) ([]*graph.Graph, bool) {
	if len((separator)) > 0 {

		// Create adjancy list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		// To differentiate connected components (subgraphs)
		parent := unionFind(copyAdjlist)

		if checkBalanced(parent, len(g.AdjList)) {
			return decomposeGraph(g, parent), true
		}
	}
	return nil, false
}

// Decompose graph into convex balanced subgraphs
// but checking first and last node for usefulness and observation 7 for easier convexity check
func OneShortestPathBalancedConvexDecomposition(g *graph.Graph, separator []int, ctx context.Context) ([]*graph.Graph, bool) {
	if len(separator) > 0 {
		// Create adjancy list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		// To differentiate connected components (subgraphs)
		parent := unionFind(copyAdjlist)

		nodes := []int{separator[0], separator[len(separator)-1]} // check first and last node
		if separator[0] == separator[len(separator)-1] {          // edge-case if separator length == 1
			nodes = nodes[:1]
		}

		// Check if start and end node are useful for decomposition
		for i, node := range nodes {
			visited := make(map[int]struct{})
			for _, neighbor := range g.AdjList[node] {
				if _, exists := parent[neighbor]; exists {
					visited[parent[neighbor]] = struct{}{}
				}
			}
			// separator node has only one component as neighbor, restore node
			if len(visited) <= 1 {
				neighbors := []int{}
				for _, neighbor := range g.AdjList[node] {
					// filter for non-separator neighbors
					if _, exists := parent[neighbor]; exists {
						neighbors = append(neighbors, neighbor)
					}
				}
				copyAdjlist[node] = neighbors
				for _, neighbor := range neighbors {
					copyAdjlist[neighbor] = append(copyAdjlist[neighbor], node)
				}
				// Delete node from separator set
				if i == 0 {
					separator = separator[1:]
				} else {
					separator = separator[:len(separator)-1]
				}
				if len(neighbors) > 0 {
					parent[node] = parent[neighbors[0]] // Add node to parent map
				}

			}
		}
		if checkBalanced(parent, len(g.AdjList)) {
			if degreeFour(g, separator) {
				return decomposeGraph(g, parent), true
			}
			if checkObservationAndConvexity(g, copyAdjlist, parent, separator, ctx) {
				return decomposeGraph(g, parent), true
			}
		}
	}

	return nil, false
}

// Check wether or not the given separator set leads to a valid solution
func CheckOneShortestPathBalancedConvex(g *graph.Graph, separator []int, ctx context.Context) bool {
	if len((separator)) > 0 {

		//Create adjancy list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		// To differentiate connected components (subgraphs)
		parent := unionFind(copyAdjlist)

		if checkBalanced(parent, len(g.AdjList)) {
			if degreeFour(g, separator) {
				return true
			}
			if checkObservationAndConvexity(g, copyAdjlist, parent, separator, ctx) {
				return true
			}
		}
	}
	return false
}

func CheckBalanced(g *graph.Graph, separator []int) bool {
	if len((separator)) > 0 {

		// Create adjancy list of all subgraphes combined into list
		copyAdjlist := g.CopyAdjlist()
		for _, node := range separator {
			graph.RemoveNode(copyAdjlist, node)
		}

		return checkBalanced(unionFind(copyAdjlist), len(g.AdjList))
	}

	return false
}

// Decomposes graph if it already has more than one connected component
func DecomposeInputComponents(g *graph.Graph) ([]*graph.Graph, bool) {
	parent := unionFind(g.AdjList)
	components := make(map[int]bool)

	if countComponents(parent, components) > 1 {
		return decomposeGraph(g, parent), true
	}

	return nil, false
}

// decompose graph into subgraphs
func decomposeGraph(g *graph.Graph, parent map[int]int) []*graph.Graph {
	// Coordinates for computing size of grid [][]int
	yTop, yLow := g.Height, -1
	xLeft, xRight := g.Width, -1

	sizeNodes := make(map[int][]int) // key = connected components, values = coordinates for size

	// Initialize values for comparison
	for key, val := range parent {
		if key == val {
			sizeNodes[key] = make([]int, 0, 4)
			sizeNodes[key] = append(sizeNodes[key], yTop, yLow, xLeft, xRight)
		}
	}

	// Iterate through grid and store coordinates for rectangle creation (grid [][]int)
	for y := range g.Height {
		for x := range g.Width {
			node := g.Grid[y][x]
			// if node is -1 non passable node it shouldnt exist in parent (-1 can't be a nodeid)
			root, exists := parent[node]
			if exists {
				// Sequential if statements cause a subgraph can be for example one node, a row of nodes, a column...
				// access coordinates of connected component via sizeNodes[root]
				if y < sizeNodes[root][0] {
					sizeNodes[root][0] = y
				}
				if y > sizeNodes[root][1] {
					sizeNodes[root][1] = y
				}

				if x < sizeNodes[root][2] {
					sizeNodes[root][2] = x
				}

				if x > sizeNodes[root][3] {
					sizeNodes[root][3] = x
				}
			}
		}
	}

	// create children array for original graph
	subgraphes := make([]*graph.Graph, 0, len(sizeNodes))
	// key is root of component
	for key := range sizeNodes {
		yTop, yLow = sizeNodes[key][0], sizeNodes[key][1]
		xLeft, xRight = sizeNodes[key][2], sizeNodes[key][3]
		height := yLow - yTop + 1
		width := xRight - xLeft + 1

		// create subgraph as own graph object, reserve appropiate memory
		subgraph := graph.NewGraph(height, width)
		subgrid := make([][]int, height)
		for y := range height {
			subgrid[y] = make([]int, width)
		}

		// Coordinates in subgraph for iteration
		ySub := 0
		xSub := 0

		// Iterate through original grid and copy section
		for y := yTop; y <= yLow; y++ {
			for x := xLeft; x <= xRight; x++ {
				node := g.Grid[y][x]
				root, exists := parent[node]
				// check for -1
				if exists {
					// check if nodeid should be in current subgraph (connected component)
					if root == key {
						subgrid[ySub][xSub] = node
					} else {
						subgrid[ySub][xSub] = -1
					}
				} else {
					subgrid[ySub][xSub] = -1
				}
				xSub++
			}
			xSub = 0
			ySub++
		}
		subgraph.Grid = subgrid
		subgraph.BuildAdjlist()
		subgraphes = append(subgraphes, subgraph)
	}
	return subgraphes
}
