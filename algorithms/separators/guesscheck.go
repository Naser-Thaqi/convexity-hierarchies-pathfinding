package separators

import (
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
)

func GuessAndCheck(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	passableNodes := []int{}

	for key := range g.AdjList {
		passableNodes = append(passableNodes, key)
	}

	n := len(passableNodes)

	//using bit manipulation to test every possible subset
	for i := range 1 << n {
		candidateSet := []int{}

		for j := range n {
			if (i & (1 << j)) != 0 {
				candidateSet = append(candidateSet, passableNodes[j])
			}
		}
		select {
		case <-ctx.Done():
			// method cancelled
			return nil, false
		default:
			// proceed
		}
		//testing for alpha balancing and convexity
		convexComponents, valid := graphdecomp.BalancedConvexDecomposition(g, candidateSet, ctx)
		if valid {
			return convexComponents, true
		}
	}

	return nil, false
}
