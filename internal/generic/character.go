package generic

import (
	"context"
	"log/slog"
	"time"
	"unsafe"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Character struct {
	name  string
	state oas.CharacterSchema
	point game.Point

	cli    *api.Client
	logger *slog.Logger
}

func NewCharacter(name string, client *api.Client) *Character {
	c := &Character{
		name: name,

		cli:    client,
		logger: slog.Default().With(ycloggingslog.Stream, name),
	}

	return c
}

func (c *Character) Init() {
	for {
		resp, err := c.cli.GetCharacterCharactersNameGet(context.Background(), oas.GetCharacterCharactersNameGetParams{Name: c.name})
		if err != nil {
			c.logger.Warn("fail to init character data", "error", err)
			continue
		}

		switch v := resp.(type) {
		case *oas.CharacterResponseSchema:
			c.syncState(unsafe.Pointer(&v.Data))
			if t, ok := c.state.CooldownExpiration.Get(); ok {
				if cooldown := time.Until(t); cooldown > 0 {
					c.logger.Warn("cooldown at init: " + cooldown.String())
					time.Sleep(cooldown)
				}
			}

		default:
			c.logger.Warn("fail to init character data: got unexpected response", "response", v)
			continue
		}

		break
	}
}

func (c *Character) syncState(p unsafe.Pointer) {
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

	itemCount.ResetAll()
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

func (c *Character) Health() (int, int) {
	return c.state.Hp, c.state.MaxHp
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

func (c *Character) Equiped() map[oas.EquipSchemaSlot]string {
	return map[oas.EquipSchemaSlot]string{
		oas.EquipSchemaSlotWeapon:    c.state.WeaponSlot,
		oas.EquipSchemaSlotBodyArmor: c.state.BodyArmorSlot,
		oas.EquipSchemaSlotLegArmor:  c.state.LegArmorSlot,
		oas.EquipSchemaSlotShield:    c.state.ShieldSlot,
		oas.EquipSchemaSlotHelmet:    c.state.HelmetSlot,
		oas.EquipSchemaSlotBoots:     c.state.BootsSlot,
		oas.EquipSchemaSlotRing1:     c.state.Ring1Slot,
		oas.EquipSchemaSlotRing2:     c.state.Ring2Slot,
		oas.EquipSchemaSlotAmulet:    c.state.AmuletSlot,
		oas.EquipSchemaSlotArtifact1: c.state.Artifact1Slot,
		oas.EquipSchemaSlotArtifact2: c.state.Artifact2Slot,
		oas.EquipSchemaSlotArtifact3: c.state.Artifact3Slot,
		oas.EquipSchemaSlotUtility1:  c.state.Utility1Slot,
		oas.EquipSchemaSlotUtility2:  c.state.Utility2Slot,
	}
}

func (c *Character) Utilities() map[oas.EquipSchemaSlot]int {
	return map[oas.EquipSchemaSlot]int{
		oas.EquipSchemaSlotUtility1: c.state.Utility1SlotQuantity,
		oas.EquipSchemaSlotUtility2: c.state.Utility2SlotQuantity,
	}
}

func (c *Character) Location() game.Point {
	return game.Point{
		X: c.state.X,
		Y: c.state.Y,
	}
}

type Task struct {
	Code    string
	Current int
	Total   int
	Type    string
}

var (
	NoTask = Task{}
)

func (c *Character) Task() Task {
	if c.state.Task == "" {
		return NoTask
	}

	return Task{
		Code:    c.state.Task,
		Current: c.state.TaskProgress,
		Total:   c.state.TaskTotal,
		Type:    c.state.TaskType,
	}
}

func (c *Character) Log(msg string, args ...any) {
	c.logger.Debug(msg, args...)
}

func (c *Character) Name() string {
	return c.name
}
