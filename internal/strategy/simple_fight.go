package strategy

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/macro"
)

var (
	maxFoodCount = 60
)

type simpleFight struct {
	character *generic.Character
	game      *game.Game

	monster         string
	deposit         []string
	depositGold     bool
	food            []string
	allowSwitchGear bool

	currentMonster string
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

func (s *simpleFight) UseFood(food ...string) *simpleFight {
	s.food = append(s.food, food...)
	return s
}

func (s *simpleFight) AllowSwitchGear() *simpleFight {
	s.allowSwitchGear = true
	return s
}

func (s *simpleFight) Name() string {
	return "fight with " + s.monster
}

func (s *simpleFight) Do(ctx context.Context) error {
	if s.character.HealthPercent() < 60 {
		if err := s.heal(ctx); err != nil {
			return fmt.Errorf("heal: %w", err)
		}
	}

	if s.character.InventoryFull() {
		macro.Deposit(ctx, s.character, s.game, s.deposit...)

		if s.depositGold {
			macro.DepositGold(ctx, s.character, s.game)
		}
	}

	monster, err := s.game.Find(s.monster, s.character.Location())
	if err != nil {
		return fmt.Errorf("get map: %w", err)
	}

	if err := s.switchGear(ctx, monster.Name); err != nil {
		return fmt.Errorf("swith tools: %w", err)
	}

	err = s.character.Move(ctx, monster)
	if err != nil {
		return fmt.Errorf("move: %w", err)
	}

	_, err = s.character.Fight(ctx)
	if err != nil {
		return fmt.Errorf("fight: %w", err)
	}

	return nil
}

func (s *simpleFight) heal(ctx context.Context) error {
	inventory := s.character.Inventory()

	for _, food := range s.food {
		if inventory[food] == 0 {
			continue
		}

		item, err := s.game.GetItem(food)
		if err != nil {
			continue
		}

		healEffect := 0
		for _, ef := range item.Effects {
			if ef.Name == "heal" {
				healEffect = ef.Value
			}
		}

		if healEffect == 0 {
			continue
		}

		hp, maxHp := s.character.Health()

		useCount := int(math.Ceil(float64(maxHp-hp) / float64(healEffect)))

		err = s.character.Use(ctx, food, min(useCount, inventory[food]))
		if err != nil {
			return fmt.Errorf("use: %w", err)
		}
	}

	hp, maxHp := s.character.Health()

	if hp < maxHp {
		bank := s.game.BankItems()

		space, _ := s.character.InventorySpace()
		limit := min(maxFoodCount, space)

		for _, food := range s.food {
			if bank[food] == 0 {
				continue
			}

			err := s.character.Move(ctx, s.game.BankLocation(s.character.Location()))
			if err != nil {
				return fmt.Errorf("move: %w", err)
			}

			count := min(bank[food], limit)
			err = s.character.Withdraw(ctx, food, count)
			if err != nil {
				return fmt.Errorf("withdraw: %w", err)
			}

			limit -= count
		}

		return s.character.Rest(ctx)
	}

	return nil
}

func (s *simpleFight) switchGear(ctx context.Context, monster string) error {
	if !s.allowSwitchGear {
		return nil
	}

	if s.currentMonster == monster {
		// already weared best gear
		return nil
	}

	s.game.LockBank()
	defer s.game.UnlockBank()

	start := time.Now()
	gear := macro.GetBestGearForMonster(s.character, s.game, monster)
	s.character.Log(fmt.Sprintf("choose best gear for monster %s: %v", monster, time.Since(start)), slog.Any("items", gear))

	if err := macro.Wear(ctx, s.character, s.game, gear); err != nil {
		return fmt.Errorf("wear: %w", err)
	}

	s.currentMonster = monster
	return nil
}
