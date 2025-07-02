package separators

import (
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
)

// Decompose Graph along orthogonals of obstacles
func HoleCutting(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	parent := findObstacleComponents(g)                                       //get parent map
	innerObstacles := getInnerObstaclesRoots(g, parent)                       //return set of roots that are inner components
	innerObstacleNodeSets := getInnerObstacleNodeSets(parent, innerObstacles) //return a map key = root, and values are nodes connected to same component

	// empty for g.c.
	parent = nil
	innerObstacles = nil

	separator := []int{}
	// gather all orthogonal separator node of each in inner obstacle component
	for _, set := range innerObstacleNodeSets {
		centralBoundaryNodes := getCentralBoundaryNodesCoords(g, set)
		separator = append(separator, getSeparatorOfObstacle(g, centralBoundaryNodes)...)
	}
	// give ctx to balanced convex decomposition, to abort during costly bfs convexity check
	return graphdecomp.BalancedConvexDecomposition(g, separator, ctx)
}

// returns nodes orthogonal to the obstacle till the orthogonal nodes reach another obstacle or end of grid
// appends the original nodeids into separator set
func getSeparatorOfObstacle(g *graph.Graph, centralBoundaryNodesCoords [8]int) []int {
	separators := []int{}

	nx, ny := centralBoundaryNodesCoords[0], centralBoundaryNodesCoords[1]
	end := false

	//top node
	for !end {
		ny--
		if ny >= 0 {
			if g.Grid[ny][nx] != -1 {
				separators = append(separators, g.Grid[ny][nx])
			} else {
				end = true
			}
		} else {
			end = true
		}
	}

	nx, ny = centralBoundaryNodesCoords[2], centralBoundaryNodesCoords[3]

	//bottom node
	end = false
	for !end {
		ny++
		if ny < g.Height {
			if g.Grid[ny][nx] != -1 {
				separators = append(separators, g.Grid[ny][nx])
			} else {
				end = true
			}
		} else {
			end = true
		}
	}

	nx, ny = centralBoundaryNodesCoords[4], centralBoundaryNodesCoords[5]
	//left node
	end = false
	for !end {
		nx--
		if nx >= 0 {
			if g.Grid[ny][nx] != -1 {
				separators = append(separators, g.Grid[ny][nx])
			} else {
				end = true
			}
		} else {
			end = true
		}
	}

	nx, ny = centralBoundaryNodesCoords[6], centralBoundaryNodesCoords[7]
	//right node
	end = false
	for !end {
		nx++
		if nx < g.Width {
			if g.Grid[ny][nx] != -1 {
				separators = append(separators, g.Grid[ny][nx])
			} else {
				end = true
			}
		} else {
			end = true
		}
	}

	return separators
}

// returns central node of each boundary
func getCentralBoundaryNodesCoords(g *graph.Graph, obstacleNodeSet []int) [8]int {
	//first node as reference point for further computation
	xLeft, yTop := graph.CoordinatesFromNodeID(obstacleNodeSet[0], g.Width)
	xRight, yLow := graph.CoordinatesFromNodeID(obstacleNodeSet[0], g.Width)

	//compute width and height of obstacle
	for i := range obstacleNodeSet {
		nodeId := obstacleNodeSet[i]
		x, y := graph.CoordinatesFromNodeID(nodeId, g.Width)
		if yTop > y {
			yTop = y
		}
		if yLow < y {
			yLow = y
		}
		if xLeft > x {
			xLeft = x
		}
		if xRight < x {
			xRight = x
		}
	}

	pmx := (xLeft + xRight) / 2 //middle x coordinate
	pmy := (yLow + yTop) / 2    //middle y coordinate

	centralBoundaryNodesCoords := [8]int{
		pmx, g.Height, //top node
		pmx, -1, //bottom node
		g.Width, pmy, //left node
		-1, pmy, //right
	}

	// find for every obstacle the appropriate nodes that contain pmx, pmy
	for i := range obstacleNodeSet {
		nodeId := obstacleNodeSet[i]
		x, y := graph.CoordinatesFromNodeID(nodeId, g.Width)
		if x == pmx && y < centralBoundaryNodesCoords[1] {
			centralBoundaryNodesCoords[1] = y
		}
		if x == pmx && y > centralBoundaryNodesCoords[3] {
			centralBoundaryNodesCoords[3] = y
		}
		if x < centralBoundaryNodesCoords[4] && y == pmy {
			centralBoundaryNodesCoords[4] = x
		}
		if x > centralBoundaryNodesCoords[6] && y == pmy {
			centralBoundaryNodesCoords[6] = x
		}
	}

	return centralBoundaryNodesCoords
}

// return a map where roots are keys and they store every node connected to them
func getInnerObstacleNodeSets(parent map[int]int, innerObstacles map[int]struct{}) map[int][]int {
	obstacleNodeSets := make(map[int][]int)

	for key, val := range parent {
		// check for innerObstacles
		if _, exists := innerObstacles[val]; exists {
			// gather nodes connected to root node
			obstacleNodeSets[val] = append(obstacleNodeSets[val], key)
		}
	}

	return obstacleNodeSets
}

// Returns Set of roots that are inner obstacles
func getInnerObstaclesRoots(g *graph.Graph, parent map[int]int) map[int]struct{} {
	obstacles := make(map[int]struct{})

	for key, val := range parent {
		if key == val {
			obstacles[key] = struct{}{}
		}
	}

	//first and last row
	for x := range g.Width {
		if g.Grid[0][x] == -1 {
			root := find(parent, graph.NodeID(x, 0, g.Width))
			delete(obstacles, root)
		}
		if g.Grid[g.Height-1][x] == -1 {
			root := find(parent, graph.NodeID(x, g.Height-1, g.Width))
			delete(obstacles, root)
		}
	}
	//first and last column except first and last cell
	for y := 1; y < g.Height-1; y++ {
		if g.Grid[y][0] == -1 {
			root := find(parent, graph.NodeID(0, y, g.Width))
			delete(obstacles, root)
		}
		if g.Grid[y][g.Width-1] == -1 {
			root := find(parent, graph.NodeID(g.Width-1, y, g.Width))
			delete(obstacles, root)
		}
	}

	return obstacles
}

// returns map of nodeids which stores root node (connected components)
// find obstacle components via union find
// every nodeid gets new nodeid in parent map, for calculating position
func findObstacleComponents(g *graph.Graph) map[int]int {
	parent := make(map[int]int)

	//iterate through whole grid, initialize parent map
	for y := range g.Height {
		for x := range g.Width {
			//filter for obstacle nodes
			if g.Grid[y][x] == -1 {
				//give nodeid based on current grid
				nodeID := graph.NodeID(x, y, g.Width)
				parent[nodeID] = nodeID
			}
		}
	}

	for y := range g.Height {
		for x := range g.Width {
			nodeID := g.Grid[y][x]
			if nodeID == -1 {
				for _, dir := range graph.Directions {
					nx := x + dir[0]
					ny := y + dir[1]
					if nx < 0 || ny < 0 || nx >= g.Width || ny >= g.Height {
						continue
					}
					if g.Grid[ny][nx] == -1 {
						//compute nodeids of current obstacle and neighbor and combine them
						union(parent, graph.NodeID(x, y, g.Width), graph.NodeID(nx, ny, g.Width))
					}
				}
			}
		}
	}
	return parent
}
