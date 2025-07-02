package separators

import (
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
)

// Decompose graph by compressing the grid and apply one shortest path.
// After finding a possible solution, decompress the possible solution and recheck in original graph
func TwoShortestPath(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	// build compressed grid
	gc := graph.NewGraph(g.Height/3, g.Width/3)
	gc.Grid = compressGrid(g)
	gc.BuildAdjlist()

	boundaryNodesC := extractBoundaryNodes(gc)
	// Test every boundary node of compressed grid till one offers a valid solution
	for node := range boundaryNodesC {
		prevMap := bfsPaths(gc.AdjList, node)
		paths := createPaths(prevMap, boundaryNodesC, node) // compute shortest paths
		//try every path till one succeeds
	Outerloop:
		for _, candidate := range paths {
			select {
			case <-ctx.Done():
				return nil, false
			default:
				// proceed
			}
			validC := graphdecomp.CheckOneShortestPathBalancedConvex(gc, candidate, ctx) // check validity in compressed grid
			if validC {                                                                  // valid solution in compressed grid
				separator := decompressPath(g, candidate)
				if separator == nil {
					continue
				}
				outerPaths := getOuterPaths(g, decompressBlocks(g, candidate), separator)
				for _, outerPath := range outerPaths {
					if !graphdecomp.CheckBalanced(g, outerPath) {
						continue Outerloop
					}
				}
				if convexComponents, valid := graphdecomp.BalancedDecomposition(g, separator); valid {
					return convexComponents, valid
				}
			}
		}
	}

	return nil, false
}

// get outer paths p1 and p2
func getOuterPaths(g *graph.Graph, separatorO, innerPath []int) map[int][]int {
	// make visited map for deletion
	visited := make(map[int]struct{}, len(innerPath))
	for _, node := range innerPath {
		if _, exists := visited[node]; !exists {
			visited[node] = struct{}{}
		}
	}
	// delete all appearances of inner path nodes in decompressed path
	head := 0
	for head < len(separatorO) {
		if _, exists := visited[separatorO[head]]; exists {
			separatorO[head], separatorO[len(separatorO)-1] = separatorO[len(separatorO)-1], separatorO[head]
			separatorO = separatorO[:len(separatorO)-1]
		} else {
			head++
		}
	}
	visited = nil

	parent := unionFindPaths(g.AdjList, separatorO)
	outerPaths := make(map[int][]int, 2)
	for node, root := range parent {
		outerPaths[root] = append(outerPaths[root], node)
	}
	return outerPaths
}

func unionFindPaths(adjlist map[int][]int, separatorO []int) map[int]int {
	parent := make(map[int]int)

	// initialize map with every node pointing to itself
	for _, node := range separatorO {
		parent[node] = node
	}
	// unite every key with its neighbor nodes
	for _, node := range separatorO {
		for _, neighbor := range adjlist[node] {
			if _, exists := parent[neighbor]; exists {
				union(parent, node, neighbor)
			}
		}
	}

	for node := range parent {
		parent[node] = find(parent, node)
	}

	return parent
}

// Decompresses separator nodes into all original nodes
func decompressBlocks(g *graph.Graph, separatorC []int) []int {
	separatorO := make([]int, 0, len(separatorC)*9)
	surround := [8][2]int{
		//x,y
		{-1, -1}, // top left
		{0, -1},  // top
		{1, -1},  // top right
		{-1, 0},  // left
		{1, 0},   // right
		{-1, 1},  // bottom left
		{0, 1},   // bottom
		{1, 1},   // bottom right
	}

	widthC := g.Width / 3 // compute compressed grid's width

	for i := range separatorC {
		xC, yC := graph.CoordinatesFromNodeID(separatorC[i], widthC) // should be middle node
		xO, yO := xC*3+1, yC*3+1
		separatorO = append(separatorO, g.Grid[yO][xO]) // append middle node
		for _, dir := range surround {                  // append all nodes surrounding the middle node
			nx, ny := xO+dir[0], yO+dir[1] //compute neighbor node (octile)
			separatorO = append(separatorO, g.Grid[ny][nx])
		}
	}

	return separatorO
}

