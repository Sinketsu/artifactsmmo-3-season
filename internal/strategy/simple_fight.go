package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

type simpleFight struct {
	character *generic.Character
	game      *game.Game

	monster     string
	deposit     []string
	depositGold bool
}

func SimpleFight(character *generic.Character, game *game.Game) *simpleFight {
	return &simpleFight{
		character: character,
		game:      game,
	}
}

func (s *simpleFight) With(monster string) *simpleFight {
	s.monster = monster
	return s
}

func (s *simpleFight) Deposit(items ...string) *simpleFight {
	s.deposit = items
	return s
}

func (s *simpleFight) DepositGold() *simpleFight {
	s.depositGold = true
	return s
}

func (s *simpleFight) Name() string {
	return "fight with " + s.monster
}

func (s *simpleFight) Do(ctx context.Context) error {
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

	monster, err := s.game.Map().Get(ctx, s.monster)
	if err != nil {
		return fmt.Errorf("get map: %w", err)
	}

	// TODO choose the gear

	err = s.character.Move(ctx, monster)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Fight(ctx)
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	if s.character.HealthPercent() < 60 {
		if err := s.character.Rest(ctx); err != nil {
			return fmt.Errorf("rest: %w", err)
		}
	}

	return nil
}
