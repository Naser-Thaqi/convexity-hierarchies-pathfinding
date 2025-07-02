package algorithms

import "testing"

func TestBreadthFirstSearch(t *testing.T) {
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

	dist := BreadthFirstSearch(adjList, 1, 1)
	if dist != 0 {
		t.Errorf("Distance from node 1 to itself should be 0, got %d", dist)
	}

	// Test: Distance from node 1 to node 20
	dist = BreadthFirstSearch(adjList, 1, 20)
	if dist != 5 {
		t.Errorf("Distance from node 1 to node 20 should be 5, got %d", dist)
	}

	// Test: Distance from node 1 to node 23
	dist = BreadthFirstSearch(adjList, 1, 23)
	if dist != 8 {
		t.Errorf("Distance from node 1 to node 23 should be 8, got %d", dist)
	}

	// Test: Node 25 doesn't exist (unreachable)
	dist = BreadthFirstSearch(adjList, 1, 25)
	if dist != -1 {
		t.Errorf("Distance from node 1 to unreachable node 25 should be -1, got %d", dist)
	}

}
