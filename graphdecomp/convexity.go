package graphdecomp

import (
	"bachelor-project/graph"
	"context"
)

// check original graph and subgraph for convexity with help of their adjancy lists
// expects already adjlist of subgraph and parent map (union find)
func checkConvexity(g *graph.Graph, adjlist map[int][]int, parent map[int]int, ctx context.Context) bool {
	boundaryNodes := extractBorderNodesOfComponents(adjlist, parent)
	// check every connected component
	for _, boundaryList := range boundaryNodes {
		for _, node := range boundaryList {
			select {
			case <-ctx.Done():
				// method cancelled
				return false
			default:
				// proceed
			}
			// start bfs for every boundary node per component
			distSub, maxDepth := bfs(adjlist, node)                //compute distance to each other node and max depth of bfs
			filterDistancesForBoundary(adjlist, distSub)           //filter dist map only for boundary nodes
			if !isNodeConvex(g.AdjList, distSub, node, maxDepth) { // check for same distance between nodes in original graph
				return false
			}
		}
	}

	return true
}

func checkObservationAndConvexity(g *graph.Graph, adjlist map[int][]int, parent map[int]int, separator []int, ctx context.Context) bool {
	boundaryNodes := extractBorderNodesOfComponents(adjlist, parent)
	adjacentNodes := getAdjacentNodesOfSeparator(g, separator, parent)

	// Adjlist for coordinates
	coordAdjlist := make(map[int]int, len(g.AdjList))
	coordID := 0
	for y := range g.Height {
		for x := range g.Width {
			nodeID := g.Grid[y][x]
			if nodeID != -1 {
				coordAdjlist[nodeID] = coordID
			}
			coordID++
		}
	}

	for root, boundaryList := range boundaryNodes {
		// check if observation 7 applies to skip convexity check with bfs
		if !checkObservation(g, coordAdjlist, adjacentNodes[root]) {
			for _, node := range boundaryList {
				select {
				case <-ctx.Done():
					return false
				default:
					// proceed
				}
				// start bfs for every boundary node per component
				distSub, maxDepth := bfs(adjlist, node)                // compute distance to each other node and max depth of bfs
				filterDistancesForBoundary(adjlist, distSub)           // filter dist map only for boundary nodes
				if !isNodeConvex(g.AdjList, distSub, node, maxDepth) { // check for same distance between nodes in original graph
					return false
				}
			}
		}
	}

	return true
}

// Checks if every node in path has degree 4 except first and last node
func degreeFour(g *graph.Graph, path []int) bool {
	for i := 1; i < len(path)-1; i++ {
		if len(g.AdjList[path[i]]) != 4 {
			return false
		}
	}
	return true
}

func checkObservation(g *graph.Graph, coordAdjlist map[int]int, path []int) bool {
	if len(path) < 2 {
		return true
	}

	var xDir, yDir int // +1 für rechts/unten, -1 für links/oben

	for i := 0; i < len(path)-1; i++ {
		x1, y1 := graph.CoordinatesFromNodeID(coordAdjlist[path[i]], g.Width)
		x2, y2 := graph.CoordinatesFromNodeID(coordAdjlist[path[i+1]], g.Width)

		dx := x2 - x1
		dy := y2 - y1

		// check x direction
		if dx != 0 {
			dir := dx / abs(dx)
			if xDir == 0 {
				xDir = dir
			} else if xDir != dir {
				return false // conflict in x-direction
			}
		}

		// check y direction
		if dy != 0 {
			dir := dy / abs(dy)
			if yDir == 0 {
				yDir = dir
			} else if yDir != dir {
				return false // conflict in y-direction
			}
		}
	}

	return true
}

// compute absolute value
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getAdjacentNodesOfSeparator(g *graph.Graph, separators []int, parent map[int]int) map[int][]int {
	adjacentNodes := make(map[int][]int)

	for _, node := range separators {
		for _, neighbor := range g.AdjList[node] {
			if root, exists := parent[neighbor]; exists {
				adjacentNodes[root] = append(adjacentNodes[root], neighbor)
			}
		}
	}

	return adjacentNodes
}

// filter dist map for boundary nodes
func filterDistancesForBoundary(adjlist map[int][]int, dist map[int]int) {
	for key := range dist {
		if len(adjlist[key]) == 4 {
			delete(dist, key)
		}
	}
}

// Returns boundary nodes (deg(v)<4) of a given Graph for each connected component
func extractBorderNodesOfComponents(adjlist map[int][]int, parent map[int]int) map[int][]int {
	boundaryNodes := make(map[int][]int) // key = root of connected component, values = all nodes in same connected component

	for node := range adjlist {
		if len(adjlist[node]) < 4 {
			boundaryNodes[parent[node]] = append(boundaryNodes[parent[node]], node)
		}
	}
	return boundaryNodes
}

// Returns one-to-many distance relationship via bfs
func bfs(adjlist map[int][]int, start int) (map[int]int, int) {
	dist := make(map[int]int) // store distances to each node
	visited := make(map[int]struct{})

	queue := []int{start}
	dist[start] = 0             // distance from start to start is 0
	visited[start] = struct{}{} // set true
	head := 0                   // pointer for avoiding sclice copys

	for head < len(queue) {
		current := queue[head]
		head++
		// visit all neighbors
		for _, neighbor := range adjlist[current] {
			// check if neighbor was already visited
			if _, exists := visited[neighbor]; !exists {
				visited[neighbor] = struct{}{}     // set true
				dist[neighbor] = dist[current] + 1 // store distance
				queue = append(queue, neighbor)    // push to queue to visit its neighbors later
			}
		}
	}
	maxDepth := 0
	for _, distance := range dist {
		if maxDepth < distance {
			maxDepth = distance
		}
	}

	return dist, maxDepth
}

// Checks if distances of a node to all other border nodes in subgraph are similar in original graph
func isNodeConvex(originalAdj map[int][]int, distSub map[int]int, start int, maxDepth int) bool {
	visited := make(map[int]struct{})
	queue := []int{start}
	visited[start] = struct{}{} // set true

	depth := 0                   // store current depth for early abortion
	subVisitedCount := 1         // trace how many bordernodes of subgraph were visited
	subNodeCount := len(distSub) // number of bordernodes that needs to be visited

	head := 0 // pointer

	for head < len(queue) && depth <= maxDepth {
		levelSize := len(queue) - head
		for range levelSize {
			current := queue[head]
			head++

			// visit all neighbors
			for _, neighbor := range originalAdj[current] {
				// visit non-visited neighbors
				if _, exists := visited[neighbor]; !exists {
					visited[neighbor] = struct{}{} // set true
					// check if visited node is a border node in subgraph
					if subDist, inSub := distSub[neighbor]; inSub {
						subVisitedCount++
						if depth+1 < subDist {
							return false
						}
						// early abort condition, if all border nodes of subgraph were visited
						if subVisitedCount == subNodeCount {
							return true
						}
					}
					queue = append(queue, neighbor)
				}
			}
		}
		depth++
	}

	return true
}
