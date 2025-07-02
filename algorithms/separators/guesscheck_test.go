package separators

import (
	"bachelor-project/graph"
	"context"
	"testing"
	"time"
)

func TestGuessAndCheck(t *testing.T) {
	grid := [][]int{
		{0, 1, 2, 3},
		{4, -1, -1, 7},
		{8, 9, 10, 11},
	}

	g := graph.NewGraph(3, 4)
	g.Grid = grid
	g.BuildAdjlist()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	subgraphs, ok := GuessAndCheck(g, ctx)
	defer cancel()
	if !ok {
		t.Errorf("Expected GuessAndCheck to return a valid decomposition, but it failed")
	}
	if len(subgraphs) == 0 {
		t.Errorf("Expected non-empty subgraph list, but got 0")
	}
}
