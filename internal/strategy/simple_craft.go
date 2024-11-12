package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
	"golang.org/x/exp/maps"
)

type simpleCraft struct {
	character *generic.Character
	game      *game.Game

	items []string
}

func SimpleCraft(character *generic.Character, game *game.Game) *simpleCraft {
	return &simpleCraft{
		character: character,
		game:      game,
	}
}

func (s *simpleCraft) Items(items ...string) *simpleCraft {
	s.items = append(s.items, items...)
	return s
}

func (s *simpleCraft) Name() string {
	return fmt.Sprintf("craft anything of %v", s.items)
}

func (s *simpleCraft) Do(ctx context.Context) error {
	macro.Recycle(ctx, s.character, s.game, s.items...)

	// deposit all items
	macro.Deposit(ctx, s.character, s.game, maps.Keys(s.character.Inventory())...)

	macro.CraftFromBank(ctx, s.character, s.game, s.items...)

	time.Sleep(1 * time.Second)
	return nil
}
