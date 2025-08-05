package pathing

import (
	"container/heap"
	"fmt"
	"math"
	"palbaseiq/pkg/types"
)

// Node represents a node in the pathfinding graph
type Node struct {
	Position types.Position
	Cost     float64
	Priority float64
	Parent   *Node
	Index    int // for heap operations
}

// Path represents a path between two points
type Path struct {
	Nodes    []types.Position
	Distance float64
	Cost     float64
}

// Graph represents the pathfinding graph for the base
type Graph struct {
	Base      *types.Base
	Nodes     map[string]*Node
	Edges     map[string][]Edge
	Heuristic HeuristicFunction
}

// Edge represents a connection between two nodes
type Edge struct {
	From   types.Position
	To     types.Position
	Cost   float64
	Weight float64
}

// HeuristicFunction defines the heuristic function for A* pathfinding
type HeuristicFunction func(from, to types.Position) float64

// NewGraph creates a new pathfinding graph for the base
func NewGraph(base *types.Base) *Graph {
	return &Graph{
		Base:      base,
		Nodes:     make(map[string]*Node),
		Edges:     make(map[string][]Edge),
		Heuristic: ManhattanDistance,
	}
}

// GetNodeKey returns a unique key for a position
func GetNodeKey(pos types.Position) string {
	return fmt.Sprintf("%d,%d,%d", pos.X, pos.Y, pos.Z)
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(pos types.Position) {
	key := GetNodeKey(pos)
	if _, exists := g.Nodes[key]; !exists {
		g.Nodes[key] = &Node{
			Position: pos,
			Cost:     math.Inf(1),
		}
	}
}

// AddEdge adds an edge between two positions
func (g *Graph) AddEdge(from, to types.Position, cost float64) {
	fromKey := GetNodeKey(from)

	// Add nodes if they don't exist
	g.AddNode(from)
	g.AddNode(to)

	// Add edge
	edge := Edge{
		From:   from,
		To:     to,
		Cost:   cost,
		Weight: cost,
	}

	g.Edges[fromKey] = append(g.Edges[fromKey], edge)
}

// GetNeighbors returns all valid neighbors of a position
func (g *Graph) GetNeighbors(pos types.Position) []types.Position {
	var neighbors []types.Position

	// Define the 6 possible directions (up, down, left, right, forward, backward)
	directions := []types.Position{
		{0, 1, 0},  // up
		{0, -1, 0}, // down
		{-1, 0, 0}, // left
		{1, 0, 0},  // right
		{0, 0, -1}, // forward
		{0, 0, 1},  // backward
	}

	for _, dir := range directions {
		neighbor := types.Position{
			X: pos.X + dir.X,
			Y: pos.Y + dir.Y,
			Z: pos.Z + dir.Z,
		}

		// Check if neighbor is valid and not occupied
		if g.Base.IsPositionValid(neighbor) && !g.Base.IsPositionOccupied(neighbor) {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

// BuildGraph builds the complete graph from the base
func (g *Graph) BuildGraph() {
	// Clear existing graph
	g.Nodes = make(map[string]*Node)
	g.Edges = make(map[string][]Edge)

	// Add all free positions as nodes
	freePositions := g.Base.GetFreePositions()
	for _, pos := range freePositions {
		g.AddNode(pos)
	}

	// Add edges between adjacent free positions
	for _, pos := range freePositions {
		neighbors := g.GetNeighbors(pos)
		for _, neighbor := range neighbors {
			// Calculate edge cost based on distance and terrain
			cost := g.CalculateEdgeCost(pos, neighbor)
			g.AddEdge(pos, neighbor, cost)
		}
	}
}

// CalculateEdgeCost calculates the cost of moving between two positions
func (g *Graph) CalculateEdgeCost(from, to types.Position) float64 {
	baseCost := from.Distance(to)

	// Add penalties for vertical movement (climbing/descending)
	if from.Y != to.Y {
		baseCost *= 1.5 // Vertical movement is more expensive
	}

	// Add penalties for proximity to walls or other obstacles
	obstaclePenalty := g.CalculateObstaclePenalty(to)

	return baseCost + obstaclePenalty
}

// CalculateObstaclePenalty calculates penalty for being near obstacles
func (g *Graph) CalculateObstaclePenalty(pos types.Position) float64 {
	penalty := 0.0

	// Check in a 3x3x3 area around the position
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				checkPos := types.Position{
					X: pos.X + dx,
					Y: pos.Y + dy,
					Z: pos.Z + dz,
				}

				if g.Base.IsPositionValid(checkPos) && g.Base.IsPositionOccupied(checkPos) {
					// Calculate distance-based penalty
					distance := math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
					if distance > 0 {
						penalty += 0.1 / distance
					}
				}
			}
		}
	}

	return penalty
}

// FindPath finds the shortest path between two positions using A* algorithm
func (g *Graph) FindPath(start, end types.Position) (*Path, error) {
	if !g.Base.IsPositionValid(start) || !g.Base.IsPositionValid(end) {
		return nil, fmt.Errorf("invalid start or end position")
	}

	if g.Base.IsPositionOccupied(start) || g.Base.IsPositionOccupied(end) {
		return nil, fmt.Errorf("start or end position is occupied")
	}

	// Initialize open and closed sets
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	closedSet := make(map[string]bool)

	// Initialize start node
	startKey := GetNodeKey(start)
	startNode := &Node{
		Position: start,
		Cost:     0,
		Priority: g.Heuristic(start, end),
	}

	heap.Push(openSet, startNode)

	// Keep track of all nodes for path reconstruction
	allNodes := make(map[string]*Node)
	allNodes[startKey] = startNode

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)
		currentKey := GetNodeKey(current.Position)

		// Check if we reached the goal
		if current.Position == end {
			return g.ReconstructPath(current), nil
		}

		closedSet[currentKey] = true

		// Check neighbors
		neighbors := g.GetNeighbors(current.Position)
		for _, neighborPos := range neighbors {
			neighborKey := GetNodeKey(neighborPos)

			if closedSet[neighborKey] {
				continue
			}

			// Calculate tentative cost
			tentativeCost := current.Cost + g.CalculateEdgeCost(current.Position, neighborPos)

			// Get or create neighbor node
			neighbor, exists := allNodes[neighborKey]
			if !exists {
				neighbor = &Node{
					Position: neighborPos,
					Cost:     math.Inf(1),
				}
				allNodes[neighborKey] = neighbor
			}

			if tentativeCost < neighbor.Cost {
				neighbor.Parent = current
				neighbor.Cost = tentativeCost
				neighbor.Priority = tentativeCost + g.Heuristic(neighborPos, end)

				if !exists {
					heap.Push(openSet, neighbor)
				} else {
					heap.Fix(openSet, neighbor.Index)
				}
			}
		}
	}

	return nil, fmt.Errorf("no path found between %s and %s", start, end)
}

