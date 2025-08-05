package optimizer

import (
	"math"
	"math/rand"
	"palbaseiq/pkg/pathing"
	"palbaseiq/pkg/types"
	"sort"
	"time"
)

// PlacementOptimizer handles the optimization of item placement in the base
type PlacementOptimizer struct {
	Base  *types.Base
	Graph *pathing.Graph
}

// OptimizationConfig holds configuration for the optimization process
type OptimizationConfig struct {
	MaxIterations     int
	Temperature       float64
	CoolingRate       float64
	MinTemperature    float64
	RandomSeed        int64
	PathfindingWeight float64
	EfficiencyWeight  float64
	CompactnessWeight float64
}

// DefaultConfig returns a default optimization configuration
func DefaultConfig() *OptimizationConfig {
	return &OptimizationConfig{
		MaxIterations:     1000,
		Temperature:       100.0,
		CoolingRate:       0.95,
		MinTemperature:    0.1,
		RandomSeed:        time.Now().UnixNano(),
		PathfindingWeight: 0.4,
		EfficiencyWeight:  0.3,
		CompactnessWeight: 0.3,
	}
}

// NewPlacementOptimizer creates a new placement optimizer
func NewPlacementOptimizer(base *types.Base) *PlacementOptimizer {
	graph := pathing.NewGraph(base)
	return &PlacementOptimizer{
		Base:  base,
		Graph: graph,
	}
}

// PlacementScore represents the score of a placement configuration
type PlacementScore struct {
	TotalScore       float64
	PathfindingScore float64
	EfficiencyScore  float64
	CompactnessScore float64
	Details          map[string]float64
}

// OptimizePlacement optimizes the placement of items in the base
func (po *PlacementOptimizer) OptimizePlacement(items []*types.Item, config *OptimizationConfig) (*types.Base, *PlacementScore, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Set random seed
	rand.Seed(config.RandomSeed)

	// Create a copy of the base for optimization
	optimizedBase := po.Base.Clone()

	// Build the pathfinding graph
	po.Graph.Base = optimizedBase
	po.Graph.BuildGraph()

	// Sort items by priority (higher priority first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	// Initial placement using greedy algorithm
	po.placeItemsGreedy(optimizedBase, items)

	// Optimize using simulated annealing
	bestBase := optimizedBase.Clone()
	bestScore := po.evaluatePlacement(optimizedBase, items, config)

	temperature := config.Temperature

	for iteration := 0; iteration < config.MaxIterations; iteration++ {
		// Create a new candidate by perturbing the current placement
		candidateBase := optimizedBase.Clone()
		po.perturbPlacement(candidateBase, items)

		// Evaluate the candidate
		candidateScore := po.evaluatePlacement(candidateBase, items, config)

		// Accept or reject based on simulated annealing
		if po.shouldAccept(bestScore.TotalScore, candidateScore.TotalScore, temperature) {
			optimizedBase = candidateBase

			// Update best if this is better
			if candidateScore.TotalScore > bestScore.TotalScore {
				bestBase = candidateBase.Clone()
				bestScore = candidateScore
			}
		}

		// Cool down
		temperature *= config.CoolingRate
		if temperature < config.MinTemperature {
			break
		}
	}

	return bestBase, bestScore, nil
}

// placeItemsGreedy places items using a greedy algorithm
func (po *PlacementOptimizer) placeItemsGreedy(base *types.Base, items []*types.Item) {
	for _, item := range items {
		bestPosition := po.findBestPosition(base, item)
		if bestPosition != nil {
			item.Position = *bestPosition
			base.PlaceItem(item)
		}
	}
}

// findBestPosition finds the best position for an item
func (po *PlacementOptimizer) findBestPosition(base *types.Base, item *types.Item) *types.Position {
	var bestPosition *types.Position
	bestScore := math.Inf(-1)

	// Try different positions
	freePositions := base.GetFreePositions()
	for _, pos := range freePositions {
		// Check if item can be placed here
		testItem := &types.Item{
			ID:       item.ID,
			Type:     item.Type,
			Position: pos,
			Bounds:   item.Bounds,
			Rotation: item.Rotation,
			Priority: item.Priority,
		}

		if base.CanPlaceItem(testItem) {
			score := po.evaluateItemPosition(base, testItem)
			if score > bestScore {
				bestScore = score
				bestPosition = &pos
			}
		}
	}

	return bestPosition
}

