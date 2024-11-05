package game

import "github.com/Sinketsu/artifactsmmo-3-season/internal/api"

type Game struct {
	maps *maps
}

func New(client *api.Client) *Game {
	return &Game{
		maps: newMaps(client),
	}
}

func (g *Game) Map() *maps {
	return g.maps
}
