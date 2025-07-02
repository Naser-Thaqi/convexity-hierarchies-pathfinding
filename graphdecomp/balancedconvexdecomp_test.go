package graphdecomp

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"context"
	"testing"
	"time"
)

func makeTestGraph() *graph.Graph {
	g := graph.NewGraph(3, 3)
	g.Grid = [][]int{
		{0, 1, -1},
		{2, -1, 3},
		{-1, 4, 5},
	}
	g.Width = 3
	g.Height = 3
	g.AdjList = map[int][]int{
		0: {1, 2},
		1: {0},
		2: {0},
		3: {5},
		4: {5},
		5: {3, 4},
	}

	return g
}

func TestCheckOneShortestPathBalancedConvex(t *testing.T) {
	grid := [][]int{
		{11, -1, 13, 14, 15},
		{17, 18, 19, 20, 21},
		{23, -1, 25, 26, 27},
		{29, 30, 31, 32, 33},
		{35, 36, 37, -1, 39},
	}
	g := graph.NewGraph(5, 5)
	g.Grid = grid
	g.BuildAdjlist()
	separator := []int{13, 19, 25, 31, 37}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if CheckOneShortestPathBalancedConvex(g, []int{}, ctx) {
		t.Errorf("Expected false for empty separator")
	}

	original := config.Alpha
	config.Alpha = 1
	res := CheckOneShortestPathBalancedConvex(g, separator, ctx)
	if res == false {
		t.Errorf("Expected true for separator %v with alpha=1, got false", separator)
	}

	config.Alpha = 0.2
	res = CheckOneShortestPathBalancedConvex(g, separator, ctx)
	if res == true {
		t.Errorf("Expected false for separator %v with alpha=0.2, got true", separator)
	}
	config.Alpha = original
}

func TestBalancedConvexDecomposition(t *testing.T) {
	g := makeTestGraph()

	separator := []int{}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	subgraphs, ok := BalancedConvexDecomposition(g, separator, ctx)
	if !ok && subgraphs != nil {
		t.Errorf("Expected decomposition to fail and return nil subgraphs")
	}

	separator = append(separator, 1)
	subgraphs, ok = BalancedConvexDecomposition(g, separator, ctx)
	if ok {
		if len(subgraphs) == 0 {
			t.Errorf("Expected at least one subgraph when decomposition successful")
		}
	}
}

func TestDecomposeGraph(t *testing.T) {
	candidates := []int{1, 2} // separator nodes
	for _, c := range candidates {
		g := makeTestGraph()
		graph.RemoveNode(g.AdjList, c)

		parent := unionFind(g.AdjList)

		subgraphs := decomposeGraph(g, parent)

		if len(subgraphs) == 0 {
			t.Fatal("Expected at least one subgraph decomposition")
		}

		for _, sg := range subgraphs {
			//fmt.Println("subgraph", c, sg.Grid)
			if sg.Height <= 0 || sg.Width <= 0 {
				t.Errorf("Subgraph has invalid dimensions: Height=%d, Width=%d", sg.Height, sg.Width)
			}
		}
	}
}

func TestOneShortestPathDecomposition(t *testing.T) {
	grid := [][]int{
		{11, -1, 13, 14, 15},
		{17, 18, 19, 20, 21},
		{23, 24, 25, 26, 27},
		{29, 30, 31, 32, 33},
		{35, 36, 37, -1, 39},
	}
	g := graph.NewGraph(5, 5)
	g.Grid = grid
	g.BuildAdjlist()
	separator := []int{13, 19, 25, 31, 37}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	subgraphs, ok := OneShortestPathBalancedConvexDecomposition(g, separator, ctx)
	if !ok {
		if len(subgraphs) != 2 {
			t.Errorf("Expected two subgraphs when decomposition fails, but got %d", len(subgraphs))
		}
	}
	for _, node := range []int{13, 37} {
		if _, exists := subgraphs[0].AdjList[node]; !exists {
			if _, exists := subgraphs[1].AdjList[node]; !exists {
				t.Errorf("Expected node %d to appear in at least one subgraph", node)
			}
		}
	}
}

func TestDecomposeInputComponents(t *testing.T) {
	grid := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	g := graph.NewGraph(3, 3)
	g.Grid = grid
	g.BuildAdjlist()

	subgraphs, ok := DecomposeInputComponents(g)
	if ok {
		t.Errorf("Expected no subgraphs, but got %d", len(subgraphs))
	}

	grid = [][]int{
		{1, -1, 3},
		{4, -1, 6},
		{7, -1, 9},
	}
	g = graph.NewGraph(3, 3)
	g.Grid = grid
	g.BuildAdjlist()
	subgraphs, ok = DecomposeInputComponents(g)
	if !ok {
		t.Errorf("Expected two subgraphs, but got %d", len(subgraphs))
	}

}
