package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Subaru string = "Subaru"
)

func (c *liveCharacter) subaruStrategy() Strategy {
	skills := c.character.Skills()
	items := []string{}

	needIron := false

	if skills[string(oas.CraftSchemaSkillWeaponcrafting)] < 20 {
		items = append(items, "iron_sword", "fire_bow")
		needIron = true
	}

	if skills[string(oas.CraftSchemaSkillGearcrafting)] < 20 {
		items = append(items, "leather_armor", "iron_armor", "iron_helm")
		needIron = true
	}

	if skills[string(oas.CraftSchemaSkillJewelrycrafting)] < 20 {
		items = append(items, "iron_ring", "air_ring")
		needIron = true
	}

	if needIron {
		c.game.IntercomSet(c.character.Name(), "need_iron")
	} else {
		c.game.IntercomUnSet(c.character.Name(), "need_iron")
	}

	return strategy.SimpleCraft(c.character, c.game).Items(items...)
}
