package separators

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"context"
	"testing"
	"time"
)

func TestRowColumn(t *testing.T) {
	original := config.Alpha
	testCases := []struct {
		name       string
		grid       [][]int
		width      int
		height     int
		expected   bool
		separators []int
		alpha      float64
	}{
		{
			name:     "Empty",
			grid:     [][]int{},
			width:    0,
			height:   0,
			expected: false,
			alpha:    1,
		},
		{
			name: "Possible",
			grid: [][]int{
				{0, 1, 2, 4, 5},
				{6, -1, 8, -1, 10},
				{11, 12, 13, -1, 15},
				{16, -1, 18, -1, 20},
			},
			width:      5,
			height:     4,
			expected:   true,
			separators: []int{4},
			alpha:      1,
		},
		{ //2x2 grid, first and last row/column not passable
			name: "Grid too smal",
			grid: [][]int{
				{0, 1},
				{2, 3},
			},
			width:      2,
			height:     2,
			expected:   false,
			separators: []int{},
			alpha:      1,
		},
		{ //3x3 grid, every node is passable
			name: "No valid balancing",
			grid: [][]int{
				{0, 1, 2},
				{3, 4, 5},
				{6, 7, 8},
			},
			height:     3,
			width:      3,
			expected:   false,
			separators: []int{},
			alpha:      0.2,
		},
	}

	// test every testcase
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := graph.NewGraph(tc.height, tc.width)
			g.Grid = tc.grid
			g.BuildAdjlist()
			config.Alpha = tc.alpha

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

			components, ok := RowColumn(g, ctx)
			defer cancel()

			if ok != tc.expected {
				t.Errorf("Test %s failed: expected %v, got %v", tc.name, tc.expected, ok)
			}

			if ok && len(components) == 0 {
				t.Errorf("Test %s: expected components but got none", tc.name)
			}
			if ok {
				for _, subg := range components {
					for _, sep := range tc.separators {
						if _, exists := subg.AdjList[sep]; exists {
							t.Errorf("Test %s: separator %d still in subgraph", tc.name, sep)
						}
					}
				}
			}
		})
	}
	config.Alpha = original
}
