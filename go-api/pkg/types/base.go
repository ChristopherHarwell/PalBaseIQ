package types

import (
	"fmt"
	"math"
)

// Position represents a 3D coordinate in the base
//
// It remains unchanged from the original implementation, providing
// convenience methods for distance calculations and string
// representations. See the original code for details on Euclidean
// and Manhattan distance helpers.
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
//
// Width, Height and Depth define how many grid cells an item occupies
// along each axis when placed in a base.
type BoundingBox struct {
	Width, Height, Depth int
}

// Volume returns the total volume of the bounding box
func (bb BoundingBox) Volume() int {
	return bb.Width * bb.Height * bb.Depth
}

// StructureCategory enumerates the high level categories used by
// Palworld.gg to group structures.  These values map directly to
// the options in the "Structure Type" dropdown on the website and
// should not be renamed without also updating any corresponding
// constants.
//
// The names follow the site’s capitalization and spacing to ease
// serialization and comparison.  See
// https://palworld.gg/structures for the authoritative list.
type StructureCategory string

const (
	StructureCategoryFood           StructureCategory = "Food"
	StructureCategoryFoundation     StructureCategory = "Foundation"
	StructureCategoryDefense        StructureCategory = "Defense"
	StructureCategoryInfrastructure StructureCategory = "Infrastructure"
	StructureCategoryStorage        StructureCategory = "Storage"
	StructureCategoryPals           StructureCategory = "Pals"
	StructureCategoryLight          StructureCategory = "Light"
	StructureCategoryProduction     StructureCategory = "Production"
	StructureCategoryFurniture      StructureCategory = "Furniture"
	StructureCategoryOther          StructureCategory = "Other"
)

// StructureName identifies a specific placeable item within a
// StructureCategory.  This type replaces the monolithic StructureType
// used previously, allowing items to be grouped by category without
// conflating the category and the specific item.  When adding new
// structures, prefer using the canonical name from Palworld.gg and
// group it with the appropriate category in StructureDefinitions.
type StructureName string

const (
	// Food structures
	StructureNameCampfire         StructureName = "campfire"
	StructureNameCookingPot       StructureName = "cooking_pot"
	StructureNameColdFoodBox      StructureName = "cold_food_box"
	StructureNameElectricKitchen  StructureName = "electric_kitchen"
	StructureNameBerryPlantation  StructureName = "berry_plantation"
	StructureNameCarrotPlantation StructureName = "carrot_plantation"

	// Foundation/defense structures
	StructureNameStoneDefensiveWall  StructureName = "stone_defensive_wall"
	StructureNameMetalDefensiveWall  StructureName = "metal_defensive_wall"
	StructureNameWoodenDefensiveWall StructureName = "wooden_defensive_wall"
	StructureNameGlassWallAndDoor    StructureName = "glass_wall_and_door"
	StructureNameGlassFence          StructureName = "glass_fence"
	StructureNameGlassSlantedRoof    StructureName = "glass_slanted_roof"

	// Product/production structures
	StructureNameProductionAssemblyLineII     StructureName = "production_assembly_line_ii"
	StructureNameAdvancedCivilizationWorkshop StructureName = "advanced_civilization_workshop"
	StructureNameGoldCoinAssemblyLine         StructureName = "gold_coin_assembly_line"

	// Furniture
	StructureNameJapanesePaperLantern  StructureName = "japanese_paper_lantern"
	StructureNameRedMetalBarrel        StructureName = "red_metal_barrel"
	StructureNameBlueMetalBarrel       StructureName = "blue_metal_barrel"
	StructureNameGreenMetalBarrel      StructureName = "green_metal_barrel"
	StructureNameAntiqueBathtub        StructureName = "antique_bathtub"
	StructureNameFreePalAllianceBanner StructureName = "free_pal_alliance_banner"

	// Storage
	StructureNameWoodenBarrel         StructureName = "wooden_barrel"
	StructureNameItemRetrievalMachine StructureName = "item_retrieval_machine"

	// Pals-related structures
	StructureNameMonitoringStand     StructureName = "monitoring_stand"
	StructureNamePalboxControlDevice StructureName = "palbox_control_device"
	StructureNamePalBed              StructureName = "pal_bed"
	StructureNamePalSphereWorkbench  StructureName = "pal_sphere_workbench"
	StructureNamePalbox              StructureName = "palbox"

	// Other existing structures from the original code
	StructureNameFoodBox                   StructureName = "food_box"
	StructureNameFoodPlot                  StructureName = "food_plot"
	StructureNamePowerGenerator            StructureName = "power_generator"
	StructureNameAccumulator               StructureName = "accumulator"
	StructureNameOuterWall                 StructureName = "outer_wall"
	StructureNameWorkbench                 StructureName = "workbench"
	StructureNameStorage                   StructureName = "storage"
	StructureNameFurnace                   StructureName = "furnace"
	StructureNameMedievalMedicineWorkbench StructureName = "medieval_medicine_workbench"
	StructureNameElectricMedicineWorkbench StructureName = "electric_medicine_workbench"
	StructureNameAdvancedMedicineWorkbench StructureName = "advanced_medicine_workbench"
	StructureNameBreedingFarm              StructureName = "breeding_farm"
	StructureNameIncubator                 StructureName = "incubator"
)

