package types

import (
	"fmt"
	"math"
)

// Position represents a 3D coordinate in the base
type Position struct {
	X, Y, Z int
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.X, p.Y, p.Z)
}

// Distance calculates the Euclidean distance between two positions
func (p Position) Distance(other Position) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	dz := float64(p.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// ManhattanDistance calculates the Manhattan distance between two positions
func (p Position) ManhattanDistance(other Position) int {
	return abs(p.X-other.X) + abs(p.Y-other.Y) + abs(p.Z-other.Z)
}

// BoundingBox represents the dimensions of an item
type BoundingBox struct {
	Width, Height, Depth int
}

// Volume returns the total volume of the bounding box
func (bb BoundingBox) Volume() int {
	return bb.Width * bb.Height * bb.Depth
}

// ItemType represents different types of items that can be placed in the base
type ItemType string

const (
	ItemTypePalbox             ItemType = "palbox"
	ItemTypePalBed             ItemType = "pal_bed"
	ItemTypeFoodBox            ItemType = "food_box"
	ItemTypeFoodPlot           ItemType = "food_plot"
	ItemTypePowerGenerator     ItemType = "power_generator"
	ItemTypeAccumulator        ItemType = "accumulator"
	ItemTypeOuterWall          ItemType = "outer_wall"
	ItemTypeWorkbench          ItemType = "workbench"
	ItemTypeStorage            ItemType = "storage"
	ItemTypeFurnace            ItemType = "furnace"
	ItemTypeCookingPot         ItemType = "cooking_pot"
	ItemTypeMedicineWorkbench  ItemType = "medicine_workbench"
	ItemTypeBreedingFarm       ItemType = "breeding_farm"
	ItemTypeIncubator          ItemType = "incubator"
	ItemTypePalSphereWorkbench ItemType = "pal_sphere_workbench"
)

// Item represents a placeable item in the base
type Item struct {
	ID       string
	Type     ItemType
	Position Position
	Bounds   BoundingBox
	Rotation int // 0, 90, 180, 270 degrees
	Priority int // Higher priority items are placed first
}

// String returns a string representation of the item
func (i Item) String() string {
	return fmt.Sprintf("%s[%s] at %s", i.Type, i.ID, i.Position)
}

// GetOccupiedPositions returns all positions occupied by this item
func (i Item) GetOccupiedPositions() []Position {
	positions := make([]Position, 0, i.Bounds.Volume())

	for x := 0; x < i.Bounds.Width; x++ {
		for y := 0; y < i.Bounds.Height; y++ {
			for z := 0; z < i.Bounds.Depth; z++ {
				positions = append(positions, Position{
					X: i.Position.X + x,
					Y: i.Position.Y + y,
					Z: i.Position.Z + z,
				})
			}
		}
	}

	return positions
}

// Intersects checks if this item intersects with another item
func (i Item) Intersects(other Item) bool {
	// Check if bounding boxes overlap
	return i.Position.X < other.Position.X+other.Bounds.Width &&
		i.Position.X+i.Bounds.Width > other.Position.X &&
		i.Position.Y < other.Position.Y+other.Bounds.Height &&
		i.Position.Y+i.Bounds.Height > other.Position.Y &&
		i.Position.Z < other.Position.Z+other.Bounds.Depth &&
		i.Position.Z+i.Bounds.Depth > other.Position.Z
}

// Base represents the entire base layout
type Base struct {
	Width  int
	Height int
	Depth  int
	Items  map[string]*Item
	Grid   [][][]bool // 3D grid representing occupied spaces
}

// NewBase creates a new base with the specified dimensions
func NewBase(width, height, depth int) *Base {
	// Initialize 3D grid
	grid := make([][][]bool, width)
	for x := range grid {
		grid[x] = make([][]bool, height)
		for y := range grid[x] {
			grid[x][y] = make([]bool, depth)
		}
	}

	return &Base{
		Width:  width,
		Height: height,
		Depth:  depth,
		Items:  make(map[string]*Item),
		Grid:   grid,
	}
}

// IsPositionValid checks if a position is within the base bounds
func (b *Base) IsPositionValid(pos Position) bool {
	return pos.X >= 0 && pos.X < b.Width &&
		pos.Y >= 0 && pos.Y < b.Height &&
		pos.Z >= 0 && pos.Z < b.Depth
}

// IsPositionOccupied checks if a position is occupied by any item
func (b *Base) IsPositionOccupied(pos Position) bool {
	if !b.IsPositionValid(pos) {
		return true // Invalid positions are considered occupied
	}
	return b.Grid[pos.X][pos.Y][pos.Z]
}

// CanPlaceItem checks if an item can be placed at the given position
func (b *Base) CanPlaceItem(item *Item) bool {
	// Check if all positions the item would occupy are valid and unoccupied
	for _, pos := range item.GetOccupiedPositions() {
		if b.IsPositionOccupied(pos) {
			return false
		}
	}
	return true
}

// PlaceItem places an item in the base
func (b *Base) PlaceItem(item *Item) error {
	if !b.CanPlaceItem(item) {
		return fmt.Errorf("cannot place item %s at position %s", item.ID, item.Position)
	}

	// Mark all occupied positions as occupied
	for _, pos := range item.GetOccupiedPositions() {
		b.Grid[pos.X][pos.Y][pos.Z] = true
	}

	b.Items[item.ID] = item
	return nil
}

// RemoveItem removes an item from the base
func (b *Base) RemoveItem(itemID string) error {
	item, exists := b.Items[itemID]
	if !exists {
		return fmt.Errorf("item %s not found", itemID)
	}

	// Mark all occupied positions as unoccupied
	for _, pos := range item.GetOccupiedPositions() {
		b.Grid[pos.X][pos.Y][pos.Z] = false
	}

	delete(b.Items, itemID)
	return nil
}

// GetItemAtPosition returns the item at the given position, if any
func (b *Base) GetItemAtPosition(pos Position) *Item {
	for _, item := range b.Items {
		for _, itemPos := range item.GetOccupiedPositions() {
			if itemPos == pos {
				return item
			}
		}
	}
	return nil
}

// GetOccupiedPositions returns all occupied positions in the base
func (b *Base) GetOccupiedPositions() []Position {
	var positions []Position
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			for z := 0; z < b.Depth; z++ {
				if b.Grid[x][y][z] {
					positions = append(positions, Position{X: x, Y: y, Z: z})
				}
			}
		}
	}
	return positions
}

// GetFreePositions returns all free positions in the base
func (b *Base) GetFreePositions() []Position {
	var positions []Position
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			for z := 0; z < b.Depth; z++ {
				if !b.Grid[x][y][z] {
					positions = append(positions, Position{X: x, Y: y, Z: z})
				}
			}
		}
	}
	return positions
}

// GetOccupancyPercentage returns the percentage of occupied space
func (b *Base) GetOccupancyPercentage() float64 {
	total := b.Width * b.Height * b.Depth
	occupied := len(b.GetOccupiedPositions())
	return float64(occupied) / float64(total) * 100
}

// Clone creates a deep copy of the base
func (b *Base) Clone() *Base {
	clone := NewBase(b.Width, b.Height, b.Depth)

	// Copy items
	for id, item := range b.Items {
		cloneItem := &Item{
			ID:       item.ID,
			Type:     item.Type,
			Position: item.Position,
			Bounds:   item.Bounds,
			Rotation: item.Rotation,
			Priority: item.Priority,
		}
		clone.Items[id] = cloneItem
	}

	// Copy grid
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			for z := 0; z < b.Depth; z++ {
				clone.Grid[x][y][z] = b.Grid[x][y][z]
			}
		}
	}

	return clone
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
