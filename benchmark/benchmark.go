package benchmark

import (
	"bachelor-project/algorithms"
	"bachelor-project/algorithms/separators"
	"bachelor-project/config"
	"bachelor-project/graph"
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CombinedBenchmarkConvexNormal(mapDir, scenDir, outputBasePath string) {
	mapFiles, err := filepath.Glob(filepath.Join(mapDir, "*.map"))
	if err != nil {
		fmt.Printf("Error reading map directory: %v\n", err)
		return
	}

	// prepare csv files
	buildTimeCsvPath := filepath.Join(filepath.Dir(outputBasePath), "build-time.csv")
	distanceTimeDir := filepath.Join(filepath.Dir(outputBasePath), "distance")

	err = os.MkdirAll(filepath.Dir(buildTimeCsvPath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating build-time CSV folder: %v\n", err)
		return
	}
	err = os.MkdirAll(distanceTimeDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating distance CSV folder: %v\n", err)
		return
	}

	// Build-Time csv file
	buildTimeFile, err := os.Create(buildTimeCsvPath)
	if err != nil {
		fmt.Printf("Error creating build-time CSV file: %v\n", err)
		return
	}
	defer buildTimeFile.Close()

	buildWriter := csv.NewWriter(buildTimeFile)
	defer buildWriter.Flush()
	buildWriter.Write([]string{"Map Name", "Time1 normal (ms)", "Time2 convex (ms)", "Number of Subgraphs"})
	buildWriter.Flush()

	for _, mapPath := range mapFiles {
		mapName := strings.TrimSuffix(filepath.Base(mapPath), ".map")
		scenPath := filepath.Join(scenDir, mapName+".map"+".scen")
		fmt.Printf("Processing map: %s\n", mapName)

		// build time
		start1 := time.Now()
		graph.LoadGraphFromFile(mapPath)
		time1 := time.Since(start1).Milliseconds()

		start2 := time.Now()
		g := graph.LoadGraphFromFile(mapPath)
		algorithms.BuildConvexHierarchy(g)
		time2 := time.Since(start2).Milliseconds()

		countSubgraphs := countLeaves(g)

		buildWriter.Write([]string{
			mapName,
			fmt.Sprintf("%d", time1),
			fmt.Sprintf("%d", time2),
			fmt.Sprintf("%d", countSubgraphs),
		})
		buildWriter.Flush()

		fmt.Println("Start scenario")

		// distance time
		if _, err := os.Stat(scenPath); os.IsNotExist(err) {
			fmt.Printf("Scenario file missing for map %s\n", mapName)
			continue
		}

		scenarios := LoadScenario(scenPath)
		if len(scenarios) == 0 {
			fmt.Printf("No valid scenarios in: %s\n", scenPath)
			continue
		}

		csvDistFilePath := filepath.Join(distanceTimeDir, mapName+".csv")
		distFile, err := os.Create(csvDistFilePath)
		if err != nil {
			fmt.Printf("Error creating distance CSV file: %v\n", err)
			continue
		}
		writer := csv.NewWriter(distFile)
		writer.Write([]string{"Instance", "Bucket", "Time1 normal (ms)", "Time2 convex (ms)", "Distance 1", "Distance 2", "search size normal", "search size convex", "Time to find subgraph Convex", "Time to find Distance in subgraph Convex", "Count subgraphs"})
		writer.Flush()

		for i, s := range scenarios {
			mapWidth := s[0]
			startX, startY := s[1], s[2]
			goalX, goalY := s[3], s[4]
			bucket := s[5]
			x := graph.NodeID(startX, startY, mapWidth)
			y := graph.NodeID(goalX, goalY, mapWidth)

			startTimeNormal := time.Now()
			distanceNormal := algorithms.BreadthFirstSearch(g.AdjList, x, y)
			runTimeNormal := time.Since(startTimeNormal).Milliseconds()
			searchSpaceSizeNormal := len(g.AdjList)

			startTimeConvex := time.Now()
			subgraph := algorithms.FindSmallestConvexComponent(g, x, y)
			distanceConvex := algorithms.BreadthFirstSearch(subgraph.AdjList, x, y)
			runTimeConvex := time.Since(startTimeConvex).Milliseconds()

			startTimeConvexFindSubgraph := time.Now()
			subgraph = algorithms.FindSmallestConvexComponent(g, x, y)
			runTimeConvexFindSubgraph := time.Since(startTimeConvexFindSubgraph).Milliseconds()

			startTimeConvexFindDistance := time.Now()
			algorithms.BreadthFirstSearch(subgraph.AdjList, x, y)
			runTimeConvexFindDistance := time.Since(startTimeConvexFindDistance).Milliseconds()

			searchSpaceSizeConvex := len(subgraph.AdjList)

			writer.Write([]string{
				fmt.Sprintf("Instance-%d", i+1),
				fmt.Sprintf("%d", bucket),
				fmt.Sprintf("%d", runTimeNormal),
				fmt.Sprintf("%d", runTimeConvex),
				fmt.Sprintf("%d", distanceNormal),
				fmt.Sprintf("%d", distanceConvex),
				fmt.Sprintf("%d", searchSpaceSizeNormal),
				fmt.Sprintf("%d", searchSpaceSizeConvex),
				fmt.Sprintf("%d", runTimeConvexFindSubgraph),
				fmt.Sprintf("%d", runTimeConvexFindDistance),
				fmt.Sprintf("%d", countSubgraphs),
			})
			writer.Flush()
		}
		distFile.Close()
	}

	fmt.Println("Combined benchmark completed. Output saved to:")
	fmt.Println(" - Build times:", buildTimeCsvPath)
	fmt.Println(" - Distance times per map in:", distanceTimeDir)
}

func BuildGraphBenchmarkConvexNormal(directory string, csvFilePath string) {
	// search for all .map files in folder
	mapFiles, err := filepath.Glob(filepath.Join(directory, "*.map"))
	if err != nil {
		fmt.Printf("Error reading map directory: %v\n", err)
		return
	}

	// create csv folder
	err = os.MkdirAll(filepath.Dir(csvFilePath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating CSV folder: %v\n", err)
		return
	}

	// create csv file
	csvFile, err := os.Create(csvFilePath + ".csv")
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// write header
	writer.Write([]string{"Instance", "Time1 normal (ms)", "Time2 convex (ms)", "Number of Subgraphs"})

	// iterate over maps
	for _, mapPath := range mapFiles {
		mapName := strings.TrimSuffix(filepath.Base(mapPath), ".map")
		fmt.Printf("Processing map: %s\n", mapName)

		// normal building
		start1 := time.Now()
		graph.LoadGraphFromFile(mapPath) // read map
		time1 := time.Since(start1).Milliseconds()

		// convex building
		start2 := time.Now()
		g := graph.LoadGraphFromFile(mapPath) // read map
		algorithms.BuildConvexHierarchy(g)    // build hierarichal structure
		time2 := time.Since(start2).Milliseconds()

		countSubgraphs := countLeaves(g) // of convex building

		writer.Write([]string{
			mapName,
			fmt.Sprintf("%d", time1),
			fmt.Sprintf("%d", time2),
			fmt.Sprintf("%d", countSubgraphs),
		})
	}
	fmt.Println("Benchmarking completed. Results saved to:", csvFilePath)
}

// Distance normal/convex, Time normal/convex, Time finding convex subgraph, time bfs in convex subgraph, count subgraphs
func FindDistanceTimeNormalConvex(mapDir, scenDir, csvFilePath string) {
	mapFiles, err := filepath.Glob(filepath.Join(mapDir, "*.map"))
	if err != nil {
		fmt.Printf("Error reading map directory: %v\n", err)
		return
	}

	// create csv folder
	err = os.MkdirAll(filepath.Dir(csvFilePath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating CSV folder: %v\n", err)
		return
	}

	for _, mapPath := range mapFiles {
		mapName := strings.TrimSuffix(filepath.Base(mapPath), ".map")
		scenPath := filepath.Join(scenDir, mapName+".scen")

		// create csv file for corresponding map
		csvFile, err := os.Create(filepath.Join(filepath.Dir(csvFilePath), mapName+".csv"))
		if err != nil {
			fmt.Printf("Error creating CSV file: %v\n", err)
			return
		}

		writer := csv.NewWriter(csvFile)
		// write Header
		writer.Write([]string{"Instance", "Bucket", "Time1 normal (ms)", "Time2 convex (ms)", "Distance 1", "Distance 2", "search size normal", "search size convex", "Time to find subgraph Convex", "Time to find Distance in subgraph Convex", "Count subgraphs"})

		// Check scen file
		if _, err := os.Stat(scenPath); os.IsNotExist(err) {
			fmt.Printf("Scenario file missing for map %s\n", mapName)
			continue
		}

		scenarios := LoadScenario(scenPath)
		if len(scenarios) == 0 {
			fmt.Printf("No valid scenarios found in: %s\n", scenPath)
			continue
		}

		g := graph.LoadGraphFromFile(mapPath)
		algorithms.BuildConvexHierarchy(g)

		for i, s := range scenarios {
			mapWidth := s[0]
			startX, startY := s[1], s[2]
			goalX, goalY := s[3], s[4]
			bucket := s[5]
			x := graph.NodeID(startX, startY, mapWidth)
			y := graph.NodeID(goalX, goalY, mapWidth)

			startTimeNormal := time.Now()
			distanceNormal := algorithms.BreadthFirstSearch(g.AdjList, x, y)
			runTimeNormal := time.Since(startTimeNormal).Milliseconds()
			searchSpaceSizeNormal := len(g.AdjList)

			startTimeConvex := time.Now()
			subgraph := algorithms.FindSmallestConvexComponent(g, x, y)
			distanceConvex := algorithms.BreadthFirstSearch(subgraph.AdjList, x, y)
			runTimeConvex := time.Since(startTimeConvex).Milliseconds()

			startTimeConvexFindSubgraph := time.Now()
			subgraph = algorithms.FindSmallestConvexComponent(g, x, y)
			runTimeConvexFindSubgraph := time.Since(startTimeConvexFindSubgraph).Milliseconds()

			startTimeConvexFindDistance := time.Now()
			algorithms.BreadthFirstSearch(subgraph.AdjList, x, y)
			runTimeConvexFindDistance := time.Since(startTimeConvexFindDistance).Milliseconds()

			searchSpaceSizeConvex := len(subgraph.AdjList)
			countSubgraphs := countLeaves(g)

			// write into csv file
			writer.Write([]string{
				fmt.Sprintf("Instance-%d", i+1),
				fmt.Sprintf("%d", bucket),
				fmt.Sprintf("%d", runTimeNormal),
				fmt.Sprintf("%d", runTimeConvex), // find subgraph + distance
				fmt.Sprintf("%d", distanceNormal),
				fmt.Sprintf("%d", distanceConvex),
				fmt.Sprintf("%d", searchSpaceSizeNormal),
				fmt.Sprintf("%d", searchSpaceSizeConvex),
				fmt.Sprintf("%d", runTimeConvexFindSubgraph),
				fmt.Sprintf("%d", runTimeConvexFindDistance),
				fmt.Sprintf("%d", countSubgraphs),
			})
		}
		writer.Flush()
		csvFile.Close()
	}

	fmt.Println("Benchmarking completed. Results saved to:", csvFilePath)
}

func GetSizeOfGraph(directory string, csvFilePath string) {
	// find all .map in folder
	mapFiles, err := filepath.Glob(filepath.Join(directory, "*.map"))
	if err != nil {
		fmt.Printf("Error reading map directory: %v\n", err)
		return
	}

	// create csv folder
	err = os.MkdirAll(filepath.Dir(csvFilePath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating CSV folder: %v\n", err)
		return
	}

	// create csv file
	csvFile, err := os.Create(csvFilePath + ".csv")
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// write header
	writer.Write([]string{"Instance", "Map Name", "Number of Nodes", "Number of Edges"})

	for i, mapPath := range mapFiles {
		mapName := strings.TrimSuffix(filepath.Base(mapPath), ".map")
		fmt.Printf("Processing map: %s\n", mapName)

		g := graph.LoadGraphFromFile(mapPath)
		if g == nil {
			fmt.Printf("Error loading graph: %s\n", mapPath)
			continue
		}

		nodes := len(g.AdjList)
		edges := 0
		for _, neighbors := range g.AdjList {
			edges += len(neighbors)
		}
		edges /= 2

		// write into csv file
		writer.Write([]string{
			fmt.Sprintf("Instance-%d", i+1),
			mapName,
			fmt.Sprintf("%d", nodes),
			fmt.Sprintf("%d", edges),
		})
	}

	fmt.Println("Graph size analysis completed. Output saved to:", csvFilePath)
}

type sep func(g *graph.Graph, ctx context.Context) ([]*graph.Graph, bool)

func BenchEveryHeuristic(directory string, csvFilePath string) {
	// search all .map in folder
	mapFiles, err := filepath.Glob(filepath.Join(directory, "*.map"))
	if err != nil {
		fmt.Printf("Error reading map directory: %v\n", err)
		return
	}

	// create csv folder
	err = os.MkdirAll(filepath.Dir(csvFilePath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating CSV folder: %v\n", err)
		return
	}

	// create csv file
	csvFile, err := os.Create(csvFilePath + ".csv")
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// write header
	writer.Write([]string{"Instance", "Map Name",
		"Kaffpa Separator", "OSP Separator", "TSP Separator", "Row/Column Separator", "Holecutting Separator", "GuessAndCheck separator",
		"Kaffpa imbalance", "OSP imbalance", "TSP imbalance", "Row/Column imbalance", "Holecutting imbalance", "GuessCheck imbalance",
	})

	for i, mapPath := range mapFiles {
		mapName := strings.TrimSuffix(filepath.Base(mapPath), ".map")
		fmt.Printf("Processing map: %s\n", mapName)

		separatorSizes := []int{-1, -1, -1, -1, -1, -1}
		imbalancedRatio := []float64{-1.0, -1.0, -1.0, -1.0, -1.0, -1.0}
		heuristics := []sep{separators.KaFFPaSeparator, separators.OneShortestPath, separators.TwoShortestPath, separators.RowColumn, separators.HoleCutting, separators.GuessAndCheck}
		if i > -1 {
			for j, sepF := range heuristics {

				g := graph.LoadGraphFromFile(mapPath)
				if g == nil {
					fmt.Printf("Error loading graph: %s\n", mapPath)
					continue
				}

				resultChan := make(chan []*graph.Graph, 1) // channel for result
				ctx, cancel := context.WithTimeout(context.Background(), config.Time)

				sepFunc := sepF
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
					fmt.Println("heuristic done")
					// return only positive result
					if res != nil {
						g.Childs = res
						separatorSizes[j] = getSeparatorSize(g)
						imbalancedRatio[j] = getImblancedRatio(g)
						continue
					} // else: try another heuristic
				case <-ctx.Done():
					fmt.Println("heuristic done")
					cancel()
					continue
				}

			}
		}
		// write one line into csv file
		writer.Write([]string{
			fmt.Sprintf("Instance-%d", i+1),
			mapName,
			fmt.Sprintf("%d", separatorSizes[0]),
			fmt.Sprintf("%d", separatorSizes[1]),
			fmt.Sprintf("%d", separatorSizes[2]),
			fmt.Sprintf("%d", separatorSizes[3]),
			fmt.Sprintf("%d", separatorSizes[4]),
			fmt.Sprintf("%d", separatorSizes[5]),
			fmt.Sprintf("%.4f", imbalancedRatio[0]),
			fmt.Sprintf("%.4f", imbalancedRatio[1]),
			fmt.Sprintf("%.4f", imbalancedRatio[2]),
			fmt.Sprintf("%.4f", imbalancedRatio[3]),
			fmt.Sprintf("%.4f", imbalancedRatio[4]),
			fmt.Sprintf("%.4f", imbalancedRatio[5]),
		})

		writer.Flush()
	}

	fmt.Println("Graph size analysis completed. Output saved to:", csvFilePath)
}

func LoadScenario(filePath string) [][6]int {
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

	var scen [][6]int
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		mapWidth, _ := strconv.Atoi(parts[2])
		// mapHeight, _ := strconv.Atoi(parts[3])
		startX, _ := strconv.Atoi(parts[4])
		startY, _ := strconv.Atoi(parts[5])
		goalX, _ := strconv.Atoi(parts[6])
		goalY, _ := strconv.Atoi(parts[7])
		bucket, _ := strconv.Atoi(parts[0])

		coords := [6]int{mapWidth, startX, startY, goalX, goalY, bucket}
		scen = append(scen, coords)
	}

	return scen
}

// Count number of subgraphs
func countLeaves(g *graph.Graph) int {
	if len(g.Childs) == 0 {
		return 1
	}

	count := 0
	for _, child := range g.Childs {
		count += countLeaves(child)
	}
	return count
}

func getSeparatorSize(g *graph.Graph) int {
	size := len(g.AdjList)
	for _, subgraph := range g.Childs {
		size -= len(subgraph.AdjList)
	}
	return size
}

func getImblancedRatio(g *graph.Graph) float64 {
	if len(g.Childs) == 0 {
		return -1.0
	}

	minSize := int(^uint(0) >> 1) // max int
	maxSize := 0

	for _, subgraph := range g.Childs {
		size := len(subgraph.AdjList)
		if size < minSize {
			minSize = size
		}
		if size > maxSize {
			maxSize = size
		}
	}

	if maxSize == 0 {
		return -1.0
	}

	return float64(minSize) / float64(maxSize)
}
