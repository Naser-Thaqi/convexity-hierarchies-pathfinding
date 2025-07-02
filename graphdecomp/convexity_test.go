package graphdecomp

import (
	"bachelor-project/graph"
	"context"
	"reflect"
	"testing"
	"time"
)

func makeTestGraph1() *graph.Graph {
	return &graph.Graph{
		AdjList: map[int][]int{
			0:  {1, 4},
			1:  {0, 5},
			3:  {7},
			4:  {0, 5, 8},
			5:  {1, 4, 9},
			7:  {3, 11},
			8:  {4, 9, 12},
			9:  {5, 8, 10, 13},
			10: {9, 11},
			11: {7, 10, 15},
			12: {8, 13},
			13: {9, 12},
			15: {11},
		},
	}
}

func TestExtractBorderNodesOfComponents(t *testing.T) {
	adjList := map[int][]int{
		1:  {5},
		4:  {5},
		5:  {1, 4, 6, 9},
		6:  {5},
		9:  {5},
		11: {15},
		12: {},
		14: {15},
		15: {11, 14},
	}

	parent := map[int]int{
		1:  4,
		4:  4,
		5:  4,
		6:  4,
		9:  4,
		11: 11,
		12: 12,
		14: 11,
		15: 11,
	}

	boundaries := extractBorderNodesOfComponents(adjList, parent)

	if len(boundaries) != 3 {
		t.Errorf("Expected 3 components, got %d", len(boundaries))
	}
	for root, nodes := range boundaries {
		for _, n := range nodes {
			if len(adjList[n]) >= 4 {
				t.Errorf("Node %d in component %d is not boundary (deg >=4)", n, root)
			}
		}
	}
	// 1, 3, 4 expected boundary lengths
	visited := []bool{false, false, false}
	for root := range boundaries {
		if len(boundaries[root]) == 4 {
			visited[2] = true
		}
		if len(boundaries[root]) == 3 {
			visited[1] = true
		}
		if len(boundaries[root]) == 1 {
			visited[0] = true
		}
	}
	for i := range visited {
		if visited[i] == false {
			t.Errorf("Node has not expected length of component border nodes")
		}
	}

}

func TestBFS(t *testing.T) {
	/*
		@   1   @  @  @
		5 , 6 , 7  @  @
		10, @,  12 13 14
		15, @   17 @  19
		20  21  @  23 24
	*/
	adjList := map[int][]int{
		1:  {6},
		5:  {6, 10},
		6:  {1, 5, 7},
		7:  {6, 12},
		10: {5, 15},
		12: {7, 13, 17},
		13: {12, 14},
		14: {13, 19},
		15: {10, 20},
		17: {12},
		19: {14, 24},
		20: {15, 21},
		21: {20},
		23: {24},
		24: {19, 23},
	}

	dist, maxDepth := bfs(adjList, 1)

	if dist[1] != 0 {
		t.Errorf("Distance to start node 0 should be 0, got %d", dist[0])
	}

	if dist[21] != 6 {
		t.Errorf("Distance to node 2 should be 2, got %d", dist[2])
	}

	if maxDepth != 8 {
		t.Errorf("Max depth should be 2, got %d", maxDepth)
	}
}

func TestFilterDistancesForBoundary(t *testing.T) {
	/*
		0  1  @  3
		4  5  6  7
		8  9  10 11
		12 13 14  15
	*/
	originalAdj := map[int][]int{
		0:  {1, 4},
		1:  {0, 5},
		3:  {7},
		4:  {0, 5, 8},
		5:  {1, 4, 6, 9},
		6:  {5, 7, 10},
		7:  {3, 6, 11},
		8:  {4, 9, 12},
		9:  {5, 8, 10, 13},
		10: {6, 9, 11, 14},
		11: {7, 10, 15},
		12: {8, 13},
		13: {9, 12, 14},
		14: {10, 13, 15},
		15: {11, 14},
	}

	distSub := map[int]int{
		0:  0,
		1:  1,
		3:  5,
		4:  1,
		5:  2,
		6:  3,
		7:  4,
		8:  2,
		9:  3,
		10: 4,
		11: 5,
		12: 3,
		13: 4,
		14: 5,
		15: 6,
	}

	removed := []int{5, 9, 10}
	shouldRemain := []int{0, 1, 3, 4, 6, 7, 8, 11, 12, 13, 14, 15}

	filterDistancesForBoundary(originalAdj, distSub)

	for _, node := range removed {
		if _, ok := distSub[node]; ok {
			t.Errorf("Node %d should have been removed from distSub but is still present", node)
		}
	}
	for _, node := range shouldRemain {
		if _, ok := distSub[node]; !ok {
			t.Errorf("Node %d should still be present but was removed", node)
		}
	}

}

