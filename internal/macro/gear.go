package macro

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/simulator"
	combinations "github.com/mxschmitt/golang-combinations"
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

func GetBestGearForMonster(character *generic.Character, game *game.Game, code string) []oas.ItemSchema {
	monster, err := game.GetMonster(code)
	if err != nil {
		slog.Error("fail to get monster: " + code)
		return nil
	}

	bank := game.BankItems()
	all := map[oas.EquipSchemaSlot][]oas.ItemSchema{}

	for itemCode, q := range bank {
		item, err := game.GetItem(itemCode)
		if err != nil {
			slog.Error("fail to get item: " + itemCode)
			continue
		}

		if item.Level > character.Level() || item.Subtype == "tool" {
			continue
		}

		switch item.Type {
		case "ring":
			all[oas.EquipSchemaSlot(item.Type)] = append(all[oas.EquipSchemaSlot(item.Type)], item)
			if q > 1 {
				all[oas.EquipSchemaSlot(item.Type)] = append(all[oas.EquipSchemaSlot(item.Type)], item)
			}
		default:
			all[oas.EquipSchemaSlot(item.Type)] = append(all[oas.EquipSchemaSlot(item.Type)], item)
		}
	}

	for slot, code := range character.Equiped() {
		if code == "" || slot == oas.EquipSchemaSlotUtility1 || slot == oas.EquipSchemaSlotUtility2 {
			continue
		}

		// assume all equiped items are exist
		item, _ := game.GetItem(code)
		all[oas.EquipSchemaSlot(item.Type)] = append(all[oas.EquipSchemaSlot(item.Type)], item)
	}

	for slot := range all {
		slices.SortFunc(all[slot], func(a, b oas.ItemSchema) int {
			return strings.Compare(a.Code, b.Code)
		})
		if slot == "ring" {
			all[slot] = compact(all[slot], 2)
		} else {
			all[slot] = slices.CompactFunc(all[slot], func(a, b oas.ItemSchema) bool {
				return a.Code == b.Code
			})
		}
	}

	weapons := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotWeapon], 1))
	bodyArmors := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotBodyArmor], 1))
	legsArmors := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotLegArmor], 1))
	shields := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotShield], 1))
	amulets := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotAmulet], 1))
	boots := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotBoots], 1))
	rings := omitEmpty(combinations.Combinations(all["ring"], 2))
	helmets := omitEmpty(combinations.Combinations(all[oas.EquipSchemaSlotHelmet], 1))
	artifacts := omitEmpty(combinations.Combinations(all["artifact"], 3))

	sim := simulator.New()

	best := []oas.ItemSchema(nil)
	bestTime := 9999
	bestRemainingHp := 0

	for _, w := range weapons {
		for _, ba := range bodyArmors {
			for _, la := range legsArmors {
				for _, s := range shields {
					for _, a := range amulets {
						for _, r := range rings {
							for _, b := range boots {
								for _, h := range helmets {
									for _, ar := range artifacts {
										items := flatten(w, ba, la, s, a, r, b, h, ar)
										result := sim.Fight(character, items, monster)

										if !result.Win {
											continue
										}

										if result.Seconds < bestTime {
											bestTime = result.Seconds
											best = items
											bestRemainingHp = result.RemainingCharacterHp
											continue
										}
										if result.Seconds == bestTime && result.RemainingCharacterHp > bestRemainingHp {
											bestTime = result.Seconds
											best = items
											bestRemainingHp = result.RemainingCharacterHp
											continue
										}
										// may be additional criteria
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return best
}

func omitEmpty(items [][]oas.ItemSchema) [][]oas.ItemSchema {
	if len(items) == 0 {
		return [][]oas.ItemSchema{{{}}}
	}

	return items
}

func flatten(sets ...[]oas.ItemSchema) (result []oas.ItemSchema) {
	for _, s := range sets {
		result = append(result, s...)
	}

	return
}

func compact(items []oas.ItemSchema, n int) []oas.ItemSchema {
	if len(items) < n+1 {
		return items
	}

	result := make([]oas.ItemSchema, 0)
	last := ""
	lastCount := 0

	for i := 0; i < len(items); i++ {
		if last != items[i].Code {
			result = append(result, items[i])
			last = items[i].Code
			lastCount = 1
		} else {
			if lastCount < n {
				result = append(result, items[i])
				lastCount++
			}
		}
	}

	return result
}

func Wear(ctx context.Context, character *generic.Character, game *game.Game, items []oas.ItemSchema) error {
	if len(items) == 0 {
		return fmt.Errorf("wear: empty item list")
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
