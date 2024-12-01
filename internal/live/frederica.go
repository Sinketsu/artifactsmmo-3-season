package live

import (
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
)

const (
	Frederica string = "Frederica"
)

func (c *liveCharacter) fredericaStrategy() Strategy {
	fight := strategy.SimpleFight(c.character, c.game).
		Deposit("skeleton_bone", "skeleton_skull", "pig_skin", "ogre_eye", "ogre_skin", "wooden_club", "cyclops_eye").
		Deposit("red_cloth", "gingerbread", "gift").
		DepositGold().
		AllowSwitchGear().
		UseFood("cooked_bass", "cooked_salmon", "gingerbread")

	return fight.With("gingerbread")
}
