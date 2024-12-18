package strategy

import (
	"context"
	"fmt"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type simpleCraft struct {
	character *generic.Character
	game      *game.Game

	items []string
	buy   map[string]int
}

func SimpleCraft(character *generic.Character, game *game.Game) *simpleCraft {
	return &simpleCraft{
		character: character,
		game:      game,

		buy: make(map[string]int),
	}
}

func (s *simpleCraft) Items(items ...string) *simpleCraft {
	s.items = append(s.items, items...)
	return s
}

func (s *simpleCraft) Buy(items map[string]int) *simpleCraft {
	for k, v := range items {
		s.buy[k] = v
	}

	return s
}

func (s *simpleCraft) Name() string {
	return fmt.Sprintf("craft anything of %v", s.items)
}

func (s *simpleCraft) Do(ctx context.Context) error {
	s.checkGE(ctx)

	macro.Recycle(ctx, s.character, s.game, s.items...)

	// deposit all items
	macro.Deposit(ctx, s.character, s.game)

	macro.CraftFromBank(ctx, s.character, s.game, s.items...)

	time.Sleep(1 * time.Second)
	return nil
}

func (s *simpleCraft) checkGE(ctx context.Context) {
	needSync := false

	for item, price := range s.buy {
		orders := s.game.GEOrders(item)

		for _, order := range orders {
			if order.Price <= price {
				macro.Buy(ctx, s.character, s.game, order.Id, order.Quantity, order.Price*order.Quantity)
				needSync = true
			}
		}
	}

	if needSync {
		s.game.SyncGE()
	}
}
