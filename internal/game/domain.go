package game

import "github.com/Sinketsu/artifactsmmo-3-season/internal/api"

type Game struct {
	maps  *maps
	items *items
}

func New(client *api.Client) *Game {
	return &Game{
		maps:  newMaps(client),
		items: newItems(client),
	}
}

func (g *Game) Map() *maps {
	return g.maps
}

func (g *Game) Item() *items {
	return g.items
}
