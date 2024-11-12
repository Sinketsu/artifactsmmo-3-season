package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Subaru string = "Subaru"
)

func (c *liveCharacter) subaruStrategy() Strategy {
	items := []string{}

	// jewerlycrafting
	if c.character.Skills()[string(oas.CraftSchemaSkillJewelrycrafting)] < 10 {
		items = append(items, "copper_ring", "life_amulet")
	}

	// cooking
	if c.character.Skills()[string(oas.CraftSchemaSkillCooking)] < 20 {
		items = append(items, "cheese")
	}

	return strategy.SimpleCraft(c.character, c.game).Items(items...)
}