// evaluateItemPosition evaluates how good a position is for an item
func (po *PlacementOptimizer) evaluateItemPosition(base *types.Base, item *types.Item) float64 {
	score := 0.0

	// Prefer positions near the center for important items
	if item.Type == types.ItemTypePalbox {
		center := types.Position{X: base.Width / 2, Y: 0, Z: base.Depth / 2}
		distance := item.Position.Distance(center)
		score += 100.0 / (1.0 + distance)
	}

	// Prefer positions near related items
	score += po.evaluateProximityToRelatedItems(base, item)

	// Prefer positions that don't block paths
	score += po.evaluatePathAccessibility(base, item)

	return score
}

// evaluateProximityToRelatedItems evaluates proximity to related items
func (po *PlacementOptimizer) evaluateProximityToRelatedItems(base *types.Base, item *types.Item) float64 {
	score := 0.0

	// Define related item types
	relatedItems := po.getRelatedItemTypes(item.Type)

	for _, existingItem := range base.Items {
		if relatedItems[existingItem.Type] {
			distance := item.Position.Distance(existingItem.Position)
			score += 10.0 / (1.0 + distance)
		}
	}

	return score
}

// getRelatedItemTypes returns item types that are related to the given type
func (po *PlacementOptimizer) getRelatedItemTypes(itemType types.ItemType) map[types.ItemType]bool {
	related := make(map[types.ItemType]bool)

	switch itemType {
	case types.ItemTypeFoodBox:
		related[types.ItemTypeFoodPlot] = true
		related[types.ItemTypeCookingPot] = true
	case types.ItemTypeFoodPlot:
		related[types.ItemTypeFoodBox] = true
		related[types.ItemTypeCookingPot] = true
	case types.ItemTypePowerGenerator:
		related[types.ItemTypeAccumulator] = true
		related[types.ItemTypeWorkbench] = true
	case types.ItemTypeWorkbench:
		related[types.ItemTypePowerGenerator] = true
		related[types.ItemTypeStorage] = true
	case types.ItemTypeStorage:
		related[types.ItemTypeWorkbench] = true
		related[types.ItemTypeFurnace] = true
	}

	return related
}

// evaluatePathAccessibility evaluates how well an item placement maintains path accessibility
func (po *PlacementOptimizer) evaluatePathAccessibility(base *types.Base, item *types.Item) float64 {
	score := 0.0

	// Check if placement creates isolated areas
	isolatedPenalty := po.calculateIsolationPenalty(base, item)
	score -= isolatedPenalty

	// Check if placement blocks important paths
	blockingPenalty := po.calculateBlockingPenalty(base, item)
	score -= blockingPenalty

	return score
}

// calculateIsolationPenalty calculates penalty for creating isolated areas
func (po *PlacementOptimizer) calculateIsolationPenalty(base *types.Base, item *types.Item) float64 {
	// This is a simplified calculation
	// In a full implementation, you would use flood fill or connected components
	// to detect isolated areas
	return 0.0
}

// calculateBlockingPenalty calculates penalty for blocking important paths
func (po *PlacementOptimizer) calculateBlockingPenalty(base *types.Base, item *types.Item) float64 {
	penalty := 0.0

	// Check if item blocks access to important items
	for _, existingItem := range base.Items {
		if existingItem.Type == types.ItemTypePalbox {
			// Check if path to Palbox is blocked
			path, err := po.Graph.FindPath(item.Position, existingItem.Position)
			if err != nil {
				penalty += 50.0 // High penalty for blocking Palbox access
			} else {
				// Lower penalty for longer paths
				penalty += path.Cost * 0.1
			}
		}
	}

	return penalty
}

// perturbPlacement creates a perturbation of the current placement
func (po *PlacementOptimizer) perturbPlacement(base *types.Base, items []*types.Item) {
	// Randomly select an item to move
	if len(items) == 0 {
		return
	}

	itemIndex := rand.Intn(len(items))
	item := items[itemIndex]

	// Remove the item
	base.RemoveItem(item.ID)

	// Find a new position
	newPosition := po.findBestPosition(base, item)
	if newPosition != nil {
		item.Position = *newPosition
		base.PlaceItem(item)
	}
}