// Decompresses middle path into its original nodes and returns only the separating path that is sandwiched by two others
func decompressPath(g *graph.Graph, separatorC []int) []int {
	// trivial case, would be handled by one-shortest-path heuristic before
	// and can't determine which nodes to choose of 3x3 block as separator nodes
	if len(separatorC) < 2 {
		return nil
	}
	separatorO := make([]int, 0, len(separatorC)*3)
	width := g.Width / 3
	// 0=top,1=bottom,2=left,3=right

	// path is made of 3 nodes per 3x3 blocks, 3 nodes equals one compressed node
	// first cGrid-node, add "left end" nodes that are passable but in non full 3x3 blocks
	x, y := graph.CoordinatesFromNodeID(separatorC[0], width)
	xi, yi := graph.CoordinatesFromNodeID(separatorC[1], width)
	for _, dir := range graph.Directions {
		nx := x + dir[0]
		ny := y + dir[1]
		if nx == xi && ny == yi {
			x, y = x*3+1, y*3+1       //compute position in original grid
			x, y = x-dir[0], y-dir[1] //compute opposite direction
			separatorO = append(separatorO, g.Grid[y][x])
			break
		}
	}

	// one gridC-node equals 3 nodes in decompressed path
	// append middle,next node of current gridC-node and left node of next gridC-node
	// except last node
	for i := range len(separatorC) - 1 {
		x, y := graph.CoordinatesFromNodeID(separatorC[i], width) //gridC coordinates
		separatorO = append(separatorO, g.Grid[1+y*3][1+x*3])

		xi, yi := graph.CoordinatesFromNodeID(separatorC[i+1], width) //gridC coordinates of next node
		for _, dir := range graph.Directions {
			nx := x + dir[0]
			ny := y + dir[1]
			if nx == xi && ny == yi {
				separatorO = append(separatorO, g.Grid[1+y*3+dir[1]][1+x*3+dir[0]])
				separatorO = append(separatorO, g.Grid[1+y*3+dir[1]+dir[1]][1+x*3+dir[0]+dir[0]])
				break
			}
		}
	}

	// last gridC-node and right end
	x, y = graph.CoordinatesFromNodeID(separatorC[len(separatorC)-1], width)
	separatorO = append(separatorO, g.Grid[1+y*3][1+x*3])
	xi, yi = graph.CoordinatesFromNodeID(separatorC[len(separatorC)-2], width)
	for _, dir := range graph.Directions {
		nx := x + dir[0]
		ny := y + dir[1]
		if nx == xi && ny == yi {
			x, y = x*3+1, y*3+1       //compute position in original grid
			x, y = x-dir[0], y-dir[1] //compute opposite direction
			separatorO = append(separatorO, g.Grid[y][x])
			break
		}
	}

	return separatorO
}

// Compresses grid by replacing 3x3 blocks that consists only of passable nodes into one node, otherwise -1 (non-passable node)
func compressGrid(g *graph.Graph) [][]int {
	surround := [8][2]int{
		//x,y
		{-1, -1}, // top left
		{0, -1},  // top
		{1, -1},  // top right
		{-1, 0},  // left
		{1, 0},   // right
		{-1, 1},  // bottom left
		{0, 1},   // bottom
		{1, 1},   // bottom right
	}

	// compute size of compressed grid
	// rounding down leaves out non full 3x3 blocks at the end of row or column
	widthC := g.Width / 3
	heightC := g.Height / 3
	gridC := make([][]int, heightC)

	// coordinates for center nodes of 3x3 block in original grid
	y := 1
	x := 1
	nodeId := 0 // compute new nodeid

	//iterate through original grid, but only the center node of 3x3 blocks
	for yC := range heightC {
		gridC[yC] = make([]int, widthC)
		for xC := range widthC {
			if g.Grid[y][x] != -1 {
				// check all 8 directions of center node
				for _, dir := range surround {
					nx := x + dir[0]
					ny := y + dir[1]
					if g.Grid[ny][nx] == -1 {
						gridC[yC][xC] = -1
						break // one non passable node suffices to declare a 3x3 block as not passable
					}
				}
				// Check if current node in compressed grid got -1 of for loop
				if gridC[yC][xC] != -1 {
					gridC[yC][xC] = nodeId
				}
				// else: current node is not passable,
			} else {
				gridC[yC][xC] = -1
			}
			x += 3   // move to next block in row
			nodeId++ // move to next node in compressed grid
		}
		x = 1  // move to beginning of a row
		y += 3 // move to next row
	}
	return gridC
}
