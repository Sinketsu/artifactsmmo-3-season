package live

import "github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"

const (
	Emilia string = "Emilia"
)

func (c *liveCharacter) emiliaStrategy() Strategy {
	return strategy.SimpleFight(c.character, c.game).
		With("chicken").
		Deposit("feather", "raw_chicken", "egg").
		DepositGold()
}
