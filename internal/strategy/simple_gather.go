package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type simpleGather struct {
	character *generic.Character
	game      *game.Game

	spot        string
	craft       []string
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

func (s *simpleGather) Craft(items ...string) *simpleGather {
	s.craft = items
	s.deposit = append(s.deposit, items...)
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
		macro.CraftFromInventory(ctx, s.character, s.game, s.craft...)

		macro.Deposit(ctx, s.character, s.deposit...)

		if s.depositGold {
			macro.DepositGold(ctx, s.character)
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
