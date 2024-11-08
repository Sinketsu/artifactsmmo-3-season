package macro

import (
	"context"
	"log/slog"
	"math"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

func CraftFromInventory(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	for _, code := range items {
		item, err := game.Item().Get(ctx, code)
		if err != nil {
			slog.Error("fail get item "+code, slog.Any("error", err))
			continue
		}

		if !item.Craft.IsSet() {
			slog.Warn("item " + code + " is not craftable...")
			continue
		}

		q := math.MaxInt64
		for _, res := range item.Craft.Value.CraftSchema.Items {
			q = min(q, character.Inventory()[code]/res.Quantity)
		}

		if q == 0 {
			continue
		}

		workshop, err := game.Map().Get(ctx, string(item.Craft.Value.CraftSchema.Skill.Value))
		if err != nil {
			slog.Error("fail to find workshop", slog.Any("error", err))
			continue
		}

		if err := character.Move(ctx, workshop); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			continue
		}

		if _, err := character.Craft(ctx, code, q); err != nil {
			slog.Error("fail to craft", slog.Any("error", err))
		}
	}
}

func Recycle(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	for _, code := range items {
		q := character.Inventory()[code]

		if q == 0 {
			continue
		}

		item, err := game.Item().Get(ctx, code)
		if err != nil {
			slog.Error("fail get item "+code, slog.Any("error", err))
			continue
		}

		if !item.Craft.IsSet() {
			slog.Warn("item " + code + " is not recyclable...")
			continue
		}

		workshop, err := game.Map().Get(ctx, string(item.Craft.Value.CraftSchema.Skill.Value))
		if err != nil {
			slog.Error("fail to find workshop", slog.Any("error", err))
			continue
		}

		if err := character.Move(ctx, workshop); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			continue
		}

		if _, err := character.Recycle(ctx, code, q); err != nil {
			slog.Error("fail to recycle", slog.Any("error", err))
		}
	}
}

func Deposit(ctx context.Context, character *generic.Character, items ...string) {
	if len(items) == 0 {
		return
	}

	if err := character.Move(ctx, game.Bank); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	for code, q := range character.Inventory() {
		if err := character.Deposit(ctx, code, q); err != nil {
			slog.Error("fail to deposit item "+code, slog.Any("error", err))
			continue
		}
	}
}

func DepositGold(ctx context.Context, character *generic.Character) {
	if err := character.Move(ctx, game.Bank); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	if q := character.Gold(); q > 0 {
		if err := character.DepositGold(ctx, q); err != nil {
			slog.Error("fail to deposit gold", slog.Any("error", err))
		}
	}
}

func CraftFromBank(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	if len(items) == 0 {
		return
	}

	bank := game.Bank().Items()

	for _, code := range items {
		item, err := game.Item().Get(ctx, code)
		if err != nil {
			slog.Error("fail get item "+code, slog.Any("error", err))
			continue
		}

		if !item.Craft.IsSet() {
			slog.Warn("item " + code + " is not craftable...")
			continue
		}

		q := math.MaxInt64
		resourcesForOneCraft := 0
		for _, res := range item.Craft.Value.CraftSchema.Items {
			q = min(q, bank[code]/res.Quantity)
			resourcesForOneCraft += res.Quantity
		}

		inventorySpace, _ := character.InventorySpace()
		q = min(q, inventorySpace/resourcesForOneCraft)

		if q == 0 {
			continue
		}

		bankLocation, _ := game.Map().Get(ctx, "bank")
		if err := character.Move(ctx, bankLocation); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			return
		}

		for _, res := range item.Craft.Value.CraftSchema.Items {
			if err := character.Withdraw(ctx, res.Code, res.Quantity*q); err != nil {
				slog.Error("fail to withdraw", slog.Any("error", err))
				return
			}
		}

		workshop, err := game.Map().Get(ctx, string(item.Craft.Value.CraftSchema.Skill.Value))
		if err != nil {
			slog.Error("fail to find workshop", slog.Any("error", err))
			continue
		}

		if err := character.Move(ctx, workshop); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			continue
		}

		if _, err := character.Craft(ctx, code, q); err != nil {
			slog.Error("fail to craft", slog.Any("error", err))
		}
	}
}
