package macro

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

func SwitchTools(ctx context.Context, character *generic.Character, game *game.Game, spot game.Point) error {
	game.LockBank()
	defer game.UnlockBank()

	start := time.Now()
	gear := GetBestGearForResource(character, game, spot.Name)
	character.Log(fmt.Sprintf("choose best gear for resource %s: %v", spot.Name, time.Since(start)), slog.Any("items", gear))

	if err := Wear(ctx, character, game, gear); err != nil {
		return fmt.Errorf("wear: %w", err)
	}

	return nil
}
