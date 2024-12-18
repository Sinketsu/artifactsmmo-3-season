package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Ram string = "Ram"
)

func (c *liveCharacter) ramStrategy() Strategy {
	quests := strategy.TasksItems(c.character, c.game).
		// AllowEvents("strange_apparition", "magic_apparition").
		Cancel("steel", "hardwood_plank", "strange_ore", "magic_wood", "strangold", "magical_plank",
			"cooked_salmon", "cooked_bass", "cooked_trout")

	return quests
}
