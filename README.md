# PalBaseIQ - Palworld Base Optimization System

A sophisticated Go-based optimization system for Palworld base layout and item placement. This system uses advanced algorithms including A* pathfinding and simulated annealing to create optimal base layouts that maximize efficiency, accessibility, and space utilization.

## Features

### 🏗️ **3D Base Management**
- Full 3D coordinate system supporting Palworld's base dimensions
- Support for Palbox placement with 2x2x2 or 2x3x2 tile footprints
- Height limit support up to 16 tiles vertically
- Collision detection and spatial validation

### 🧭 **Advanced Pathfinding**
- A* algorithm implementation for optimal pathfinding
- 6-directional movement (up, down, left, right, forward, backward)
- Obstacle avoidance and terrain penalties
- Path cost optimization for Pal movement efficiency

### 🎯 **Intelligent Item Placement**
- Priority-based placement system
- Related item proximity optimization
- Workflow efficiency analysis
- Compactness and space utilization scoring

### 🔧 **Optimization Algorithms**
- Simulated annealing for global optimization
- Greedy initial placement for fast convergence
- Multi-objective scoring (pathfinding, efficiency, compactness)
- Configurable optimization parameters

## Supported Item Types

- **Palbox** - Central base hub (2x2x2 or 2x3x2)
- **Pal Beds** - Resting areas for Pals
- **Food Box** - Food storage
- **Food Plots** - Farming areas
- **Power Generator** - Energy production
- **Accumulator** - Energy storage
- **Workbench** - Crafting stations
- **Storage** - Item storage
- **Furnace** - Smelting operations
- **Cooking Pot** - Food preparation
- **Medicine Workbench** - Medical crafting
- **Breeding Farm** - Pal breeding
- **Incubator** - Egg hatching
- **Pal Sphere Workbench** - Sphere crafting

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd PalBaseIQ

# Install dependencies
go mod tidy

# Build the application
go build -o palbaseiq cmd/main.go

# Run the application
./palbaseiq
```

## Usage

### Basic Usage

```go
package main

import (
    "palbaseiq/pkg/types"
    "palbaseiq/pkg/optimizer"
)

func main() {
    // Create a base with dimensions
    base := types.NewBase(20, 16, 20) // width x height x depth
    
    // Define items to place
    items := []*types.Item{
        {
            ID: "palbox_1",
            Type: types.ItemTypePalbox,
            Bounds: types.BoundingBox{Width: 2, Height: 2, Depth: 2},
            Priority: 100, // Highest priority
        },
        // Add more items...
    }
    
    // Create optimizer
    opt := optimizer.NewPlacementOptimizer(base)
    
    // Configure optimization
    config := optimizer.DefaultConfig()
    config.MaxIterations = 1000
    
    // Run optimization
    optimizedBase, score, err := opt.OptimizePlacement(items, config)
    if err != nil {
        panic(err)
    }
    
    // Use optimized base
    fmt.Printf("Optimization score: %.2f\n", score.TotalScore)
}
```

### Advanced Configuration

```go
config := &optimizer.OptimizationConfig{
    MaxIterations:     2000,
    Temperature:       150.0,
    CoolingRate:       0.98,
    MinTemperature:    0.05,
    PathfindingWeight: 0.5,  // Emphasize pathfinding
    EfficiencyWeight:  0.3,  // Balance efficiency
    CompactnessWeight: 0.2,  // Less emphasis on compactness
}
```

## Architecture

### Core Components

1. **Types Package** (`pkg/types/`)
   - Base and Item data structures
   - 3D position and bounding box management
   - Spatial validation and collision detection

2. **Pathfinding Package** (`pkg/pathing/`)
   - A* algorithm implementation
   - Graph construction and traversal
   - Path cost calculation and optimization

3. **Optimizer Package** (`pkg/optimizer/`)
   - Simulated annealing optimization
   - Multi-objective scoring system
   - Placement evaluation and improvement

### Key Algorithms

#### A* Pathfinding
- **Heuristic Functions**: Manhattan distance, Euclidean distance
- **Cost Calculation**: Distance + terrain penalties + obstacle proximity
- **Path Optimization**: Minimizes total movement cost for Pals

#### Simulated Annealing
- **Temperature Schedule**: Exponential cooling with configurable parameters
- **Perturbation Strategy**: Random item relocation with greedy repositioning
- **Acceptance Criteria**: Boltzmann probability for uphill moves

#### Multi-Objective Scoring
- **Pathfinding Score**: Accessibility and movement efficiency
- **Efficiency Score**: Related item proximity and workflow optimization
- **Compactness Score**: Space utilization and layout density

## Performance Considerations

### Optimization Parameters
- **MaxIterations**: 500-2000 (trade-off between quality and speed)
- **Temperature**: 100-200 (higher = more exploration)
- **CoolingRate**: 0.95-0.99 (slower = more thorough search)

### Base Size Guidelines
- **Small Base** (10x16x10): ~100-500 iterations
- **Medium Base** (20x16x20): ~500-1000 iterations
- **Large Base** (30x16x30): ~1000-2000 iterations

## Example Output

```
PalBaseIQ - Palworld Base Optimization System
=============================================
Starting base optimization...
Base dimensions: 20x16x20
Items to place: 42
Optimization iterations: 500

Optimization Results:
====================
Total Score: 847.32
Pathfinding Score: 234.56
Efficiency Score: 312.78
Compactness Score: 299.98
Occupancy: 15.2%

Optimized Item Placements:
==========================
palbox: (10, 0, 10) (Priority: 100)
pal_bed_1: (5, 0, 5) (Priority: 90)
food_box_1: (8, 0, 8) (Priority: 80)
power_generator_1: (12, 0, 12) (Priority: 85)
...

Pathfinding Analysis:
====================
Palbox location: (10, 0, 10)
Path to food_box: 2.45 cost (5 steps)
Path to power_generator: 1.87 cost (3 steps)
Path to workbench: 3.12 cost (7 steps)
Average path cost: 2.48
Reachable items: 4/4
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Palworld's base building mechanics
- Uses established algorithms from computational geometry and optimization
- Built with Go for performance and simplicity 