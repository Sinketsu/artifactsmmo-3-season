package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Emilia string = "Emilia"
)

func (c *liveCharacter) emiliaStrategy() Strategy {
	fight := strategy.SimpleFight(c.character, c.game).
		DepositGold().
		AllowSwitchGear().
		AllowEvents("portal_demon", "bandit_camp", "snowman").
		// event resources
		Deposit("demon_horn", "piece_of_obsidian", "bandit_armor", "lizard_skin", "carrot", "snowman_hat", "gift").
		UseFood("cooked_bass", "cooked_wolf_meat")

	return fight.With("death_knight").Deposit("death_knight_sword", "red_cloth")
}
