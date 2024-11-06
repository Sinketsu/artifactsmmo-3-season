package generic

import (
	"log/slog"
	"unsafe"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
	ycloggingslog "github.com/Sinketsu/yc-logging-slog"
)

type Character struct {
	name  string
	state oas.CharacterSchema

	cli    *api.Client
	logger *slog.Logger
}

func NewCharacter(name string, client *api.Client) *Character {
	return &Character{
		name: name,

		cli:    client,
		logger: slog.Default().With(ycloggingslog.Stream, name),
	}
}

func (c *Character) syncState(p unsafe.Pointer) {
	// tricky hack, because `ogen` generates different models for Character state from different methods instead of reusing one. But fields are the same - so we can cast it
	c.state = *(*oas.CharacterSchema)(p)
}

func (c *Character) InventoryFull() bool {
	slots := 0
	items := 0

	for _, slot := range c.state.Inventory {
		if slot.Code != "" {
			slots++
			items += slot.Quantity
		}
	}

	return c.state.InventoryMaxItems-5 < items || slots >= 19
}

func (c *Character) InInventory(code string) int {
	for _, slot := range c.state.Inventory {
		if slot.Code == code {
			return slot.Quantity
		}
	}
	return 0
}

func (c *Character) Gold() int {
	return c.state.Gold
}

func (c *Character) HealthPercent() float64 {
	return float64(c.state.Hp) / float64(c.state.MaxHp) * 100
}
