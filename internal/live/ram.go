package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Ram string = "Ram"
)

func (c *liveCharacter) ramStrategy() Strategy {
	skills := c.character.Skills()

	gather := strategy.SimpleGather(c.character, c.game).
		DepositGold().
		AllowSwitchTools().
		Deposit(
			"iron_ore",
			"topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire",
			"coal", "birch_wood", "sap", "trout", "algae", "nettle_leaf", "health_potion",
			"bass", "algae", "glowstem_leaf",
		).
		AllowEvents("Strange Apparition", "Magic Apparition").
		Deposit("strange_ore", "diamond_stone", "diamond", "magic_wood", "magic_sap").
		Craft("gold", "dead_wood_plank")

	switch {
	case skills[string(oas.ResourceSchemaSkillAlchemy)] < 40:
		return strategy.SimpleCraft(c.character, c.game).Items("health_potion")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 40 || c.game.GetAchievment("Expert Fisherman").Current < c.game.GetAchievment("Expert Fisherman").Total:
		return gather.Spot("bass_fishing_spot")
	case c.game.GetAchievment("Expert Alchemist").Current < c.game.GetAchievment("Expert Alchemist").Total:
		return gather.Spot("glowstem")
	case skills[string(oas.ResourceSchemaSkillMining)] < 40:
		return gather.Spot("gold_rocks")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 40:
		return gather.Spot("dead_tree")
	}

	return strategy.Empty()
}
