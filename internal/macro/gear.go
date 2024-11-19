package macro

import (
	"context"
	"log/slog"
	"slices"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

func GetBestGearForResource(character *generic.Character, game *game.Game, code string) []oas.ItemSchema {
	resource, err := game.GetResource(code)
	if err != nil {
		slog.Error("fail to get resource: " + code)
		return nil
	}

	bank := game.BankItems()
	candidates := []oas.ItemSchema{}

	for itemCode := range bank {
		item, err := game.GetItem(itemCode)
		if err != nil {
			slog.Error("fail to get item: " + itemCode)
			continue
		}

		if item.Subtype == "tool" && item.Level <= character.Level() {
			candidates = append(candidates, item)
		}
	}

	if weapon := character.Equiped()[oas.EquipSchemaSlotWeapon]; weapon != "" {
		item, err := game.GetItem(weapon)
		if err != nil {
			slog.Error("fail to get item: " + weapon)
		} else {
			if item.Subtype == "tool" {
				candidates = append(candidates, item)
			}
		}
	}

	if len(candidates) == 0 || len(candidates) == 1 {
		return candidates
	}

	slices.SortFunc(candidates, func(a, b oas.ItemSchema) int {
		aEffect, bEffect := 0, 0

		for _, effect := range a.Effects {
			if effect.Name == string(resource.Skill) {
				aEffect = effect.Value
				break
			}
		}

		for _, effect := range b.Effects {
			if effect.Name == string(resource.Skill) {
				bEffect = effect.Value
				break
			}
		}

		return aEffect - bEffect
	})

	return candidates[:1]
}

func Wear(ctx context.Context, character *generic.Character, game *game.Game, items []oas.ItemSchema) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		// TODO support multi slots
		if item.Type == "ring" || item.Type == "artifact" || item.Code == "utility" {
			continue
		}

		current := character.Equiped()[oas.EquipSchemaSlot(item.Type)]

		if current == item.Code {
			continue
		}

		if current != "" {
			if err := character.Move(ctx, game.BankLocation(character.Location())); err != nil {
				return err
			}

			if err := character.Unequip(ctx, oas.UnequipSchemaSlot(item.Type), 1); err != nil {
				return err
			}

			if err := character.Deposit(ctx, current, 1); err != nil {
				return err
			}
		}

		if err := character.Withdraw(ctx, item.Code, 1); err != nil {
			return err
		}

		if err := character.Equip(ctx, item.Code, oas.EquipSchemaSlot(item.Type), 1); err != nil {
			return err
		}
	}

	return nil
}
