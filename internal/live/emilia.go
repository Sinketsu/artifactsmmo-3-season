package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Emilia string = "Emilia"
)

func (c *liveCharacter) emiliaStrategy() Strategy {
	// quests := strategy.TasksMonsters(c.character, c.game).
	// 	AllowEvents("portal_demon", "bandit_camp", "snowman").
	// 	UseFood("cooked_salmon", "maple_syrup", "carrot", "cooked_wolf_meat").
	// 	Cancel("hellhound", "goblin_wolfrider", "lich", "bat", "goblin", "cultist_acolyte", "imp")

	fight := strategy.SimpleFight(c.character, c.game).
		AllowEvents("portal_demon", "bandit_camp", "snowman").
		UseFood("cooked_salmon", "maple_syrup", "carrot", "cooked_wolf_meat", "gingerbread").
		DepositGold().
		AllowSwitchGear().
		With("gingerbread")

	return fight
}
