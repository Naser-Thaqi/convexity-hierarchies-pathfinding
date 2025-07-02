package algorithms

// computes the distance between two nodes in a graph with BFS
func BreadthFirstSearch(adjlist map[int][]int, startID, endID int) int {
	// base case: start and end are the same node
	if startID == endID {
		return 0
	}

	visited := make(map[int]struct{}, len(adjlist))
	queue := make([]int, 0, len(adjlist))
	queue = append(queue, startID)
	visited[startID] = struct{}{} // set true
	head := 0                     // pointer
	depth := 0                    // depth for distance

	for head < len(queue) {
		levelSize := len(queue) - head
		depth++

		// visit every node of current lebel
		for range levelSize {
			current := queue[head]
			head++

			// visit every neighbor of current node
			for _, neighbor := range adjlist[current] {
				if neighbor == endID {
					return depth
				}
				// visit not visited nodes later
				if _, exists := visited[neighbor]; !exists {
					visited[neighbor] = struct{}{} // set true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// return -1 if endID is not reachable from startID
	return -1
}