// ReconstructPath reconstructs the path from the goal node
func (g *Graph) ReconstructPath(goalNode *Node) *Path {
	var positions []types.Position
	current := goalNode

	for current != nil {
		positions = append([]types.Position{current.Position}, positions...)
		current = current.Parent
	}

	// Calculate total distance and cost
	distance := 0.0
	cost := 0.0
	for i := 1; i < len(positions); i++ {
		dist := positions[i-1].Distance(positions[i])
		distance += dist
		cost += g.CalculateEdgeCost(positions[i-1], positions[i])
	}

	return &Path{
		Nodes:    positions,
		Distance: distance,
		Cost:     cost,
	}
}

// FindOptimalPath finds the optimal path considering multiple factors
func (g *Graph) FindOptimalPath(start, end types.Position, constraints []PathConstraint) (*Path, error) {
	// For now, just use the basic A* algorithm
	// In the future, this could implement more sophisticated pathfinding
	// that considers multiple constraints and objectives
	return g.FindPath(start, end)
}

// PathConstraint represents a constraint for pathfinding
type PathConstraint interface {
	IsValid(path *Path) bool
	GetCost(path *Path) float64
}

// ManhattanDistance is a heuristic function using Manhattan distance
func ManhattanDistance(from, to types.Position) float64 {
	return float64(from.ManhattanDistance(to))
}

// EuclideanDistance is a heuristic function using Euclidean distance
func EuclideanDistance(from, to types.Position) float64 {
	return from.Distance(to)
}

// PriorityQueue implementation for A* algorithm
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}
