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
		AllowSwitchTools().
		AllowEvents("strange_apparition", "magic_apparition").
		// default side resources from mining
		Deposit("topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire").
		// default side resources from woodcutting
		Deposit("sap", "apple", "maple_sap").
		// default side resources from fishing
		Deposit("algae").
		// event resources
		Deposit("strange_ore", "diamond_stone", "diamond", "magic_wood", "magic_sap")

	switch {
	case skills[string(oas.ResourceSchemaSkillMining)] < 40:
		return gather.Spot("gold_rocks").Craft("gold").Deposit("salmon")
	case skills[string(oas.ResourceSchemaSkillWoodcutting)] < 40:
		return gather.Spot("dead_tree").Deposit("gold").Craft("dead_wood_plank")
	case c.game.GetAchievment("Professional Fisherman").Current < c.game.GetAchievment("Professional Fisherman").Total:
		return gather.Spot("salmon_fishing_spot").Deposit("dead_wood_plank", "salmon")
	}

	return strategy.Empty()
}
