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
		Deposit(
			"topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire",
			"gudgeon", "algae", "shrimp", "sap", "apple",
		).
		Craft("copper", "iron", "ash_plank", "spruce_plank", "small_health_potion")

	switch {
	case skills[string(oas.ResourceSchemaSkillMining)] < 10 || c.game.GetAchievment("Amateur Miner").Current < c.game.GetAchievment("Amateur Miner").Total:
		return gather.Spot("copper_rocks")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 10 || c.game.GetAchievment("Amateur Lumberjack").Current < c.game.GetAchievment("Amateur Lumberjack").Total:
		return gather.Spot("ash_tree")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 10 || c.game.GetAchievment("Amateur Fisherman").Current < c.game.GetAchievment("Amateur Fisherman").Total:
		return gather.Spot("gudgeon_fishing_spot")
	case skills[string(oas.ResourceSchemaSkillAlchemy)] < 10:
		return gather.Spot("sunflower")
	case skills[string(oas.ResourceSchemaSkillMining)] < 20 || c.game.GetAchievment("Miner").Current < c.game.GetAchievment("Miner").Total:
		return gather.Spot("iron_rocks")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 20 || c.game.GetAchievment("Lumberjack").Current < c.game.GetAchievment("Lumberjack").Total:
		return gather.Spot("spruce_tree")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 20 || c.game.GetAchievment("Fisherman").Current < c.game.GetAchievment("Fisherman").Total:
		return gather.Spot("shrimp_fishing_spot")
	case skills[string(oas.ResourceSchemaSkillAlchemy)] < 20 || c.game.GetAchievment("Amateur Alchemist").Current < c.game.GetAchievment("Amateur Alchemist").Total:
		return gather.Spot("sunflower")
	default:
		return gather.Spot("copper_rocks")
	}
}
