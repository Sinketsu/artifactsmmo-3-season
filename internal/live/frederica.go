package live

import "github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"

const (
	Frederica string = "Frederica"
)

func (c *liveCharacter) fredericaStrategy() Strategy {
	fight := strategy.SimpleFight(c.character, c.game).
		Deposit("raw_wolf_meat", "wolf_bone", "wolf_hair", "blue_slimeball", "apple", "raw_chicken", "egg", "feather").
		DepositGold()

	return fight.With("blue_slime")
}
