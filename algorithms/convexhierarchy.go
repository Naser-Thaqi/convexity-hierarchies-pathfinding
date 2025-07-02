package algorithms

import (
	"bachelor-project/algorithms/separators"
	"bachelor-project/config"
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
)

type sep func(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool)

// Create convex subgraphes
func BuildConvexHierarchy(g *graph.Graph) {
	// prevent undefined behavior, by decomposing graph if it already has components
	childs, ok := graphdecomp.DecomposeInputComponents(g)

	if ok {
		g.Childs = childs
	} else {
		g.Childs = pipeline(g)
	}

	stack := []*graph.Graph{}
	for i := len(g.Childs) - 1; i >= 0; i-- {
		stack = append(stack, g.Childs[i])
	}

	// Build Tree preorder iterative
	for len(stack) > 0 {
		c := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		c.Childs = pipeline(c)
		c.Grid = nil

		for i := len(c.Childs) - 1; i >= 0; i-- {
			stack = append(stack, c.Childs[i])
		}
	}
}

// pipeline for using several heuristics to compute convex subgraphs
func pipeline(g *graph.Graph) []*graph.Graph {
	if len(g.AdjList) < 3 {
		return nil
	}
	// separators.KaFFPaSeparator,  separators.OneShortestPath, separators.TwoShortestPath, separators.RowColumn, separators.HoleCutting
	sepFuncs := []sep{separators.OneShortestPath, separators.TwoShortestPath, separators.RowColumn, separators.HoleCutting}
	/*
	   pipeline:
	   kaffpa
	   separating shortest path
	   two separating shortest path
	   each row and column
	   hole cutting
	*/
	// try every function (heuristic) in array
	for _, sepFunc := range sepFuncs {
		ctx, cancel := context.WithTimeout(context.Background(), config.Time)

		resultChan := make(chan []*graph.Graph, 1) // channel for result
		go func(f sep) {
			graphs, ok := f(g, ctx)
			if ok {
				resultChan <- graphs
			} else {
				resultChan <- nil
			}
		}(sepFunc)

		select {
		case res := <-resultChan:
			cancel()
			// return only positive result
			if res != nil {
				return res
			} // else: try another heuristic
		case <-ctx.Done():
			cancel()
		}
	}
	// no heuristic found a valid alpha balanced convex decomposition
	return nil
}

// Returns smallest convex component that has start and end -node in a single adjacency list
func FindSmallestConvexComponent(g *graph.Graph, startNode, endNode int) *graph.Graph {
	// Check if key (node) exists
	_, startExists := g.AdjList[startNode]
	_, endExists := g.AdjList[endNode]

	if !(startExists && endExists) {
		return nil
	}

	// Search childs
	for _, child := range g.Childs {
		found := FindSmallestConvexComponent(child, startNode, endNode)
		if found != nil {
			// Child contains both nodes in a smaller adjlist
			return found
		}
	}

	return g
}
