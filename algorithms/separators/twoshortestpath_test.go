package separators

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"context"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestCompressGrid(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10, 11},
		{12, 13, 14, 15, 16, 17},
		{18, 19, 20, 21, 22, 23},
		{24, 25, 26, 27, 28, 29},
		{30, -1, 32, 33, 34, 35},
		{36, 37, 38, 39, 40, 41},
	}
	expected := [][]int{
		{0, 1},
		{-1, 3},
	}

	g := graph.NewGraph(7, 6)
	g.Grid = grid
	cGrid := compressGrid(g)
	widthC := g.Width / 3
	heightC := g.Height / 3
	for y := range heightC {
		for x := range widthC {
			if cGrid[y][x] != expected[y][x] {
				t.Errorf("mismatch at compressed cell (%d, %d): got %d, want %d", y, x, cGrid[y][x], expected[y][x])
			}
		}
	}
}

func TestDecompressPath(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8},
		{9, 10, 11, 12, 13, 14, 15, 16, 17},
		{18, 19, 20, 21, 22, 23, 24, 25, 26},
		{27, 28, 29, 30, 31, 32, 33, 34, 35},
		{36, 37, 38, 39, 40, 41, 42, 43, 44},
		{45, -1, 47, 48, 49, 50, 51, 52, 53},
		{54, 55, 56, 57, 58, 59, 60, 61, 62},
	}

	g := graph.NewGraph(7, 9)
	g.Grid = grid
	// cGrid := [][]int{
	// 	{0, 1, 2},
	// 	{-1, 4, 5},
	// }

	tests := []struct {
		name       string
		separatorC []int
		expected   []int
	}{
		{
			name:       "horizontal path, no holes",
			separatorC: []int{0, 1, 2},
			expected:   []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
		},
		{
			name:       "vertical path, with holes",
			separatorC: []int{1, 4},
			expected:   []int{4, 13, 22, 31, 40, 49},
		},
	}

	for _, tc := range tests {
		separator := decompressPath(g, tc.separatorC)
		if !reflect.DeepEqual(separator, tc.expected) {
			t.Errorf("%s: got %v, want %v", tc.name, separator, tc.expected)
		}
	}
}

func TestDecompressBlocks(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8},
		{9, 10, 11, 12, 13, 14, 15, 16, 17},
		{18, 19, 20, 21, 22, 23, 24, 25, 26},
		{27, 28, 29, 30, 31, 32, 33, 34, 35},
		{36, 37, 38, 39, 40, 41, 42, 43, 44},
		{45, -1, 47, 48, 49, 50, 51, 52, 53},
		{54, 55, 56, 57, 58, 59, 60, 61, 62},
	}

	g := graph.NewGraph(7, 9)
	g.Grid = grid
	// cGrid := [][]int{
	// 	{0, 1, 2},
	// 	{-1, 4, 5},
	// }

	separatorC := []int{1, 4}
	expected := []int{3, 4, 5, 12, 13, 14, 21, 22, 23, 30, 31, 32, 39, 40, 41, 48, 49, 50}
	separator := decompressBlocks(g, separatorC)
	sort.Ints(separator)
	if !reflect.DeepEqual(separator, expected) {
		t.Errorf(": got %v, want %v", separator, expected)
	}
}

func TestGetOuterPaths(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8},
		{9, 10, 11, 12, 13, 14, 15, 16, 17},
		{18, 19, 20, 21, 22, 23, 24, 25, 26},
		{27, 28, 29, 30, 31, 32, 33, 34, 35},
		{36, 37, 38, 39, 40, 41, 42, 43, 44},
		{45, -1, 47, 48, 49, 50, 51, 52, 53},
		{54, 55, 56, 57, 58, 59, 60, 61, 62},
	}

	g := graph.NewGraph(7, 9)
	g.Grid = grid
	g.BuildAdjlist()
	// cGrid := [][]int{
	// 	{0, 1, 2},
	// 	{-1, 4, 5},
	// }
	separatorO := []int{3, 4, 5, 12, 13, 14, 21, 22, 23, 30, 31, 32, 39, 40, 41, 48, 49, 50}
	innerPath := []int{4, 13, 22, 31, 40, 49}
	expected1 := []int{3, 12, 21, 30, 39, 48}
	expected2 := []int{5, 14, 23, 32, 41, 50}
	outerPaths := getOuterPaths(g, separatorO, innerPath)
	keys := []int{}
	for key := range outerPaths {
		sort.Ints(outerPaths[key])
		keys = append(keys, key)
	}
	if len(outerPaths) != 2 {
		t.Errorf("expected at least 2 paths, got %d", len(outerPaths))
	}
	if !reflect.DeepEqual(outerPaths[keys[0]], expected1) {
		if !reflect.DeepEqual(outerPaths[keys[0]], expected2) {
			t.Errorf("expected path 1 to one expected path, got %d", outerPaths[0])
		}
	}
	if !reflect.DeepEqual(outerPaths[keys[1]], expected1) {
		if !reflect.DeepEqual(outerPaths[keys[1]], expected2) {
			t.Errorf("expected path 2 to one expected path, got %d", outerPaths[1])
		}
	}

}

func TestTwoShortestPath(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8},
		{9, 10, 11, 12, 13, 14, 15, 16, 17},
		{18, 19, 20, 21, 22, 23, 24, 25, 26},
		{27, 28, 29, 30, 31, 32, 33, 34, 35},
		{36, 37, 38, 39, 40, 41, 42, 43, 44},
		{45, 46, 47, 48, 49, 50, 51, 52, 53},
	}
	// cGrid := [][]int{
	// 	{0, 1, 2},
	// 	{-1, 4, 5},
	// }

	g := graph.NewGraph(6, 9)
	g.Grid = grid
	g.BuildAdjlist()
	original := config.Alpha
	config.Alpha = 2.0 / 3.0
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	components, valid := TwoShortestPath(g, ctx)

	defer cancel()

	if !valid {
		t.Fatalf("expected true, got false")
	}

	if len(components) < 2 {
		t.Fatalf("expected at least 2 components, got %d", len(components))
	}

	for i, comp := range components {
		if comp == nil {
			t.Errorf("component %d is nil", i)
			continue
		}
		if len(comp.AdjList) == 0 {
			t.Errorf("component %d has empty or nil adjacency list", i)
		}
	}
	config.Alpha = original
}