func TestIsNodeConvex(t *testing.T) {
	/*
		0  1  @  3
		4  5  @  7
		8  9  10 11
		12 13 @  15
	*/
	originalAdj := map[int][]int{
		0:  {1, 4},
		1:  {0, 5},
		3:  {7},
		4:  {0, 5, 8},
		5:  {1, 4, 9},
		7:  {3, 11},
		8:  {4, 9, 12},
		9:  {5, 8, 10, 13},
		10: {9, 11},
		11: {7, 10, 15},
		12: {8, 13},
		13: {9, 12},
		15: {11},
	}

	// distSub from bfs in a subgraph
	distSub := map[int]int{
		0: 2,
		1: 1,
		4: 1,
		5: 0,
	}
	start := 5
	maxDepth := 2

	if !isNodeConvex(originalAdj, distSub, start, maxDepth) {
		t.Errorf("Expected node to be convex")
	}

	// Separator is NodeID 9,0,1
	distSubFalse := map[int]int{
		4:  3,
		5:  4,
		8:  2,
		12: 1,
		13: 0,
	}
	start = 13
	maxDepth = 4

	if isNodeConvex(originalAdj, distSubFalse, start, maxDepth) {
		t.Errorf("Expected node to be non-convex")
	}
}

func TestCheckConvexity(t *testing.T) {
	/*
		0  1  @  3
		4  5  @  7
		8  9  10 11
		12 13 @  15
	*/
	testcases := []struct {
		name     string
		adjlist  map[int][]int
		expected bool
	}{{
		// separator is 9
		name: "Non-convex",
		adjlist: map[int][]int{
			0:  {1, 4},
			1:  {0, 5},
			3:  {7},
			4:  {0, 5, 8},
			5:  {1, 4},
			7:  {3, 11},
			8:  {4, 12},
			10: {11},
			11: {7, 10, 15},
			12: {8, 13},
			13: {12},
			15: {11},
		},
		expected: false,
	}, {
		// separator is 10
		name: "Convex",
		adjlist: map[int][]int{
			0:  {1, 4},
			1:  {0, 5},
			3:  {7},
			4:  {0, 5, 8},
			5:  {1, 4, 9},
			7:  {3, 11},
			8:  {4, 9, 12},
			9:  {5, 8, 13},
			11: {7, 15},
			12: {8, 13},
			13: {9, 12},
			15: {11},
		},
		expected: true,
	},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			g := makeTestGraph1()

			parent := unionFind(tc.adjlist)

			result := checkConvexity(g, tc.adjlist, parent, ctx)
			if result != tc.expected {
				t.Errorf("Test case '%s' failed: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

}

func TestDegreeFour(t *testing.T) {
	// Graph where all intermediate nodes have degree 4
	g := &graph.Graph{
		AdjList: map[int][]int{
			1: {2},
			2: {1, 3, 4, 5},
			3: {2},
			4: {2},
			5: {2},
		},
	}

	path := []int{1, 2, 3}
	if !degreeFour(g, path) {
		t.Error("Expected degreeFour to return true")
	}

	// Change node 2 to have fewer neighbors
	g.AdjList[2] = []int{1, 3}
	if degreeFour(g, path) {
		t.Error("Expected degreeFour to return false")
	}
}

func TestCheckObservation(t *testing.T) {
	grid := [][]int{
		{10, 11, 12, 13, 14, 15},
		{16, 17, 18, 19, 20, 21},
		{22, 23, 24, 25, 26, 27},
	}
	g := graph.NewGraph(3, 6)
	g.Grid = grid
	g.BuildAdjlist()

	// one direction x and y
	path1 := []int{12, 18, 19, 25}
	coordAdj := map[int]int{
		12: 2,
		18: 8,
		19: 9,
		25: 15,
	}
	if !checkObservation(g, coordAdj, path1) {
		t.Error("Expected monotonePath to return true for monotone path")
	}

	// Mixed direction in x and y
	path2 := []int{11, 17, 18, 19, 13}
	coordAdj = map[int]int{
		11: 1,
		17: 7,
		18: 8,
		19: 9,
		13: 3,
	}
	if checkObservation(g, coordAdj, path2) {
		t.Error("Expected monotonePath to return false for turning path")
	}
}

func TestGetAdjacentNodesOfSeparator(t *testing.T) {
	grid := [][]int{
		{10, 11, 12, 13, 14},
		{16, 17, 18, 19, 20},
		{22, 23, 24, 25, 26},
	}
	g := graph.NewGraph(3, 5)
	g.Grid = grid
	g.BuildAdjlist()
	separator := []int{12, 18, 24}
	parent := map[int]int{
		10: 10, 11: 10, 16: 10, 17: 10, 22: 10, 23: 10,
		13: 13, 14: 13, 19: 13, 20: 13, 25: 13, 26: 13,
	}

	result := getAdjacentNodesOfSeparator(g, separator, parent)
	expected := map[int][]int{
		10: {11, 17, 23},
		13: {13, 19, 25},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCheckObservationAndConvexity(t *testing.T) {
	grid := [][]int{
		{10, 11, 12, 13, 14},
		{16, 17, 18, 19, 20},
		{22, 23, 24, 25, 26},
	}
	g := graph.NewGraph(3, 5)
	g.Grid = grid
	g.BuildAdjlist()
	separator := []int{12, 18, 24}
	copy := g.CopyAdjlist()
	for _, node := range separator {
		graph.RemoveNode(copy, node)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parent := map[int]int{
		10: 10, 11: 10, 16: 10, 17: 10, 22: 10, 23: 10,
		13: 13, 14: 13, 19: 13, 20: 13, 25: 13, 26: 13,
	}
	result := checkObservationAndConvexity(g, copy, parent, separator, ctx)

	if !result {
		t.Errorf("Expected true got false")
	}

}
