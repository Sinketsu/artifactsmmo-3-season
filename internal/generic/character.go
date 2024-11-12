package generic

import (
	"context"
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
	c := &Character{
		name: name,

		cli:    client,
		logger: slog.Default().With(ycloggingslog.Stream, name),
	}

	c.init()

	return c
}

func (c *Character) init() {
	for {
		resp, err := c.cli.GetCharacterCharactersNameGet(context.Background(), oas.GetCharacterCharactersNameGetParams{Name: c.name})
		if err != nil {
			c.logger.Warn("fail to init character data", "error", err)
			continue
		}

		switch v := resp.(type) {
		case *oas.CharacterResponseSchema:
			c.syncState(unsafe.Pointer(&v.Data))
		default:
			c.logger.Warn("fail to init character data: got unexpected response", "response", v)
			continue
		}

		break
	}
}

func (c *Character) syncState(p unsafe.Pointer) {
	for _, slot := range c.state.Inventory {
		if slot.Code != "" {
			itemCount.Reset(slot.Code)
		}
	}

	// tricky hack, because `ogen` generates different models for Character state from different methods instead of reusing one. But fields are the same - so we can cast it
	c.state = *(*oas.CharacterSchema)(p)

	characterLevel.Set(int64(c.state.Level), c.name)

	skillLevel.Set(int64(c.state.AlchemyLevel), c.name, string(oas.CraftSchemaSkillAlchemy))
	skillLevel.Set(int64(c.state.MiningLevel), c.name, string(oas.CraftSchemaSkillMining))
	skillLevel.Set(int64(c.state.WoodcuttingLevel), c.name, string(oas.CraftSchemaSkillWoodcutting))
	skillLevel.Set(int64(c.state.FishingLevel), c.name, string(oas.ResourceSchemaSkillFishing))
	skillLevel.Set(int64(c.state.CookingLevel), c.name, string(oas.CraftSchemaSkillCooking))
	skillLevel.Set(int64(c.state.WeaponcraftingLevel), c.name, string(oas.CraftSchemaSkillWeaponcrafting))
	skillLevel.Set(int64(c.state.GearcraftingLevel), c.name, string(oas.CraftSchemaSkillGearcrafting))
	skillLevel.Set(int64(c.state.JewelrycraftingLevel), c.name, string(oas.CraftSchemaSkillJewelrycrafting))

	goldCount.Set(int64(c.state.Gold), c.name)

	for _, slot := range c.state.Inventory {
		if slot.Code != "" {
			itemCount.Set(int64(slot.Quantity), c.name, slot.Code)
		}
	}
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
		string(oas.CraftSchemaSkillWoodcutting):     c.state.WoodcuttingLevel,
		string(oas.CraftSchemaSkillMining):          c.state.MiningLevel,
		string(oas.ResourceSchemaSkillFishing):      c.state.FishingLevel,
	}
}
