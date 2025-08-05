package main

import (
	"fmt"
	"log"
	"palbaseiq/pkg/optimizer"
	"palbaseiq/pkg/types"
)

func main() {
	fmt.Println("PalBaseIQ - Palworld Base Optimization System")
	fmt.Println("=============================================")

	// Create a base with dimensions based on your layout
	// Assuming a 20x16x20 base (width x height x depth)
	base := types.NewBase(20, 16, 20)

	// Define items to place in the base
	items := createBaseItems()

	// Create the placement optimizer
	opt := optimizer.NewPlacementOptimizer(base)

	// Configure optimization parameters
	config := optimizer.DefaultConfig()
	config.MaxIterations = 500 // Reduced for faster demo
	config.PathfindingWeight = 0.4
	config.EfficiencyWeight = 0.3
	config.CompactnessWeight = 0.3

	fmt.Println("Starting base optimization...")
	fmt.Printf("Base dimensions: %dx%dx%d\n", base.Width, base.Height, base.Depth)
	fmt.Printf("Items to place: %d\n", len(items))
	fmt.Printf("Optimization iterations: %d\n", config.MaxIterations)

	// Run optimization
	optimizedBase, score, err := opt.OptimizePlacement(items, config)
	if err != nil {
		log.Fatalf("Optimization failed: %v", err)
	}

	// Display results
	fmt.Println("\nOptimization Results:")
	fmt.Println("====================")
	fmt.Printf("Total Score: %.2f\n", score.TotalScore)
	fmt.Printf("Pathfinding Score: %.2f\n", score.PathfindingScore)
	fmt.Printf("Efficiency Score: %.2f\n", score.EfficiencyScore)
	fmt.Printf("Compactness Score: %.2f\n", score.CompactnessScore)
	fmt.Printf("Occupancy: %.1f%%\n", optimizedBase.GetOccupancyPercentage())

	// Display item placements
	fmt.Println("\nOptimized Item Placements:")
	fmt.Println("==========================")
	for _, item := range optimizedBase.Items {
		fmt.Printf("%s: %s (Priority: %d)\n", item.Type, item.Position, item.Priority)
	}

	// Analyze pathfinding
	analyzePathfinding(optimizedBase, opt)

	fmt.Println("\nOptimization complete!")
}

// createBaseItems creates a list of items to place in the base
func createBaseItems() []*types.Item {
	items := []*types.Item{
		// Palbox (highest priority - must be placed first)
		{
			ID:       "palbox_1",
			Type:     types.ItemTypePalbox,
			Bounds:   types.BoundingBox{Width: 2, Height: 2, Depth: 2},
			Priority: 100,
		},

		// Pal Beds (32 beds as shown in your layout)
		{
			ID:       "pal_bed_1",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},
		{
			ID:       "pal_bed_2",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},
		{
			ID:       "pal_bed_3",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},
		{
			ID:       "pal_bed_4",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},
		{
			ID:       "pal_bed_5",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},
		{
			ID:       "pal_bed_6",
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		},

		// Food Box
		{
			ID:       "food_box_1",
			Type:     types.ItemTypeFoodBox,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 80,
		},

		// Food Plots (4 plots as shown in your layout)
		{
			ID:       "food_plot_1",
			Type:     types.ItemTypeFoodPlot,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 75,
		},
		{
			ID:       "food_plot_2",
			Type:     types.ItemTypeFoodPlot,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 75,
		},
		{
			ID:       "food_plot_3",
			Type:     types.ItemTypeFoodPlot,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 75,
		},
		{
			ID:       "food_plot_4",
			Type:     types.ItemTypeFoodPlot,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 75,
		},

		// Power Generator
		{
			ID:       "power_generator_1",
			Type:     types.ItemTypePowerGenerator,
			Bounds:   types.BoundingBox{Width: 1, Height: 2, Depth: 1},
			Priority: 85,
		},

		// Accumulator
		{
			ID:       "accumulator_1",
			Type:     types.ItemTypeAccumulator,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 80,
		},

		// Additional items for a more complete base
		{
			ID:       "workbench_1",
			Type:     types.ItemTypeWorkbench,
			Bounds:   types.BoundingBox{Width: 2, Height: 1, Depth: 1},
			Priority: 70,
		},
		{
			ID:       "storage_1",
			Type:     types.ItemTypeStorage,
			Bounds:   types.BoundingBox{Width: 1, Height: 2, Depth: 1},
			Priority: 65,
		},
		{
			ID:       "furnace_1",
			Type:     types.ItemTypeFurnace,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 60,
		},
		{
			ID:       "cooking_pot_1",
			Type:     types.ItemTypeCookingPot,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 70,
		},
	}

	// Add more Pal Beds to reach 32 total
	for i := 7; i <= 32; i++ {
		items = append(items, &types.Item{
			ID:       fmt.Sprintf("pal_bed_%d", i),
			Type:     types.ItemTypePalBed,
			Bounds:   types.BoundingBox{Width: 1, Height: 1, Depth: 1},
			Priority: 90,
		})
	}

	return items
}

