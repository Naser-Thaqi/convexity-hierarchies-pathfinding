package algorithms

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"fmt"
	"strings"
	"testing"
)

func TestBuildConvexHierarchy(t *testing.T) {
	grid := [][]int{
		{-1, -1, 2, 3, -1},
		{-1, 6, 7, 8, 9},
		{10, 11, 12, 13, 14},
		{15, 16, 17, 18, 19},
		{-1, 21, 22, 23, -1},
	}
	g := graph.NewGraph(5, 5)
	g.Grid = grid
	g.BuildAdjlist()

	original := config.KaFFPaPath
	config.KaFFPaPath = "../KaHIP/build/kaffpa"
	BuildConvexHierarchy(g)
	config.KaFFPaPath = original

	if g.Childs == nil {
		t.Error("Expected non-nil Childs after hierarchy build")
	}
	// PrintHierarchyAdjLists(g, 0)
	// PrintHierarchyGrids(g, 0)
}

func PrintHierarchyGrids(g *graph.Graph, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%sGraph at level %d\n", indent, level)

	if g.Grid != nil {
		fmt.Printf("%sGrid:\n", indent)
		for _, row := range g.Grid {
			fmt.Printf("%s  %v\n", indent, row)
		}
	} else {
		fmt.Printf("%sGrid: <nil>\n", indent)
	}

	for _, child := range g.Childs {
		PrintHierarchyGrids(child, level+1)
	}
}

func PrintHierarchyAdjLists(g *graph.Graph, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%sGraph at level %d: %d nodes\n", indent, level, len(g.AdjList))
	for k, v := range g.AdjList {
		fmt.Printf("%s  %d: %v\n", indent, k, v)
	}
	for _, child := range g.Childs {
		PrintHierarchyAdjLists(child, level+1)
	}
}

func TestFindSmallestConvexComponent(t *testing.T) {
	parent := &graph.Graph{
		AdjList: map[int][]int{
			0: {},
			1: {2},
			2: {1},
			3: {},
			4: {},
		},
	}
	child1 := &graph.Graph{
		AdjList: map[int][]int{
			3: {},
		},
	}
	child2 := &graph.Graph{
		AdjList: map[int][]int{
			0: {0},
			1: {2},
			2: {1},
		},
	}
	child3 := &graph.Graph{
		AdjList: map[int][]int{
			0: {0},
		},
	}
	child4 := &graph.Graph{
		AdjList: map[int][]int{
			1: {2},
			2: {1},
		},
	}
	parent.Childs = []*graph.Graph{child1, child2}
	child2.Childs = []*graph.Graph{child3, child4}
	/*
			parent
		child1	    child2
				child3  child4
	*/

	result := FindSmallestConvexComponent(parent, 1, 2)
	if result == nil || len(result.AdjList) != 2 {
		t.Error("Expected to find convex component containing 1 and 2")
	}

	result = FindSmallestConvexComponent(parent, 1, 99)
	if result != nil {
		t.Error("Expected nil for non-existent end node")
	}
}
