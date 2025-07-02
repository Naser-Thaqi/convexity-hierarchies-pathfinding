# convexity-hierarchies-pathfinding
Heuristics and hierarchical graph decomposition for grid-based pathfinding. Based on the paper "Convexity Hierarchies in Grid Networks (2023)"

# General
This explains the structure of the folders and files, and how to run the code.

---

# Table of Contents

1. [Folder and File Structure](#folder-and-file-structure)
2. [Dependencies](#dependencies)
3. [How to Run](#how-to-run)

---

# Folder and File Structure

- **`main.go`**: The main program to execute everything

- **`config/`**:
  - `config.go`: Contains configuration for alpha, timeout for heuristic, and relative path to kaffpa

- **`algorithms/`**:  
  - `bfs.go`: Breadth-First search implementation
  - `convexhierarchy.go`: Build convex hierarchical structure
  - `bfs_test.go`: Test functions of bfs.go
  - `convexhierarchy_test.go`:  Test functions of convexhierarchy.go
  - **`separators/`**: All heuristics to compute alpha balanced convex decompositions. Every heuristic has its own name_test.go file
  - `guesscheck.go`: heuristic
  - `KaFFPa.go`: heuristic
  - `oneshortestpath.go`: heuristic
  - `twoshortestpath.go`: heuristic
  - `rowcolumn.go`: heurisitc
  - `holecutting.go`: heuristic
  - `unionfind.go`: helper methods for heuristics

- **`graph/`**:
  - `graph.go`: Own implementation of a graph class (structure) and helper methods

- **`graphdecomp/`**: Core graph decomposition logic ,Every file has its own name_test.go file
  - `balanced.go`
  - `convexity.go`
  - `balancedconvexdecomp.go`

- **`benchmark/`**:  
  - `benchmark.go`: Code for evaluating the performance of the algorithms and writing into CSV file.  
  - **`map/`**: Contains all benchmark maps
  - **`scen/`**: Contains all scen files to corresponding map files
  - **`output/`** Contains all csv files that were generated

---

## Dependencies

Ensure you have Go installed on your system.
Ensure you have KaFFPa installed and the path is configured in config/config.go
Ensure you have choosen the correct alpha in config if you want to change the value
For benchmarks you need https://www.movingai.com/benchmarks/formats.html for map and scen files.
In map/ and scen/ every "benchmark" folder needs  "-scen" or "-map" to their name.
map files end with ".map" and scen files with ".map.scen"

---

# How to Run
If you have KAFFPa installed, you can add separator.KaFFpaSeparator into the pipeline in algorithms/convexhierarchy.go and
outcomment TestKaFFPaSeparator(t *testing.T) in algorithms/separators/KaFFPa_test.go.

Use the following command to run the program in main.go:
b + Number for various benchmarks
c for convexity hierarchy and t for traditional 

```bash
go run main.go <b1/b2/b3/b4/b5>
go run main.go <c> < filepath map > <filepath scen >
go run main.go <t> < filepath map > <filepath scen >
```
Use the following command to run all tests (open console in main folder):
 ```bash
go run test -v ./...
```

Use the following command to run a single test (open console where <name>.test.go is):
 ```bash
go run test -v <Functionname>
```

Use the following command to open code coverage visualization (open console in main folder):
```bash
go tool cover -html=coverage.out
```
