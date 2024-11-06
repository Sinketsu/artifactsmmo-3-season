package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

type simpleGather struct {
	character *generic.Character
	game      *game.Game

	spot        string
	deposit     []string
	depositGold bool
}

func SimpleGather(character *generic.Character, game *game.Game) *simpleGather {
	return &simpleGather{
		character: character,
		game:      game,
	}
}

func (s *simpleGather) Spot(spot string) *simpleGather {
	s.spot = spot
	return s
}

func (s *simpleGather) Deposit(items ...string) *simpleGather {
	s.deposit = items
	return s
}

func (s *simpleGather) DepositGold() *simpleGather {
	s.depositGold = true
	return s
}

func (s *simpleGather) Name() string {
	return "gather " + s.spot
}

func (s *simpleGather) Do(ctx context.Context) error {
	if s.character.InventoryFull() {
		if err := s.character.Move(ctx, game.Bank); err != nil {
			return fmt.Errorf("move: %w", err)
		}

		for _, item := range s.deposit {
			if q := s.character.InInventory(item); q > 0 {
				if err := s.character.Deposit(ctx, item, q); err != nil {
					return fmt.Errorf("deposit: %w", err)
				}
			}
		}

		if s.depositGold && s.character.Gold() > 0 {
			if err := s.character.DepositGold(ctx, s.character.Gold()); err != nil {
				return fmt.Errorf("deposit gold: %w", err)
			}
		}
	}

	spot, err := s.game.Map().Get(ctx, s.spot)
	if err != nil {
		return fmt.Errorf("get map: %w", err)
	}

	// TODO choose the gear

	err = s.character.Move(ctx, spot)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Gather(ctx)
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	return nil
}
