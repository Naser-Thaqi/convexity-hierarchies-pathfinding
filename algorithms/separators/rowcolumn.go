package separators

import (
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
	"sort"
	"sync"
)

// Decompose graph via valid alpha balanced convexity separator which are a whole column or row
func RowColumn(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	row_candidates := make([][]int, 0, g.Height)
	column_candidates := make([][]int, 0, g.Width)

	var wg sync.WaitGroup
	wg.Add(2)

	//compute row candidates and sort after length
	go func() {
		defer wg.Done()
		if g.Height > 2 {
			//skip first and last row
			for y := 1; y < g.Height-1; y++ {
				found := false //only append arrays with content
				row := []int{}
				for x := range g.Width {
					nodeID := g.Grid[y][x]
					if nodeID != -1 {
						row = append(row, nodeID)
						found = true
					}
				}
				if found {
					row_candidates = append(row_candidates, row)
				}
			}
			sort.Slice(row_candidates, func(i, j int) bool {
				return len(row_candidates[i]) < len(row_candidates[j])
			})
		}

	}()

	//compute column candidates and sort after lenghts
	go func() {
		defer wg.Done()
		if g.Width > 2 {
			//skip first and last column
			for x := 1; x < g.Width-1; x++ {
				found := false //only append arrays with content
				column := []int{}
				for y := range g.Height {
					nodeID := g.Grid[y][x]
					if nodeID != -1 {
						column = append(column, nodeID)
						found = true
					}
				}
				if found {
					column_candidates = append(column_candidates, column)
				}
			}
			sort.Slice(column_candidates, func(i, j int) bool {
				return len(column_candidates[i]) < len(column_candidates[j])
			})
		}
	}()

	//wait for computation of candidates
	wg.Wait()

	//merge-wise iteration for trying candidates for minimal length
	i := 0
	j := 0

	for i < len(row_candidates) && j < len(column_candidates) {
		if len(row_candidates[i]) <= len(column_candidates[j]) {
			convexComponents, valid := graphdecomp.BalancedConvexDecomposition(g, row_candidates[i], ctx)
			if valid {
				return convexComponents, true
			}
			i++
		} else {
			convexComponents, valid := graphdecomp.BalancedConvexDecomposition(g, column_candidates[j], ctx)
			if valid {
				return convexComponents, true
			}
			j++
		}
	}
	for i < len(row_candidates) {
		convexComponents, valid := graphdecomp.BalancedConvexDecomposition(g, row_candidates[i], ctx)
		if valid {
			return convexComponents, true
		}
		i++
	}
	for j < len(column_candidates) {

		convexComponents, valid := graphdecomp.BalancedConvexDecomposition(g, column_candidates[j], ctx)
		if valid {
			return convexComponents, true
		}
		j++
	}

	return nil, false
}
