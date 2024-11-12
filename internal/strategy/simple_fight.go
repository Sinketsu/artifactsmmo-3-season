package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
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
	s.deposit = append(s.deposit, items...)
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
		macro.Deposit(ctx, s.character, s.game, s.deposit...)

		if s.depositGold {
			macro.DepositGold(ctx, s.character, s.game)
		}
	}

	monster, err := s.game.Find(ctx, s.monster)
	if err != nil {
		return fmt.Errorf("get map: %w", err)
	}

	// TODO choose the gear

	err = s.character.Move(ctx, monster)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if s.character.HealthPercent() < 60 {
		if err := s.character.Rest(ctx); err != nil {
			return fmt.Errorf("rest: %w", err)
		}
	}

	_, err = s.character.Fight(ctx)
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	return nil
}
