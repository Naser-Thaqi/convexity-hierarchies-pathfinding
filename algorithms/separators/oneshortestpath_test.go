package separators

import (
	"bachelor-project/graph"
	"context"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestExtractBoundaryNodes(t *testing.T) {
	grid := [][]int{
		{-1, -1, 2, 3, 4, 5},
		{-1, 7, 8, 9, 10, 11},
		{12, 13, 14, 15, 16, 17},
		{18, 19, 20, -1, 22, 23},
		{-1, 25, 26, 27, 28, 29},
	}
	g := graph.NewGraph(5, 6)
	g.Grid = grid
	g.BuildAdjlist()
	expected := []int{2, 3, 4, 5, 7, 11, 12, 15, 17, 18, 20, 22, 23, 25, 26, 27, 28, 29}
	boundaryNodes := extractBoundaryNodes(g)
	if len(boundaryNodes) != len(expected) {
		t.Errorf("Expected %d boundary nodes, but got %d", len(expected), len(boundaryNodes))
	}
	for _, exp := range expected {
		found := slices.Contains(boundaryNodes, exp)
		if !found {
			t.Errorf("Expected node %d in boundaryNodes, but not found", exp)
		}
	}
}

func TestBfsPaths(t *testing.T) {
	adjlist := map[int][]int{
		0: {1, 2},
		1: {0, 3},
		2: {0, 3},
		3: {1, 2},
	}

	result := bfsPaths(adjlist, 0)
	expected := map[int]int{
		1: 0,
		2: 0,
		3: 1, // or 2, both would be correct, depending on traversal
	}

	for k, v := range expected {
		if result[k] != v && result[k] != 2 {
			t.Errorf("For node %d expected previous %d or 2, got %d", k, v, result[k])
		}
	}
}

func TestCreatePaths(t *testing.T) {
	// backtrace map from BFS
	prevMap := map[int]int{
		4: 3,
		3: 1,
		1: 0,
	}
	boundary := []int{4, 0, 5}
	start := 0

	result := createPaths(prevMap, boundary, start)
	expected := map[int][]int{
		4: {4, 3, 1, 0},
		0: {0},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected path %v, got %v", expected, result)
	}
}

func TestReducePaths(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10, 11},
		{12, 13, 14, 15, 16, 17},
		{18, -1, -1, -1, -1, 23},
		{24, 25, 26, 27, 28, 29},
	}

	g := graph.NewGraph(5, 6)
	g.Grid = grid
	g.BuildAdjlist()
	path := []int{18, 12, 13, 14, 15, 16, 17, 23}
	paths := map[int][]int{
		18: path,
	}
	reduced := reducePaths(g, paths)
	expected := [][]int{{18, 23}}

	if !reflect.DeepEqual(reduced, expected) {
		t.Errorf("Expected %v, got %v", expected, reduced)
	}
}

func TestOneShortestPath(t *testing.T) {
	g := graph.NewGraph(3, 3)
	g.Grid = [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}
	g.BuildAdjlist()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	components, valid := OneShortestPath(g, ctx)
	defer cancel()
	if !valid {
		t.Log("Expected valid decomposition for trivial 3x3")
	} else {
		if len(components) != 2 {
			t.Errorf("Expected 2 components, got %d", len(components))
		}
	}

	g2 := graph.NewGraph(2, 2)
	g2.Grid = [][]int{
		{15, 16},
		{-1, 21},
	}
	g2.BuildAdjlist()

	//single node path 16:16, should be a valid separator
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	components2, valid2 := OneShortestPath(g2, ctx)
	defer cancel()
	if !valid2 {
		t.Error("Expected valid decomposition for custom 2x2 grid")
	} else {
		if len(components2) != 2 {
			t.Errorf("Expected 2 components, got %d", len(components2))
		}
	}
}
