package live

import (
	"context"
	"log/slog"
	"time"

	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/game"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/generic"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/strategy"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Strategy interface {
	Name() string
	Do(ctx context.Context) error
}

type selector func() Strategy

type liveCharacter struct {
	character *generic.Character
	game      *game.Game
	strategy  Strategy
	selector  selector
	logger    *slog.Logger
}

func Character(name string, client *api.Client, game *game.Game) *liveCharacter {
	c := &liveCharacter{
		character: generic.NewCharacter(name, client),
		game:      game,
		selector:  func() Strategy { return strategy.Empty() },
		logger:    slog.Default().With(ycloggingslog.Stream, name),
	}

	switch name {
	case Ram:
		c.selector = c.ramStrategy
	case Rem:
		c.selector = c.remStrategy
	case Emilia:
		c.selector = c.emiliaStrategy
	case Frederica:
		c.selector = c.fredericaStrategy
	case Subaru:
		c.selector = c.subaruStrategy
	}

	return c
}

func (c *liveCharacter) Live(ctx context.Context) {
	c.character.Init()

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

func (c *liveCharacter) do(ctx context.Context) error {
	new := c.getStrategy()

	if c.strategy == nil || c.strategy.Name() != new.Name() {
		c.logger.Info("change strategy to " + new.Name())
		c.strategy = new
	}

	return c.strategy.Do(ctx)
}

func (c *liveCharacter) getStrategy() Strategy {
	return c.selector()
}
