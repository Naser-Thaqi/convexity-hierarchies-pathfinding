package separators

import (
	"bachelor-project/graph"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// create small graph instance
func buildTestGraph() *graph.Graph {
	g := graph.NewGraph(3, 3)
	g.Grid = [][]int{
		{3, 4, 5},
		{6, 7, 8},
		{9, 10, 11},
	}
	/* in metis node form
	1, 2, 3
	4, 5, 6
	7, 8, 9
	*/
	g.BuildAdjlist()
	return g
}

func TestWriteMetisGraphMapped(t *testing.T) {
	g := buildTestGraph()

	// gather original nodeids (old now)
	oldIDs := []int{}
	for id := range g.AdjList {
		oldIDs = append(oldIDs, id)
	}
	sort.Ints(oldIDs)

	idMap := make(map[int]int, len(oldIDs))
	reverseMap := make(map[int]int, len(oldIDs))
	for newID, oldID := range oldIDs {
		idMap[oldID] = newID + 1
		reverseMap[newID+1] = oldID
	}

	testdataDir := "KaFFPa-Test-Data"
	err := os.MkdirAll(testdataDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create testdata dir: %v", err)
	}

	// Open file handle, pass to function
	path := filepath.Join(testdataDir, "testgraph.metis")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer f.Close()

	err = writeMetisGraphMapped(g, f, idMap, reverseMap)
	if err != nil {
		t.Fatalf("writeMetisGraphMapped failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Output file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Output file is empty")
	}
}

func TestReadPartitionFile(t *testing.T) {
	path := "KaFFPa-Test-Data/test_partition.out"

	content := "0\n0\n0\n0\n0\n0\n1\n1\n1\n"
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test partition file: %v", err)
	}
	defer os.Remove(path)

	partitions, err := readPartitionFile(path)
	if err != nil {
		t.Fatalf("readPartitionFile failed: %v", err)
	}

	expected := []int{0, 0, 0, 0, 0, 0, 1, 1, 1}
	if len(partitions) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(partitions))
	}

	for i := range partitions {
		if partitions[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], partitions[i])
		}
	}
}

func TestFindSeparator(t *testing.T) {
	g := buildTestGraph()
	partitions := []int{0, 0, 0, 0, 0, 0, 1, 1, 1}

	idMap := map[int]int{
		3:  1,
		4:  2,
		5:  3,
		6:  4,
		7:  5,
		8:  6,
		9:  7,
		10: 8,
		11: 9,
	}
	separator := findSeparator(g, partitions, idMap)

	expected := map[int]struct{}{
		6:  {},
		7:  {},
		8:  {},
		9:  {},
		10: {},
		11: {},
	}
	if len(separator) != len(expected) {
		t.Fatalf("Expected %d separator nodes, got %d", len(expected), len(separator))
	}
	for _, node := range separator {
		if _, exists := expected[node]; !exists {
			t.Errorf("Unexpected separator node: %d", node)
		}
	}
}

// func TestKaFFPaSeparator(t *testing.T) {
// 	original := config.KaFFPaPath
// 	config.KaFFPaPath = "../../KaHIP/build/kaffpa"

// 	g := buildTestGraph()
// 	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*60)
// 	defer cancel()

// 	subgraphs, ok := KaFFPaSeparator(g, ctxTimeout)

// 	if !ok {
// 		t.Fatal("KaFFPa could not find valid separators")
// 	}

// 	config.KaFFPaPath = original

// 	if len(subgraphs) < 2 {
// 		t.Fatalf("Expected at least 2 subgraphs, got %d", len(subgraphs))
// 	}

// 	for i, sg := range subgraphs {
// 		if len(sg.AdjList) == 0 {
// 			t.Errorf("Subgraph %d is empty", i)
// 		}
// 	}
// }
