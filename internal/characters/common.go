package characters

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
)

type Strategy interface {
	Name() string
	Do(ctx context.Context) error
}

type commonCharacter struct {
	character *generic.Character
	game      *game.Game
	strategy  Strategy
	logger    *slog.Logger
}

func (c *commonCharacter) Live(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := c.do(ctx)
		if err != nil {
			c.logger.Error(err.Error())
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *commonCharacter) do(ctx context.Context) error {
	new := c.getStrategy()

	if c.strategy == nil || c.strategy.Name() != new.Name() {
		c.logger.Info("change strategy to " + new.Name())
		c.strategy = new
	}

	return c.strategy.Do(ctx)
}

func (c *commonCharacter) getStrategy() Strategy {
	panic("unimplemented")
}
