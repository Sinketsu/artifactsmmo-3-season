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
		AllowEvents("rosenblood").
		UseFood("cooked_salmon").
		WithGearLevelDelta(5)

	return fight.With("hellhound")
}
