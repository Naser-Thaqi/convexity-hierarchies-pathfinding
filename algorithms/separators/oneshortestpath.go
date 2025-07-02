package separators

import (
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
	"math/rand"
	"sort"
)

// Decompose graph by removing a shortest path between two boundary nodes
func OneShortestPath(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	bordernodes := extractBoundaryNodes(g)
	//randomize bordernodes testing order
	rand.Shuffle(len(bordernodes), func(i, j int) {
		bordernodes[i], bordernodes[j] = bordernodes[j], bordernodes[i]
	})

	// try every border node till one path succeeds
	for _, node := range bordernodes {
		select {
		case <-ctx.Done():
			return nil, false
		default:
			// proceed
		}
		prevMap := bfsPaths(g.AdjList, node)             // compute prev map of starting node
		paths := createPaths(prevMap, bordernodes, node) // compute a path to every other border node
		reducedPaths := reducePaths(g, paths)            // reduce computed paths
		paths = nil
		sort.Slice(reducedPaths, func(i, j int) bool { // sort every path in ascending length
			return len(reducedPaths[i]) < len(reducedPaths[j])
		})
		// try every path till one succeeds
		for _, candidate := range reducedPaths {
			select {
			case <-ctx.Done():
				return nil, false
			default:
				// proceed
			}
			convexComponents, valid := graphdecomp.OneShortestPathBalancedConvexDecomposition(g, candidate, ctx)
			if valid {
				bordernodes = nil
				return convexComponents, valid
			}
		}
	}
	return nil, false
}

// Reduces given paths p to p* paths
// Delete all subsequences of nodes with degree less than 4 except first and last node
func reducePaths(g *graph.Graph, paths map[int][]int) [][]int {
	reducedPaths := [][]int{}
	//iterate through all paths
	for _, path := range paths {
		if len(path) < 2 {
			reducedPaths = append(reducedPaths, path)
			continue
		}
		// build reduced path
		reducedP := []int{}
		obstacleStart := -1
		intervalEnd := -1
		// iterate through path
		for i := range path {
			degree := len(g.AdjList[path[i]])
			if degree != 4 {
				// check if node is starting point of subsequence
				if obstacleStart == -1 {
					obstacleStart = i
				}
				intervalEnd = i
			} else {
				if obstacleStart != -1 {
					// if subsequence is one node in length, keep in path
					if obstacleStart == i-1 {
						reducedP = append(reducedP, path[obstacleStart])
					} else {
						// keep start and end node of subsequence
						reducedP = append(reducedP, path[obstacleStart], path[i-1])
					}
					// subsequence finished
					obstacleStart = -1
					intervalEnd = -1
				}
				// current degree == 4, so append in path
				reducedP = append(reducedP, path[i])
			}
		}

		if obstacleStart != -1 {
			// path ends during obstacle subsequence
			if obstacleStart == intervalEnd {
				reducedP = append(reducedP, path[obstacleStart])
			} else {
				// keep first and last node
				reducedP = append(reducedP, path[obstacleStart], path[intervalEnd])
			}
		}
		//append reduced path into result
		reducedPaths = append(reducedPaths, reducedP)
	}
	return reducedPaths
}

// Returns set of paths of start node to each other boundarynode
func createPaths(prev map[int]int, boundaryNodes []int, start int) map[int][]int {
	fullPaths := make(map[int][]int)

	fullPaths[start] = []int{start} // add path startnode to itself as path

	//iterate through all boundary nodes
	for _, node := range boundaryNodes {
		// check if boundary node was reachable of startnode (same connected component)
		if _, exists := prev[node]; exists {
			path := []int{node} //build path
			current := node
			// traverse map back to start
			for current != start {
				parent := prev[current]
				path = append(path, parent)
				current = parent
			}
			fullPaths[node] = path
		}
	}

	return fullPaths
}

// Visit each node in graph with bfs and store previous visited node in a map
func bfsPaths(adjlist map[int][]int, start int) map[int]int {
	prev := make(map[int]int)
	visited := make(map[int]struct{})
	queue := []int{start}

	visited[start] = struct{}{} // set true
	head := 0

	for head < len(queue) {
		current := queue[head]
		head++

		//visit neighbors
		for _, neighbor := range adjlist[current] {
			//put non-visited nodes into queue
			if _, exists := visited[neighbor]; !exists {
				visited[neighbor] = struct{}{}
				// set current node to previous node of neighbor
				prev[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return prev
}

// returns set of boundarynodes
func extractBoundaryNodes(g *graph.Graph) []int {
	boundaryNodes := []int{}

	for node, neighbors := range g.AdjList {
		if len(neighbors) != 4 {
			boundaryNodes = append(boundaryNodes, node)
		}
	}
	return boundaryNodes
}
