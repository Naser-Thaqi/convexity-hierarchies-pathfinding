package graph

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// For visiting neighbors
var Directions = [4][2]int{
	{-1, 0}, // north
	{1, 0},  // south
	{0, -1}, // west
	{0, 1},  // east
}

type Graph struct {
	AdjList map[int][]int // adjacency list
	Grid    [][]int
	Childs  []*Graph
	Height  int
	Width   int
}

// Create new graph object
func NewGraph(height int, width int) *Graph {
	return &Graph{
		AdjList: make(map[int][]int), // creates empty adjlist
		Height:  height,
		Width:   width,
	}
}

// Add edge to the adjacency list of graph object
func (g *Graph) AddEdge(v, w int) {
	if _, exists := g.AdjList[v]; !exists {
		g.AdjList[v] = make([]int, 0, 4)
	}
	g.AdjList[v] = append(g.AdjList[v], w) // add w to the adjacency list of v
}

// removes a node and all its incident edges from the graph (through adjacecy list)
// the node is removed from the adjaceny list and also from all of the adjacency lists of all its neighbors
func RemoveNode(adjlist map[int][]int, node int) {
	for _, neighbor := range adjlist[node] { // iterate through the neighbors of the node
		adjlist[neighbor] = remove(adjlist[neighbor], node) // remove the node from the adjacency list of each neighbor
	}
	delete(adjlist, node) // delete the node itself from the adjacency list
}

// remove is a helper function to remove an element from a slice
// it replaces the element to be removed with the last element in the slice
// and then reslices to exclude the last element, reducing the slice size by one
func remove(slice []int, elem int) []int {
	for i, v := range slice { // iterate through the slice
		if v == elem { // find the element to remove
			slice[i] = slice[len(slice)-1] // replace it with the last element
			return slice[:len(slice)-1]    // reslice to exclude the last element
		}
	}
	return slice
}

func RestoreNode(adjlist, restoredAdj map[int][]int, nodeid int) {
	restoredAdj[nodeid] = adjlist[nodeid]
	for _, neighbor := range restoredAdj[nodeid] {
		restoredAdj[neighbor] = append(restoredAdj[neighbor], nodeid)
	}
}

// Compute nodeid for given position of grid
func NodeID(x, y, width int) int {
	return y*width + x
}

// Compute x,y coordinates of given nodeid in current grid
func CoordinatesFromNodeID(id, width int) (x, y int) {
	y = id / width
	x = id % width
	return x, y
}

// Load graph from file and build grid and adjacency list
func LoadGraphFromFile(filePath string) *Graph {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: could not open file: %v\n", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil
	}
	scanner.Scan()

	heightLine := scanner.Text()
	parts := strings.Fields(heightLine)
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	// read third line width
	if !scanner.Scan() {
		return nil
	}
	widthLine := scanner.Text()
	parts = strings.Fields(widthLine)
	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	// check for word map
	if !scanner.Scan() {
		return nil
	}
	if scanner.Text() != "map" {
		return nil
	}

	// read grid from fifth line
	grid := make([][]int, height)
	for y := 0; y < height && scanner.Scan(); y++ {
		line := scanner.Text()

		grid[y] = make([]int, width)

		for x := range width {
			// ".", "G" are passable nodes
			if line[x] == '.' || line[x] == 'G' {
				grid[y][x] = NodeID(x, y, width)
			} else {
				grid[y][x] = -1
			}
		}
	}

	//create graph object
	graph := NewGraph(height, width)

	graph.Grid = grid

	graph.BuildAdjlist()

	return graph
}

// buildy adjacency list
func (g *Graph) BuildAdjlist() {
	for y := range g.Height {
		for x := range g.Width {
			// if node is not passable skip (continue)
			if g.Grid[y][x] == -1 {
				continue
			}
			g.AdjList[g.Grid[y][x]] = []int{}
			// check neighbors
			for _, dir := range Directions {
				nx := x + dir[0]
				ny := y + dir[1]
				// check for out of grid nodes
				if nx < 0 || ny < 0 || nx >= g.Width || ny >= g.Height {
					continue
				}
				// check for passable neighbors
				if g.Grid[ny][nx] == -1 {
					continue
				}
				//add edge
				g.AddEdge(g.Grid[y][x], g.Grid[ny][nx])
			}
		}
	}
}

// Copy adjacency list of a given graph object and return the copy
func (g *Graph) CopyAdjlist() map[int][]int {
	copy := make(map[int][]int, len(g.AdjList))
	for node, neighbors := range g.AdjList {
		copy[node] = slices.Clone(neighbors)
	}
	return copy
}
