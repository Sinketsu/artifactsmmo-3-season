package game

import (
	"os"

	"github.com/Sinketsu/artifactsmmo-3-season/gen/oas"
	"github.com/Sinketsu/artifactsmmo-3-season/internal/api"
)

type Game struct {
	maps        *mapService
	items       *itemService
	resources   *resourceService
	monsters    *monsterService
	bank        *bankService
	achievments *achievmentService
	events      *eventService
	intercom    *intercomService
	ge          *geService
}

func New(client *api.Client) *Game {
	g := &Game{
		maps:        newMapService(client),
		items:       newItemService(client),
		resources:   newResourceService(client),
		monsters:    newMonsterService(client),
		bank:        newBankService(client),
		achievments: newAchievmentService(client, os.Getenv("SERVER_ACCOUNT")),
		events:      newEventService(client),
		intercom:    newIntercomService(),
		ge:          newGEService(client),
	}

	return g
}

func (g *Game) BankLocation(closestTo Point) Point {
	location, _ := g.maps.get("bank", closestTo)
	return location
}

func (g *Game) GrandExchangeLocation(closestTo Point) Point {
	location, _ := g.maps.get("grand_exchange", closestTo)
	return location
}

func (g *Game) TaskMasterItemsLocation(closestTo Point) Point {
	location, _ := g.maps.get("items", closestTo)
	return location
}

func (g *Game) TaskMasterMonstersLocation(closestTo Point) Point {
	location, _ := g.maps.get("monsters", closestTo)
	return location
}

func (g *Game) Find(code string, closestTo Point) (Point, error) {
	return g.maps.get(code, closestTo)
}

func (g *Game) GetItem(code string) (oas.ItemSchema, error) {
	return g.items.get(code)
}

func (g *Game) BankItems() map[string]int {
	return g.bank.Items()
}

func (g *Game) SyncBank() {
	g.bank.sync()
}

func (g *Game) LockBank() {
	g.bank.lock()
}

func (g *Game) UnlockBank() {
	g.bank.unlock()
}

func (g *Game) GetEvent(code string) (Point, error) {
	return g.events.get(code)
}

func (g *Game) GetAchievment(name string) achievment {
	return g.achievments.get(name)
}

func (g *Game) GetResource(code string) (oas.ResourceSchema, error) {
	return g.resources.get(code)
}

func (g *Game) GetMonster(code string) (oas.MonsterSchema, error) {
	return g.monsters.get(code)
}

func (g *Game) IntercomSet(character string, name string) {
	g.intercom.Set(character, name)
}

func (g *Game) IntercomUnSet(character string, name string) {
	g.intercom.UnSet(character, name)
}

func (g *Game) IntercomGet(character string, name string) bool {
	return g.intercom.Get(character, name)
}

func (g *Game) GEOrders(item string) []Order {
	return g.ge.Get(item)
}

func (g *Game) SyncGE() {
	g.ge.sync()
}
