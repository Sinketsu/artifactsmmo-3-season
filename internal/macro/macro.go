package macro

import (
	"context"
	"log/slog"
	"math"
	"slices"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

func CraftFromInventory(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	for _, code := range items {
		item, err := game.GetItem(code)
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
			q = min(q, character.Inventory()[res.Code]/res.Quantity)
		}

		if q == 0 {
			continue
		}

		workshop, err := game.Find(string(item.Craft.Value.CraftSchema.Skill.Value), character.Location())
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

		item, err := game.GetItem(code)
		if err != nil {
			slog.Error("fail get item "+code, slog.Any("error", err))
			continue
		}

		if !item.Craft.IsSet() ||
			item.Type == string(oas.GetAllItemsItemsGetTypeConsumable) ||
			item.Type == string(oas.GetAllItemsItemsGetTypeUtility) {
			// these types are not recyclable
			continue
		}

		workshop, err := game.Find(string(item.Craft.Value.CraftSchema.Skill.Value), character.Location())
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

func Deposit(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	if len(items) == 0 {
		return
	}

	if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	needSync := false
	defer func() {
		if needSync {
			game.SyncBank()
		}
	}()

	for code, q := range character.Inventory() {
		if !slices.Contains(items, code) {
			continue
		}

		if err := character.Deposit(ctx, code, q); err != nil {
			slog.Error("fail to deposit item "+code, slog.Any("error", err))
			continue
		}
		needSync = true
	}
}

func DepositGold(ctx context.Context, character *generic.Character, game *game.Game) {
	if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	needSync := false
	defer func() {
		if needSync {
			game.SyncBank()
		}
	}()

	if q := character.Gold(); q > 0 {
		if err := character.DepositGold(ctx, q); err != nil {
			slog.Error("fail to deposit gold", slog.Any("error", err))
		}
		needSync = true
	}
}

func CraftFromBank(ctx context.Context, character *generic.Character, game *game.Game, items ...string) {
	if len(items) == 0 {
		return
	}

	bank := game.BankItems()
	needSync := false
	defer func() {
		if needSync {
			game.SyncBank()
		}
	}()

	for _, code := range items {
		item, err := game.GetItem(code)
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
			q = min(q, bank[res.Code]/res.Quantity)
			resourcesForOneCraft += res.Quantity
		}

		inventorySpace, _ := character.InventorySpace()
		q = min(q, inventorySpace/resourcesForOneCraft)

		if q == 0 {
			continue
		}

		if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			return
		}
		needSync = true

		for _, res := range item.Craft.Value.CraftSchema.Items {
			if err := character.Withdraw(ctx, res.Code, res.Quantity*q); err != nil {
				slog.Error("fail to withdraw", slog.Any("error", err))
				return
			}
		}

		workshop, err := game.Find(string(item.Craft.Value.CraftSchema.Skill.Value), character.Location())
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

func AcceptItemTask(ctx context.Context, character *generic.Character, game *game.Game) {
	if err := character.Move(ctx, game.TaskMasterItemsLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	if _, err := character.AcceptNewTask(ctx); err != nil {
		slog.Error("fail to accept new task", slog.Any("error", err))
		return
	}
}

func TradeItemTask(ctx context.Context, character *generic.Character, game *game.Game, code string, quantity int) {
	if err := character.Move(ctx, game.TaskMasterItemsLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	if err := character.TaskTrade(ctx, code, quantity); err != nil {
		slog.Error("fail to trade task", slog.Any("error", err))
		return
	}
}

func CompleteTask(ctx context.Context, character *generic.Character, game *game.Game) {
	if err := character.Move(ctx, game.TaskMasterItemsLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	if _, err := character.CompleteTask(ctx); err != nil {
		slog.Error("fail to complete task", slog.Any("error", err))
		return
	}
}

func Buy(ctx context.Context, character *generic.Character, game *game.Game, id string, quantity int, totalPrice int) {
	if character.Gold() < totalPrice {
		if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
			slog.Error("fail to move", slog.Any("error", err))
			return
		}

		if err := character.WithdrawGold(ctx, totalPrice-character.Gold()); err != nil {
			slog.Error("fail to withdraw gold", slog.Any("error", err))
			return
		}
	}

	if err := character.Move(ctx, game.GrandExchangeLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	order, err := character.Buy(ctx, id, quantity)
	if err != nil {
		slog.Error("fail to buy", slog.Any("error", err))
		return
	}

	if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
		slog.Error("fail to move", slog.Any("error", err))
		return
	}

	if err := character.Deposit(ctx, order.Code, order.Quantity); err != nil {
		slog.Error("fail to deposit", slog.Any("error", err))
	}
}
