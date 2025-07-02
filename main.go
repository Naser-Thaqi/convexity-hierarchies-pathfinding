package main

import (
	"bachelor-project/algorithms"
	"bachelor-project/benchmark"
	"bachelor-project/config"
	"bachelor-project/graph"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	/*
		os.Args[i] execution hint:
		1 = "t" traditional, "c" convexity hierarchies, b1/b2/b3/b4/b5 for specifig benchmark function
		2 = benchmark folder bg512,maze,random,room,street
		3 = must be defined if t or c is choosen => specific file
		4 =
		use t and c only
	*/
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <mode> <benchmarkFolder> [alpha]")
		fmt.Println("Modes:")
		fmt.Println("  t  = traditional BFS")
		fmt.Println("  c  = convex benchmark")
		fmt.Println("  b1 = FindDistanceBenchmarkNormal")
		fmt.Println("  b2 = BuildGraphBenchmarkConvexNormal")
		fmt.Println("  b3 = FindDistanceTimeNormalConvex")
		fmt.Println("  b4 = GetSizeOfGraph")
		return
	}
	mode := os.Args[1]

	if len(os.Args) > 2 {
		parsedAlpha, err := strconv.ParseFloat(os.Args[3], 64)
		if err == nil {
			config.Alpha = parsedAlpha
		}
	}

	//add benchmark names
	benchmarks := []string{}

	switch mode {
	case "t":
		mapPath := os.Args[2]
		scenPath := os.Args[3]
		g := graph.LoadGraphFromFile(mapPath)
		scenarios := benchmark.LoadScenario(scenPath)
		for _, s := range scenarios {
			mapWidth := s[0]
			startX, startY := s[1], s[2]
			goalX, goalY := s[3], s[4]
			x := graph.NodeID(startX, startY, mapWidth)
			y := graph.NodeID(goalX, goalY, mapWidth)

			startTime := time.Now()
			distance := algorithms.BreadthFirstSearch(g.AdjList, x, y)
			runTime := time.Since(startTime).Milliseconds()
			fmt.Println(x, y, distance, runTime)
		}

	case "c":
		mapPath := os.Args[2]
		scenPath := os.Args[3]
		g := graph.LoadGraphFromFile(mapPath)
		scenarios := benchmark.LoadScenario(scenPath)
		for _, s := range scenarios {
			mapWidth := s[0]
			startX, startY := s[1], s[2]
			goalX, goalY := s[3], s[4]
			x := graph.NodeID(startX, startY, mapWidth)
			y := graph.NodeID(goalX, goalY, mapWidth)

			startTime := time.Now()
			subgraph := algorithms.FindSmallestConvexComponent(g, x, y)
			distance := algorithms.BreadthFirstSearch(subgraph.AdjList, x, y)
			runTime := time.Since(startTime).Milliseconds()
			fmt.Println(x, y, distance, runTime)
		}

	case "b1":
		fmt.Println("Running every heuristic for all benchmarks...")
		for _, b := range benchmarks {
			mapFolderName := b + "-map"
			mapDir := filepath.Join("benchmark", "map", mapFolderName)
			output := filepath.Join("benchmark", "output", b)
			everyHeuristic := filepath.Join(output + "-individually-heuristic" + b)
			benchmark.BenchEveryHeuristic(mapDir, everyHeuristic)
		}
	case "b2":
		fmt.Println("Running building normal vs convex hierarchy benchmark for all benchmarks...")
		for _, b := range benchmarks {
			mapFolderName := b + "-map"
			mapDir := filepath.Join("benchmark", "map", mapFolderName)
			output := filepath.Join("benchmark", "output", b)
			buildTime := filepath.Join(output + "-build-time-analysis" + b)
			benchmark.BuildGraphBenchmarkConvexNormal(mapDir, buildTime)
		}
	case "b3":
		fmt.Println("Running normal vs convex time comparison for all benchmarks...")
		for _, b := range benchmarks {
			mapFolderName := b + "-map"
			mapDir := filepath.Join("benchmark", "map", mapFolderName)
			scenDir := filepath.Join("benchmark", "scen", mapFolderName)
			output := filepath.Join("benchmark", "output", b)
			findDistanceTime := filepath.Join(output + "-distance-time-analysis" + b)
			benchmark.FindDistanceTimeNormalConvex(mapDir, scenDir, findDistanceTime)
		}
	case "b4":
		fmt.Println("Running graph size analysis for all benchmarks...")
		for _, b := range benchmarks {
			mapFolderName := b + "-map"
			mapDir := filepath.Join("benchmark", "map", mapFolderName)
			output := filepath.Join("benchmark", "output", b)
			graphAnalysisCsv := filepath.Join(output+"-graph-analysis", b)
			benchmark.GetSizeOfGraph(mapDir, graphAnalysisCsv)
		}
	case "b5":
		fmt.Println("Running combined benchmark (build + distance) for all benchmarks...")
		for _, b := range benchmarks {
			mapFolderName := b + "-map"
			mapDir := filepath.Join("benchmark", "map", mapFolderName)
			scenDir := filepath.Join("benchmark", "scen", b+"-scen")
			outputBasePath := filepath.Join("benchmark", "output", b, "combined-benchmark")

			fmt.Printf("Benchmarking: %s\n", b)
			benchmark.CombinedBenchmarkConvexNormal(mapDir, scenDir, outputBasePath)
		}
	default:
		fmt.Println("Unknown mode:", mode)
	}
}