// shouldAccept determines if a candidate should be accepted in simulated annealing
func (po *PlacementOptimizer) shouldAccept(currentScore, candidateScore, temperature float64) bool {
	if candidateScore > currentScore {
		return true
	}

	// Calculate acceptance probability
	delta := candidateScore - currentScore
	probability := math.Exp(delta / temperature)

	return rand.Float64() < probability
}

// evaluatePlacement evaluates the overall quality of a placement
func (po *PlacementOptimizer) evaluatePlacement(base *types.Base, items []*types.Item, config *OptimizationConfig) *PlacementScore {
	score := &PlacementScore{
		Details: make(map[string]float64),
	}

	// Evaluate pathfinding efficiency
	pathfindingScore := po.evaluatePathfinding(base, items)
	score.PathfindingScore = pathfindingScore

	// Evaluate efficiency (proximity of related items)
	efficiencyScore := po.evaluateEfficiency(base, items)
	score.EfficiencyScore = efficiencyScore

	// Evaluate compactness
	compactnessScore := po.evaluateCompactness(base)
	score.CompactnessScore = compactnessScore

	// Calculate weighted total score
	score.TotalScore = config.PathfindingWeight*pathfindingScore +
		config.EfficiencyWeight*efficiencyScore +
		config.CompactnessWeight*compactnessScore

	// Store detailed scores
	score.Details["pathfinding"] = pathfindingScore
	score.Details["efficiency"] = efficiencyScore
	score.Details["compactness"] = compactnessScore

	return score
}

// evaluatePathfinding evaluates the pathfinding efficiency of the placement
func (po *PlacementOptimizer) evaluatePathfinding(base *types.Base, items []*types.Item) float64 {
	score := 0.0

	// Find the Palbox
	var palbox *types.Item
	for _, item := range base.Items {
		if item.Type == types.ItemTypePalbox {
			palbox = item
			break
		}
	}

	if palbox == nil {
		return 0.0
	}

	// Evaluate paths from Palbox to all other items
	for _, item := range base.Items {
		if item.ID == palbox.ID {
			continue
		}

		path, err := po.Graph.FindPath(palbox.Position, item.Position)
		if err == nil {
			// Shorter paths are better
			score += 100.0 / (1.0 + path.Cost)
		} else {
			// Penalty for unreachable items
			score -= 50.0
		}
	}

	return score
}

// evaluateEfficiency evaluates the efficiency of item placement
func (po *PlacementOptimizer) evaluateEfficiency(base *types.Base, items []*types.Item) float64 {
	score := 0.0

	for _, item := range base.Items {
		relatedItems := po.getRelatedItemTypes(item.Type)

		for _, otherItem := range base.Items {
			if item.ID == otherItem.ID {
				continue
			}

			if relatedItems[otherItem.Type] {
				distance := item.Position.Distance(otherItem.Position)
				score += 20.0 / (1.0 + distance)
			}
		}
	}

	return score
}

// evaluateCompactness evaluates how compact the placement is
func (po *PlacementOptimizer) evaluateCompactness(base *types.Base) float64 {
	// Calculate the bounding box of all items
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	minZ, maxZ := math.Inf(1), math.Inf(-1)

	for _, item := range base.Items {
		for _, pos := range item.GetOccupiedPositions() {
			minX = math.Min(minX, float64(pos.X))
			maxX = math.Max(maxX, float64(pos.X))
			minY = math.Min(minY, float64(pos.Y))
			maxY = math.Max(maxY, float64(pos.Y))
			minZ = math.Min(minZ, float64(pos.Z))
			maxZ = math.Max(maxZ, float64(pos.Z))
		}
	}

	// Calculate volume of bounding box
	volume := (maxX - minX) * (maxY - minY) * (maxZ - minZ)

	// Calculate total item volume
	totalItemVolume := 0.0
	for _, item := range base.Items {
		totalItemVolume += float64(item.Bounds.Volume())
	}

	// Compactness is the ratio of item volume to bounding box volume
	if volume > 0 {
		return totalItemVolume / volume
	}

	return 0.0
}