// StructureDefinition captures metadata for a structure, including
// its canonical name, high-level category, human-readable description,
// build work (abstract work units), and material costs (by material name).
//
// Use canonical names from Palworld.gg for both name and category fields.
type StructureDefinition struct {
	Name         StructureName
	Category     StructureCategory
	Description  string
	BuildWork    int
	MaterialCost map[string]int
}

// StructureDefinitions maps each StructureName to its StructureDefinition.
// When adding new structures, append new entries here.
var StructureDefinitions = map[StructureName]StructureDefinition{
	// Food
	StructureNameCampfire:         {Name: StructureNameCampfire, Category: StructureCategoryFood},
	StructureNameCookingPot:       {Name: StructureNameCookingPot, Category: StructureCategoryFood},
	StructureNameColdFoodBox:      {Name: StructureNameColdFoodBox, Category: StructureCategoryFood},
	StructureNameElectricKitchen:  {Name: StructureNameElectricKitchen, Category: StructureCategoryFood},
	StructureNameBerryPlantation:  {Name: StructureNameBerryPlantation, Category: StructureCategoryFood},
	StructureNameCarrotPlantation: {Name: StructureNameCarrotPlantation, Category: StructureCategoryFood},

	// Foundation/Defense
	StructureNameStoneDefensiveWall:  {Name: StructureNameStoneDefensiveWall, Category: StructureCategoryFoundation},
	StructureNameMetalDefensiveWall:  {Name: StructureNameMetalDefensiveWall, Category: StructureCategoryFoundation},
	StructureNameWoodenDefensiveWall: {Name: StructureNameWoodenDefensiveWall, Category: StructureCategoryFoundation},
	StructureNameGlassWallAndDoor:    {Name: StructureNameGlassWallAndDoor, Category: StructureCategoryFoundation},
	StructureNameGlassFence:          {Name: StructureNameGlassFence, Category: StructureCategoryFoundation},
	StructureNameGlassSlantedRoof:    {Name: StructureNameGlassSlantedRoof, Category: StructureCategoryFoundation},

	// Product/production
	StructureNameProductionAssemblyLineII:     {Name: StructureNameProductionAssemblyLineII, Category: StructureCategoryProduction},
	StructureNameAdvancedCivilizationWorkshop: {Name: StructureNameAdvancedCivilizationWorkshop, Category: StructureCategoryProduction},
	StructureNameGoldCoinAssemblyLine:         {Name: StructureNameGoldCoinAssemblyLine, Category: StructureCategoryProduction},

	// Furniture
	StructureNameJapanesePaperLantern:  {Name: StructureNameJapanesePaperLantern, Category: StructureCategoryFurniture},
	StructureNameRedMetalBarrel:        {Name: StructureNameRedMetalBarrel, Category: StructureCategoryFurniture},
	StructureNameBlueMetalBarrel:       {Name: StructureNameBlueMetalBarrel, Category: StructureCategoryFurniture},
	StructureNameGreenMetalBarrel:      {Name: StructureNameGreenMetalBarrel, Category: StructureCategoryFurniture},
	StructureNameAntiqueBathtub:        {Name: StructureNameAntiqueBathtub, Category: StructureCategoryFurniture},
	StructureNameFreePalAllianceBanner: {Name: StructureNameFreePalAllianceBanner, Category: StructureCategoryFurniture},

	// Storage
	StructureNameWoodenBarrel:         {Name: StructureNameWoodenBarrel, Category: StructureCategoryStorage},
	StructureNameItemRetrievalMachine: {Name: StructureNameItemRetrievalMachine, Category: StructureCategoryStorage},

	// Pals
	StructureNameMonitoringStand:     {Name: StructureNameMonitoringStand, Category: StructureCategoryPals},
	StructureNamePalboxControlDevice: {Name: StructureNamePalboxControlDevice, Category: StructureCategoryPals},
	StructureNamePalBed:              {Name: StructureNamePalBed, Category: StructureCategoryPals},
	StructureNamePalSphereWorkbench:  {Name: StructureNamePalSphereWorkbench, Category: StructureCategoryPals},
	StructureNamePalbox:              {Name: StructureNamePalbox, Category: StructureCategoryPals},

	// Other miscellaneous items from original code
	StructureNameFoodBox:                   {Name: StructureNameFoodBox, Category: StructureCategoryFood},
	StructureNameFoodPlot:                  {Name: StructureNameFoodPlot, Category: StructureCategoryFood},
	StructureNamePowerGenerator:            {Name: StructureNamePowerGenerator, Category: StructureCategoryInfrastructure},
	StructureNameAccumulator:               {Name: StructureNameAccumulator, Category: StructureCategoryInfrastructure},
	StructureNameOuterWall:                 {Name: StructureNameOuterWall, Category: StructureCategoryFoundation},
	StructureNameWorkbench:                 {Name: StructureNameWorkbench, Category: StructureCategoryProduction},
	StructureNameStorage:                   {Name: StructureNameStorage, Category: StructureCategoryStorage},
	StructureNameFurnace:                   {Name: StructureNameFurnace, Category: StructureCategoryProduction},
	StructureNameMedievalMedicineWorkbench: {Name: StructureNameMedievalMedicineWorkbench, Category: StructureCategoryProduction},
	StructureNameElectricMedicineWorkbench: {Name: StructureNameElectricMedicineWorkbench, Category: StructureCategoryProduction},
	StructureNameAdvancedMedicineWorkbench: {Name: StructureNameAdvancedMedicineWorkbench, Category: StructureCategoryProduction},
	StructureNameBreedingFarm:              {Name: StructureNameBreedingFarm, Category: StructureCategoryPals},
	StructureNameIncubator:                 {Name: StructureNameIncubator, Category: StructureCategoryPals},
}

// Item represents a placeable item in the base.
//
// It now uses StructureName instead of the old StructureType so that
// the specific item can be cross-referenced with its category via
// StructureDefinitions.  The rest of the fields are unchanged.
type Item struct {
	ID       string
	Type     StructureName
	Position Position
	Bounds   BoundingBox
	Rotation int // 0, 90, 180, 270 degrees
	Priority int // Higher priority items are placed first
}

// String returns a string representation of the item, including its
// structure type and identifier.  The category can be obtained by
// looking up the Type in StructureDefinitions.
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
//
// The structure and behaviour of Base are unchanged from the
// original implementation.  It uses a 3D boolean grid to mark
// occupied positions and maintains a map of items by their ID.
// New items created with the updated StructureName types are fully
// compatible.
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

// abs is a helper function for ManhattanDistance.  It is defined
// outside of any type to avoid creating methods on basic types.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
