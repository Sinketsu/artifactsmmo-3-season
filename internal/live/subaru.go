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

	// cooking
	if c.character.Skills()[string(oas.CraftSchemaSkillCooking)] < 10 {
		items = append(items, "cooked_chicken", "cooked_gudgeon")

		if c.character.Skills()[string(oas.CraftSchemaSkillCooking)] >= 5 {
			items = append(items, "fried_eggs")
		}
	}

	// weaponcrafting
	if c.character.Skills()[string(oas.CraftSchemaSkillWeaponcrafting)] < 10 {
		items = append(items, "copper_dagger")
	}

	// gearcrafting
	if c.character.Skills()[string(oas.CraftSchemaSkillGearcrafting)] < 10 {
		items = append(items, "copper_helmet", "wooden_shield")
	}

	// jewerlycrafting
	if c.character.Skills()[string(oas.CraftSchemaSkillJewelrycrafting)] < 10 {
		items = append(items, "copper_ring")
	}

	return strategy.SimpleCraft(c.character, c.game).Items(items...)
}
