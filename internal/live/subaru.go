package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Subaru string = "Subaru"
)

var (
	buyMap = map[string]int{
		// resources
		"yellow_slimeball": 10,
		"red_slimeball":    10,
		"blue_slimeball":   10,
		"green_slimeball":  10,

		// for recycle
		"copper_dagger":        40,
		"fire_staff":           20,
		"sticky_dagger":        20,
		"sticky_sword":         20,
		"water_bow":            20,
		"greater_wooden_staff": 20,
		"fire_bow":             20,
		"copper_boots":         20,
		"iron_dagger":          20,
		"iron_sword":           20,
		"copper_helmet":        20,
		"copper_legs_armor":    20,
		"feather_coat":         20,
		"copper_armor":         20,
		"iron_boots":           20,
		"leather_boots":        20,
		"leather_armor":        20,
		"iron_armor":           20,
		"iron_legs_armor":      20,
		"iron_helm":            20,
		"leather_hat":          20,
		"wooden_shield":        20,
	}
)

func (c *liveCharacter) subaruStrategy() Strategy {
	// skills := c.character.Skills()

	craft := strategy.SimpleCraft(c.character, c.game).
		Buy(buyMap).
		Items("cooked_bass", "cooked_salmon")

	// if skills[string(oas.CraftSchemaSkillGearcrafting)] < 30 {
	// 	craft = craft.Items("steel_armor", "steel_legs_armor", "steel_boots", "steel_helm", "skeleton_armor", "tromatising_mask")
	// }

	return craft
}
