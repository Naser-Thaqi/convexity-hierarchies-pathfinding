package separators

import (
	"bachelor-project/graph"
	"context"
	"testing"
	"time"
)

func TestFindObstacleComponents(t *testing.T) {
	grid := [][]int{
		{0, -1, 2, 3, 4},
		{5, -1, 7, -1, 9},
		{10, 11, -1, -1, 14},
		{-1, 16, 17, 18, 19},
	}

	g := graph.NewGraph(4, 5)
	g.Grid = grid
	parent := findObstacleComponents(g)
	roots := make(map[int]struct{})
	for node := range parent {
		root := find(parent, node)
		roots[root] = struct{}{}
	}
	if len(roots) != 3 {
		t.Errorf("Expected 3 obstacle components, but found %d", len(roots))
	}
}

func TestGetInnerObstaclesRoots(t *testing.T) {
	grid := [][]int{
		{0, -1, 2, 3, 4},
		{5, -1, 7, -1, 9},
		{10, 11, -1, -1, 14},
		{-1, 16, 17, 18, 19},
	}
	parent := map[int]int{
		1:  1,
		6:  1,
		8:  12,
		12: 12,
		13: 12,
		15: 15,
	}

	g := graph.NewGraph(4, 5)
	g.Grid = grid
	obstacles := getInnerObstaclesRoots(g, parent)
	if len(obstacles) != 1 {
		t.Errorf("Expected 1 obstacle component, but found %d", len(obstacles))
	}
	if _, exists := obstacles[12]; !exists {
		t.Errorf("Expected root 12")
	}
}

func TestGetInnerObstacleNodeSets(t *testing.T) {
	parent := map[int]int{
		1:  1,
		6:  1,
		8:  12,
		12: 12,
		13: 12,
		15: 15,
	}
	innerObstacles := map[int]struct{}{
		12: {},
	}
	innerNodeSet := getInnerObstacleNodeSets(parent, innerObstacles)
	for _, set := range innerNodeSet {
		if len(set) != 3 {
			t.Errorf("Expected 3 nodes, but found %d", len(set))
		}
	}
}

func TestGetSeparatorOfObstacle(t *testing.T) {
	grid := [][]int{
		{0, 1, -1, 3, 4},
		{5, 6, 7, 8, 9},
		{10, -1, -1, -1, 14},
		{15, 16, -1, 18, 19},
		{20, 21, 22, 23, 24},
	}
	g := graph.NewGraph(5, 5)
	g.Grid = grid
	obstacleNodeSet := []int{11, 12, 13, 17} //nodeids of obstacle nodes

	centralBoundaryNodesCoords := getCentralBoundaryNodesCoords(g, obstacleNodeSet)
	separators := getSeparatorOfObstacle(g, centralBoundaryNodesCoords)

	expected := map[int]struct{}{
		7:  {},
		22: {},
		10: {},
		14: {},
	}

	for _, node := range separators {
		_, exists := expected[node]
		if !exists {
			t.Errorf("Node is not in expected Solution")
		}
	}
}

func TestHoleCutting(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4},
		{5, -1, 7, 8, 9},
		{10, 11, 12, 13, 14},
		{15, -1, 17, -1, 19},
		{20, 21, 22, 23, 24},
	}
	g := graph.NewGraph(5, 5)
	g.Grid = grid
	g.BuildAdjlist()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	components, ok := HoleCutting(g, ctx)
	defer cancel()

	if !ok {
		t.Errorf("HoleCutting failed: expected successful decomposition")
	}
	if len(components) != 9 {
		t.Errorf("HoleCutting returned wrong number of components: expected 9, got %d", len(components))
	}
}
