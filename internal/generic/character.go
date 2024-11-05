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

type Point struct {
	Name string
	X, Y int
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
