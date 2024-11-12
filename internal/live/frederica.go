package live

import "github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"

const (
	Frederica string = "Frederica"
)

func (c *liveCharacter) fredericaStrategy() Strategy {
	return strategy.SimpleFight(c.character, c.game).
		With("cow").
		Deposit("feather", "raw_chicken", "egg", "green_slimeball", "apple", "red_slimeball", "raw_beef", "milk_bucket", "cowhide").
		DepositGold()
}
