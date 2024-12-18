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
	keep        []string
	depositGold bool
	events      []string

	current string
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
	s.craft = append(s.craft, items...)
	return s
}

func (s *simpleGather) Keep(items ...string) *simpleGather {
	s.keep = append(s.keep, items...)
	return s
}

func (s *simpleGather) DepositGold() *simpleGather {
	s.depositGold = true
	return s
}

func (s *simpleGather) AllowEvents(events ...string) *simpleGather {
	s.events = append(s.events, events...)
	return s
}

func (s *simpleGather) Name() string {
	return fmt.Sprintf("gather %s and do events %v if possible", s.spot, s.events)
}

func (s *simpleGather) Do(ctx context.Context) error {
	if s.character.InventoryFull() {
		macro.CraftFromInventory(ctx, s.character, s.game, s.craft...)

		macro.Deposit(ctx, s.character, s.game, s.keep...)

		if s.depositGold {
			macro.DepositGold(ctx, s.character, s.game)
		}
	}

	var spot game.Point
	for _, event := range s.events {
		if event, err := s.game.GetEvent(event); err == nil {
			spot = event
			break
		}
	}

	if spot.Name == "" {
		var err error
		spot, err = s.game.Find(s.spot, s.character.Location())
		if err != nil {
			return fmt.Errorf("get map: %w", err)
		}
	}

	if s.current != spot.Name {
		if err := macro.SwitchTools(ctx, s.character, s.game, spot); err != nil {
			return fmt.Errorf("switch tools: %w", err)
		}
		s.current = spot.Name
	}

	err := s.character.Move(ctx, spot)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Gather(ctx)
	if err != nil {
		return fmt.Errorf("gather: %w", err)
	}

	return nil
}
