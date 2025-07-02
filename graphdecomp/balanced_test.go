package graphdecomp

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"testing"
)

func TestCheckBalanced(t *testing.T) {
	config.Alpha = 0.5

	testCases := []struct {
		name      string
		adjList   map[int][]int
		expectBal bool
	}{
		{ // one component, one line
			name: "Not balanced, one component",
			adjList: map[int][]int{
				0: {1},
				1: {0, 2},
				2: {1, 3},
				3: {2, 4},
				4: {3}},
			expectBal: false,
		},
		{
			name: "not balanced, two components",
			adjList: map[int][]int{
				0:  {1},
				1:  {0, 2},
				2:  {1, 3},
				3:  {2, 4},
				4:  {3},
				10: {11},
				11: {10},
			},
			expectBal: false,
		},
		{
			name: "balanced, two components",
			adjList: map[int][]int{
				0: {1},
				1: {0, 2},
				2: {1},
				6: {7},
				7: {6, 8},
				8: {7},
			},
			expectBal: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &graph.Graph{AdjList: tc.adjList}
			parent := unionFind(g.AdjList)
			balanced := checkBalanced(parent, len(g.AdjList))
			if balanced != tc.expectBal {
				t.Errorf("Expected balanced=%v, got %v", tc.expectBal, balanced)
			}
		})
	}
}
