# PalBaseIQ Setup Guide

## Prerequisites

### Installing Go
If Go is not installed on your system, you can install it using one of these methods:

#### macOS (using Homebrew)
```bash
brew install go
```

#### macOS (manual installation)
1. Download Go from https://golang.org/dl/
2. Extract to `/usr/local/go`
3. Add to PATH: `export PATH=$PATH:/usr/local/go/bin`

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
```

#### Windows
1. Download Go from https://golang.org/dl/
2. Run the installer
3. Follow the installation wizard

### Verify Installation
```bash
go version
```

## Building and Running PalBaseIQ

Once Go is installed, follow these steps:

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Build the Application
```bash
go build -o palbaseiq cmd/main.go
```

### 3. Run the Application
```bash
./palbaseiq
```

## Project Structure

```
PalBaseIQ/
├── cmd/
│   └── main.go              # Main application entry point
├── pkg/
│   ├── types/
│   │   └── base.go          # Core data structures and types
│   ├── pathing/
│   │   └── graph.go         # A* pathfinding implementation
│   └── optimizer/
│       └── placement.go     # Simulated annealing optimization
├── go.mod                   # Go module dependencies
├── README.md               # Comprehensive documentation
└── SETUP.md               # This setup guide
```

## Key Features Implemented

### 1. 3D Base Management
- **Position System**: 3D coordinates (X, Y, Z) for item placement
- **Bounding Box**: Volume-based collision detection
- **Palbox Support**: 2x2x2 or 2x3x2 tile footprints
- **Height Limits**: Up to 16 tiles vertically

### 2. Advanced Pathfinding
- **A* Algorithm**: Optimal pathfinding between any two points
- **6-Directional Movement**: Up, down, left, right, forward, backward
- **Obstacle Avoidance**: Automatic pathfinding around placed items
- **Cost Optimization**: Terrain penalties and proximity costs

### 3. Intelligent Optimization
- **Simulated Annealing**: Global optimization algorithm
- **Multi-Objective Scoring**: Pathfinding, efficiency, and compactness
- **Priority System**: Important items (like Palbox) placed first
- **Related Item Proximity**: Workflow optimization

### 4. Supported Item Types
- Palbox (central hub)
- Pal Beds (32 beds as in your layout)
- Food Box and Food Plots
- Power Generator and Accumulator
- Workbench, Storage, Furnace
- Cooking Pot and other crafting stations

## Example Usage

The main application demonstrates:
1. Creating a 20x16x20 base
2. Placing 42 items including 32 Pal Beds
3. Optimizing placement using simulated annealing
4. Analyzing pathfinding efficiency
5. Displaying optimization scores and results

## Configuration Options

You can customize the optimization by modifying:
- **MaxIterations**: Number of optimization iterations (500-2000)
- **Temperature**: Initial temperature for simulated annealing (100-200)
- **CoolingRate**: How fast temperature decreases (0.95-0.99)
- **Weights**: Balance between pathfinding, efficiency, and compactness

## Expected Output

When you run the application, you should see:
- Base optimization progress
- Final placement scores
- Item positions and priorities
- Pathfinding analysis
- Occupancy percentage

## Troubleshooting

### Common Issues

1. **"go: command not found"**
   - Install Go using the instructions above
   - Ensure Go is in your PATH

2. **"module not found"**
   - Run `go mod tidy` to download dependencies
   - Check that you're in the correct directory

3. **Build errors**
   - Ensure all files are in the correct package structure
   - Check that imports match the package names

### Performance Notes

- **Small bases** (10x16x10): ~100-500 iterations
- **Medium bases** (20x16x20): ~500-1000 iterations  
- **Large bases** (30x16x30): ~1000-2000 iterations

The optimization time scales with base size and iteration count.

## Next Steps

Once the basic system is running, you can:
1. Modify item types and dimensions
2. Adjust optimization parameters
3. Add new scoring criteria
4. Implement visualization features
5. Add web interface or API

## Support

If you encounter issues:
1. Check that Go is properly installed
2. Verify all files are in the correct locations
3. Ensure dependencies are downloaded
4. Review the error messages for specific issues 