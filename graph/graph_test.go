package graph

import (
	"reflect"
	"testing"
)

func TestAddEdge(t *testing.T) {
	g := NewGraph(3, 3)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)

	expected := []int{2, 3}
	if !reflect.DeepEqual(g.AdjList[1], expected) {
		t.Errorf("AddEdge failed: expected %v, got %v", expected, g.AdjList[1])
	}
}

func TestRemoveNode(t *testing.T) {
	adjList := map[int][]int{
		1: {2, 3},
		2: {1},
		3: {1},
	}
	RemoveNode(adjList, 1)

	if _, exists := adjList[1]; exists {
		t.Errorf("RemoveNode failed: node 1 still exists in adjList")
	}
	if len(adjList[2]) > 0 || len(adjList[3]) > 0 {
		t.Errorf("RemoveNode failed: node 1 still exists in neighbors")
	}
}

func TestRestoreNode(t *testing.T) {
	adjList := map[int][]int{
		1: {2, 3},
		2: {1},
		3: {1},
	}
	restoredAdj := map[int][]int{
		2: {},
		3: {},
	}
	RestoreNode(adjList, restoredAdj, 1)
	if _, exists := restoredAdj[1]; !exists {
		t.Errorf("RestoreNode failed: node 1 is not in adjList")
	}
	if len(restoredAdj[2]) == 0 || len(restoredAdj[3]) == 0 {
		t.Errorf("RestoreNode failed: node 1 is not in neighbors")
	}
}

func TestNodeID(t *testing.T) {
	/*
		0, -1, 2, 3, 4, 5,
		-1, 7, 8, 9, 10, -1
	*/
	width := 6
	x, y := 4, 1 // coordinates are index position
	expected := 10
	if result := NodeID(x, y, width); result != expected {
		t.Errorf("NodeID failed: expected %d, got %d", expected, result)
	}
}

func TestCoordinatesFromNodeID(t *testing.T) {
	/*
		0, -1, 2, 3, 4, 5,
		-1, 7, 8, 9, 10, -1
	*/
	width := 6
	id := 10
	expectedX, expectedY := 4, 1
	x, y := CoordinatesFromNodeID(id, width)
	if x != expectedX || y != expectedY {
		t.Errorf("CoordinatesFromNodeID failed: expected (%d,%d), got (%d,%d)", expectedX, expectedY, x, y)
	}
}

func TestBuildAdjlist(t *testing.T) {
	g := NewGraph(3, 3)
	g.Grid = [][]int{
		{0, 1, 2},
		{3, -1, 5},
		{-1, 7, -1},
	}
	g.BuildAdjlist()

	expected := map[int][]int{
		0: {1, 3},
		1: {0, 2},
		2: {1, 5},
		3: {0},
		5: {2},
		7: {},
	}

	for k, v := range expected {
		if !reflect.DeepEqual(g.AdjList[k], v) {
			t.Errorf("BuildAdjlist failed at node %d: expected %v, got %v", k, v, g.AdjList[k])
		}
	}
}

func TestCopyAdjlist(t *testing.T) {
	original := &Graph{
		AdjList: map[int][]int{
			1: {2, 3},
			2: {1},
		},
	}

	copy := original.CopyAdjlist()

	if !reflect.DeepEqual(original.AdjList, copy) {
		t.Errorf("CopyAdjlist failed: expected %v, got %v", original.AdjList, copy)
	}

	// Modify the copy and check whether or not the original is affected
	copy[1][0] = 99
	if original.AdjList[1][0] == 99 {
		t.Errorf("CopyAdjlist failed: changes to copy affected the original")
	}
}

func TestLoadGraphFromFile(t *testing.T) {
	graph := LoadGraphFromFile("graph_test.txt")
	if graph == nil {
		t.Fatal("LoadGraphFromFile returned nil")
	}

	expectedNodes := []int{0, 1, 2, 3, 5, 7}

	// Check for all expected nodes
	for _, node := range expectedNodes {
		if _, exists := graph.AdjList[node]; !exists {
			t.Errorf("Expected node %d in adjacency list, but it was missing", node)
		}
	}

	// Check that no unexpected nodes exist
	if len(graph.AdjList) != len(expectedNodes) {
		t.Errorf("Adjacency list contains unexpected nodes: got %d, expected %d", len(graph.AdjList), len(expectedNodes))
	}

	expectedGrid := [][]int{
		{0, 1, 2},
		{3, -1, 5},
		{-1, 7, -1},
	}
	if !reflect.DeepEqual(graph.Grid, expectedGrid) {
		t.Errorf("Grid does not match expected grid.\nExpected: %v\nGot: %v", expectedGrid, graph.Grid)
	}

	// Expected adjancy list
	expectedAdj := map[int][]int{
		0: {1, 3},
		1: {0, 2},
		2: {1, 5},
		3: {0},
		5: {2},
		7: {},
	}

	for node, expectedNeighbors := range expectedAdj {
		actualNeighbors := graph.AdjList[node]
		if !reflect.DeepEqual(actualNeighbors, expectedNeighbors) {
			t.Errorf("For node %d, expected neighbors %v, got %v", node, expectedNeighbors, actualNeighbors)
		}
	}
}
