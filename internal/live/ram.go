package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Ram string = "Ram"
)

func (c *liveCharacter) ramStrategy() Strategy {
	gather := strategy.SimpleGather(c.character, c.game).DepositGold()

	if c.game.BankItems()["cooked_salmon"] < 1000 {
		return gather.Spot("salmon_fishing_spot")
	}

	return strategy.TasksMonsters(c.character, c.game)
}
