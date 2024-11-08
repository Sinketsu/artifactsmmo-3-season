package generic

import (
	"log/slog"
	"unsafe"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Character struct {
	name  string
	state oas.CharacterSchema

	cli    *api.Client
	logger *slog.Logger
}

func NewCharacter(name string, client *api.Client) *Character {
	return &Character{
		name: name,

		cli:    client,
		logger: slog.Default().With(ycloggingslog.Stream, name),
	}
}

func (c *Character) syncState(p unsafe.Pointer) {
	// tricky hack, because `ogen` generates different models for Character state from different methods instead of reusing one. But fields are the same - so we can cast it
	c.state = *(*oas.CharacterSchema)(p)

	characterLevel.Set(int64(c.state.Level), c.name)

	skillLevel.Set(int64(c.state.AlchemyLevel), c.name, "alchemy")
	skillLevel.Set(int64(c.state.MiningLevel), c.name, "mining")
	skillLevel.Set(int64(c.state.WoodcuttingLevel), c.name, "woodcutting")
	skillLevel.Set(int64(c.state.FishingLevel), c.name, "fishing")
	skillLevel.Set(int64(c.state.CookingLevel), c.name, "cooking")
	skillLevel.Set(int64(c.state.GearcraftingLevel), c.name, "gearcrafting")
	skillLevel.Set(int64(c.state.JewelrycraftingLevel), c.name, "jewelrycrafting")
}

func (c *Character) InventorySpace() (freeSpace int, freeSlots int) {
	freeSpace = c.state.InventoryMaxItems
	freeSlots = 20

	for _, count := range c.Inventory() {
		freeSlots--
		freeSpace -= count
	}

	return
}

func (c *Character) InventoryFull() bool {
	space, slots := c.InventorySpace()

	return space < 5 || slots < 3
}

func (c *Character) Inventory() map[string]int {
	result := map[string]int{}
	for _, slot := range c.state.Inventory {
		if slot.Code != "" {
			result[slot.Code] = slot.Quantity
		}
	}

	return result
}

func (c *Character) Gold() int {
	return c.state.Gold
}

func (c *Character) HealthPercent() float64 {
	return float64(c.state.Hp) / float64(c.state.MaxHp) * 100
}

func (c *Character) Level() int {
	return c.state.Level
}

func (c *Character) Skills() map[string]int {
	return map[string]int{
		string(oas.CraftSchemaSkillAlchemy):         c.state.AlchemyLevel,
		string(oas.CraftSchemaSkillCooking):         c.state.CookingLevel,
		string(oas.CraftSchemaSkillWeaponcrafting):  c.state.WeaponcraftingLevel,
		string(oas.CraftSchemaSkillGearcrafting):    c.state.GearcraftingLevel,
		string(oas.CraftSchemaSkillJewelrycrafting): c.state.JewelrycraftingLevel,
		string(oas.CraftSchemaSkillWoodcutting):     c.state.WeaponcraftingLevel,
		string(oas.CraftSchemaSkillMining):          c.state.MiningLevel,
		string(oas.ResourceSchemaSkillFishing):      c.state.FishingLevel,
	}
}