// analyzePathfinding analyzes the pathfinding efficiency of the optimized base
func analyzePathfinding(base *types.Base, optimizer *optimizer.PlacementOptimizer) {
	fmt.Println("\nPathfinding Analysis:")
	fmt.Println("=====================")

	// Find the Palbox
	var palbox *types.Item
	for _, item := range base.Items {
		if item.Type == types.ItemTypePalbox {
			palbox = item
			break
		}
	}

	if palbox == nil {
		fmt.Println("No Palbox found!")
		return
	}

	fmt.Printf("Palbox location: %s\n", palbox.Position)

	// Analyze paths to key items
	keyItems := []types.ItemType{
		types.ItemTypeFoodBox,
		types.ItemTypePowerGenerator,
		types.ItemTypeWorkbench,
		types.ItemTypeStorage,
	}

	totalPathCost := 0.0
	reachableItems := 0

	for _, itemType := range keyItems {
		for _, item := range base.Items {
			if item.Type == itemType {
				path, err := optimizer.Graph.FindPath(palbox.Position, item.Position)
				if err == nil {
					fmt.Printf("Path to %s: %.2f cost (%d steps)\n", itemType, path.Cost, len(path.Nodes))
					totalPathCost += path.Cost
					reachableItems++
				} else {
					fmt.Printf("Path to %s: UNREACHABLE\n", itemType)
				}
				break // Only check first item of each type
			}
		}
	}

	if reachableItems > 0 {
		avgPathCost := totalPathCost / float64(reachableItems)
		fmt.Printf("Average path cost: %.2f\n", avgPathCost)
		fmt.Printf("Reachable items: %d/%d\n", reachableItems, len(keyItems))
	}
}

// visualizeBase creates a simple text visualization of the base
func visualizeBase(base *types.Base) {
	fmt.Println("\nBase Visualization (Top-down view at Y=0):")
	fmt.Println("=========================================")

	// Create a 2D representation of the base at ground level
	grid := make([][]string, base.Width)
	for x := range grid {
		grid[x] = make([]string, base.Depth)
		for z := range grid[x] {
			grid[x][z] = "."
		}
	}

	// Mark occupied positions
	for _, item := range base.Items {
		for _, pos := range item.GetOccupiedPositions() {
			if pos.Y == 0 && pos.X >= 0 && pos.X < base.Width && pos.Z >= 0 && pos.Z < base.Depth {
				// Use different symbols for different item types
				symbol := "X"
				switch item.Type {
				case types.ItemTypePalbox:
					symbol = "P"
				case types.ItemTypePalBed:
					symbol = "B"
				case types.ItemTypeFoodBox:
					symbol = "F"
				case types.ItemTypeFoodPlot:
					symbol = "G"
				case types.ItemTypePowerGenerator:
					symbol = "E"
				case types.ItemTypeAccumulator:
					symbol = "A"
				case types.ItemTypeWorkbench:
					symbol = "W"
				case types.ItemTypeStorage:
					symbol = "S"
				}
				grid[pos.X][pos.Z] = symbol
			}
		}
	}

	// Print the grid
	for z := 0; z < base.Depth; z++ {
		for x := 0; x < base.Width; x++ {
			fmt.Print(grid[x][z], " ")
		}
		fmt.Println()
	}

	fmt.Println("\nLegend:")
	fmt.Println("P = Palbox")
	fmt.Println("B = Pal Bed")
	fmt.Println("F = Food Box")
	fmt.Println("G = Food Plot")
	fmt.Println("E = Power Generator")
	fmt.Println("A = Accumulator")
	fmt.Println("W = Workbench")
	fmt.Println("S = Storage")
	fmt.Println(". = Empty space")
}
