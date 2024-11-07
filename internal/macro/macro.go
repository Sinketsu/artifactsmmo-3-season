package macro

import (
	"context"
	"log/slog"
	"math"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

func CraftAll(ctx context.Context, character *generic.Character, game *game.Game, allowBank bool, items ...string) error {
	for _, code := range items {
		item, err := game.Item().Get(ctx, code)
		if err != nil {
			slog.Error("fail get item", slog.Any("error", err))
			continue
		}

		if !item.Craft.IsSet() {
			slog.Warn("item " + code + " is not craftable...")
			continue
		}

		q := math.MaxInt64
		for _, res := range item.Craft.Value.CraftSchema.Items {
			q = min(q, character.InInventory(code)/res.Quantity)
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

	return nil
}
