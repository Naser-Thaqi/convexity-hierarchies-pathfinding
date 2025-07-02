package separators

import (
	"bachelor-project/config"
	"bachelor-project/graph"
	"bachelor-project/graphdecomp"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Decompose graph into balanced convex components using KaFFPa partitions and separator nodes
func KaFFPaSeparator(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool) {
	tmpInputFile, err := os.CreateTemp("", "kaffpa_input_*.graph")
	if err != nil {
		fmt.Println("Error creating temp input file:", err)
		return nil, false
	}
	defer os.Remove(tmpInputFile.Name())
	defer tmpInputFile.Close()

	tmpOutputFile, err := os.CreateTemp("", "kaffpa_output_*.out")
	if err != nil {
		fmt.Println("Error creating temp output file:", err)
		return nil, false
	}
	defer os.Remove(tmpOutputFile.Name())
	defer tmpOutputFile.Close()

	// gather original nodeids (old now)
	oldIDs := []int{}
	for id := range g.AdjList {
		oldIDs = append(oldIDs, id)
	}
	sort.Ints(oldIDs) // for stabile order, and backtracking of original/old ids

	idMap := make(map[int]int, len(oldIDs))
	reverseMap := make(map[int]int, len(oldIDs))
	for newID, oldID := range oldIDs {
		idMap[oldID] = newID + 1 // nodeid begins with 1
		reverseMap[newID+1] = oldID
	}

	// write graph into correct format for application of KaFFPa
	err = writeMetisGraphMapped(g, tmpInputFile, idMap, reverseMap)
	if err != nil {
		fmt.Println("Error while writing metis file:", err)
		return nil, false
	}

	// parameters for KaFFPa
	k := 2         //number of partitions
	imbalance := 0 // partition can differ by 33% of Nodes/k
	for {
		select {
		case <-ctx.Done():
			return nil, false
		default:
			// proceed
		}

		cmd := exec.Command(config.KaFFPaPath, tmpInputFile.Name(),
			"--k="+strconv.Itoa(k),
			"--imbalance="+strconv.Itoa(imbalance),
			"--preconfiguration=strong",
			"--output_filename="+tmpOutputFile.Name(),
		)

		cmd.Stdout = nil
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println("KaFFPa Error on execution:", err)
			continue
		}

		partitions, err := readPartitionFile(tmpOutputFile.Name())
		if err != nil {
			fmt.Println("Error reading the partition:", err)
			continue
		}

		select {
		case <-ctx.Done():
			return nil, false
		default:
			// proceed
		}

		separator := findSeparator(g, partitions, idMap)

		// check convexity
		subgraphs, isConvex := graphdecomp.ConvexDecomposition(g, separator, ctx)
		if isConvex {
			return subgraphs, true
		}
		//k++
		//imbalance = int(float64(k)*config.Alpha*100 - 100.0)
	}
}

// write current graph into a file in metis format
func writeMetisGraphMapped(g *graph.Graph, file *os.File, idMap, reverseMap map[int]int) error {
	// count edges
	edgeNum := 0
	for _, neighbors := range g.AdjList {
		edgeNum += len(neighbors)
	}
	edgeNum /= 2

	n := len(idMap) // number of nodes
	m := edgeNum    // number of edges
	fmt.Fprintf(file, "%d %d\n", n, m)
	// write all neighbors of nodes 1 to n
	for i := 1; i <= n; i++ {
		oldID := reverseMap[i]
		for _, nb := range g.AdjList[oldID] {
			fmt.Fprintf(file, "%d ", idMap[nb])
		}
		fmt.Fprintln(file)
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

// read partition of file
func readPartitionFile(path string) ([]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("file does not exist")
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	partitions := make([]int, len(lines))
	for i, line := range lines {
		partitions[i], err = strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			return nil, err
		}
	}
	return partitions, nil
}

// find separators by choosing nodes that have neighbors of 2 different partitions
func findSeparator(g *graph.Graph, partitions []int, idMap map[int]int) []int {
	separators := []int{}
	for node, neighbors := range g.AdjList {
		if len(neighbors) < 2 {
			continue
		}
		for _, neighbor := range neighbors {
			found := partitions[idMap[neighbors[0]]-1]
			if found != partitions[idMap[neighbor]-1] {
				separators = append(separators, node)
				break
			}
		}
	}

	return separators
}
