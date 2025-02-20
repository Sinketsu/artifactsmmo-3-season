package strategy

import (
	"context"
	"fmt"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

type simpleFight struct {
	character *generic.Character
	game      *game.Game

	monster         string
	keep            []string
	depositGold     bool
	food            []string
	allowSwitchGear bool
	gearLevelDelta  int
	events          []string
	potions         []string

	current string
}

func SimpleFight(character *generic.Character, game *game.Game) *simpleFight {
	return &simpleFight{
		character: character,
		game:      game,

		gearLevelDelta: 15,
	}
}

func (s *simpleFight) With(monster string) *simpleFight {
	s.monster = monster
	return s
}

func (s *simpleFight) Keep(items ...string) *simpleFight {
	s.keep = append(s.keep, items...)
	return s
}

func (s *simpleFight) DepositGold() *simpleFight {
	s.depositGold = true
	return s
}

func (s *simpleFight) UseFood(food ...string) *simpleFight {
	s.food = append(s.food, food...)
	s.keep = append(s.keep, food...)
	return s
}

func (s *simpleFight) AllowSwitchGear() *simpleFight {
	s.allowSwitchGear = true
	return s
}

func (s *simpleFight) AllowEvents(events ...string) *simpleFight {
	s.events = append(s.events, events...)
	return s
}

func (s *simpleFight) AllowUsePotions(potions ...string) *simpleFight {
	s.potions = append(s.potions, potions...)
	return s
}

func (s *simpleFight) WithGearLevelDelta(delta int) *simpleFight {
	s.gearLevelDelta = delta
	return s
}

func (s *simpleFight) Name() string {
	return "fight with " + s.monster
}

func (s *simpleFight) Do(ctx context.Context) error {
	if err := macro.Heal(ctx, s.character, s.game, s.food...); err != nil {
		return fmt.Errorf("heal: %w", err)
	}

	if s.character.InventoryFull() {
		macro.Deposit(ctx, s.character, s.game, s.keep...)

		if s.depositGold {
			macro.DepositGold(ctx, s.character, s.game)
		}
	}

	var monster game.Point
	for _, event := range s.events {
		if event, err := s.game.GetEvent(event); err == nil {
			monster = event
			break
		}
	}

	if monster.Name == "" {
		var err error
		monster, err = s.game.Find(s.monster, s.character.Location())
		if err != nil {
			return fmt.Errorf("get map: %w", err)
		}
	}

	// temporary part
	if len(s.potions) == 0 {
		for slot, count := range s.character.Utilities() {
			if count > 0 {
				if err := s.character.Move(ctx, s.game.BankLocation(s.character.Location())); err != nil {
					return fmt.Errorf("move: %w", err)
				}

				item := s.character.Equiped()[slot]

				if err := s.character.Unequip(ctx, oas.UnequipSchemaSlot(slot), count); err != nil {
					return fmt.Errorf("unequip utility: %w", err)
				}

				if err := s.character.Deposit(ctx, item, count); err != nil {
					return fmt.Errorf("deposit utility: %w", err)
				}
			}
		}
	}

	if s.allowSwitchGear && s.current != monster.Name {
		if err := macro.SwitchGear(ctx, s.character, s.game, monster, s.gearLevelDelta, s.potions...); err != nil {
			return fmt.Errorf("switch gear: %w", err)
		}
		s.current = monster.Name
	}

	err := s.character.Move(ctx, monster)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Fight(ctx)
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	return nil
}
