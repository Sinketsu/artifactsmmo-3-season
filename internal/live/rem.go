package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Rem string = "Rem"
)

func (c *liveCharacter) remStrategy() Strategy {
	skills := c.character.Skills()

	gather := strategy.SimpleGather(c.character, c.game).
		DepositGold().
		Craft("copper", "iron", "ash_plank", "spruce_plank", "small_health_potion").
		Deposit(
			"topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire",
			"gudgeon", "algae", "shrimp", "sap", "apple",
		)

	switch {
	case skills[string(oas.ResourceSchemaSkillMining)] < 10:
		return gather.Spot("copper_rocks")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 10:
		return gather.Spot("gudgeon_fishing_spot")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 10:
		return gather.Spot("ash_tree")
	case skills[string(oas.ResourceSchemaSkillAlchemy)] < 20:
		return gather.Spot("sunflower")
	case skills[string(oas.ResourceSchemaSkillMining)] < 20:
		return gather.Spot("iron_rocks")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 20:
		return gather.Spot("shrimp_fishing_spot")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 20:
		return gather.Spot("spruce_tree")
	default:
		return strategy.Empty()
	}
}
