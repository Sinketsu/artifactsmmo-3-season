package macro

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

const (
	maxFoodCount = 60
)

func Heal(ctx context.Context, character *generic.Character, game *game.Game, food ...string) error {
	if character.HealthPercent() >= 60 {
		return nil
	}

	inventory := character.Inventory()

	for _, meal := range food {
		if inventory[meal] == 0 {
			continue
		}

		item, err := game.GetItem(meal)
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

		hp, maxHp := character.Health()

		useCount := int(math.Ceil(float64(maxHp-hp) / float64(healEffect)))

		err = character.Use(ctx, meal, min(useCount, inventory[meal]))
		if err != nil {
			return fmt.Errorf("use: %w", err)
		}
	}

	hp, maxHp := character.Health()

	if hp < maxHp {
		bank := game.BankItems()

		space, _ := character.InventorySpace()
		limit := min(maxFoodCount, space)

		for _, meal := range food {
			if bank[meal] == 0 {
				continue
			}

			err := character.Move(ctx, game.BankLocation(character.Location()))
			if err != nil {
				return fmt.Errorf("move: %w", err)
			}

			count := min(bank[meal], limit)
			err = character.Withdraw(ctx, meal, count)
			if err != nil {
				return fmt.Errorf("withdraw: %w", err)
			}

			limit -= count
		}

		return character.Rest(ctx)
	}

	return nil
}

func SwitchGear(ctx context.Context, character *generic.Character, game *game.Game, monster game.Point) error {
	game.LockBank()
	defer game.UnlockBank()

	start := time.Now()
	gear := GetBestGearForMonster(character, game, monster.Name)
	character.Log(fmt.Sprintf("choose best gear for monster %s: %v", monster.Name, time.Since(start)), slog.Any("items", gear))

	if err := Wear(ctx, character, game, gear); err != nil {
		return fmt.Errorf("wear: %w", err)
	}

	return nil
}
