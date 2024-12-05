package strategy

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type simpleGather struct {
	character *generic.Character
	game      *game.Game

	spot            string
	craft           []string
	deposit         []string
	depositGold     bool
	allowSwithTools bool
	events          []string

	currentSpot string
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
	s.deposit = append(s.deposit, items...)
	return s
}

func (s *simpleGather) Deposit(items ...string) *simpleGather {
	s.deposit = append(s.deposit, items...)
	return s
}

func (s *simpleGather) DepositGold() *simpleGather {
	s.depositGold = true
	return s
}

func (s *simpleGather) AllowSwitchTools() *simpleGather {
	s.allowSwithTools = true
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

		macro.Deposit(ctx, s.character, s.game, s.deposit...)

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

	if err := s.switchTools(ctx, spot); err != nil {
		return fmt.Errorf("swith tools: %w", err)
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

func (s *simpleGather) switchTools(ctx context.Context, spot game.Point) error {
	if !s.allowSwithTools {
		return nil
	}

	if s.currentSpot == spot.Name {
		// already weared best tool
		return nil
	}

	s.game.LockBank()
	defer s.game.UnlockBank()

	start := time.Now()
	gear := macro.GetBestGearForResource(s.character, s.game, spot.Name)
	s.character.Log(fmt.Sprintf("choose best tool for resource %s: %v", spot.Name, time.Since(start)), slog.Any("items", gear))

	if err := macro.Wear(ctx, s.character, s.game, gear); err != nil {
		return fmt.Errorf("wear: %w", err)
	}

	s.currentSpot = spot.Name
	return nil
}
