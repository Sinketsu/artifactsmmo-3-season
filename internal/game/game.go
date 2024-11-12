package game

import (
	"context"
	"os"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
)

type Game struct {
	maps        *mapService
	items       *itemService
	bank        *bankService
	achievments *achievmentService
	events      *eventService

	Bank          Point
	GrandExchange Point
}

func New(client *api.Client) *Game {
	g := &Game{
		maps:        newMapService(client),
		items:       newItemService(client),
		bank:        newBankService(client),
		achievments: newAchievmentService(client, os.Getenv("SERVER_ACCOUNT")),
		events:      newEventService(client),
	}

	g.Bank, _ = g.Find(context.Background(), "bank")
	g.GrandExchange, _ = g.Find(context.Background(), "grand_exchange")

	return g
}

func (g *Game) Find(ctx context.Context, code string) (Point, error) {
	return g.maps.get(ctx, code)
}

func (g *Game) GetItem(ctx context.Context, code string) (oas.ItemSchema, error) {
	return g.items.get(ctx, code)
}

func (g *Game) BankItems() map[string]int {
	return g.bank.Items()
}

func (g *Game) SyncBank() {
	g.bank.sync()
}

func (g *Game) GetEvent(code string) (Point, error) {
	return g.events.get(code)
}

func (g *Game) GetAchievment(name string) achievment {
	return g.achievments.get(name)
}
