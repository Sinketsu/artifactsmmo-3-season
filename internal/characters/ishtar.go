package characters

import (
	"log/slog"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type ishtar struct {
	*commonCharacter
}

func NewIshtar(client *api.Client, game *game.Game) *ishtar {
	name := "Ishtar"

	return &ishtar{
		commonCharacter: &commonCharacter{
			character: generic.NewCharacter(name, client),
			game:      game,
			logger:    slog.Default().With(ycloggingslog.Stream, name),
		},
	}
}

func (c *ishtar) getStrategy() Strategy {
	skills := c.character.Skills()

	gather := strategy.SimpleGather(c.character, c.game).
		DepositGold().
		Craft("copper", "iron", "ash_plank", "spruce_plank", "small_health_potion").
		Deposit(
			"topaz_stone", "topaz", "emerald_stone", "emerald", "ruby_stone", "ruby", "sapphire_stone", "sapphire",
			"gudgeon", "algae", "shrimp", "sap", "apple", "sunflower",
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
