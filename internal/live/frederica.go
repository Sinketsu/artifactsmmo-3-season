package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Frederica string = "Frederica"
)

func (c *liveCharacter) fredericaStrategy() Strategy {
	fight := strategy.SimpleFight(c.character, c.game).
		DepositGold().
		AllowSwitchGear().
		AllowEvents("portal_demon", "bandit_camp", "snowman").
		// event resources
		Deposit("demon_horn", "piece_of_obsidian", "bandit_armor", "lizard_skin", "carrot", "snowman_hat", "gift").
		UseFood("cooked_bass", "cooked_wolf_meat")

	return fight.With("imp").Deposit("demoniac_dust", "piece_of_obsidian")
}
