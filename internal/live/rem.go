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

	if c.game.IntercomGet("Subaru", "need_iron") {
		return strategy.SimpleGather(c.character, c.game).DepositGold().Deposit("topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire", "air_boost_potion", "fire_boost_potion").AllowSwitchTools().Craft("iron").Spot("iron_rocks")
	}

	gather := strategy.SimpleGather(c.character, c.game).
		DepositGold().
		AllowSwitchTools().
		Deposit(
			"topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire",
			"coal", "birch_wood", "sap", "bass", "algae", "nettle_leaf",
		)

	switch {
	case skills[string(oas.ResourceSchemaSkillMining)] < 30 || c.game.GetAchievment("Expert Miner").Current < c.game.GetAchievment("Expert Miner").Total:
		return gather.Spot("coal_rocks")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 30 || c.game.GetAchievment("Intermediate Lumberjack").Current < c.game.GetAchievment("Intermediate Lumberjack").Total:
		return gather.Spot("birch_tree")
	case skills[string(oas.ResourceSchemaSkillFishing)] < 30 || c.game.GetAchievment("Intermediate Fisherman").Current < c.game.GetAchievment("Intermediate Fisherman").Total:
		return gather.Spot("trout_fishing_spot")
	case skills[string(oas.ResourceSchemaSkillAlchemy)] < 30 || c.game.GetAchievment("Intermediate Alchemist").Current < c.game.GetAchievment("Intermediate Alchemist").Total:
		return gather.Spot("nettle")
	default:
		return gather.Spot("iron_rocks")
	}
}
